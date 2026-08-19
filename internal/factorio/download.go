package factorio

import (
	"context"
	"encoding/json"
	"errors"
	"factorio/internal/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	grove "github.com/StevenAlexanderJohnson/grove"
)

type FactorioVersionRelease struct {
	Alpha     string `json:"alpha"`
	Demo      string `json:"demo"`
	Expansion string `json:"expansion"`
	Headless  string `json:"headless"`
}

type FactorioVersionResponse struct {
	Experimental FactorioVersionRelease `json:"experimental"`
	Stable       FactorioVersionRelease `json:"stable"`
}

// CompareVersions compares two dot-separated version strings (e.g., "1.1.107" and "2.0.28").
// It returns -1 if v1 < v2, 0 if v1 == v2, and 1 if v1 > v2.
func CompareVersions(v1, v2 string) int {
	v1 = strings.TrimSpace(strings.TrimPrefix(v1, "v"))
	v2 = strings.TrimSpace(strings.TrimPrefix(v2, "v"))

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) && parts1[i] != "" {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) && parts2[i] != "" {
			n2, _ = strconv.Atoi(parts2[i])
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}

	return 0
}

// ReadDiskVersion reads and trims the version string from the given version.txt path.
// If the file does not exist, it returns an empty string without an error.
func ReadDiskVersion(versionFilePath string) (string, error) {
	if versionFilePath == "" {
		versionFilePath = "version.txt"
	}

	data, err := os.ReadFile(versionFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read version file %q: %w", versionFilePath, err)
	}

	return strings.TrimSpace(string(data)), nil
}

// UpdateVersionFile writes the given version string to the version.txt file on disk.
func UpdateVersionFile(versionFilePath string, version string) error {
	if versionFilePath == "" {
		versionFilePath = "version.txt"
	}

	dir := filepath.Dir(versionFilePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for version file: %w", err)
		}
	}

	content := strings.TrimSpace(version) + "\n"
	if err := os.WriteFile(versionFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write version file %q: %w", versionFilePath, err)
	}

	return nil
}

// FetchLatestReleases queries the Factorio API for the latest version metadata.
func FetchLatestReleases(ctx context.Context) (*FactorioVersionResponse, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://factorio.com/api/latest-releases", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to check factorio version: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request to check factorio version failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code while checking factorio version: %d", resp.StatusCode)
	}

	var versionResponse FactorioVersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResponse); err != nil {
		return nil, fmt.Errorf("failed to parse factorio version response: %w", err)
	}

	return &versionResponse, nil
}

// NeedsDownload checks the version.txt file on disk against the latest release from the Factorio API.
// It returns true if the disk version is missing or smaller than the latest version.
func NeedsDownload(ctx context.Context, versionFilePath string) (bool, string, error) {
	releases, err := FetchLatestReleases(ctx)
	if err != nil {
		return false, "", err
	}

	latestVersion := releases.Stable.Headless
	if latestVersion == "" {
		return false, "", errors.New("latest stable headless version not found in API response")
	}

	diskVersion, err := ReadDiskVersion(versionFilePath)
	if err != nil {
		return false, "", err
	}

	if diskVersion == "" {
		return true, latestVersion, nil
	}

	if CompareVersions(diskVersion, latestVersion) < 0 {
		return true, latestVersion, nil
	}

	return false, latestVersion, nil
}

// DownloadFactorioServer downloads the Factorio headless server archive, extracts it,
// and updates version.txt on disk with the downloaded version.
func DownloadFactorioServer(ctx context.Context, downloadUrl string, destPath string, extractDir string, versionFilePath string, version string) error {
	client := &http.Client{
		// Big timeout period just in case that it's hosted on a slow machine or slow internet
		Timeout: 30 * time.Minute,
	}

	if downloadUrl == "" {
		downloadUrl = "https://factorio.com/get-download/stable/headless/linux64"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request to download factorio: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request to download factorio failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code while downloading factorio: %d", resp.StatusCode)
	}

	destDir := filepath.Dir(destPath)
	if destDir != "" && destDir != "." {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for destination archive: %w", err)
		}
	}

	out, err := os.Create(destPath + ".tmp")
	if err != nil {
		return fmt.Errorf("failed to create destination file for factorio: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(destPath + ".tmp")
		return fmt.Errorf("failed during downloading factorio: %w", err)
	}

	if err := os.Rename(destPath+".tmp", destPath); err != nil {
		_ = os.Remove(destPath + ".tmp")
		return fmt.Errorf("failed during renamed temp factorio file: %w", err)
	}

	if extractDir == "" {
		extractDir = filepath.Dir(destPath)
	}

	cmd := exec.CommandContext(ctx, "tar", "-xJf", destPath, "-C", extractDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar extraction failed: %w (output: %s)", err, string(output))
	}

	if version != "" {
		if err := UpdateVersionFile(versionFilePath, version); err != nil {
			return fmt.Errorf("failed to update version file after download: %w", err)
		}
	}

	return nil
}

// EnsureUpdated checks if a new version is available and downloads/extracts it if needed.
// It returns true if an update was performed, along with the latest version.
func EnsureUpdated(ctx context.Context, cfg config.FactorioConfig, logger grove.ILogger) (bool, string, error) {
	needsDownload, latestVersion, err := NeedsDownload(ctx, cfg.VersionFilePath)
	if err != nil {
		return false, "", fmt.Errorf("failed to check for factorio updates: %w", err)
	}

	if !needsDownload {
		if logger != nil {
			logger.Infof("Factorio is already up-to-date (version %s).", latestVersion)
		}
		return false, latestVersion, nil
	}

	if logger != nil {
		logger.Infof("New Factorio version available: %s. Starting download and extraction...", latestVersion)
	}

	downloadURL := cfg.DownloadURL
	if downloadURL == "" {
		downloadURL = fmt.Sprintf("https://factorio.com/get-download/%s/headless/linux64", latestVersion)
	}

	tempArchive := cfg.TempArchivePath
	if tempArchive == "" {
		tempArchive = filepath.Join(os.TempDir(), "factorio.tar.xz")
	}

	extractDir := cfg.ExtractDir
	if extractDir == "" {
		extractDir = "/opt"
	}

	if err := DownloadFactorioServer(ctx, downloadURL, tempArchive, extractDir, cfg.VersionFilePath, latestVersion); err != nil {
		return false, latestVersion, fmt.Errorf("failed to download and update factorio: %w", err)
	}

	// Clean up temp archive file after successful extraction
	_ = os.Remove(tempArchive)

	if logger != nil {
		logger.Infof("Factorio successfully updated to version %s.", latestVersion)
	}

	return true, latestVersion, nil
}
