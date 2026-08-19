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

func DownloadFactorioServer(ctx context.Context, downloadUrl string, destPath string, fatalChan chan<- error) {
	client := &http.Client{
		// Big timeout period just in case that it's hosted on a slow machine or slow internet
		Timeout: 30 * time.Minute,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		fatalChan <- fmt.Errorf("failed to create request to download factorio: %w", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fatalChan <- fmt.Errorf("http request to download factorio failed: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fatalChan <- fmt.Errorf("unexpected status code while downloading factorio: %d", resp.StatusCode)
		return
	}

	out, err := os.Create(destPath + ".tmp")
	if err != nil {
		fatalChan <- fmt.Errorf("failed to create destination file for factorio: %w", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(destPath + ".tmp")
		fatalChan <- fmt.Errorf("failed during downloading factorio: %w", err)
		return
	}

	if err := os.Rename(destPath+".tmp", destPath); err != nil {
		_ = os.Remove(destPath + ".tmp")
		fatalChan <- fmt.Errorf("failed during renamed temp factorio file: %w", err)
		return
	}

	cmd := exec.CommandContext(ctx, "tar", "-xJf", destPath, "-C", destPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fatalChan <- fmt.Errorf("tar extraction failed: %w (output: %s)", err, string(output))
		return
	}
}
