package cloudwatch

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/stretchr/testify/suite"
)

type writerTestSuite struct {
	suite.Suite

	api                   *mockAPI
	ctx                   context.Context
	groupName, streamName string
	sut                   io.WriteCloser
}

func (w *writerTestSuite) SetupTest() {
	w.api = new(mockAPI)
	w.ctx = context.Background()
	w.groupName = "groupName"
	w.streamName = "streamName"

	w.api.On(
		"DescribeLogStreamsWithContext",
		w.ctx,
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(w.groupName),
			LogStreamNamePrefix: aws.String(w.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.DescribeLogStreamsOutput{LogStreams: nil}, nil)

	w.api.On(
		"CreateLogStreamWithContext",
		w.ctx,
		&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(w.groupName),
			LogStreamName: aws.String(w.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.CreateLogStreamOutput{}, nil)

	group := NewGroup(w.api, w.groupName)
	writer, err := group.Create(w.ctx, w.streamName, freezeTime(time.Unix(1, 0)))
	w.NoError(err)
	w.sut = writer
}

func (w *writerTestSuite) TestLifecycle() {
	w.api.On(
		"PutLogEventsWithContext",
		w.ctx,
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{Message: aws.String("Hello\n"), Timestamp: aws.Int64(1000)},
				{Message: aws.String("World"), Timestamp: aws.Int64(1000)},
			},
			LogGroupName:  aws.String(w.groupName),
			LogStreamName: aws.String(w.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.PutLogEventsOutput{}, nil)

	n, err := io.WriteString(w.sut, "Hello\nWorld")
	w.NoError(err)
	w.Equal(11, n)
	w.NoError(w.sut.Close())
}

func (w *writerTestSuite) TestWriteRejected() {
	w.api.On(
		"PutLogEventsWithContext",
		w.ctx,
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{Message: aws.String("Hello\n"), Timestamp: aws.Int64(1000)},
				{Message: aws.String("World"), Timestamp: aws.Int64(1000)},
			},
			LogGroupName:  aws.String(w.groupName),
			LogStreamName: aws.String(w.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.PutLogEventsOutput{
		RejectedLogEventsInfo: &cloudwatchlogs.RejectedLogEventsInfo{
			TooOldLogEventEndIndex: aws.Int64(2),
		},
	}, nil)

	_, err := io.WriteString(w.sut, "Hello\nWorld")
	w.NoError(err)

	const expectedError = "log messages were rejected"
	w.EqualError(w.sut.(*writerImpl).flushAll(), expectedError)

	_, err = io.WriteString(w.sut, "Hello")
	w.EqualError(err, expectedError)
}

func (w *writerTestSuite) TestNewline() {
	w.api.On(
		"PutLogEventsWithContext",
		w.ctx,
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{Message: aws.String("Hello\n"), Timestamp: aws.Int64(1000)},
			},
			LogGroupName:  aws.String(w.groupName),
			LogStreamName: aws.String(w.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.PutLogEventsOutput{}, nil)

	n, err := io.WriteString(w.sut, "Hello\n")
	w.NoError(err)
	w.Equal(6, n)

	w.NoError(w.sut.Close())
}

func TestWriter(t *testing.T) {
	suite.Run(t, new(writerTestSuite))
}
