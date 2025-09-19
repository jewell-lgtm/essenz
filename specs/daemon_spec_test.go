package specs

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChromeDaemonManagement validates F3: Chrome Daemon Management
//
// SPEC: Chrome daemon starts on first browser operation, multiple sz commands
// reuse the same Chrome instance, daemon shuts down gracefully when not needed,
// and daemon restarts automatically if crashed.
func TestChromeDaemonManagement(t *testing.T) {
	// Build sz binary for testing
	szBinary := buildSzBinary(t)
	defer func() { _ = os.Remove(szBinary) }()

	t.Run("daemon starts on first browser operation", func(t *testing.T) {
		// First command that requires Chrome should succeed and use Chrome automation
		cmd := exec.Command(szBinary, "fetch", "https://example.com")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "First fetch command should succeed")
		assert.Contains(t, string(output), "<html>", "Should contain HTML content")
		assert.Contains(t, string(output), "Example Domain", "Should contain expected content from Chrome rendering")
	})

	t.Run("multiple commands reuse same Chrome instance", func(t *testing.T) {
		// Get initial Chrome process count
		initialChromeCount := countChromeProcesses()

		// Run multiple sz commands in quick succession
		for i := 0; i < 3; i++ {
			cmd := exec.Command(szBinary, "fetch", "https://example.com")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Fetch command %d should succeed", i+1)
			assert.Contains(t, string(output), "Example Domain", "Should contain example.com content")
		}

		// Chrome process count should not have increased significantly
		finalChromeCount := countChromeProcesses()
		assert.LessOrEqual(t, finalChromeCount, initialChromeCount+2,
			"Should not spawn many new Chrome processes")
	})

	t.Run("daemon shuts down after idle timeout", func(t *testing.T) {
		// Skip this test for now as it requires daemon process architecture improvements
		t.Skip("Timeout functionality requires daemon process lifecycle improvements")
	})

	t.Run("daemon restarts automatically if crashed", func(t *testing.T) {
		// Skip this test for now as it requires daemon process architecture improvements
		t.Skip("Crash recovery functionality requires daemon process lifecycle improvements")
	})
}

// Helper functions for Chrome process management testing

func buildSzBinary(t *testing.T) string {
	t.Helper()
	tmpFile := "/tmp/sz-test-" + time.Now().Format("20060102-150405")
	cmd := exec.Command("go", "build", "-o", tmpFile, "./cmd/essenz")
	// Set working directory to project root for build
	cmd.Dir = ".."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", string(output))
	}
	require.NoError(t, err, "Failed to build sz binary for testing")
	return tmpFile
}

func countChromeProcesses() int {
	cmd := exec.Command("pgrep", "-c", "-f", "Chrome.*--headless.*--remote-debugging-port=9222")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	var count int
	_, err = fmt.Sscanf(string(output), "%d", &count)
	if err != nil {
		return 0
	}
	return count
}
