package cloudwatch

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	iface "github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/stretchr/testify/mock"
)

type mockAPI struct {
	mock.Mock
	iface.CloudWatchLogsAPI
}

func (m *mockAPI) PutLogEventsWithContext(ctx aws.Context, input *cloudwatchlogs.PutLogEventsInput, opts ...request.Option) (*cloudwatchlogs.PutLogEventsOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*cloudwatchlogs.PutLogEventsOutput), args.Error(1)
}

func (m *mockAPI) CreateLogStreamWithContext(ctx aws.Context, input *cloudwatchlogs.CreateLogStreamInput, opts ...request.Option) (*cloudwatchlogs.CreateLogStreamOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*cloudwatchlogs.CreateLogStreamOutput), args.Error(1)
}

func (m *mockAPI) GetLogEventsWithContext(ctx aws.Context, input *cloudwatchlogs.GetLogEventsInput, opts ...request.Option) (*cloudwatchlogs.GetLogEventsOutput, error) {
	args := m.Called(ctx, input, opts)
	return args.Get(0).(*cloudwatchlogs.GetLogEventsOutput), args.Error(1)
}
