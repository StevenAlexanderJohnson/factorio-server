package main

import (
	"context"
	"errors"
	"factorio/internal/config"
	"factorio/internal/factorio"
	"flag"
	"fmt"
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
		logger.Errorf("failed to ensure factorio was up to date on startup: %w", err)
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
		WithRoute("/healthz", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("healthy"))
			},
		)).
		WithRoute("POST /start", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				responseChan := make(chan error)
				messageChan <- factorio.FactorioMessage{
					Type:  factorio.FactorioStart,
					Reply: responseChan,
				}

				if err := <-responseChan; err != nil {
					if errors.Is(err, factorio.ErrServerAlreadyRunning) {
						grove.WriteErrorToResponse(w, http.StatusBadRequest, "the factorio server is already running")
						return
					}
					grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
				}
			},
		)).
		WithRoute("POST /stop", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				responseChan := make(chan error)
				messageChan <- factorio.FactorioMessage{
					Type:  factorio.FactorioStop,
					Reply: responseChan,
				}

				if err := <-responseChan; err != nil {
					if errors.Is(err, factorio.ErrServerAlreadyStopped) {
						grove.WriteErrorToResponse(w, http.StatusBadRequest, "the factorio server is already stopped")
						return
					}
					grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
				}
			},
		)).
		WithRoute("POST /update", http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				updated, version, err := factorio.EnsureUpdated(ctx, cfg.Factorio, logger)
				if err != nil {
					grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
					return
				}
				output := map[string]string{}
				if updated {
					output["message"] = fmt.Sprintf("factorio has been updated to version: %s", version)
				} else {
					output["message"] = "factorio was already up to date"
				}

				if err := grove.WriteJsonBodyToResponse(w, output); err != nil {
					grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())

				}
			},
		))

	app := grove.NewApp("factorio").WithScope("/", scope)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
