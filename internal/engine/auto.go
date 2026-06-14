package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// thinTextThreshold is the minimum visible-text length (runes) below which a
// result is considered "thin" and may trigger escalation to a heavier engine.
const thinTextThreshold = 500

// AutoEngine is a composite engine that escalates through its leaf engines until
// one returns usable content: http -> lightpanda -> chrome by default.
type AutoEngine struct {
	engines []Engine
}

// NewAutoEngine returns an auto engine. With no arguments it uses the default
// tiers (http, lightpanda, chrome).
func NewAutoEngine(engines ...Engine) *AutoEngine {
	if len(engines) == 0 {
		engines = []Engine{NewHTTPEngine(), NewLightpandaEngine(), NewChromeEngine()}
	}
	return &AutoEngine{engines: engines}
}

// Name implements Engine.
func (e *AutoEngine) Name() string { return "auto" }

// Fetch tries each engine in order, escalating on error or thin output. The last
// engine's result is always accepted (nothing heavier remains to try).
func (e *AutoEngine) Fetch(ctx context.Context, rawURL string, opts Options) (Result, error) {
	var lastErr error
	var fallback *Result // best successful-but-thin result so far
	for i, eng := range e.engines {
		last := i == len(e.engines)-1
		res, err := eng.Fetch(ctx, rawURL, opts)
		if err != nil {
			// A definitive HTTP status (404/410/5xx) won't be helped by a heavier
			// engine — propagate it rather than rendering an error page as content.
			var se *StatusError
			if errors.As(err, &se) && definitiveStatus(se.Code) {
				return Result{}, fmt.Errorf("auto engine: %w", err)
			}
			lastErr = err
			debugf(opts, "[engine:auto] %s failed: %v", eng.Name(), err)
			continue
		}
		if last || sufficient(res) {
			debugf(opts, "[engine:auto] using %s", res.Engine)
			return res, nil
		}
		// Thin but usable: remember it in case heavier engines are unavailable.
		debugf(opts, "[engine:auto] %s output thin, escalating", eng.Name())
		r := res
		fallback = &r
	}
	// A working result with real content beats failing outright.
	if fallback != nil && len([]rune(visibleText(fallback.HTML))) > 0 {
		debugf(opts, "[engine:auto] escalation exhausted, falling back to %s", fallback.Engine)
		return *fallback, nil
	}
	if lastErr != nil {
		return Result{}, fmt.Errorf("auto engine: all engines failed; last error: %w", lastErr)
	}
	if fallback != nil {
		return *fallback, nil
	}
	return Result{}, fmt.Errorf("auto engine: no engines configured")
}

// definitiveStatus reports whether an HTTP status is final (escalating to a
// heavier engine cannot turn it into content). 401/403/429 are excluded since a
// real browser or cookies may succeed.
func definitiveStatus(code int) bool {
	return code == 404 || code == 410 || code >= 500
}

// sufficient reports whether a result is good enough to accept without
// escalating. Substantial content is always sufficient; thin content escalates
// only when the page looks JavaScript-dependent (so genuinely small static pages
// are not pointlessly handed to a heavier engine).
func sufficient(res Result) bool {
	vis := len([]rune(visibleText(res.HTML)))
	if vis == 0 {
		return false // empty output is never sufficient
	}
	if vis >= thinTextThreshold {
		return true
	}
	probe := res.Raw
	if probe == "" {
		probe = res.HTML
	}
	return !looksJSDependent(probe)
}

var (
	scriptStyleRe = regexp.MustCompile(`(?is)<(script|style)\b[^>]*>.*?</(script|style)>`)
	tagRe         = regexp.MustCompile(`(?s)<[^>]+>`)
	wsRe          = regexp.MustCompile(`\s+`)
	scriptTagRe   = regexp.MustCompile(`(?i)<script[\s>]`)
	jsMarkerRe    = regexp.MustCompile(`(?i)(enable|requires?|turn on)\s+javascript|please enable js`)
)

// visibleText approximates the human-visible text of an HTML fragment by
// dropping script/style blocks and tags, then collapsing whitespace.
func visibleText(html string) string {
	s := scriptStyleRe.ReplaceAllString(html, " ")
	s = tagRe.ReplaceAllString(s, " ")
	s = wsRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

// looksJSDependent reports whether markup suggests content is rendered by
// JavaScript (has <script> tags or an "enable JavaScript" notice).
func looksJSDependent(html string) bool {
	return scriptTagRe.MatchString(html) || jsMarkerRe.MatchString(html)
}

func debugf(opts Options, format string, args ...any) {
	if opts.Debug {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
