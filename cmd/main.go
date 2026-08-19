package main

import (
	"context"
	"errors"
	"factorio/internal/factorio"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	messageChan, fatalChan := factorio.StartFactorioLoop(ctx, "/factorio/data/saves/my-world.zip")

	logger := grove.NewDefaultLogger("FactorioAPI")

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
				responseChan := make(chan error)
				messageChan <- factorio.FactorioMessage{
					Type:  factorio.FactorioUpdate,
					Reply: responseChan,
				}

				if err := <-responseChan; err != nil {
					grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
				}
			},
		))

	app := grove.NewApp("factorio").WithScope("/", scope)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
