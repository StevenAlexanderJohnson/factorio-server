package factorio

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func DownloadFactorioServer(ctx context.Context, downloadUrl string, destPath string) error {
	client := &http.Client{
		// Big timeout period just in case that it's hosted on a slow machine or slow internet
		Timeout: 30 * time.Minute,
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

	cmd := exec.CommandContext(ctx, "tar", "-xJf", destPath, "-C", destPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tar extraction failed: %w (output: %s)", err, string(output))
	}
	return nil
}
