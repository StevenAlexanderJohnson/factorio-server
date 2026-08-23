# Factorio Docker API

A lightweight HTTP API service for managing a headless Factorio inside Docker containers.

## Overview

Factorio Docker API exposes a RESTful interface to control the server lifecycle and manipulate configuration files programmatically over HTTP without requiring direct shell access to the host or container.

Key design choices:
* **Minimal Image Footprint**: Lightweight runtime image built on `debian:bookworm-slim` with compiled Go binaries.
* **Zero-Configuration Deployment**: Starts up out-of-the-box with built-in defaults or can be customized via YAML configuration.
* **Graceful Lifecycle Management**: Handles container shutdown signals (`SIGTERM`/`SIGINT`) by issuing a graceful `/quit` signal to the Factorio server process to ensure worlds are saved before exiting.

## Features & Current Capabilities

* **Server Control**: Endpoints to start (`POST /start`), stop (`POST /stop`), and update (`POST /update`) the Factorio server.
* **In-Game Command Execution**: Send commands (`POST /command`) to the running server via RCON and receive execution responses over HTTP.
* **Manual Version Updates**: When hitting the update endpoint it will check the latest Factorio releases and handles downloading and extracting updates on startup or on demand.
* **HTTP Configuration Management**: Full `GET` and `PUT` endpoints to inspect and update Factorio configuration files on disk:
  * Server Settings (`/settings/server`)
  * Map Settings (`/settings/map`)
  * Map Generation Settings (`/settings/map-gen`)
  * Admin List (`/settings/admin-list`)
  * Whitelist (`/settings/whitelist`)
  * Ban List (`/settings/ban-list`)

## Configuration

Configuration is handled via a `config.yaml` file passed with the `-config` flag. If no config file is provided, default settings are applied automatically.

```yaml
auth:
  api_key: ""
factorio:
  executable_path: "/opt/factorio/bin/x64/factorio"
  save_path: "/factorio/data/saves/my-world.zip"
  server_settings_path: "/factorio/data/server-settings.json"
  server_adminlist_path: "/factorio/data/server-adminlist.json"
  server_banlist_path: "/factorio/data/server-banlist.json"
  server_whitelist_path: "/factorio/data/server-whitelist.json"
  use_server_whitelist: false
  map_settings_path: "/factorio/data/map-settings.json"
  map_gen_settings_path: "/factorio/data/map-gen-settings.json"
  shutdown_timeout: 1m
  auto_download_on_start: true
  rcon_port: 27015
  rcon_password: ""
```

## Running with Docker

```bash
docker run -d \
  -p 34197:34197/udp \
  -p 8080:8080 \
  -v ./data:/factorio/data \
  factorio-api:latest
```

## API Endpoints

### Health Check
* `GET /healthz` - Returns health status.

### Lifecycle Management
* `POST /start` - Starts the Factorio server process.
* `POST /stop` - Gracefully stops the server process.
* `POST /update` - Triggers an update check and restarts the server if updated.

### Command Execution
* `POST /command` - Sends a command to the Factorio server over RCON and returns the output response (`{"command": "/version"}` -> `{"response": "..."}`).

### Settings Management
* `GET /settings/server` / `PUT /settings/server` - Get/Update server settings.
* `GET /settings/map` / `PUT /settings/map` - Get/Update map settings.
* `GET /settings/map-gen` / `PUT /settings/map-gen` - Get/Update map generation settings.
* `GET /settings/admin-list` / `PUT /settings/admin-list` - Get/Update admin list.
* `GET /settings/whitelist` / `PUT /settings/whitelist` - Get/Update whitelist.
* `GET /settings/ban-list` / `PUT /settings/ban-list` - Get/Update ban list.

## Roadmap

* **Save Management**: Uploading, downloading, creating, and selecting active world save files over HTTP.
* **Log Parsing & OpenTelemetry**: Parsing Factorio stdout logs to export metrics and traces via OpenTelemetry.
* **Real-time Event Streaming**: Webhooks or Server-Sent Events (SSE) for player join/leave events, chat, and server status.
* **Mod Management**: Downloading, updating, enabling, and disabling Factorio mods via API.

## Contributing

If there is a feature that you would like see added, please open an issue first to see if it aligns with path I want this project to go. I am open to adding cool features, but I would like to keep this project pretty slim and focused. I like having new features work out of the box without extra configuration. If configuration is required a sensable default should be provided. Quality of life changes should always be accepted, but may need to be prioritized.

If you don't open an issue first for dicussion and provide a PR it will almost certainly be rejected.

See the next section for how I use AI in this project. If you do open a PR I would like that AI use is disclaimed. I'm not against the use of AI, I would just like to know.

There is no strict standards that this project uses. So until that is set up, if you want to contribute try your best to follow the patterns that I've already used. I may ask for a refactor or update to your code if it doesn't match closely enough.


## AI Use Disclaimer

I do use AI to assist in development. All code in this project is either written or reviewed by me. No code makes it into the main branch without thorough review or was hand written.

Some use cases that I use AI for as follows:
  * Creating logs
  * Boilerplate code (see get and update methods in [settings](/internal/services/settings.go) for an example) 
  * Brainstorming
  * Repetitive tasks
  * Configuration parsing