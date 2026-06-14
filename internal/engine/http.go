package engine

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

// browserUA is sent by the http engine. Many sites serve degraded content or
// reject the default Go user-agent, which would force needless escalation to a
// JS engine, so we present as a mainstream desktop browser.
const browserUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

const defaultTimeout = 30 * time.Second

// HTTPEngine fetches a page with net/http and returns its raw HTML. It runs no
// JavaScript: it is the lightest tier. Content extraction (readability) is
// applied downstream by essenz's extractor, uniformly across all engines.
type HTTPEngine struct {
	// Client is the HTTP client to use. If nil, a per-fetch client honouring
	// Options.Timeout is created.
	Client *http.Client
}

// NewHTTPEngine returns an http engine using a default client.
func NewHTTPEngine() *HTTPEngine { return &HTTPEngine{} }

// Name implements Engine.
func (e *HTTPEngine) Name() string { return "http" }

// Fetch retrieves rawURL over HTTP and returns readability-extracted HTML,
// falling back to the raw HTML when the page is not article-shaped.
func (e *HTTPEngine) Fetch(ctx context.Context, rawURL string, opts Options) (Result, error) {
	client := e.Client
	if client == nil {
		timeout := opts.Timeout
		if timeout <= 0 {
			timeout = defaultTimeout
		}
		client = &http.Client{Timeout: timeout}
		if opts.InsecureTLS {
			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // opt-in, matches essenz's existing behaviour
			}
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return Result{}, fmt.Errorf("http engine: %w", err)
	}
	req.Header.Set("User-Agent", browserUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	for i := range opts.Cookies {
		c := opts.Cookies[i]
		req.AddCookie(&http.Cookie{Name: c.Name, Value: c.Value, Path: c.Path, Domain: c.Domain})
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("http engine: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Result{}, &StatusError{Code: resp.StatusCode, Status: resp.Status}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, fmt.Errorf("http engine: reading body: %w", err)
	}
	raw := string(body)

	return Result{HTML: raw, Raw: raw, Engine: e.Name()}, nil
}
