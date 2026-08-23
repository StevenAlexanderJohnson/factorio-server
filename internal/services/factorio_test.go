package services

import (
	"context"
	"errors"
	"factorio/internal/config"
	"testing"

	grove "github.com/StevenAlexanderJohnson/grove"
)

func TestFactorioService_StoppedState(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := grove.NewDefaultLogger("TestServer")
	cfg := config.FactorioConfig{}

	srv, err := NewFactorioService(ctx, logger, cfg)
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
	mgr := &factorioManager{}

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
