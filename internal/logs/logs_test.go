package logs

import (
	"context"
	"factorio/internal/config"
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

	parser, err := NewLogParser(config.LogsConfig{LogParserScript: scriptPath}, bus)
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

func TestLogParserWatcher(t *testing.T) {
	scriptV1 := `
return function(line)
    if string.find(line, "test") then
        return { type = "event_v1", data = { msg = line } }
    end
    return nil
end
`
	scriptV2 := `
return function(line)
    if string.find(line, "test") then
        return { type = "event_v2", data = { msg = line } }
    end
    return nil
end
`
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "parser.lua")
	if err := os.WriteFile(scriptPath, []byte(scriptV1), 0644); err != nil {
		t.Fatalf("failed to write test lua script: %v", err)
	}

	bus := events.NewEventBus()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	parser, err := NewLogParser(config.LogsConfig{LogParserScript: scriptPath}, bus)
	if err != nil {
		t.Fatalf("failed to create log parser: %v", err)
	}

	parser.StartWatcher(ctx, 50*time.Millisecond, nil)

	ch1, unsub1 := bus.Subscribe(ctx, "event_v1")
	defer unsub1()

	if err := parser.ParseLine("test line"); err != nil {
		t.Fatalf("ParseLine failed: %v", err)
	}

	select {
	case msg := <-ch1:
		if msg.Type != "event_v1" {
			t.Fatalf("expected event_v1, got %s", msg.Type)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for v1 event")
	}

	ch2, unsub2 := bus.Subscribe(ctx, "event_v2")
	defer unsub2()

	// Ensure modification time changes even on fast filesystems
	time.Sleep(100 * time.Millisecond)
	if err := os.WriteFile(scriptPath, []byte(scriptV2), 0644); err != nil {
		t.Fatalf("failed to update test lua script: %v", err)
	}

	// Wait for polling watcher to detect change and reload
	var reloaded bool
	for i := 0; i < 20; i++ {
		time.Sleep(50 * time.Millisecond)
		_ = parser.ParseLine("test line")
		select {
		case msg := <-ch2:
			if msg.Type == "event_v2" {
				reloaded = true
				break
			}
		default:
		}
		if reloaded {
			break
		}
	}

	if !reloaded {
		t.Fatal("watcher failed to reload updated lua script")
	}
}
