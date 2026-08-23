package config

import (
	"time"
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
