package input

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"time"
)

type YdotoolInputSender struct {
	bin        string
	socketPath string
}

func NewYdotoolInputSender(bin string, socketPath string) (InputSender, error) {
	sp := path.Join(os.Getenv("HOME"), ".ydotool_socket")

	return &YdotoolInputSender{bin: bin, socketPath: sp}, nil
}

// connects to the ydotoold Unix socket to verify the daemon is running.
func ensureYdotooldRunning(socketPath string, timeout time.Duration) error {
	// Single upfront check: if the socket doesn't exist, fail fast
	if _, statErr := os.Stat(socketPath); statErr != nil {
		if os.IsNotExist(statErr) {
			return fmt.Errorf("ydotoold socket not found at %s: %w", socketPath, statErr)
		}
		return fmt.Errorf("cannot stat ydotoold socket at %s: %w", socketPath, statErr)
	}

	deadline := time.Now().Add(timeout)
	attempt := 1
	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return fmt.Errorf("ydotoold not available: could not connect to %s within %s", socketPath, timeout)
		}

		// Use a short per-attempt timeout to avoid long hangs
		attemptTimeout := 300 * time.Millisecond
		if remaining < attemptTimeout {
			attemptTimeout = remaining
		}

		conn, err := net.DialTimeout("unixgram", socketPath, attemptTimeout)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		// Brief backoff before retrying
		sleep := 100 * time.Millisecond
		if remaining < sleep {
			sleep = remaining
		}
		log.Printf("[ydotoold] running check attempt %d failed: %v. Retrying in %s (socket: %s)", attempt, err, sleep, socketPath)
		attempt++
		time.Sleep(sleep)
	}
}

// Sends a command to ydotool using specified binary location and socket path
func (s *YdotoolInputSender) sendCommand(args ...string) error {
	cmd := exec.Command(s.bin, args...)
	cmd.Env = append(os.Environ(), "YDOTOOL_SOCKET="+s.socketPath)
	return cmd.Run()
}

func (s *YdotoolInputSender) Close() error {
	return nil
}

func (s *YdotoolInputSender) SendClick() error {
	// Left-click
	return s.sendCommand("click", "0xC0")
}

func (s *YdotoolInputSender) SendKeyPress(keyCode string) error {
	return s.sendCommand("key", keyCode)
}
