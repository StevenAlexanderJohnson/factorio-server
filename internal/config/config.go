package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Auth     Auth           `yaml:"auth"`
	Factorio FactorioConfig `yaml:"factorio"`
}

type Auth struct {
	ApiKey string `yaml:"api_key"`
}

type FactorioConfig struct {
	StartServerOnStartup bool          `yaml:"start_server_on_startup"`
	ExecutablePath       string        `yaml:"executable_path"`
	SavePath             string        `yaml:"save_path"`
	ServerSettingsPath   string        `yaml:"server_settings_path,omitempty"`
	ServerAdminListPath  string        `yaml:"server_adminlist_path,omitempty"`
	ServerBanListPath    string        `yaml:"server_banlist_path,omitempty"`
	ServerWhiteListPath  string        `yaml:"server_whitelist_path,omitempty"`
	UseServerWhitelist   bool          `yaml:"use_server_whitelist"`
	MapSettingsPath      string        `yaml:"map_settings_path,omitempty"`
	MapGenSettingsPath   string        `yaml:"map_gen_settings_path,omitempty"`
	ShutdownTimeout      time.Duration `yaml:"shutdown_timeout"`
	DownloadURL          string        `yaml:"download_url,omitempty"`
	VersionFilePath      string        `yaml:"version_file_path,omitempty"`
	ExtractDir           string        `yaml:"extract_dir,omitempty"`
	TempArchivePath      string        `yaml:"temp_archive_path,omitempty"`
	AutoDownloadOnStart  bool          `yaml:"auto_download_on_start"`
	RCONPort             int           `yaml:"rcon_port,omitempty"`
	RCONPassword         string        `yaml:"rcon_password,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		Auth: Auth{
			ApiKey: "",
		},
		Factorio: FactorioConfig{
			ExecutablePath:      "/opt/factorio/bin/x64/factorio",
			SavePath:            "/factorio/data/saves/my-world.zip",
			ServerSettingsPath:  "/factorio/data/server-settings.json",
			ServerAdminListPath: "/factorio/data/server-adminlist.json",
			ServerBanListPath:   "/factorio/data/server-banlist.json",
			ServerWhiteListPath: "/factorio/data/server-whitelist.json",
			MapSettingsPath:     "/factorio/data/map-settings.json",
			MapGenSettingsPath:  "/factorio/data/map-gen-settings.json",
			VersionFilePath:     "/factorio/version.txt",
			ExtractDir:          "/opt",
			TempArchivePath:     "/tmp/factorio.tar.xz",
			ShutdownTimeout:     1 * time.Minute,
			AutoDownloadOnStart: true,
			RCONPort:            27015,
			RCONPassword:        "",
		},
	}
}

func createDefaultApiKey() (string, error) {
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return "", fmt.Errorf("an error occurred while generating a default api key: %w", err)
	}
	apiKey := base64.StdEncoding.EncodeToString(apiKeyBytes)
	return apiKey, nil
}

// LoadConfig loads configuration from a YAML file.
// If path is empty, it checks the CONFIG_PATH env var, or defaults to "/factorio/data/config.yaml".
func LoadConfig(path string) (*Config, error) {
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
		logger := grove.NewDefaultLogger("Config")
		logger.Infof("The API Key was empty in the configuration file. A new one was generated and written to the config file.")
	}

	if cfg.Factorio.RCONPort > 0 && cfg.Factorio.RCONPassword == "" {
		rconPass, err := createDefaultApiKey()
		if err != nil {
			return nil, err
		}
		cfg.Factorio.RCONPassword = rconPass
		updatedConfig = true
		logger := grove.NewDefaultLogger("Config")
		logger.Infof("The RCON password was empty in the configuration file. A new one was generated and written to the config file.")
	}

	if updatedConfig {
		cfgBytes, err := yaml.Marshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config after generating defaults: %w", err)
		}
		if err := os.WriteFile(path, cfgBytes, os.FileMode(os.O_TRUNC|os.O_WRONLY)); err != nil {
			return nil, fmt.Errorf("failed to write config back to file after generating defaults: %w", err)
		}
	}

	return cfg, nil
}
