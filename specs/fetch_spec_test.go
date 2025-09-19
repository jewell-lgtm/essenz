package specs

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchLocalFileSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command with Local File Support")
	t.Log("GIVEN a local HTML file exists")
	t.Log("WHEN the user runs `sz fetch /path/to/file.html`")
	t.Log("THEN the output should display the file contents")

	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.html")
	testContent := "<html><body>Hello World</body></html>"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err, "Should create test file")

	// Run the fetch command
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", testFile)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch command should succeed")

	outputStr := string(output)
	assert.Contains(t, outputStr, testContent, "Output should contain file contents")
}

func TestFetchHTTPSURLSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command with HTTPS URL Support")
	t.Log("GIVEN a valid HTTPS URL")
	t.Log("WHEN the user runs `sz fetch https://example.com`")
	t.Log("THEN the output should display the HTTPS response content")

	// Create a test HTTP server
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>Test Server Response</body></html>"))
	}))
	defer server.Close()

	// Run the fetch command with the test server URL
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch command should succeed")

	outputStr := string(output)
	assert.Contains(t, outputStr, "Test Server Response", "Output should contain server response")
}

func TestFetchHTTPURLSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command with HTTP URL Support")
	t.Log("GIVEN a valid HTTP URL")
	t.Log("WHEN the user runs `sz fetch http://example.com`")
	t.Log("THEN the output should display the HTTP response content")

	// Create a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>HTTP Test Response</body></html>"))
	}))
	defer server.Close()

	// Run the fetch command with the test server URL
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch command should succeed")

	outputStr := string(output)
	assert.Contains(t, outputStr, "HTTP Test Response", "Output should contain server response")
}

func TestFetchFileNotFoundSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command Error Handling - File Not Found")
	t.Log("GIVEN a file path that does not exist")
	t.Log("WHEN the user runs `sz fetch /nonexistent/file.html`")
	t.Log("THEN the command should exit with error and show helpful message")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "/nonexistent/file.html")
	output, err := cmd.CombinedOutput()
	require.Error(t, err, "Command should fail for nonexistent file")

	outputStr := string(output)
	assert.Contains(t, outputStr, "no such file", "Should show file not found error")
}

func TestFetchInvalidURLSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command Error Handling - Invalid URL")
	t.Log("GIVEN an invalid URL")
	t.Log("WHEN the user runs `sz fetch invalid-url`")
	t.Log("THEN the command should exit with error and show helpful message")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "invalid-url")
	output, err := cmd.CombinedOutput()
	require.Error(t, err, "Command should fail for invalid URL")

	outputStr := string(output)
	// Should contain either URL parsing error or file not found error
	hasURLError := strings.Contains(outputStr, "invalid") ||
		strings.Contains(outputStr, "no such file") ||
		strings.Contains(outputStr, "unsupported protocol")
	assert.True(t, hasURLError, "Should show appropriate error message")
}

func TestFetchHelpSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command Usage Help")
	t.Log("GIVEN the user needs help with the fetch command")
	t.Log("WHEN the user runs `sz fetch --help`")
	t.Log("THEN the output should show usage examples for both files and URLs")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "--help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Help command should succeed")

	outputStr := string(output)
	assert.Contains(t, outputStr, "fetch", "Should mention fetch command")
	assert.Contains(t, outputStr, "Usage", "Should show usage information")
}
