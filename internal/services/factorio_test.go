package services

import (
	"context"
	"errors"
	"factorio/internal/config"
	"os"
	"path/filepath"
	"testing"

	grove "github.com/StevenAlexanderJohnson/grove"
)

func TestFactorioService_StoppedState(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("auth:\n  api_key: test\n"), 0644)

	cfgManager, err := config.NewConfigManager(configPath)
	if err != nil {
		t.Fatalf("failed to create config manager: %v", err)
	}

	logger := grove.NewDefaultLogger("TestServer")

	srv, err := NewFactorioService(ctx, logger, cfgManager)
	if err != nil {
		t.Fatalf("failed to create factorio service: %v", err)
	}

	if srv.IsRunning() {
		t.Fatalf("expected server not to be running initially")
	}

	err = srv.StopServer()
	if !errors.Is(err, ErrServerAlreadyStopped) {
		t.Fatalf("expected ErrServerAlreadyStopped when stopping a stopped server, got: %v", err)
	}

	_, err = srv.SendCommand("/help")
	if !errors.Is(err, ErrServerAlreadyStopped) {
		t.Fatalf("expected ErrServerAlreadyStopped when sending command to stopped server, got: %v", err)
	}
}

func TestFactorioManager_LockConcurrency(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	_ = os.WriteFile(configPath, []byte("auth:\n  api_key: test\n"), 0644)
	cfgManager, _ := config.NewConfigManager(configPath)

	mgr := &factorioManager{
		cfgManager: cfgManager,
	}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = mgr.IsRunning()
				_, _ = mgr.SendCommand("test")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
