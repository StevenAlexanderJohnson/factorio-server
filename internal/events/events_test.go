package events

import (
	"context"
	"testing"
	"time"
)

func TestEventBus_SubscribeAndPublish(t *testing.T) {
	bus := NewEventBus()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, unsub := bus.Subscribe(ctx, "chat", "join")
	defer unsub()

	bus.Publish(EventMessage{Type: "chat", Data: "hello"})
	bus.Publish(EventMessage{Type: "other", Data: "ignored"})
	bus.Publish(EventMessage{Type: "join", Data: "player1"})

	select {
	case msg := <-ch:
		if msg.Type != "chat" || msg.Data != "hello" {
			t.Fatalf("unexpected message: %+v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for chat message")
	}

	select {
	case msg := <-ch:
		if msg.Type != "join" || msg.Data != "player1" {
			t.Fatalf("unexpected message: %+v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for join message")
	}
}

func TestEventBus_SubscribeAll(t *testing.T) {
	bus := NewEventBus()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, unsub := bus.Subscribe(ctx)
	defer unsub()

	bus.Publish(EventMessage{Type: "any_event", Data: 123})

	select {
	case msg := <-ch:
		if msg.Type != "any_event" || msg.Data != 123 {
			t.Fatalf("unexpected message: %+v", msg)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for event")
	}
}
