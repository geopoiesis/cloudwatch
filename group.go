package cloudwatch

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	iface "github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
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
	client    iface.CloudWatchLogsAPI
	groupName string
}

// NewGroup returns a new Group instance.
func NewGroup(client iface.CloudWatchLogsAPI, groupName string) Group {
	return &groupImpl{
		client:    client,
		groupName: groupName,
	}
}

// Create creates a log stream in the group and returns an io.Writer for it.
func (g *groupImpl) Create(ctx context.Context, streamName string, opts ...CreateOption) (io.WriteCloser, error) {
	_, err := g.client.CreateLogStreamWithContext(ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(g.groupName),
		LogStreamName: aws.String(streamName),
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not create a log stream")
	}
	ret := &writerImpl{
		client:     g.client,
		ctx:        ctx,
		groupName:  aws.String(g.groupName),
		streamName: aws.String(streamName),
		throttle:   time.Tick(writeThrottle),
	}
	for _, opt := range opts {
		opt(ret)
	}

	go ret.start()
	return ret, nil
}

// Open returns an io.Reader to read from the log stream.
func (g *groupImpl) Open(ctx context.Context, streamName string) io.Reader {
	ret := &readerImpl{
		client:     g.client,
		ctx:        ctx,
		groupName:  aws.String(g.groupName),
		streamName: aws.String(streamName),
		throttle:   time.Tick(readThrottle),
	}

	go ret.start()
	return ret
}
