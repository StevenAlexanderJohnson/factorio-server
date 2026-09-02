package logs

import (
	"context"
	"factorio/internal/events"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogParser(t *testing.T) {
	scriptContent := `
return function(line)
    if string.find(line, "joined the game") then
        return {
            type = "player_join",
            data = {
                message = line,
                tags = {"player", "network"},
            }
        }
    end
    return nil
end
`
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "parser.lua")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0644); err != nil {
		t.Fatalf("failed to write test lua script: %v", err)
	}

	bus := events.NewEventBus()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, unsub := bus.Subscribe(ctx, "player_join")
	defer unsub()

	parser, err := NewLogParser(scriptPath, bus)
	if err != nil {
		t.Fatalf("failed to create log parser: %v", err)
	}

	if err := parser.ParseLine("2026-09-01 12:00:00 [JOIN] player1 joined the game"); err != nil {
		t.Fatalf("ParseLine failed: %v", err)
	}

	select {
	case msg := <-ch:
		if msg.Type != "player_join" {
			t.Fatalf("unexpected event type: %s", msg.Type)
		}
		dataMap, ok := msg.Data.(map[string]any)
		if !ok {
			t.Fatalf("expected map[string]any data, got %T", msg.Data)
		}
		tags, ok := dataMap["tags"].([]any)
		if !ok || len(tags) != 2 {
			t.Fatalf("expected []any with 2 items for tags, got %v", dataMap["tags"])
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for parsed event")
	}
}
