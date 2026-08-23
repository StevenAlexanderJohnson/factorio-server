package controllers

import (
	"encoding/json"
	"factorio/internal/factorio/models"
	"factorio/internal/services"
	"net/http"

	grove "github.com/StevenAlexanderJohnson/grove"
)

type SettingsController struct {
	logger          grove.ILogger
	settingsService services.SettingsService
}

func NewSettingsController(logger grove.ILogger, service services.SettingsService) *SettingsController {
	if logger == nil {
		panic("logger is required and cannot be nil")
	}
	return &SettingsController{
		logger:          logger,
		settingsService: service,
	}
}

func (s *SettingsController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /settings/server", s.handleGetServerSettings)
	mux.HandleFunc("PUT /settings/server", s.handleUpdateServerSettings)

	mux.HandleFunc("GET /settings/map", s.handleGetMapSettings)
	mux.HandleFunc("PUT /settings/map", s.handleUpdateMapSettings)

	mux.HandleFunc("GET /settings/map-gen", s.handleGetMapGenSettings)
	mux.HandleFunc("PUT /settings/map-gen", s.handleUpdateMapGenSettings)

	mux.HandleFunc("GET /settings/admin-list", s.handleGetAdminList)
	mux.HandleFunc("PUT /settings/admin-list", s.handleUpdateAdminList)

	mux.HandleFunc("GET /settings/whitelist", s.handleGetWhitelist)
	mux.HandleFunc("PUT /settings/whitelist", s.handleUpdateWhitelist)

	mux.HandleFunc("GET /settings/ban-list", s.handleGetBanList)
	mux.HandleFunc("PUT /settings/ban-list", s.handleUpdateBanList)
}

func (s *SettingsController) handleGetServerSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.settingsService.GetServerSettings()
	if err != nil {
		s.logger.Errorf("Failed to get server settings: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, settings)
}

func (s *SettingsController) handleUpdateServerSettings(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Updating server settings")
	var settings models.ServerSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		s.logger.Warningf("Failed to decode server settings request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := s.settingsService.UpdateServerSettings(settings); err != nil {
		s.logger.Errorf("Failed to update server settings: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.logger.Info("Server settings updated successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "server settings updated successfully",
	})
}

func (s *SettingsController) handleGetMapSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.settingsService.GetMapSettings()
	if err != nil {
		s.logger.Errorf("Failed to get map settings: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, settings)
}

func (s *SettingsController) handleUpdateMapSettings(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Updating map settings")
	var settings models.MapSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		s.logger.Warningf("Failed to decode map settings request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := s.settingsService.UpdateMapSettings(settings); err != nil {
		s.logger.Errorf("Failed to update map settings: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.logger.Info("Map settings updated successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "map settings updated successfully",
	})
}

func (s *SettingsController) handleGetMapGenSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.settingsService.GetMapGenSettings()
	if err != nil {
		s.logger.Errorf("Failed to get map-gen settings: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, settings)
}

func (s *SettingsController) handleUpdateMapGenSettings(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Updating map-gen settings")
	var settings models.MapGenSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		s.logger.Warningf("Failed to decode map-gen settings request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := s.settingsService.UpdateMapGenSettings(settings); err != nil {
		s.logger.Errorf("Failed to update map-gen settings: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.logger.Info("Map-gen settings updated successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "map-gen settings updated successfully",
	})
}

func (s *SettingsController) handleGetAdminList(w http.ResponseWriter, r *http.Request) {
	list, err := s.settingsService.GetAdminList()
	if err != nil {
		s.logger.Errorf("Failed to get admin list: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, list)
}

func (s *SettingsController) handleUpdateAdminList(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Updating admin list")
	var list models.AdminList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		s.logger.Warningf("Failed to decode admin list request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := s.settingsService.UpdateAdminList(list); err != nil {
		s.logger.Errorf("Failed to update admin list: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.logger.Info("Admin list updated successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "admin list updated successfully",
	})
}

func (s *SettingsController) handleGetWhitelist(w http.ResponseWriter, r *http.Request) {
	list, err := s.settingsService.GetWhitelist()
	if err != nil {
		s.logger.Errorf("Failed to get whitelist: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, list)
}

func (s *SettingsController) handleUpdateWhitelist(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Updating whitelist")
	var list models.Whitelist
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		s.logger.Warningf("Failed to decode whitelist request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := s.settingsService.UpdateWhitelist(list); err != nil {
		s.logger.Errorf("Failed to update whitelist: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.logger.Info("Whitelist updated successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "whitelist updated successfully",
	})
}

func (s *SettingsController) handleGetBanList(w http.ResponseWriter, r *http.Request) {
	list, err := s.settingsService.GetBanList()
	if err != nil {
		s.logger.Errorf("Failed to get ban list: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = grove.WriteJsonBodyToResponse(w, list)
}

func (s *SettingsController) handleUpdateBanList(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Updating ban list")
	var list models.BanList
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		s.logger.Warningf("Failed to decode ban list request body: %v", err)
		grove.WriteErrorToResponse(w, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	if err := s.settingsService.UpdateBanList(list); err != nil {
		s.logger.Errorf("Failed to update ban list: %v", err)
		grove.WriteErrorToResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.logger.Info("Ban list updated successfully")
	_ = grove.WriteJsonBodyToResponse(w, map[string]string{
		"message": "ban list updated successfully",
	})
}
