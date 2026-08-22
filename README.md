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
* **RCON Support**: Remote console protocol support for real-time command execution and query support.
* **Real-time Event Streaming**: Webhooks or Server-Sent Events (SSE) for player join/leave events, chat, and server status.
* **Mod Management**: Downloading, updating, enabling, and disabling Factorio mods via API.
