package controllers

import (
	"bytes"
	"encoding/json"
	"factorio/internal/config"
	"factorio/internal/factorio/models"
	"factorio/internal/services"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/StevenAlexanderJohnson/grove"
)

func TestSettingsController(t *testing.T) {
	tempDir := t.TempDir()

	cfg := config.FactorioConfig{
		ServerSettingsPath:  filepath.Join(tempDir, "server-settings.json"),
		MapSettingsPath:     filepath.Join(tempDir, "map-settings.json"),
		MapGenSettingsPath:  filepath.Join(tempDir, "map-gen-settings.json"),
		ServerAdminListPath: filepath.Join(tempDir, "server-adminlist.json"),
		ServerWhiteListPath: filepath.Join(tempDir, "server-whitelist.json"),
		ServerBanListPath:   filepath.Join(tempDir, "server-banlist.json"),
	}

	logger := grove.NewDefaultLogger("TestSettingsController")
	settingsSrv := services.NewSettingsService(cfg)
	ctrl := NewSettingsController(logger, settingsSrv)

	mux := http.NewServeMux()
	ctrl.RegisterRoutes(mux)

	// Test PUT /settings/server
	serverData, _ := json.Marshal(models.ServerSettings{Name: "Controller Test"})
	req := httptest.NewRequest("PUT", "/settings/server", bytes.NewReader(serverData))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d: %s", rec.Code, rec.Body.String())
	}

	// Test GET /settings/server
	req = httptest.NewRequest("GET", "/settings/server", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %d", rec.Code)
	}

	var gotServer models.ServerSettings
	_ = json.NewDecoder(rec.Body).Decode(&gotServer)
	if gotServer.Name != "Controller Test" {
		t.Errorf("expected server name 'Controller Test', got: %s", gotServer.Name)
	}
}
