package notification

import (
	"fmt"
	"sync"
	"time"

	"github.com/YurcheuskiRadzivon/booking-system/internal/models/notification"
	"github.com/google/uuid"
)

type MessageBroker struct {
	subscribers map[notification.NotificationChannel][]chan notification.NotificationEvent
	eventLog    []notification.NotificationEvent
	mu          sync.RWMutex
	wg          sync.WaitGroup
	closed      bool
}

func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscribers: make(map[notification.NotificationChannel][]chan notification.NotificationEvent),
		eventLog:    make([]notification.NotificationEvent, 0),
	}
}

func (b *MessageBroker) Subscribe(channel notification.NotificationChannel, bufferSize int) <-chan notification.NotificationEvent {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan notification.NotificationEvent, bufferSize)
	b.subscribers[channel] = append(b.subscribers[channel], ch)
	return ch
}

func (b *MessageBroker) Publish(event notification.NotificationEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return
	}

	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	b.eventLog = append(b.eventLog, event)

	if subs, ok := b.subscribers[event.Channel]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
				fmt.Printf("⚠️ Warning: subscriber channel full for %s\n", event.Channel)
			}
		}
	}

	if subs, ok := b.subscribers["all"]; ok {
		for _, ch := range subs {
			select {
			case ch <- event:
			default:
			}
		}
	}
}

func (b *MessageBroker) GetEventLog() []notification.NotificationEvent {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return append([]notification.NotificationEvent{}, b.eventLog...)
}

func (b *MessageBroker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closed = true
	for _, subs := range b.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
	b.subscribers = make(map[notification.NotificationChannel][]chan notification.NotificationEvent)
}
