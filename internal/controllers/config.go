package controllers

import (
	"factorio/internal/config"
	"net/http"

	grove "github.com/StevenAlexanderJohnson/grove"
)

type ConfigController struct {
	logger     grove.ILogger
	cfgManager *config.ConfigManager
}

func NewConfigController(logger grove.ILogger, cfgManager *config.ConfigManager) *ConfigController {
	return &ConfigController{
		logger:     logger,
		cfgManager: cfgManager,
	}
}

func (c *ConfigController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /config/factorio", c.handleGetConfig)
	mux.HandleFunc("PUT /config/factorio", c.handleFactorioUpdate)
}

func (c *ConfigController) handleGetConfig(w http.ResponseWriter, _ *http.Request) {
	if err := grove.WriteJsonBodyToResponse(w, c.cfgManager.GetConfig().Factorio); err != nil {
		c.logger.Errorf("failed to marshal config file to json: %w", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
	}
}

func (c *ConfigController) handleFactorioUpdate(w http.ResponseWriter, r *http.Request) {
	cfg, err := grove.ParseJsonBodyFromRequest[config.FactorioConfig](r)
	if err != nil {
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "unable to parse request")
		return
	}

	if err := c.cfgManager.UpdateFactorioConfig(cfg); err != nil {
		c.logger.Errorf("an error occurred while updating factorio config: %w", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
}
