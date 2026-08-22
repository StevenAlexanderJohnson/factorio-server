package controllers

import (
	"errors"
	"factorio/internal/services"
	"net/http"

	grove "github.com/StevenAlexanderJohnson/grove"
)

type FactorioController struct {
	factorioService *services.FactorioService
}

func NewFactorioController(service *services.FactorioService) *FactorioController {
	return &FactorioController{
		factorioService: service,
	}
}

func (f *FactorioController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /start", f.handleStart)
	mux.HandleFunc("POST /stop", f.handleStop)
	mux.HandleFunc("POST /update", f.handleUpdate)
}

func (f *FactorioController) handleStart(w http.ResponseWriter, r *http.Request) {
	if err := f.factorioService.StartServer(); err != nil {
		if errors.Is(err, services.ErrServerAlreadyRunning) {
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

func (f *FactorioController) handleStop(w http.ResponseWriter, r *http.Request) {
	if err := f.factorioService.StopServer(); err != nil {
		if errors.Is(err, services.ErrServerAlreadyStopped) {
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

func (f *FactorioController) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if err := f.factorioService.UpdateServer(); err != nil {
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "factorio server update completed successfully",
	})
}
