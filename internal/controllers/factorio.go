package controllers

import (
	"errors"
	"factorio/internal/services"
	"net/http"

	grove "github.com/StevenAlexanderJohnson/grove"
)

type FactorioController struct {
	logger          grove.ILogger
	factorioService *services.FactorioService
}

func NewFactorioController(logger grove.ILogger, service *services.FactorioService) *FactorioController {
	if logger == nil {
		panic("logger is required and cannot be nil")
	}
	return &FactorioController{
		logger:          logger,
		factorioService: service,
	}
}

func (f *FactorioController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /start", f.handleStart)
	mux.HandleFunc("POST /stop", f.handleStop)
	mux.HandleFunc("POST /update", f.handleUpdate)
}

func (f *FactorioController) handleStart(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("Received request to start Factorio server")
	if err := f.factorioService.StartServer(); err != nil {
		if errors.Is(err, services.ErrServerAlreadyRunning) {
			f.logger.Warning("Start server requested, but Factorio server is already running")
			grove.WriteErrorToResponse(w, http.StatusBadRequest, "the factorio server is already running")
			return
		}
		f.logger.Errorf("Failed to start Factorio server: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	f.logger.Info("Factorio server started successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "factorio server started",
	})
}

func (f *FactorioController) handleStop(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("Received request to stop Factorio server")
	if err := f.factorioService.StopServer(); err != nil {
		if errors.Is(err, services.ErrServerAlreadyStopped) {
			f.logger.Warning("Stop server requested, but Factorio server is already stopped")
			grove.WriteErrorToResponse(w, http.StatusBadRequest, "the factorio server is already stopped")
			return
		}
		f.logger.Errorf("Failed to stop Factorio server: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	f.logger.Info("Factorio server stopped successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "factorio server stopped",
	})
}

func (f *FactorioController) handleUpdate(w http.ResponseWriter, r *http.Request) {
	f.logger.Info("Received request to update Factorio server")
	if err := f.factorioService.UpdateServer(); err != nil {
		f.logger.Errorf("Failed to update Factorio server: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	f.logger.Info("Factorio server update completed successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "factorio server update completed successfully",
	})
}
