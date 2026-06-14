// Package engine provides pluggable fetch engines for retrieving and rendering
// web pages. Each engine turns a URL into HTML; essenz's renderer then processes
// that HTML uniformly regardless of which engine produced it.
//
// Engines:
//
//	http        net/http + go-readability (no JS, lightest)
//	lightpanda  headless JS via the `lightpanda fetch` subprocess
//	chrome      headless Chrome (heavy, full Web-API fallback)
//	auto        composite that escalates http -> lightpanda -> chrome
//
// Engines may compose other engines (auto delegates to the leaf engines).
package engine

import (
	"context"
	"fmt"
	"time"
)

// Cookie is a cookie to set before fetching a page.
type Cookie struct {
	Name   string
	Value  string
	Domain string
	Path   string
}

// Options configures a fetch.
type Options struct {
	// Cookies are set before navigation.
	Cookies []Cookie
	// Timeout bounds the whole fetch. Zero means the engine's default.
	Timeout time.Duration
	// Debug emits diagnostics to stderr (e.g. auto escalation decisions).
	Debug bool
	// InsecureTLS disables TLS certificate verification. Off by default; the CLI
	// enables it to match essenz's long-standing behaviour (test servers with
	// self-signed certs).
	InsecureTLS bool
}

// Result is the outcome of a successful fetch.
type Result struct {
	// HTML is the fetched/rendered markup, ready for essenz's renderer.
	HTML string
	// Raw is the markup before any extraction (e.g. the http engine's raw page
	// before readability). Used by the auto engine to detect JS-dependence.
	// Optional: when empty the auto engine falls back to HTML.
	Raw string
	// Engine is the name of the leaf engine that produced the result. For the
	// auto engine this reports the tier that actually succeeded.
	Engine string
}

// StatusError is returned by an engine when the server responds with a non-2xx
// HTTP status. The auto engine uses the code to decide whether escalation is
// worthwhile (e.g. 404/410/5xx are definitive; 401/403/429 may be solved by a
// real browser).
type StatusError struct {
	Code   int
	Status string
}

func (e *StatusError) Error() string {
	if e.Status != "" {
		return "unexpected status " + e.Status
	}
	return fmt.Sprintf("unexpected status %d", e.Code)
}

// Engine fetches a URL and returns HTML.
type Engine interface {
	// Name returns the engine's stable identifier (e.g. "http").
	Name() string
	// Fetch retrieves rawURL and returns its HTML.
	Fetch(ctx context.Context, rawURL string, opts Options) (Result, error)
}
