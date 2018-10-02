package cloudwatch

import (
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// Constraints are documented here:
// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudwatchlogs/#CloudWatchLogs.PutLogEvents
const (
	maxBatchSizeBytes  = 1048576
	maxBatchSizeEvents = 10000
	paddingSize        = 26
)

type logBatch struct {
	count, size int
	events      []*cloudwatchlogs.InputLogEvent
	next        *logBatch
}

func (l *logBatch) add(event *cloudwatchlogs.InputLogEvent) *logBatch {
	if event.Message == nil {
		return l
	}
	l.count++
	nextSize := l.size + len(*event.Message) + paddingSize
	if nextSize > maxBatchSizeBytes || l.count > maxBatchSizeEvents {
		l.next = new(logBatch)
		return l.next.add(event)

	}
	l.events = append(l.events, event)
	l.size = nextSize
	return l
}
