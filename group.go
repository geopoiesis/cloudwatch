package cloudwatch

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	iface "github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/enfipy/locker"

	"github.com/pkg/errors"
)

// Throttling and limits from http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/cloudwatch_limits.html
const (
	// The maximum rate of a GetLogEvents request is 10 requests per second per AWS account.
	readThrottle = time.Second / 10

	// The maximum rate of a PutLogEvents request is 5 requests per second per log stream.
	writeThrottle = time.Second / 5
)

// now is a function that returns the current time.Time. It's a variable so that
// it can be stubbed out in unit tests.
// var now = time.Now

type groupImpl struct {
	iface.CloudWatchLogsAPI
	groupName string
	locker    *locker.Locker
}

// NewGroup returns a new Group instance.
func NewGroup(client iface.CloudWatchLogsAPI, groupName string) Group {
	return &groupImpl{
		CloudWatchLogsAPI: client,
		groupName:         groupName,
		locker:            locker.Initialize(),
	}
}

func (g *groupImpl) Create(ctx context.Context, streamName string, opts ...CreateOption) (io.WriteCloser, error) {
	ret, err := g.create(ctx, streamName)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(ret)
	}

	go ret.start()
	return ret, nil
}

func (g *groupImpl) Name() string {
	return g.groupName
}

func (g *groupImpl) Open(ctx context.Context, streamName string) io.Reader {
	ret := &readerImpl{
		client:     g,
		ctx:        ctx,
		groupName:  aws.String(g.groupName),
		streamName: aws.String(streamName),
		throttle:   time.Tick(readThrottle),
	}

	go ret.start()
	return ret
}

func (g *groupImpl) create(ctx context.Context, streamName string) (*writerImpl, error) {
	ret := &writerImpl{
		client:     g,
		ctx:        ctx,
		events:     newEventsBuffer(),
		groupName:  aws.String(g.groupName),
		streamName: aws.String(streamName),
		throttle:   time.Tick(writeThrottle),
	}

	unlock := g.locker.Lock(streamName)
	defer unlock()

	description, err := g.DescribeLogStreamsWithContext(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(g.groupName),
		LogStreamNamePrefix: aws.String(streamName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get log stream description")
	}

	if len(description.LogStreams) > 0 {
		ret.sequenceToken = description.LogStreams[0].UploadSequenceToken
		return ret, nil
	}

	_, err = g.CreateLogStreamWithContext(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(g.groupName),
		LogStreamName: aws.String(streamName),
	})

	return ret, errors.Wrap(err, "could not create a log stream")
}
