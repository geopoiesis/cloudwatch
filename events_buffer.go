package cloudwatch

import (
	"sync"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// eventsBuffer represents a buffer of cloudwatch events that are protected by a
// mutex.
type eventsBuffer struct {
	sync.RWMutex
	head, tail *logBatch
}

func newEventsBuffer() *eventsBuffer {
	batch := new(logBatch)
	return &eventsBuffer{head: batch, tail: batch}
}

func (b *eventsBuffer) add(event *cloudwatchlogs.InputLogEvent) {
	b.Lock()
	defer b.Unlock()
	b.tail = b.tail.add(event)
}

func (b *eventsBuffer) drain() []*cloudwatchlogs.InputLogEvent {
	b.Lock()
	defer b.Unlock()

	ret := b.head.events
	if b.head == b.tail {
		b.head = new(logBatch)
		b.tail = b.head
	} else {
		b.head = b.head.next
	}
	return ret
}

func (b *eventsBuffer) hasMore() bool {
	b.RLock()
	defer b.RUnlock()
	return len(b.head.events) > 0
}
