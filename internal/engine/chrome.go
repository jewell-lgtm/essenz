package engine

import (
	"context"
	"fmt"

	"github.com/jewell-lgtm/essenz/internal/browser"
	"github.com/jewell-lgtm/essenz/internal/daemon"
	"github.com/jewell-lgtm/essenz/internal/pageready"
)

// ChromeEngine renders pages with headless Chrome via the existing browser/daemon
// stack. It is the heavy, full-Web-API fallback tier.
type ChromeEngine struct {
	// Readiness, when set, configures DOM readiness detection (LCP, frameworks,
	// selectors). nil uses the daemon's default behaviour.
	Readiness *pageready.ReadinessChecker
}

// NewChromeEngine returns a chrome engine with default readiness.
func NewChromeEngine() *ChromeEngine { return &ChromeEngine{} }

// Name implements Engine.
func (e *ChromeEngine) Name() string { return "chrome" }

// Fetch renders rawURL with headless Chrome and returns its HTML.
func (e *ChromeEngine) Fetch(ctx context.Context, rawURL string, opts Options) (Result, error) {
	if opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	client := browser.NewClient()
	defer client.Shutdown()

	if e.Readiness != nil {
		client = client.WithReadinessChecker(e.Readiness)
	}
	if len(opts.Cookies) > 0 {
		client = client.WithCookies(toDaemonCookies(opts.Cookies))
	}

	html, err := client.FetchContent(ctx, rawURL)
	if err != nil {
		return Result{}, fmt.Errorf("chrome engine: %w", err)
	}
	return Result{HTML: html, Engine: e.Name()}, nil
}

// toDaemonCookies converts engine cookies to the daemon's identical type.
func toDaemonCookies(cookies []Cookie) []daemon.Cookie {
	out := make([]daemon.Cookie, len(cookies))
	for i, c := range cookies {
		out[i] = daemon.Cookie(c)
	}
	return out
}
