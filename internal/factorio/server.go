package factorio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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

func StartFactorioLoop(ctx context.Context, savePath string) (chan<- FactorioMessage, <-chan error) {
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
						msg.Reply <- errors.New("server is already running")
						continue
					}

					c := exec.Command("/opt/factorio/bin/x64/factorio", "--start-server", savePath)
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

					// Graceful shutdown helper handles wait & force-kill if needed
					err := shutDownServer(ctx, cmd, logger, stdin, doneChan)

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
					_ = shutDownServer(ctx, cmd, logger, stdin, doneChan)
				}
				return
			}
		}
	}()

	return msgChan, fatalChan
}

func shutDownServer(ctx context.Context, cmd *exec.Cmd, logger grove.ILogger, stdin io.WriteCloser, done <-chan error) error {
	logger.Info("Sending /quit command to Factorio server...")

	if _, err := io.WriteString(stdin, "/quit\n"); err != nil {
		logger.Errorf("Failed to write /quit to stdin, force killing process: %v", err)
		_ = cmd.Process.Kill()
		return fmt.Errorf("%w: failed to write quit signal: %v", ErrFactorioServerError, err)
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	select {
	case err := <-done:
		logger.Info("Factorio server has shut down cleanly.")
		return err

	case <-shutdownCtx.Done():
		logger.Error("Factorio server took longer than 1 minute to shut down. Force killing process...")
		_ = cmd.Process.Kill()

		// Drain the wait channel after force killing
		<-done
		return fmt.Errorf("%w: shutdown timed out and process was force killed", ErrFactorioServerError)
	}
}
