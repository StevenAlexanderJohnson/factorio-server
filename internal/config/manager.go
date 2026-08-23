package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"sync"

	"github.com/StevenAlexanderJohnson/grove"
	"gopkg.in/yaml.v3"
)

func createDefaultApiKey() (string, error) {
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", fmt.Errorf("an error occurred while generating a default api key: %w", err)
	}
	apiKey := base64.StdEncoding.EncodeToString(apiKeyBytes)
	return apiKey, nil
}

type ConfigManager struct {
	lock sync.RWMutex

	logger grove.ILogger
	path   string
	cfg    Config
}

func NewConfigManager(path string) (*ConfigManager, error) {
	logger := grove.NewDefaultLogger("Config")
	if path == "" {
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			path = envPath
		} else {
			path = "/factorio/data/config.yaml"
		}
	}

	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file at %q: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config: %w", err)
	}

	updatedConfig := false

	if cfg.Auth.ApiKey == "" {
		defaultApiKey, err := createDefaultApiKey()
		if err != nil {
			return nil, err
		}
		cfg.Auth.ApiKey = defaultApiKey
		updatedConfig = true
		logger.Infof("The API Key was empty in the configuration file. A new one was generated and written to the config file.")
	}

	if cfg.Factorio.RCONPort > 0 && cfg.Factorio.RCONPassword == "" {
		rconPass, err := createDefaultApiKey()
		if err != nil {
			return nil, err
		}
		cfg.Factorio.RCONPassword = rconPass
		updatedConfig = true
		logger.Infof("The RCON password was empty in the configuration file. A new one was generated and written to the config file.")
	}

	manager := &ConfigManager{
		logger: logger,
		cfg:    *cfg,
		path:   path,
	}

	if updatedConfig {
		if err := manager.flushConfig(); err != nil {
			return nil, fmt.Errorf("failed to write generated defaults to the config file: %w", err)
		}
	}

	return manager, nil
}

func (c *ConfigManager) flushConfig() error {
	cfgBytes, err := yaml.Marshal(c.cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config into yaml bytes: %w", err)
	}
	if err := os.WriteFile(c.path, cfgBytes, os.FileMode(os.O_TRUNC|os.O_WRONLY)); err != nil {
		return fmt.Errorf("failed to flush config back to the file: %w", err)
	}
	c.logger.Infof("successfully flushed config file to %s", c.path)
	return nil
}

func (c *ConfigManager) GetConfig() Config {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cfg
}

func (c *ConfigManager) UpdateSavePath(path string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cfg.Factorio.SavePath = path
	return c.flushConfig()
}

func (c *ConfigManager) UpdateFactorioConfig(cfg FactorioConfig) error {
	c.cfg.Factorio = cfg
	return c.flushConfig()
}
