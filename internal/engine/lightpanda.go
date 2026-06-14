package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// LightpandaEngine renders pages by shelling out to the `lightpanda fetch`
// subprocess. Lightpanda executes JavaScript via V8 with a fraction of Chrome's
// footprint, but has no graphical renderer (so no screenshots / LCP).
type LightpandaEngine struct {
	// BinaryPath, when set, is used directly and skips resolution/download.
	BinaryPath string
	// Resolve locates the lightpanda binary. Defaults to resolveLightpanda
	// (ESSENZ_LIGHTPANDA_PATH -> PATH -> pinned auto-download).
	Resolve func(ctx context.Context, debug bool) (string, error)
}

// NewLightpandaEngine returns a lightpanda engine using default binary resolution.
func NewLightpandaEngine() *LightpandaEngine { return &LightpandaEngine{} }

// Name implements Engine.
func (e *LightpandaEngine) Name() string { return "lightpanda" }

// Fetch renders rawURL with lightpanda and returns its serialized HTML.
func (e *LightpandaEngine) Fetch(ctx context.Context, rawURL string, opts Options) (Result, error) {
	bin := e.BinaryPath
	if bin == "" {
		resolve := e.Resolve
		if resolve == nil {
			resolve = resolveLightpanda
		}
		var err error
		if bin, err = resolve(ctx, opts.Debug); err != nil {
			return Result{}, fmt.Errorf("lightpanda engine: %w", err)
		}
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	args := []string{
		"fetch", rawURL,
		"--dump", "html",
		"--wait-until", "networkidle",
		// Hard deadline: abort JS execution so endless scripts can't hang us.
		"--terminate-ms", strconv.FormatInt(timeout.Milliseconds(), 10),
	}

	if opts.InsecureTLS {
		args = append(args, "--insecure-disable-tls-host-verification")
	}

	if len(opts.Cookies) > 0 {
		path, cleanup, err := writeCookieFile(opts.Cookies)
		if err != nil {
			return Result{}, fmt.Errorf("lightpanda engine: %w", err)
		}
		defer cleanup()
		args = append(args, "--cookie", path)
	}

	// Give the subprocess a little slack beyond its own hard deadline.
	runCtx, cancel := context.WithTimeout(ctx, timeout+5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, bin, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("lightpanda engine: %w: %s", err, lastLines(stderr.String(), 3))
	}

	html := stdout.String()
	if strings.TrimSpace(html) == "" {
		return Result{}, fmt.Errorf("lightpanda engine: empty output: %s", lastLines(stderr.String(), 3))
	}
	if opts.Debug {
		fmt.Fprintf(os.Stderr, "[engine:lightpanda] %s -> %d bytes\n", rawURL, len(html))
	}
	return Result{HTML: html, Engine: e.Name()}, nil
}

// cdpCookie is lightpanda's --cookie JSON schema (CDP-style).
type cdpCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain,omitempty"`
	Path   string `json:"path,omitempty"`
}

// writeCookieFile serializes cookies to a temp JSON file for lightpanda --cookie.
func writeCookieFile(cookies []Cookie) (path string, cleanup func(), err error) {
	arr := make([]cdpCookie, len(cookies))
	for i, c := range cookies {
		arr[i] = cdpCookie(c)
	}
	data, err := json.Marshal(arr)
	if err != nil {
		return "", func() {}, err
	}
	f, err := os.CreateTemp("", "essenz-cookies-*.json")
	if err != nil {
		return "", func() {}, err
	}
	name := f.Name()
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(name)
		return "", func() {}, err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(name)
		return "", func() {}, err
	}
	return name, func() { _ = os.Remove(name) }, nil
}

// lastLines returns the last n non-empty lines of s joined with "; " for compact
// error messages.
func lastLines(s string, n int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	lines := strings.Split(s, "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "; ")
}
