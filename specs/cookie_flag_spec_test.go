package specs

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCookieFlagSpec(t *testing.T) {
	t.Log("SPEC: Cookie Flag Support")
	t.Log("GIVEN a server that requires authentication via cookie")
	t.Log("WHEN the user runs `sz --cookie 'session=abc123' <url>`")
	t.Log("THEN the request should include the cookie and return authenticated content")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value != "abc123" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Unauthorized"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>Authenticated Content</body></html>"))
	}))
	defer server.Close()

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go",
		"--cookie", "session=abc123",
		"--raw",
		server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should succeed with valid cookie: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Authenticated Content", "Should receive authenticated content")
	assert.NotContains(t, outputStr, "Unauthorized", "Should not be unauthorized")
}

func TestMultipleCookiesSpec(t *testing.T) {
	t.Log("SPEC: Multiple Cookie Flag Support")
	t.Log("GIVEN a server that requires multiple cookies")
	t.Log("WHEN the user runs `sz --cookie 'a=1' --cookie 'b=2' <url>`")
	t.Log("THEN all cookies should be sent")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a, errA := r.Cookie("a")
		b, errB := r.Cookie("b")
		if errA != nil || errB != nil || a.Value != "1" || b.Value != "2" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("Missing cookies"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>Both Cookies Received</body></html>"))
	}))
	defer server.Close()

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go",
		"--cookie", "a=1",
		"--cookie", "b=2",
		"--raw",
		server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should succeed with multiple cookies: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Both Cookies Received", "Should receive content when both cookies sent")
}

func TestInvalidCookieFormatSpec(t *testing.T) {
	t.Log("SPEC: Invalid Cookie Format Error")
	t.Log("GIVEN an invalid cookie format")
	t.Log("WHEN the user runs `sz --cookie 'invalid' <url>`")
	t.Log("THEN an error should be shown")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go",
		"--cookie", "invalid",
		"https://example.com")
	output, err := cmd.CombinedOutput()
	require.Error(t, err, "Command should fail with invalid cookie format")

	outputStr := string(output)
	assert.True(t,
		strings.Contains(outputStr, "invalid cookie") ||
			strings.Contains(outputStr, "expected name=value"),
		"Should show cookie format error: %s", outputStr)
}

func TestCookieWithFetchCommandSpec(t *testing.T) {
	t.Log("SPEC: Cookie Flag with Fetch Command")
	t.Log("GIVEN the fetch subcommand")
	t.Log("WHEN the user runs `sz fetch --cookie 'key=val' <url>`")
	t.Log("THEN the cookie should be included in the request")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("key")
		if err != nil || cookie.Value != "val" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("No Cookie"))
			return
		}
		_, _ = w.Write([]byte("<html><body>Fetch Cookie Works</body></html>"))
	}))
	defer server.Close()

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch",
		"--cookie", "key=val",
		"--raw",
		server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch command should succeed with cookie: %s", string(output))

	assert.Contains(t, string(output), "Fetch Cookie Works")
}
