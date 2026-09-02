package controllers

import (
	"bufio"
	"context"
	"factorio/internal/events"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	grove "github.com/StevenAlexanderJohnson/grove"
)

func TestSSEController_Subscribe(t *testing.T) {
	bus := events.NewEventBus()
	logger := grove.NewDefaultLogger("TestSSE")
	controller := NewSSEController(logger, bus)

	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest("GET", "/events?types=chat,join", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		mux.ServeHTTP(rec, req)
	}()

	// Allow goroutine to start and subscribe
	time.Sleep(50 * time.Millisecond)

	bus.Publish(events.EventMessage{Type: "chat", Data: map[string]string{"user": "alice", "text": "hello"}})
	bus.Publish(events.EventMessage{Type: "ignored", Data: "skip me"})
	bus.Publish(events.EventMessage{Type: "join", Data: map[string]string{"user": "bob"}})

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	body := rec.Body.String()
	if !strings.Contains(body, "event: chat\n") {
		t.Errorf("expected event: chat in stream, got:\n%s", body)
	}
	if !strings.Contains(body, `"alice"`) {
		t.Errorf("expected alice in stream, got:\n%s", body)
	}
	if strings.Contains(body, "skip me") {
		t.Errorf("did not expect ignored event in stream, got:\n%s", body)
	}
	if !strings.Contains(body, "event: join\n") {
		t.Errorf("expected event: join in stream, got:\n%s", body)
	}
}

func TestSSEController_SubscribeAll(t *testing.T) {
	bus := events.NewEventBus()
	logger := grove.NewDefaultLogger("TestSSE")
	controller := NewSSEController(logger, bus)

	mux := http.NewServeMux()
	controller.RegisterRoutes(mux)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest("GET", "/events", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	done := make(chan struct{})
	go func() {
		defer close(done)
		mux.ServeHTTP(rec, req)
	}()

	time.Sleep(50 * time.Millisecond)

	bus.Publish(events.EventMessage{Type: "custom", Data: 12345})

	time.Sleep(50 * time.Millisecond)
	cancel()
	<-done

	scanner := bufio.NewScanner(strings.NewReader(rec.Body.String()))
	foundCustom := false
	for scanner.Scan() {
		line := scanner.Text()
		if line == "event: custom" {
			foundCustom = true
		}
	}

	if !foundCustom {
		t.Errorf("expected custom event when subscribed with no filters, got:\n%s", rec.Body.String())
	}
}
