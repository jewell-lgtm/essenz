// Package daemon provides a client to communicate with the Chrome daemon.
package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/jewell-lgtm/essenz/internal/pageready"
)

// Client communicates with the Chrome daemon.
type Client struct {
	socketPath string
}

// NewDaemonClient creates a new daemon client.
func NewDaemonClient() *Client {
	socketPath := filepath.Join(os.TempDir(), "essenz-daemon.sock")
	return &Client{
		socketPath: socketPath,
	}
}

// FetchContent fetches content via the daemon.
func (c *Client) FetchContent(_ context.Context, url string) (string, error) {
	// Ensure daemon is running
	if !IsDaemonRunning() {
		if err := StartDaemonIfNeeded(); err != nil {
			return "", fmt.Errorf("failed to start daemon: %w", err)
		}
		// Give daemon time to start
		time.Sleep(1 * time.Second)
	}

	// Connect to daemon
	conn, err := net.DialTimeout("unix", c.socketPath, 5*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to connect to daemon: %w", err)
	}
	defer func() { _ = conn.Close() }()

	// Set connection timeout
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	// Send request
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	req := Request{
		Action: "fetch",
		URL:    url,
	}

	if err := encoder.Encode(req); err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	// Read response
	var resp Response
	if err := decoder.Decode(&resp); err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if !resp.Success {
		return "", fmt.Errorf("daemon error: %s", resp.Error)
	}

	return resp.Content, nil
}

// FetchContentWithReadiness fetches content via the daemon with DOM readiness detection.
func (c *Client) FetchContentWithReadiness(ctx context.Context, url string, _ *pageready.ReadinessChecker) (string, error) {
	// For now, implement this by falling back to regular fetch
	// TODO: Extend the daemon protocol to support readiness checking
	content, err := c.FetchContent(ctx, url)
	if err != nil {
		return "", err
	}

	// TODO: In future iterations, we'll integrate the readiness checker
	// into the daemon server for proper DOM event waiting

	return content, nil
}

// Ping checks if the daemon is responsive.
func (c *Client) Ping() error {
	conn, err := net.DialTimeout("unix", c.socketPath, 2*time.Second)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	req := Request{Action: "ping"}
	if err := encoder.Encode(req); err != nil {
		return err
	}

	var resp Response
	if err := decoder.Decode(&resp); err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf("ping failed: %s", resp.Error)
	}

	return nil
}

// Shutdown requests the daemon to shutdown.
func (c *Client) Shutdown() error {
	if !IsDaemonRunning() {
		return nil
	}

	conn, err := net.DialTimeout("unix", c.socketPath, 2*time.Second)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	encoder := json.NewEncoder(conn)
	req := Request{Action: "shutdown"}
	return encoder.Encode(req)
}
