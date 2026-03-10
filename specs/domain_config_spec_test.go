package specs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainConfigAutoAppliesCookies(t *testing.T) {
	t.Log("SPEC: Per-domain config auto-applies cookies for matching domain")
	t.Log("GIVEN a domain config file with cookies")
	t.Log("WHEN the user runs `sz <url>` for that domain")
	t.Log("THEN cookies from config are automatically sent")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("_gitlab_session")
		if err != nil || cookie.Value != "abc123" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>Authenticated via config</body></html>"))
	}))
	defer server.Close()

	// Create temp config dir with domain config
	configDir := t.TempDir()
	domainsDir := filepath.Join(configDir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))

	// Server URL host is 127.0.0.1:<port>, write config for that
	configYAML := `cookies:
  - name: _gitlab_session
    value: "abc123"
`
	// httptest server uses 127.0.0.1
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "127.0.0.1.yaml"), []byte(configYAML), 0o644))

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go",
		"--raw",
		server.URL)
	cmd.Env = append(os.Environ(), "SZ_CONFIG_DIR="+configDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Should succeed with config cookies: %s", string(output))
	assert.Contains(t, string(output), "Authenticated via config")
}

func TestDomainConfigCLICookieOverrides(t *testing.T) {
	t.Log("SPEC: CLI --cookie overrides config cookie by name")
	t.Log("GIVEN a domain config with cookie session=old")
	t.Log("WHEN the user runs `sz --cookie session=new <url>`")
	t.Log("THEN the CLI cookie value is used")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("No session"))
			return
		}
		_, _ = w.Write([]byte("<html><body>session=" + cookie.Value + "</body></html>"))
	}))
	defer server.Close()

	configDir := t.TempDir()
	domainsDir := filepath.Join(configDir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))

	configYAML := `cookies:
  - name: session
    value: "old_value"
`
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "127.0.0.1.yaml"), []byte(configYAML), 0o644))

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go",
		"--cookie", "session=new_value",
		"--raw",
		server.URL)
	cmd.Env = append(os.Environ(), "SZ_CONFIG_DIR="+configDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Should succeed: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "session=new_value", "CLI cookie should override config")
	assert.NotContains(t, outputStr, "old_value", "Config cookie should be overridden")
}

func TestDomainConfigUnmatchedDomain(t *testing.T) {
	t.Log("SPEC: No config cookies for unmatched domain")
	t.Log("GIVEN a domain config for example.com")
	t.Log("WHEN the user fetches a different domain")
	t.Log("THEN no config cookies are sent")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("secret"); err == nil {
			_, _ = w.Write([]byte("<html><body>LEAKED</body></html>"))
			return
		}
		_, _ = w.Write([]byte("<html><body>No cookies sent</body></html>"))
	}))
	defer server.Close()

	configDir := t.TempDir()
	domainsDir := filepath.Join(configDir, "essenz", "domains")
	require.NoError(t, os.MkdirAll(domainsDir, 0o755))

	// Config for a different domain
	configYAML := `cookies:
  - name: secret
    value: "should_not_leak"
`
	require.NoError(t, os.WriteFile(filepath.Join(domainsDir, "example.com.yaml"), []byte(configYAML), 0o644))

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go",
		"--raw",
		server.URL)
	cmd.Env = append(os.Environ(), "SZ_CONFIG_DIR="+configDir)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Should succeed: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "No cookies sent")
	assert.NotContains(t, outputStr, "LEAKED")
}
