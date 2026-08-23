package controllers

import (
	"bytes"
	"encoding/json"
	"factorio/internal/config"
	"factorio/internal/factorio/models"
	"factorio/internal/services"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/StevenAlexanderJohnson/grove"
)

func TestSettingsController(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	cfgYAML := fmt.Sprintf(`
factorio:
  server_settings_path: %q
  map_settings_path: %q
  map_gen_settings_path: %q
  server_adminlist_path: %q
  server_whitelist_path: %q
  server_banlist_path: %q
`,
		filepath.Join(tempDir, "server-settings.json"),
		filepath.Join(tempDir, "map-settings.json"),
		filepath.Join(tempDir, "map-gen-settings.json"),
		filepath.Join(tempDir, "server-adminlist.json"),
		filepath.Join(tempDir, "server-whitelist.json"),
		filepath.Join(tempDir, "server-banlist.json"),
	)
	_ = os.WriteFile(configPath, []byte(cfgYAML), 0644)

	cfgManager, err := config.NewConfigManager(configPath)
	if err != nil {
		t.Fatalf("failed to create config manager: %v", err)
	}

	logger := grove.NewDefaultLogger("TestSettingsController")
	settingsSrv := services.NewSettingsService(cfgManager)
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
