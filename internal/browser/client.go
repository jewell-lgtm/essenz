// Package browser provides high-level browser operations using Chrome automation.
package browser

import (
	"context"

	"github.com/essenz/essenz/internal/daemon"
)

// Client provides browser operations with automatic daemon management.
type Client struct{}

// NewClient creates a new browser client with global daemon management.
func NewClient() *Client {
	return &Client{}
}

// FetchContent fetches content from a URL using Chrome rendering via daemon.
func (c *Client) FetchContent(ctx context.Context, url string) (string, error) {
	client := daemon.NewDaemonClient()
	return client.FetchContent(ctx, url)
}

// Shutdown is a no-op since we use global daemon management.
// The global daemon will shut down automatically after idle timeout.
func (c *Client) Shutdown() {
	// No-op - global daemon manages its own lifecycle
}
