// Package browser provides high-level browser operations using Chrome automation.
package browser

import (
	"context"

	"github.com/jewell-lgtm/essenz/internal/daemon"
	"github.com/jewell-lgtm/essenz/internal/pageready"
)

// Client provides browser operations with automatic daemon management.
type Client struct {
	readinessChecker *pageready.ReadinessChecker
}

// NewClient creates a new browser client with global daemon management.
func NewClient() *Client {
	return &Client{
		readinessChecker: nil, // Default: no readiness checking
	}
}

// WithReadinessChecker configures the client to use DOM readiness detection.
func (c *Client) WithReadinessChecker(checker *pageready.ReadinessChecker) *Client {
	c.readinessChecker = checker
	return c
}

// FetchContent fetches content from a URL using Chrome rendering via daemon.
func (c *Client) FetchContent(ctx context.Context, url string) (string, error) {
	client := daemon.NewDaemonClient()

	// If we have a readiness checker, use enhanced fetch
	if c.readinessChecker != nil {
		return client.FetchContentWithReadiness(ctx, url, c.readinessChecker)
	}

	// Otherwise use basic fetch
	return client.FetchContent(ctx, url)
}

// Shutdown is a no-op since we use global daemon management.
// The global daemon will shut down automatically after idle timeout.
func (c *Client) Shutdown() {
	// No-op - global daemon manages its own lifecycle
}
