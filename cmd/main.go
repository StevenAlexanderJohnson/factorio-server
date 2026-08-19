package main

import (
	"context"
	"factorio/internal/config"
	"factorio/internal/factorio"
	"flag"
	"net/http"
	"os"

	"github.com/StevenAlexanderJohnson/grove"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML configuration file")
	flag.Parse()

	logger := grove.NewDefaultLogger("FactorioAPI")

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Errorf("failed to load configuration: %v", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updated, version, err := factorio.EnsureUpdated(ctx, cfg.Factorio, logger)
	if err != nil {
		logger.Errorf("failed to ensure factorio was up to date on startup: %v", err)
		os.Exit(1)
	}
	if updated {
		logger.Infof("updated factorio on startup, new version: %s", version)
	} else {
		logger.Infof("factorio was up to date on startup")
	}
	messageChan, fatalChan := factorio.StartFactorioLoop(ctx, cfg.Factorio)

	go func() {
		if err := <-fatalChan; err != nil {
			logger.Errorf("an error occurred while running the factorio loop: %v", err)
			cancel()
		}
	}()

	scope := grove.
		NewScope("main").
		WithRoute("/healthz", http.HandlerFunc(handleHealthz)).
		WithRoute("POST /start", handleStart(messageChan)).
		WithRoute("POST /stop", handleStop(messageChan)).
		WithRoute("POST /update", handleUpdate(ctx, messageChan, cfg.Factorio, logger))

	app := grove.NewApp("factorio").WithScope("/", scope)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

