package cloudwatch

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func TestGroupCreateFailure(t *testing.T) {
	ctx := context.Background()
	client := new(mockAPI)
	const groupName = "groupName"
	const streamName = "streamName"

	client.On(
		"CreateLogStreamWithContext",
		ctx,
		&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String(groupName),
			LogStreamName: aws.String(streamName),
		},
		[]request.Option(nil),
	).Return(&cloudwatchlogs.CreateLogStreamOutput{}, errors.New("bacon"))

	_, err := NewGroup(client, groupName).Create(ctx, streamName)
	assert.EqualError(t, err, "could not create a log stream: bacon")
}
