package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Factorio FactorioConfig `yaml:"factorio"`
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
}

func DefaultConfig() *Config {
	return &Config{
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
		},
	}
}

// LoadConfig loads configuration from a YAML file.
// If path is empty, it checks the CONFIG_PATH env var, or defaults to "config.yaml".
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			path = envPath
		} else {
			path = "config.yaml"
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

	return cfg, nil
}
