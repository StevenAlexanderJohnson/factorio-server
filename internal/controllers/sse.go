package controllers

import (
	"encoding/json"
	"factorio/internal/events"
	"fmt"
	"net/http"
	"strings"

	grove "github.com/StevenAlexanderJohnson/grove"
)

type SSEController struct {
	logger   grove.ILogger
	eventBus *events.EventBus
}

func NewSSEController(logger grove.ILogger, eventBus *events.EventBus) *SSEController {
	if logger == nil {
		panic("logger is required and cannot be nil")
	}
	if eventBus == nil {
		panic("eventBus is required and cannot be nil")
	}
	return &SSEController{
		logger:   logger,
		eventBus: eventBus,
	}
}

func (s *SSEController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /events", s.Subscribe)
}

func (s *SSEController) Subscribe(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		s.logger.Error("Streaming unsupported on ResponseWriter")
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	types := parseEventTypes(r)
	eventChan, unsub := s.eventBus.Subscribe(r.Context(), types...)
	defer unsub()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	s.logger.Infof("SSE client connected with subscriptions: %v", types)

	for {
		select {
		case <-r.Context().Done():
			s.logger.Info("SSE client disconnected")
			return
		case msg, ok := <-eventChan:
			if !ok {
				return
			}
			jsonData, err := json.Marshal(msg.Data)
			if err != nil {
				s.logger.Warningf("Failed to marshal event data: %v", err)
				continue
			}

			if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", msg.Type, jsonData); err != nil {
				s.logger.Warningf("Failed to write event to SSE stream: %v", err)
				return
			}
			flusher.Flush()
		}
	}
}

func parseEventTypes(r *http.Request) []string {
	var types []string
	query := r.URL.Query()
	for _, key := range []string{"type", "types", "event", "events"} {
		for _, val := range query[key] {
			for _, item := range strings.Split(val, ",") {
				trimmed := strings.TrimSpace(item)
				if trimmed != "" {
					types = append(types, trimmed)
				}
			}
		}
	}
	return types
}
