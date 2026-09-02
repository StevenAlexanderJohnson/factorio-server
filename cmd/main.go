package main

import (
	"context"
	"factorio/internal/config"
	"factorio/internal/controllers"
	"factorio/internal/events"
	"factorio/internal/factorio"
	"factorio/internal/logs"
	"factorio/internal/middleware"
	"factorio/internal/services"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
)

func main() {
	configPath := flag.String("config", "", "Path to YAML configuration file")
	flag.Parse()

	logger := grove.NewDefaultLogger("FactorioAPI")

	cfgManager, err := config.NewConfigManager(*configPath)
	if err != nil {
		logger.Errorf("failed to load configuration: %v", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	updated, version, err := factorio.EnsureUpdated(ctx, cfgManager.GetConfig().Factorio, logger)
	if err != nil {
		logger.Errorf("failed to ensure factorio was up to date on startup: %v", err)
		os.Exit(1)
	}
	if updated {
		logger.Infof("updated factorio on startup, new version: %s", version)
	} else {
		logger.Infof("factorio was up to date on startup")
	}

	eventBus := events.NewEventBus()

	logScriptPath := cfgManager.GetConfig().Logs.LogParserScript
	logParser, err := logs.NewLogParser(cfgManager.GetConfig().Logs, eventBus)
	if err != nil {
		logger.Warningf("Failed to initialize log parser with script %q: %v", logScriptPath, err)
	} else {
		logger.Infof("Loaded Lua log parser script from %s", logScriptPath)
		logParser.StartWatcher(ctx, 2*time.Second, logger)
	}

	factorioService, err := services.NewFactorioService(ctx, logger, cfgManager, logParser)
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

	settingsService := services.NewSettingsService(cfgManager)

	commandService := services.NewCommandService(factorioService)

	scope := grove.
		NewScope("main").
		WithMiddleware(middleware.NewAuthMiddleware(cfgManager)).
		WithRoute("/healthz", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("healthy"))
			}),
		).
		WithController(controllers.NewFactorioController(logger, factorioService)).
		WithController(controllers.NewSettingsController(logger, settingsService)).
		WithController(controllers.NewCommandController(logger, commandService)).
		WithController(controllers.NewConfigController(logger, cfgManager)).
		WithController(controllers.NewSSEController(logger, eventBus))

	app := grove.NewApp("factorio").WithScope("/", scope)

	if err := app.Run(); err != nil {
		cancel()
		panic(err)
	}
}
