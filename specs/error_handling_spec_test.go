package specs

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandling(t *testing.T) {
	szBinary := buildSzBinary(t)

	t.Run("invalid URL format", func(t *testing.T) {
		// Test with clearly invalid URL
		cmd := exec.Command(szBinary, "not-a-valid-url")
		output, err := cmd.CombinedOutput()

		// Should exit with non-zero status
		require.Error(t, err, "sz should exit with error for invalid URL")

		// Should contain error message about invalid URL
		outputStr := strings.ToLower(string(output))
		assert.True(t,
			strings.Contains(outputStr, "error") || strings.Contains(outputStr, "invalid"),
			"Output should mention error or invalid: %s", string(output))
	})

	t.Run("missing protocol", func(t *testing.T) {
		// Test with URL missing protocol
		cmd := exec.Command(szBinary, "example.com")
		output, err := cmd.CombinedOutput()

		// Should exit with non-zero status
		require.Error(t, err, "sz should exit with error for URL without protocol")

		outputStr := strings.ToLower(string(output))
		assert.True(t,
			strings.Contains(outputStr, "error") || strings.Contains(outputStr, "invalid"),
			"Output should mention error or invalid: %s", string(output))
	})

	t.Run("404 not found", func(t *testing.T) {
		// Test with known 404 URL
		cmd := exec.Command(szBinary, "https://github.com/this-page-does-not-exist-404")
		output, err := cmd.CombinedOutput()

		// Should exit with non-zero status
		require.Error(t, err, "sz should exit with error for 404 page")

		// Should contain error message about 404 or not found
		outputStr := strings.ToLower(string(output))
		assert.True(t,
			strings.Contains(outputStr, "404") || strings.Contains(outputStr, "not found"),
			"Output should mention 404 or not found: %s", string(output))
	})

	t.Run("timeout handling", func(t *testing.T) {
		// Test with timeout flag and slow endpoint
		cmd := exec.Command(szBinary, "--timeout", "2s", "https://httpstat.us/200?sleep=5000")

		start := time.Now()
		output, err := cmd.CombinedOutput()
		elapsed := time.Since(start)

		// Should exit with error
		require.Error(t, err, "sz should exit with error on timeout")

		// Should timeout within reasonable time (not wait full 5 seconds)
		assert.Less(t, elapsed, 4*time.Second, "Should timeout before request completes")

		// Should contain timeout message
		outputStr := strings.ToLower(string(output))
		assert.True(t,
			strings.Contains(outputStr, "timeout") || strings.Contains(outputStr, "deadline"),
			"Output should mention timeout or deadline: %s", string(output))
	})

	t.Run("network error - invalid host", func(t *testing.T) {
		// Test with invalid hostname that should cause DNS error
		cmd := exec.Command(szBinary, "https://this-domain-definitely-does-not-exist-12345.com")
		output, err := cmd.CombinedOutput()

		// Should exit with error
		require.Error(t, err, "sz should exit with error for invalid hostname")

		// Should contain error message
		outputStr := strings.ToLower(string(output))
		assert.True(t,
			strings.Contains(outputStr, "error") ||
				strings.Contains(outputStr, "failed") ||
				strings.Contains(outputStr, "no such host"),
			"Output should mention error: %s", string(output))
	})

	t.Run("successful fetch still works", func(t *testing.T) {
		// Verify that normal operation still works
		cmd := exec.Command(szBinary, "https://example.com")
		output, err := cmd.CombinedOutput()

		// Should succeed
		require.NoError(t, err, "sz should succeed for valid URL")

		// Should contain some content
		assert.NotEmpty(t, output, "Should return content for valid URL")
	})

	t.Run("file not found", func(t *testing.T) {
		// Test with non-existent file
		cmd := exec.Command(szBinary, "/path/to/nonexistent/file.html")
		output, err := cmd.CombinedOutput()

		// Should exit with error
		require.Error(t, err, "sz should exit with error for non-existent file")

		// Should contain error message
		outputStr := strings.ToLower(string(output))
		assert.True(t,
			strings.Contains(outputStr, "error") ||
				strings.Contains(outputStr, "no such file"),
			"Output should mention error or file not found: %s", string(output))
	})
}
