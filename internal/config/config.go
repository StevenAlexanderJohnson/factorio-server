package config

import (
	"time"
)

type Config struct {
	Auth     Auth           `yaml:"auth" json:"auth"`
	Factorio FactorioConfig `yaml:"factorio" json:"factorio"`
	Logs     LogsConfig     `yaml:"logs" json:"logs"`
}

type LogsConfig struct {
	LogParserScript string `yaml:"log_parser_script,omitempty" json:"log_parser_script,omitempty"`
}

type Auth struct {
	ApiKey string `yaml:"api_key" json:"api_key"`
}

type FactorioConfig struct {
	StartServerOnStartup bool          `yaml:"start_server_on_startup" json:"start_server_on_startup"`
	ExecutablePath       string        `yaml:"executable_path" json:"executable_path"`
	SavePath             string        `yaml:"save_path" json:"save_path"`
	ServerSettingsPath   string        `yaml:"server_settings_path,omitempty" json:"server_settings_path,omitempty"`
	ServerAdminListPath  string        `yaml:"server_adminlist_path,omitempty" json:"server_adminlist_path,omitempty"`
	ServerBanListPath    string        `yaml:"server_banlist_path,omitempty" json:"server_banlist_path,omitempty"`
	ServerWhiteListPath  string        `yaml:"server_whitelist_path,omitempty" json:"server_whitelist_path,omitempty"`
	UseServerWhitelist   bool          `yaml:"use_server_whitelist" json:"use_server_whitelist"`
	MapSettingsPath      string        `yaml:"map_settings_path,omitempty" json:"map_settings_path,omitempty"`
	MapGenSettingsPath   string        `yaml:"map_gen_settings_path,omitempty" json:"map_gen_settings_path,omitempty"`
	ShutdownTimeout      time.Duration `yaml:"shutdown_timeout" json:"shutdown_timeout"`
	DownloadURL          string        `yaml:"download_url,omitempty" json:"download_url,omitempty"`
	VersionFilePath      string        `yaml:"version_file_path,omitempty" json:"version_file_path,omitempty"`
	ExtractDir           string        `yaml:"extract_dir,omitempty" json:"extract_dir,omitempty"`
	TempArchivePath      string        `yaml:"temp_archive_path,omitempty" json:"temp_archive_path,omitempty"`
	AutoDownloadOnStart  bool          `yaml:"auto_download_on_start" json:"auto_download_on_start"`
	RCONPort             int           `yaml:"rcon_port,omitempty" json:"rcon_port,omitempty"`
	RCONPassword         string        `yaml:"rcon_password,omitempty" json:"rcon_password,omitempty"`
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
		Logs: LogsConfig{
			LogParserScript: "/factorio/data/log-parser.lua",
		},
	}
}
