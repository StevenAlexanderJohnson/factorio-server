package events

import (
	"context"
	"sync"
)

type EventMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type eventSubscriber struct {
	channel         chan<- EventMessage
	subscribedTypes map[string]struct{}
	all             bool
}

type EventBus struct {
	mu          sync.RWMutex
	subscribers map[*eventSubscriber]struct{}
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[*eventSubscriber]struct{}, 0),
	}
}

// Subscribes to a series of events from the event bus. If no types are provided, it subscribes to all events.
// Returns a channel that you can receive messages from, and an unsubscribe function.
func (e *EventBus) Subscribe(ctx context.Context, types ...string) (<-chan EventMessage, func()) {
	e.mu.Lock()
	defer e.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)

	eventChannel := make(chan EventMessage, 10)
	subscribedTypes := make(map[string]struct{})
	all := len(types) == 0
	for _, subTypes := range types {
		if subTypes == "*" {
			all = true
		}
		subscribedTypes[subTypes] = struct{}{}
	}

	subscriber := &eventSubscriber{
		channel:         eventChannel,
		subscribedTypes: subscribedTypes,
		all:             all,
	}

	e.subscribers[subscriber] = struct{}{}

	go func() {
		<-ctx.Done()
		e.mu.Lock()
		defer e.mu.Unlock()

		delete(e.subscribers, subscriber)
		close(subscriber.channel)
	}()

	return eventChannel, cancel
}

// Publishes a message to the event bus, sending a copy to all subscribers for that type.
func (e *EventBus) Publish(message EventMessage) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for sub := range e.subscribers {
		if _, ok := sub.subscribedTypes[message.Type]; ok || sub.all {
			select {
			case sub.channel <- message:
			default:
			}
		}
	}
}
