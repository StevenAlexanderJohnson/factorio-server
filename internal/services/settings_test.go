package services

import (
	"factorio/internal/config"
	"factorio/internal/factorio/models"
	"path/filepath"
	"testing"
)

func TestSettingsService_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()

	cfg := config.FactorioConfig{
		ServerSettingsPath:  filepath.Join(tempDir, "server-settings.json"),
		MapSettingsPath:     filepath.Join(tempDir, "map-settings.json"),
		MapGenSettingsPath:  filepath.Join(tempDir, "map-gen-settings.json"),
		ServerAdminListPath: filepath.Join(tempDir, "server-adminlist.json"),
		ServerWhiteListPath: filepath.Join(tempDir, "server-whitelist.json"),
		ServerBanListPath:   filepath.Join(tempDir, "server-banlist.json"),
	}

	srv := NewSettingsService(cfg)

	// Test ServerSettings
	srvSettings := models.ServerSettings{Name: "Test Server"}
	if err := srv.UpdateServerSettings(srvSettings); err != nil {
		t.Fatalf("failed to update server settings: %v", err)
	}
	gotSrvSettings, err := srv.GetServerSettings()
	if err != nil {
		t.Fatalf("failed to get server settings: %v", err)
	}
	if gotSrvSettings.Name != "Test Server" {
		t.Errorf("expected server name 'Test Server', got: %s", gotSrvSettings.Name)
	}

	// Test MapSettings
	mapSettings := models.MapSettings{MaxFailedBehaviorCount: 5}
	if err := srv.UpdateMapSettings(mapSettings); err != nil {
		t.Fatalf("failed to update map settings: %v", err)
	}
	gotMapSettings, err := srv.GetMapSettings()
	if err != nil {
		t.Fatalf("failed to get map settings: %v", err)
	}
	if gotMapSettings.MaxFailedBehaviorCount != 5 {
		t.Errorf("expected MaxFailedBehaviorCount = 5, got: %d", gotMapSettings.MaxFailedBehaviorCount)
	}

	// Test MapGenSettings
	mapGenSettings := models.MapGenSettings{Width: 100, Height: 200}
	if err := srv.UpdateMapGenSettings(mapGenSettings); err != nil {
		t.Fatalf("failed to update map gen settings: %v", err)
	}
	gotMapGenSettings, err := srv.GetMapGenSettings()
	if err != nil {
		t.Fatalf("failed to get map gen settings: %v", err)
	}
	if gotMapGenSettings.Width != 100 || gotMapGenSettings.Height != 200 {
		t.Errorf("expected Width 100 and Height 200, got: %d x %d", gotMapGenSettings.Width, gotMapGenSettings.Height)
	}

	// Test AdminList
	adminList := models.AdminList{"alice", "bob"}
	if err := srv.UpdateAdminList(adminList); err != nil {
		t.Fatalf("failed to update admin list: %v", err)
	}
	gotAdminList, err := srv.GetAdminList()
	if err != nil {
		t.Fatalf("failed to get admin list: %v", err)
	}
	if len(gotAdminList) != 2 || gotAdminList[0] != "alice" {
		t.Errorf("unexpected admin list: %v", gotAdminList)
	}

	// Test Whitelist
	whitelist := models.Whitelist{"charlie"}
	if err := srv.UpdateWhitelist(whitelist); err != nil {
		t.Fatalf("failed to update whitelist: %v", err)
	}
	gotWhitelist, err := srv.GetWhitelist()
	if err != nil {
		t.Fatalf("failed to get whitelist: %v", err)
	}
	if len(gotWhitelist) != 1 || gotWhitelist[0] != "charlie" {
		t.Errorf("unexpected whitelist: %v", gotWhitelist)
	}

	// Test BanList
	banList := models.BanList{"dave"}
	if err := srv.UpdateBanList(banList); err != nil {
		t.Fatalf("failed to update ban list: %v", err)
	}
	gotBanList, err := srv.GetBanList()
	if err != nil {
		t.Fatalf("failed to get ban list: %v", err)
	}
	if len(gotBanList) != 1 || gotBanList[0] != "dave" {
		t.Errorf("unexpected ban list: %v", gotBanList)
	}
}

func TestSettingsService_UninitializedListsReturnEmptyArray(t *testing.T) {
	tempDir := t.TempDir()

	cfg := config.FactorioConfig{
		ServerAdminListPath: filepath.Join(tempDir, "nonexistent-adminlist.json"),
		ServerWhiteListPath: filepath.Join(tempDir, "nonexistent-whitelist.json"),
		ServerBanListPath:   filepath.Join(tempDir, "nonexistent-banlist.json"),
	}

	srv := NewSettingsService(cfg)

	adminList, err := srv.GetAdminList()
	if err != nil {
		t.Fatalf("unexpected error getting uninitialized admin list: %v", err)
	}
	if adminList == nil {
		t.Errorf("expected non-nil empty AdminList slice")
	}

	whitelist, err := srv.GetWhitelist()
	if err != nil {
		t.Fatalf("unexpected error getting uninitialized whitelist: %v", err)
	}
	if whitelist == nil {
		t.Errorf("expected non-nil empty Whitelist slice")
	}

	banList, err := srv.GetBanList()
	if err != nil {
		t.Fatalf("unexpected error getting uninitialized ban list: %v", err)
	}
	if banList == nil {
		t.Errorf("expected non-nil empty BanList slice")
	}
}
