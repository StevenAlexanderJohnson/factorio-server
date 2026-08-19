package main

import (
	"context"
	"errors"
	"factorio/internal/config"
	"factorio/internal/factorio"
	"fmt"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

func startServer(messageChan chan<- factorio.FactorioMessage) error {
	responseChan := make(chan error)
	messageChan <- factorio.FactorioMessage{
		Type:  factorio.FactorioStart,
		Reply: responseChan,
	}
	return <-responseChan
}

func stopServer(messageChan chan<- factorio.FactorioMessage) error {
	responseChan := make(chan error)
	messageChan <- factorio.FactorioMessage{
		Type:  factorio.FactorioStop,
		Reply: responseChan,
	}
	return <-responseChan
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
}

func handleStart(messageChan chan<- factorio.FactorioMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := startServer(messageChan); err != nil {
			if errors.Is(err, factorio.ErrServerAlreadyRunning) {
				grove.WriteErrorToResponse(w, http.StatusBadRequest, "the factorio server is already running")
				return
			}
			grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = grove.WriteJsonBodyToResponse(w, map[string]string{
			"message": "factorio server started",
		})
	}
}

func handleStop(messageChan chan<- factorio.FactorioMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := stopServer(messageChan); err != nil {
			if errors.Is(err, factorio.ErrServerAlreadyStopped) {
				grove.WriteErrorToResponse(w, http.StatusBadRequest, "the factorio server is already stopped")
				return
			}
			grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		_ = grove.WriteJsonBodyToResponse(w, map[string]string{
			"message": "factorio server stopped",
		})
	}
}

func handleUpdate(ctx context.Context, messageChan chan<- factorio.FactorioMessage, cfg config.FactorioConfig, logger grove.ILogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wasRunning := true
		if err := stopServer(messageChan); err != nil {
			if errors.Is(err, factorio.ErrServerAlreadyStopped) {
				wasRunning = false
			} else {
				grove.WriteErrorToResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to stop server: %v", err))
				return
			}
		}

		updated, version, err := factorio.EnsureUpdated(ctx, cfg, logger)
		if err != nil {
			grove.WriteErrorToResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to update factorio: %v", err))
			return
		}

		if wasRunning {
			if err := startServer(messageChan); err != nil {
				grove.WriteErrorToResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to restart server after update: %v", err))
				return
			}
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
	}
}

