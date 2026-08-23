package controllers

import (
	"errors"
	"factorio/internal/models"
	"factorio/internal/services"
	"net/http"

	"github.com/StevenAlexanderJohnson/grove"
)

type CommandController struct {
	logger         grove.ILogger
	commandService *services.CommandService
}

func NewCommandController(logger grove.ILogger, commandService *services.CommandService) *CommandController {
	if logger == nil {
		panic("logger is required and cannot be nil")
	}
	return &CommandController{
		logger:         logger,
		commandService: commandService,
	}
}

func (c *CommandController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /command", c.handleCommand)
}

func (c *CommandController) handleCommand(w http.ResponseWriter, r *http.Request) {
	body, err := grove.ParseJsonBodyFromRequest[models.CommandRequest](r)
	if err != nil {
		c.logger.Warningf("Failed to parse command request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "")
		return
	}
	c.logger.Infof("Sending command to Factorio server: %s", body.Command)
	if err := c.commandService.SendCommand(body.Command); err != nil {
		if errors.Is(err, services.ErrFactorioServerNotRunning) {
			c.logger.Warningf("A request was received to send a command while the server was not running: %q", body.Command)
			grove.WriteErrorToResponse(w, http.StatusUnprocessableEntity, "The server is not online.")
			return
		}
		c.logger.Errorf("Failed to send command %q: %v", body.Command, err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, "")
		return
	}
	c.logger.Infof("Command sent successfully: %s", body.Command)
}
