package cloudwatch

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/stretchr/testify/suite"
)

type groupTestSuite struct {
	suite.Suite
	api                   *mockAPI
	ctx                   context.Context
	groupName, streamName string
	sut                   Group
}

func (gs *groupTestSuite) SetupTest() {
	gs.api = new(mockAPI)
	gs.ctx = context.Background()
	gs.groupName = "groupName"
	gs.streamName = "streamName"
	gs.sut = NewGroup(gs.api, gs.groupName)
}

func (gs *groupTestSuite) TestCreateWithExistingStream() {
	const sequenceToken = "sequenceToken"
	ctx := context.Background()
	gs.describingStreamsReturns([]*cloudwatchlogs.LogStream{
		{UploadSequenceToken: aws.String(sequenceToken)},
	}, nil)

	writer, err := gs.sut.Create(ctx, gs.streamName)
	gs.NoError(err)
	gs.NotNil(writer)

	gs.Equal(sequenceToken, *writer.(*writerImpl).sequenceToken)
}

func (gs *groupTestSuite) TestCreateDescribingStreamFails() {
	ctx := context.Background()
	gs.describingStreamsReturns(nil, errors.New("bacon"))

	writer, err := gs.sut.Create(ctx, gs.streamName)
	gs.EqualError(err, "couldn't get log stream description: bacon")
	gs.Nil(writer)
}

func (gs *groupTestSuite) TestCreateWithoutExistingStream() {
	ctx := context.Background()
	gs.describingStreamsReturns(nil, nil)
	gs.creatingLogStreamReturns(nil)

	writer, err := gs.sut.Create(ctx, gs.streamName)
	gs.NoError(err)
	gs.NotNil(writer)

	gs.Nil(writer.(*writerImpl).sequenceToken)
}

func (gs *groupTestSuite) TestCreateCreatingStreamFails() {
	gs.describingStreamsReturns(nil, nil)
	gs.creatingLogStreamReturns(errors.New("bacon"))

	writer, err := gs.sut.Create(gs.ctx, gs.streamName)
	gs.EqualError(err, "could not create a log stream: bacon")
	gs.Nil(writer)
}

func (gs *groupTestSuite) describingStreamsReturns(result []*cloudwatchlogs.LogStream, err error) {
	gs.api.On(
		"DescribeLogStreamsWithContext",
		gs.ctx,
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(gs.groupName),
			LogStreamNamePrefix: aws.String(gs.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.DescribeLogStreamsOutput{LogStreams: result}, err)
}

func (gs *groupTestSuite) creatingLogStreamReturns(err error) {
	gs.api.On(
		"CreateLogStreamWithContext",
		gs.ctx,
		&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(gs.groupName),
			LogStreamName: aws.String(gs.streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.CreateLogStreamOutput{}, err)
}

func TestGroup(t *testing.T) {
	suite.Run(t, new(groupTestSuite))
}
