package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"factorio/internal/config"
	"factorio/internal/factorio/models"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrFailedToReadSettings   = errors.New("failed to read settings file")
	ErrFailedToWriteSettings  = errors.New("failed to write settings file")
	ErrFailedToDecodeSettings = errors.New("failed to decode settings JSON")
	ErrFailedToEncodeSettings = errors.New("failed to encode settings JSON")
)

type SettingsService struct {
	cfgManager *config.ConfigManager
}

func NewSettingsService(cfgManager *config.ConfigManager) SettingsService {
	return SettingsService{
		cfgManager: cfgManager,
	}
}

func readJSONSettings[T any](targetPath string, examplePath string) (T, error) {
	var result T
	if targetPath == "" {
		targetPath = examplePath
	}

	fileBytes, err := os.ReadFile(targetPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && examplePath != "" {
			fileBytes, err = os.ReadFile(examplePath)
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return result, nil
				}
				return result, fmt.Errorf("%w: failed to read example file %q: %v", ErrFailedToReadSettings, examplePath, err)
			}
		} else {
			return result, fmt.Errorf("%w: failed to read file %q: %v", ErrFailedToReadSettings, targetPath, err)
		}
	}

	if len(bytes.TrimSpace(fileBytes)) == 0 {
		return result, nil
	}

	if err := json.Unmarshal(fileBytes, &result); err != nil {
		return result, fmt.Errorf("%w: failed to unmarshal %q: %v", ErrFailedToDecodeSettings, targetPath, err)
	}

	return result, nil
}

func writeJSONSettings[T any](targetPath string, data T) error {
	if targetPath == "" {
		return fmt.Errorf("%w: target file path is empty", ErrFailedToWriteSettings)
	}

	dir := filepath.Dir(targetPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("%w: failed to create directory %q: %v", ErrFailedToWriteSettings, dir, err)
		}
	}

	fileBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("%w: failed to marshal data: %v", ErrFailedToEncodeSettings, err)
	}

	if err := os.WriteFile(targetPath, fileBytes, 0644); err != nil {
		return fmt.Errorf("%w: failed to write file %q: %v", ErrFailedToWriteSettings, targetPath, err)
	}

	return nil
}

func (s *SettingsService) GetServerSettings() (models.ServerSettings, error) {
	cfg := s.cfgManager.GetConfig()
	return readJSONSettings[models.ServerSettings](cfg.Factorio.ServerSettingsPath, "/opt/factorio/data/server-settings.example.json")
}

func (s *SettingsService) UpdateServerSettings(settings models.ServerSettings) error {
	cfg := s.cfgManager.GetConfig()
	return writeJSONSettings(cfg.Factorio.ServerSettingsPath, settings)
}

func (s *SettingsService) GetMapSettings() (models.MapSettings, error) {
	cfg := s.cfgManager.GetConfig()
	return readJSONSettings[models.MapSettings](cfg.Factorio.MapSettingsPath, "/opt/factorio/data/map-settings.example.json")
}

func (s *SettingsService) UpdateMapSettings(settings models.MapSettings) error {
	cfg := s.cfgManager.GetConfig()
	return writeJSONSettings(cfg.Factorio.MapSettingsPath, settings)
}

func (s *SettingsService) GetMapGenSettings() (models.MapGenSettings, error) {
	cfg := s.cfgManager.GetConfig()
	return readJSONSettings[models.MapGenSettings](cfg.Factorio.MapGenSettingsPath, "/opt/factorio/data/map-gen-settings.example.json")
}

func (s *SettingsService) UpdateMapGenSettings(settings models.MapGenSettings) error {
	cfg := s.cfgManager.GetConfig()
	return writeJSONSettings(cfg.Factorio.MapGenSettingsPath, settings)
}

func (s *SettingsService) GetAdminList() (models.AdminList, error) {
	cfg := s.cfgManager.GetConfig()
	list, err := readJSONSettings[models.AdminList](cfg.Factorio.ServerAdminListPath, "/opt/factorio/data/server-adminlist.example.json")
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = models.AdminList{}
	}
	return list, nil
}

func (s *SettingsService) UpdateAdminList(list models.AdminList) error {
	if list == nil {
		list = models.AdminList{}
	}
	cfg := s.cfgManager.GetConfig()
	return writeJSONSettings(cfg.Factorio.ServerAdminListPath, list)
}

func (s *SettingsService) GetWhitelist() (models.Whitelist, error) {
	cfg := s.cfgManager.GetConfig()
	list, err := readJSONSettings[models.Whitelist](cfg.Factorio.ServerWhiteListPath, "/opt/factorio/data/server-whitelist.example.json")
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = models.Whitelist{}
	}
	return list, nil
}

func (s *SettingsService) UpdateWhitelist(list models.Whitelist) error {
	if list == nil {
		list = models.Whitelist{}
	}
	cfg := s.cfgManager.GetConfig()
	return writeJSONSettings(cfg.Factorio.ServerWhiteListPath, list)
}

func (s *SettingsService) GetBanList() (models.BanList, error) {
	cfg := s.cfgManager.GetConfig()
	list, err := readJSONSettings[models.BanList](cfg.Factorio.ServerBanListPath, "/opt/factorio/data/server-banlist.example.json")
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = models.BanList{}
	}
	return list, nil
}

func (s *SettingsService) UpdateBanList(list models.BanList) error {
	if list == nil {
		list = models.BanList{}
	}
	cfg := s.cfgManager.GetConfig()
	return writeJSONSettings(cfg.Factorio.ServerBanListPath, list)
}
