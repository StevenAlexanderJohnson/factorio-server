package main

import (
	"context"
	"factorio/internal/config"
	"factorio/internal/controllers"
	"factorio/internal/factorio"
	"factorio/internal/services"
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

	factorioService, err := services.NewFactorioService(ctx, logger, cfg.Factorio)
	if err != nil {
		logger.Errorf("An error occurred while creating the factorio service: %w\n", err)
		return
	}
	go func() {
		if err := <-factorioService.FatalChan(); err != nil {
			logger.Errorf("an error occurred while running the factorio loop: %v", err)
			cancel()
		}
	}()

	settingsService := services.NewSettingsService(cfg.Factorio)

	scope := grove.
		NewScope("main").
		WithRoute("/healthz", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("healthy"))
			}),
		).
		WithController(controllers.NewFactorioController(factorioService)).
		WithController(controllers.NewSettingsController(settingsService))

	app := grove.NewApp("factorio").WithScope("/", scope)

	if err := app.Run(); err != nil {
		cancel()
		panic(err)
	}
}
