package services

import (
	"bufio"
	"context"
	"errors"
	"factorio/internal/config"
	"factorio/internal/factorio"
	"factorio/internal/logs"
	"factorio/internal/rcon"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	grove "github.com/StevenAlexanderJohnson/grove"
)

var (
	ErrFactorioServerNotRunning = errors.New("unable to complete the request because the factorio server isn't running")
	ErrFactorioServerError      = errors.New("an error occurred with the Factorio server")
	ErrServerAlreadyRunning     = errors.New("the factorio server is already running")
	ErrServerAlreadyStopped     = errors.New("the factorio server is already stopped")
	ErrShutdownTimeout          = errors.New("the factorio server shutdown timed out and was forcefully killed")
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

type factorioManager struct {
	ctx        context.Context
	logger     grove.ILogger
	cfgManager *config.ConfigManager
	logParser  *logs.LogParser

	// Used for managing the factorio server executable
	lock sync.Mutex
	cmd  *exec.Cmd

	// Used to communicate with the server
	stdin    io.WriteCloser
	doneChan chan error

	// Used to send messages to the goroutine
	msgChan   <-chan FactorioMessage
	fatalChan chan<- error
}

func (f *factorioManager) IsRunning() bool {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.cmd != nil
}

func (f *factorioManager) SendCommand(cmd string) (string, error) {
	f.lock.Lock()
	c := f.cmd
	stdin := f.stdin
	cfg := f.cfgManager.GetConfig().Factorio
	f.lock.Unlock()

	if c == nil {
		return "", ErrServerAlreadyStopped
	}

	if cfg.RCONPort > 0 && cfg.RCONPassword != "" {
		addr := fmt.Sprintf("127.0.0.1:%d", cfg.RCONPort)
		resp, err := rcon.ExecuteCommand(addr, cfg.RCONPassword, cmd)
		if err == nil {
			return resp, nil
		}
		f.logger.Errorf("RCON command failed: %v, falling back to stdin", err)
	}

	if stdin == nil {
		return "", ErrServerAlreadyStopped
	}

	if _, err := io.WriteString(stdin, cmd+"\n"); err != nil {
		return "", fmt.Errorf("%w: failed to send command: %v", ErrFactorioServerError, err)
	}
	return "", nil
}

func (f *factorioManager) setServerState(c *exec.Cmd, stdin io.WriteCloser, doneChan chan error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.cmd = c
	f.stdin = stdin
	f.doneChan = doneChan
}

func (f *factorioManager) clearServerState() {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.cmd = nil
	f.stdin = nil
	f.doneChan = nil
}

func (f *factorioManager) getServerState() (*exec.Cmd, io.WriteCloser, chan error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	return f.cmd, f.stdin, f.doneChan
}

func (f *factorioManager) startFactorioLoop(cfgManager *config.ConfigManager) {
	for {
		_, _, currentDoneChan := f.getServerState()

		select {
		case msg := <-f.msgChan:
			switch msg.Type {

			case FactorioStart:
				if f.IsRunning() {
					f.logger.Warning("An attempt to start the Factorio server was received while it was already running.")
					msg.Reply <- ErrServerAlreadyRunning
					continue
				}

				cfg := cfgManager.GetConfig()
				execPath := cfg.Factorio.ExecutablePath
				if execPath == "" {
					execPath = "/opt/factorio/bin/x64/factorio"
				}

				if err := f.ensureSaveFile(execPath, cfg.Factorio.SavePath); err != nil {
					f.logger.Errorf("Failed to ensure save file: %v", err)
					msg.Reply <- fmt.Errorf("%w: failed to ensure save file: %v", ErrFactorioServerError, err)
					continue
				}

				args := []string{"--start-server", cfg.Factorio.SavePath}
				if cfg.Factorio.ServerSettingsPath != "" {
					if _, err := os.Stat(cfg.Factorio.ServerSettingsPath); err == nil {
						args = append(args, "--server-settings", cfg.Factorio.ServerSettingsPath)
					}
				}
				if cfg.Factorio.ServerAdminListPath != "" {
					if _, err := os.Stat(cfg.Factorio.ServerAdminListPath); err == nil {
						args = append(args, "--server-adminlist", cfg.Factorio.ServerAdminListPath)
					}
				}
				if cfg.Factorio.ServerBanListPath != "" {
					if _, err := os.Stat(cfg.Factorio.ServerBanListPath); err == nil {
						args = append(args, "--server-banlist", cfg.Factorio.ServerBanListPath)
					}
				}
				if cfg.Factorio.UseServerWhitelist {
					args = append(args, "--use-server-whitelist")
					if cfg.Factorio.ServerWhiteListPath != "" {
						if _, err := os.Stat(cfg.Factorio.ServerWhiteListPath); err == nil {
							args = append(args, "--server-whitelist", cfg.Factorio.ServerWhiteListPath)
						}
					}
				}
				if cfg.Factorio.RCONPort > 0 {
					args = append(args, "--rcon-port", fmt.Sprintf("%d", cfg.Factorio.RCONPort))
					if cfg.Factorio.RCONPassword != "" {
						args = append(args, "--rcon-password", cfg.Factorio.RCONPassword)
					}
				}

				c := exec.Command(execPath, args...)

				pStdout, err := c.StdoutPipe()
				if err != nil {
					msg.Reply <- fmt.Errorf("%w: failed to get stdout pipe: %v", ErrFactorioServerError, err)
					continue
				}

				pStderr, err := c.StderrPipe()
				if err != nil {
					msg.Reply <- fmt.Errorf("%w: failed to get stderr pipe: %v", ErrFactorioServerError, err)
					continue
				}

				pStdin, err := c.StdinPipe()
				if err != nil {
					msg.Reply <- fmt.Errorf("%w: failed to get stdin pipe: %v", ErrFactorioServerError, err)
					continue
				}

				if err := c.Start(); err != nil {
					msg.Reply <- fmt.Errorf("%w: failed to start process: %v", ErrFactorioServerError, err)
					continue
				}

				go f.scanPipe(pStdout, os.Stdout)
				go f.scanPipe(pStderr, os.Stderr)

				doneChan := make(chan error, 1)

				// Server started successfully: initialize process state
				f.setServerState(c, pStdin, doneChan)

				// Single background goroutine strictly for waiting on process exit
				go func(proc *exec.Cmd, ch chan error) {
					ch <- proc.Wait()
				}(c, doneChan)

				f.logger.Info("Factorio server started successfully.")
				msg.Reply <- nil

			case FactorioStop:
				if !f.IsRunning() {
					f.logger.Warning("An attempt to stop the Factorio server was received while it wasn't running.")
					msg.Reply <- ErrServerAlreadyStopped
					continue
				}

				cfg := cfgManager.GetConfig().Factorio
				timeout := cfg.ShutdownTimeout
				if timeout <= 0 {
					timeout = 1 * time.Minute
				}

				_, _, doneChan := f.getServerState()

				// Graceful shutdown helper handles wait & force-kill if needed
				err := f.shutDownServer(doneChan, timeout)

				// Reset state so server can be started again
				f.clearServerState()

				msg.Reply <- err

			case FactorioUpdate:
				updated, version, err := factorio.EnsureUpdated(f.ctx, f.cfgManager.GetConfig().Factorio, f.logger)
				if err != nil {
					f.logger.Errorf("Failed to update Factorio server: %v", err)
					msg.Reply <- fmt.Errorf("%w: failed to update server: %v", ErrFactorioServerError, err)
					continue
				}
				if updated {
					f.logger.Infof("Factorio server updated to version: %s", version)
				} else {
					f.logger.Info("Factorio server is already up to date")
				}
				msg.Reply <- nil
			}

		case err := <-currentDoneChan:
			// Handles unannounced process crashes or manual in-game /quit commands
			f.logger.Errorf("Factorio server process exited unexpectedly: %v", err)
			f.clearServerState()

		case <-f.ctx.Done():
			f.logger.Warning("Received cancellation signal. Shutting down Factorio actor loop...")
			if f.IsRunning() {
				cfg := cfgManager.GetConfig().Factorio
				timeout := cfg.ShutdownTimeout
				if timeout <= 0 {
					timeout = 1 * time.Minute
				}
				_, _, doneChan := f.getServerState()
				_ = f.shutDownServer(doneChan, timeout)
			}
			return
		}
	}
}

func (f *factorioManager) shutDownServer(done chan error, timeout time.Duration) error {
	f.logger.Info("Sending /quit command to Factorio server...")

	cmd, stdin, _ := f.getServerState()
	if cmd == nil || stdin == nil {
		return ErrServerAlreadyStopped
	}

	if _, err := io.WriteString(stdin, "/quit\n"); err != nil {
		f.logger.Errorf("Failed to write /quit to stdin, force killing process: %v", err)
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = stdin.Close()
		return fmt.Errorf("%w: failed to write quit signal: %v", ErrFactorioServerError, err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(f.ctx), timeout)
	defer cancel()

	select {
	case err := <-done:
		f.logger.Info("Factorio server has shut down cleanly.")
		return err

	case <-shutdownCtx.Done():
		f.logger.Errorf("Factorio server took longer than %v to shut down. Force killing process...", timeout)
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}

		if done != nil {
			<-done
		}
		return fmt.Errorf("%w: shutdown timed out and process was force killed", ErrFactorioServerError)
	}
}

// ensureSaveFile checks if the specified save file exists, and runs `factorio --create` to generate one if missing.
func (f *factorioManager) ensureSaveFile(execPath string, savePath string) error {
	if _, err := os.Stat(savePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check save file %q: %w", savePath, err)
	}

	f.logger.Infof("Save file %q not found. Creating a new world map...", savePath)

	dir := filepath.Dir(savePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for save file %q: %w", dir, err)
		}
	}

	cmd := exec.CommandContext(f.ctx, execPath, "--create", savePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create factorio map: %w", err)
	}

	f.logger.Infof("Factorio map created successfully at %q.", savePath)
	return nil
}

func (f *factorioManager) scanPipe(r io.Reader, out io.Writer) {
	tee := io.TeeReader(r, out)
	scanner := bufio.NewScanner(tee)
	for scanner.Scan() {
		line := scanner.Text()
		if f.logParser != nil {
			if err := f.logParser.ParseLine(line); err != nil {
				f.logger.Warningf("Failed to parse log line: %v", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		f.logger.Errorf("an error occurred while scanning: %w", err)
	}
}

type FactorioService struct {
	*factorioManager

	msgChan   chan<- FactorioMessage
	fatalChan <-chan error
}

func NewFactorioService(ctx context.Context, logger grove.ILogger, cfgManager *config.ConfigManager, logParser *logs.LogParser) (*FactorioService, error) {
	msgChan := make(chan FactorioMessage, 1)
	fatalChan := make(chan error, 1)

	factorioManager := &factorioManager{
		logger:     logger,
		ctx:        ctx,
		cfgManager: cfgManager,
		logParser:  logParser,
		msgChan:    msgChan,
		fatalChan:  fatalChan,
	}
	go factorioManager.startFactorioLoop(cfgManager)

	if cfgManager.GetConfig().Factorio.StartServerOnStartup {
		reply := make(chan error, 1)
		msgChan <- FactorioMessage{
			Type:  FactorioStart,
			Reply: reply,
		}
		if err := <-reply; err != nil {
			return nil, err
		}
	}

	return &FactorioService{
		factorioManager: factorioManager,
		msgChan:         msgChan,
		fatalChan:       fatalChan,
	}, nil
}

func (f *FactorioService) FatalChan() <-chan error {
	return f.fatalChan
}

func (f *FactorioService) StartServer() error {
	replyChan := make(chan error, 1)
	f.msgChan <- FactorioMessage{
		Type:  FactorioStart,
		Reply: replyChan,
	}
	return <-replyChan
}

func (f *FactorioService) StopServer() error {
	replyChan := make(chan error, 1)
	f.msgChan <- FactorioMessage{
		Type:  FactorioStop,
		Reply: replyChan,
	}
	return <-replyChan
}

func (f *FactorioService) UpdateServer() error {
	replyChan := make(chan error, 1)

	f.msgChan <- FactorioMessage{
		Type:  FactorioStop,
		Reply: replyChan,
	}
	err := <-replyChan
	wasRunning := err == nil
	if err != nil && !errors.Is(err, ErrServerAlreadyStopped) {
		return err
	}

	f.msgChan <- FactorioMessage{
		Type:  FactorioUpdate,
		Reply: replyChan,
	}
	if err := <-replyChan; err != nil {
		return err
	}

	if wasRunning {
		f.msgChan <- FactorioMessage{
			Type:  FactorioStart,
			Reply: replyChan,
		}
		if err := <-replyChan; err != nil {
			return err
		}
	}
	return nil
}
