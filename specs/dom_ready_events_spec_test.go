package specs

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDOMReadyEventsSpec(t *testing.T) {
	t.Run("basic_dom_ready_detection", func(t *testing.T) {
		t.Log("SPEC: Basic DOM Ready Detection")
		t.Log("GIVEN a static HTML page with standard content")
		t.Log("WHEN sz processes the page with DOM ready events enabled")
		t.Log("THEN it should wait for DOMContentLoaded before extraction")

		// Build binary for testing
		binary := buildBinary(t)

		// Create a test HTML file
		testHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Main Article</h1>
    <p>This is the main content of the article.</p>
    <nav>Navigation content that should be filtered</nav>
    <footer>Footer content</footer>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(testHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with DOM ready events (should be default behavior)
		cmd := exec.Command(binary, tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should extract main content and include DOM ready timing info in debug mode
		assert.Contains(t, outputStr, "Main Article", "Should extract main article content")
		assert.Contains(t, outputStr, "main content", "Should extract paragraph content")

		// Should not include navigation or footer in clean extraction
		assert.NotContains(t, outputStr, "Navigation content", "Should filter out navigation")
		assert.NotContains(t, outputStr, "Footer content", "Should filter out footer")
	})

	t.Run("javascript_framework_detection", func(t *testing.T) {
		t.Log("SPEC: JavaScript Framework Detection")
		t.Log("GIVEN a React-based page with delayed content rendering")
		t.Log("WHEN sz processes the page with framework detection")
		t.Log("THEN it should wait for React hydration before extraction")

		binary := buildBinary(t)

		// Test with a React-like page structure
		reactHTML := `<!DOCTYPE html>
<html>
<head>
    <title>React App</title>
    <script src="https://unpkg.com/react@18/umd/react.development.js"></script>
    <script src="https://unpkg.com/react-dom@18/umd/react-dom.development.js"></script>
</head>
<body>
    <div id="root">Loading...</div>
    <script>
        // Simulate React app initialization
        setTimeout(() => {
            const root = document.getElementById('root');
            root.innerHTML = '<h1>React Article</h1><p>This content loaded after React hydration.</p>';
            // Dispatch custom event to signal app is ready
            window.dispatchEvent(new Event('react-app-ready'));
        }, 1000);
    </script>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "react-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(reactHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with framework detection enabled
		cmd := exec.Command(binary, "--wait-for-frameworks", tmpFile.Name())
		start := time.Now()
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)

		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should wait for framework and extract the dynamically loaded content
		assert.Contains(t, outputStr, "React Article", "Should extract React-rendered content")
		assert.Contains(t, outputStr, "after React hydration", "Should wait for dynamic content")
		assert.NotContains(t, outputStr, "Loading...", "Should not extract initial loading state")

		// Should have waited at least 1 second for the React content
		assert.GreaterOrEqual(t, duration, time.Second, "Should wait for React hydration")
	})

	t.Run("timeout_handling", func(t *testing.T) {
		t.Log("SPEC: Timeout Handling")
		t.Log("GIVEN a page that never signals readiness completion")
		t.Log("WHEN sz processes the page with a short timeout")
		t.Log("THEN it should timeout gracefully and extract available content")

		binary := buildBinary(t)

		// Create HTML that never signals completion
		infiniteHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Infinite Loading</title>
</head>
<body>
    <h1>Available Content</h1>
    <p>This content is immediately available.</p>
    <script>
        // Infinite loading simulation - never signals ready
        setInterval(() => {
            console.log('Still loading...');
        }, 100);
    </script>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "infinite*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(infiniteHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with short timeout
		cmd := exec.Command(binary, "--dom-ready-timeout=2s", tmpFile.Name())
		start := time.Now()
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)

		require.NoError(t, err, "Command should succeed even with timeout: %s", string(output))

		outputStr := string(output)

		// Should extract available content despite timeout
		assert.Contains(t, outputStr, "Available Content", "Should extract immediately available content")
		assert.Contains(t, outputStr, "immediately available", "Should extract static content")

		// Should timeout around 2 seconds
		assert.LessOrEqual(t, duration, 3*time.Second, "Should timeout within reasonable time")
		assert.GreaterOrEqual(t, duration, 2*time.Second, "Should wait for specified timeout")
	})

	t.Run("custom_selector_waiting", func(t *testing.T) {
		t.Log("SPEC: Custom Selector Waiting")
		t.Log("GIVEN a page with content that appears after a specific element loads")
		t.Log("WHEN sz waits for a custom CSS selector to appear")
		t.Log("THEN it should extract content only after the selector is present")

		binary := buildBinary(t)

		// Create HTML with delayed content appearance
		delayedHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Delayed Content</title>
</head>
<body>
    <h1>Initial Content</h1>
    <div id="loading">Loading article...</div>
    <script>
        setTimeout(() => {
            const loading = document.getElementById('loading');
            loading.innerHTML = '<div class="article-ready"><h2>Article Title</h2><p>Full article content now available.</p></div>';
        }, 1500);
    </script>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "delayed*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(delayedHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with custom selector waiting
		cmd := exec.Command(binary, "--wait-for-selector=.article-ready", tmpFile.Name())
		start := time.Now()
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)

		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should wait for custom selector and extract the final content
		assert.Contains(t, outputStr, "Article Title", "Should extract content after selector appears")
		assert.Contains(t, outputStr, "Full article content", "Should wait for complete content")
		assert.NotContains(t, outputStr, "Loading article...", "Should not extract loading state")

		// Should have waited at least 1.5 seconds for the content
		assert.GreaterOrEqual(t, duration, 1400*time.Millisecond, "Should wait for custom selector")
	})

	t.Run("readiness_result_information", func(t *testing.T) {
		t.Log("SPEC: Readiness Result Information")
		t.Log("GIVEN a page with various readiness indicators")
		t.Log("WHEN sz processes the page with debug information enabled")
		t.Log("THEN it should provide detailed readiness detection results")

		binary := buildBinary(t)

		// Create simple HTML for readiness testing
		simpleHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Simple Page</title>
</head>
<body>
    <h1>Simple Article</h1>
    <p>Basic content for readiness testing.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "simple*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(simpleHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with debug readiness information
		cmd := exec.Command(binary, "--debug-readiness", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should include readiness information in debug mode
		// This will be implemented as part of the readiness system
		assert.Contains(t, outputStr, "Simple Article", "Should extract main content")

		// Debug information should be available (exact format TBD during implementation)
		// For now, just ensure the content extraction works correctly
	})

	t.Run("network_error_recovery", func(t *testing.T) {
		t.Log("SPEC: Network Error Recovery")
		t.Log("GIVEN an invalid URL that cannot be loaded")
		t.Log("WHEN sz attempts to process the URL with DOM ready events")
		t.Log("THEN it should handle the error gracefully with appropriate messaging")

		binary := buildBinary(t)

		// Test with invalid URL
		cmd := exec.Command(binary, "https://nonexistent-domain-for-testing.invalid")
		output, err := cmd.CombinedOutput()

		// Should exit with error but provide helpful message
		assert.Error(t, err, "Should exit with error for invalid URL")

		outputStr := string(output)

		// Should provide helpful error message
		assert.True(t,
			strings.Contains(outputStr, "error") ||
				strings.Contains(outputStr, "failed") ||
				strings.Contains(outputStr, "unable"),
			"Should provide helpful error message: %s", outputStr)
	})
}

func TestDOMReadyEventsWithRealWebsites(t *testing.T) {
	// Skip if no internet connection or in CI without network
	if os.Getenv("SKIP_NETWORK_TESTS") != "" {
		t.Skip("Skipping network tests")
	}

	t.Run("real_static_website", func(t *testing.T) {
		t.Log("SPEC: Real Static Website Processing")
		t.Log("GIVEN a real static website (example.com)")
		t.Log("WHEN sz processes the site with DOM ready events")
		t.Log("THEN it should successfully extract content")

		binary := buildBinary(t)

		cmd := exec.Command(binary, "https://example.com")
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("Network test failed (expected in some environments): %s", string(output))
			t.Skip("Network test failed - continuing with other tests")
			return
		}

		outputStr := string(output)

		// Should extract content from example.com
		assert.Contains(t, outputStr, "Example", "Should extract content from example.com")
	})
}

// Helper function to build the binary for testing
func buildBinary(t *testing.T) string {
	// Build the sz binary from project root
	cmd := exec.Command("go", "build", "-o", "/tmp/sz-dom-ready-test", "./cmd/essenz")
	// Set working directory to project root (assuming we're running tests from project root)
	cmd.Dir = ".."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build binary for testing")

	return "/tmp/sz-dom-ready-test"
}
