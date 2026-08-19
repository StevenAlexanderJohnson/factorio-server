package factorio

import (
	"context"
	"errors"
	"factorio/internal/config"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	grove "github.com/StevenAlexanderJohnson/grove"
)


var (
	ErrFactorioServerError  = errors.New("an error occurred with the Factorio server")
	ErrServerAlreadyRunning = errors.New("the factorio server is already running")
	ErrServerAlreadyStopped = errors.New("the factorio server is already stopped")
	ErrShutdownTimeout      = errors.New("the factorio server shutdown timed out and was forcefully killed")
)

type FactorioMessageType int

const (
	FactorioStart FactorioMessageType = iota
	FactorioStop
	FactorioUpdate
)

type FactorioMessage struct {
	Type  FactorioMessageType
	Reply chan<- error
}

func StartFactorioLoop(ctx context.Context, cfg config.FactorioConfig) (chan<- FactorioMessage, <-chan error) {
	msgChan := make(chan FactorioMessage, 1)
	fatalChan := make(chan error, 1)

	go func() {
		logger := grove.NewDefaultLogger("FactorioServer")

		var (
			cmd      *exec.Cmd
			stdin    io.WriteCloser
			doneChan chan error // nil when stopped; active when running
		)

		for {
			select {
			case msg := <-msgChan:
				switch msg.Type {

				case FactorioStart:
					if cmd != nil {
						logger.Warning("An attempt to start the Factorio server was received while it was already running.")
						msg.Reply <- ErrServerAlreadyRunning
						continue
					}

					execPath := cfg.ExecutablePath
					if execPath == "" {
						execPath = "/opt/factorio/bin/x64/factorio"
					}

					if err := ensureSaveFile(ctx, execPath, cfg.SavePath, logger); err != nil {
						logger.Errorf("Failed to ensure save file: %v", err)
						msg.Reply <- fmt.Errorf("%w: failed to ensure save file: %v", ErrFactorioServerError, err)
						continue
					}

					args := []string{"--start-server", cfg.SavePath}
					if cfg.ServerSettingsPath != "" {
						args = append(args, "--server-settings", cfg.ServerSettingsPath)
					}
					if cfg.ServerAdminListPath != "" {
						args = append(args, "--server-adminlist", cfg.ServerAdminListPath)
					}
					if cfg.ServerBanListPath != "" {
						args = append(args, "--server-banlist", cfg.ServerBanListPath)
					}
					if cfg.ServerWhiteListPath != "" {
						args = append(args, "--server-whitelist", cfg.ServerWhiteListPath)
					}

					c := exec.Command(execPath, args...)
					c.Stdout = os.Stdout
					c.Stderr = os.Stderr


					pStdin, err := c.StdinPipe()
					if err != nil {
						msg.Reply <- fmt.Errorf("%w: failed to get stdin pipe: %v", ErrFactorioServerError, err)
						continue
					}

					if err := c.Start(); err != nil {
						msg.Reply <- fmt.Errorf("%w: failed to start process: %v", ErrFactorioServerError, err)
						continue
					}

					// Server started successfully: initialize process state
					cmd = c
					stdin = pStdin
					doneChan = make(chan error, 1)

					// Single background goroutine strictly for waiting on process exit
					go func(proc *exec.Cmd, ch chan error) {
						ch <- proc.Wait()
					}(cmd, doneChan)

					logger.Info("Factorio server started successfully.")
					msg.Reply <- nil

				case FactorioStop:
					if cmd == nil {
						logger.Warning("An attempt to stop the Factorio server was received while it wasn't running.")
						msg.Reply <- ErrServerAlreadyStopped
						continue
					}

					timeout := cfg.ShutdownTimeout
					if timeout <= 0 {
						timeout = 1 * time.Minute
					}

					// Graceful shutdown helper handles wait & force-kill if needed
					err := shutDownServer(ctx, cmd, logger, stdin, doneChan, timeout)

					// Reset state so server can be started again
					cmd = nil
					stdin = nil
					doneChan = nil

					msg.Reply <- err

				case FactorioUpdate:
					// Update logic can be plugged in here safely!
					msg.Reply <- errors.New("update not implemented")
				}

			case err := <-doneChan:
				// Handles unannounced process crashes or manual in-game /quit commands
				logger.Errorf("Factorio server process exited unexpectedly: %v", err)
				cmd = nil
				stdin = nil
				doneChan = nil

			case <-ctx.Done():
				logger.Warning("Received cancellation signal. Shutting down Factorio actor loop...")
				if cmd != nil {
					timeout := cfg.ShutdownTimeout
					if timeout <= 0 {
						timeout = 1 * time.Minute
					}
					_ = shutDownServer(ctx, cmd, logger, stdin, doneChan, timeout)
				}
				return
			}
		}
	}()

	return msgChan, fatalChan
}

func shutDownServer(ctx context.Context, cmd *exec.Cmd, logger grove.ILogger, stdin io.WriteCloser, done <-chan error, timeout time.Duration) error {
	logger.Info("Sending /quit command to Factorio server...")

	if _, err := io.WriteString(stdin, "/quit\n"); err != nil {
		logger.Errorf("Failed to write /quit to stdin, force killing process: %v", err)
		_ = cmd.Process.Kill()
		return fmt.Errorf("%w: failed to write quit signal: %v", ErrFactorioServerError, err)
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case err := <-done:
		logger.Info("Factorio server has shut down cleanly.")
		return err

	case <-shutdownCtx.Done():
		logger.Errorf("Factorio server took longer than %v to shut down. Force killing process...", timeout)
		_ = cmd.Process.Kill()

		// Drain the wait channel after force killing
		<-done
		return fmt.Errorf("%w: shutdown timed out and process was force killed", ErrFactorioServerError)
	}
}

// ensureSaveFile checks if the specified save file exists, and runs `factorio --create` to generate one if missing.
func ensureSaveFile(ctx context.Context, execPath string, savePath string, logger grove.ILogger) error {
	if _, err := os.Stat(savePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check save file %q: %w", savePath, err)
	}

	logger.Infof("Save file %q not found. Creating a new world map...", savePath)

	dir := filepath.Dir(savePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for save file %q: %w", dir, err)
		}
	}

	cmd := exec.CommandContext(ctx, execPath, "--create", savePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create factorio map: %w", err)
	}

	logger.Infof("Factorio map created successfully at %q.", savePath)
	return nil
}


