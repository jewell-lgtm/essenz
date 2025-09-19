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

func TestDefaultReaderViewWithURLSpec(t *testing.T) {
	t.Log("SPEC: Default Reader View with URL")
	t.Log("GIVEN a URL with article content mixed with navigation and ads")
	t.Log("WHEN the user runs `sz https://example.com` (no subcommand)")
	t.Log("THEN the output should contain only the main article content as clean markdown")

	// Create complex HTML with main content and noise
	htmlContent := `<!DOCTYPE html>
<html>
<head>
	<title>Essential Article Title</title>
</head>
<body>
	<header>
		<nav>
			<a href="/home">Home</a>
			<a href="/about">About</a>
		</nav>
	</header>

	<aside class="sidebar">
		<div class="ad">Advertisement Content</div>
		<div class="related">Related Links</div>
	</aside>

	<main>
		<article>
			<h1>Essential Article Title</h1>
			<p class="byline">By Jane Doe, Published January 15, 2024</p>
			<p>This is the essential content that should be extracted by default. The tool should automatically process this into clean markdown.</p>
			<p>Another paragraph with <strong>important information</strong> that users want to read without distractions.</p>
			<h2>Key Points</h2>
			<ul>
				<li>First important point</li>
				<li>Second critical insight</li>
			</ul>
		</article>
	</main>

	<footer>
		<p>Copyright 2024. All rights reserved.</p>
		<div class="social-links">Share on social media</div>
	</footer>
</body>
</html>`

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	// Run sz command with just URL (no subcommand) - should default to reader view
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Default sz command with URL should succeed")

	outputStr := string(output)

	// Should contain main content as markdown
	assert.Contains(t, outputStr, "# Essential Article Title", "Should extract article title as H1")
	assert.Contains(t, outputStr, "essential content", "Should extract main content")
	assert.Contains(t, outputStr, "## Key Points", "Should extract section headers as H2")
	assert.Contains(t, outputStr, "**important information**", "Should format bold text as markdown")
	assert.Contains(t, outputStr, "- First important point", "Should format lists as markdown")

	// Should NOT contain navigation, ads, or footer content
	assert.NotContains(t, outputStr, "Advertisement Content", "Should filter out ads")
	assert.NotContains(t, outputStr, "Related Links", "Should filter out sidebar")
	assert.NotContains(t, outputStr, "Home", "Should filter out navigation")
	assert.NotContains(t, outputStr, "Copyright 2024", "Should filter out footer")
	assert.NotContains(t, outputStr, "social media", "Should filter out social links")
}

func TestDefaultReaderViewWithLocalFileSpec(t *testing.T) {
	t.Log("SPEC: Default Reader View with Local File")
	t.Log("GIVEN a local HTML file with article content")
	t.Log("WHEN the user runs `sz /path/to/file.html` (no subcommand)")
	t.Log("THEN the output should extract and format the content as markdown")

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "article.html")
	htmlContent := `<!DOCTYPE html>
<html>
<body>
	<div class="header">Site Navigation</div>
	<div class="content">
		<h1>Local Article Title</h1>
		<p>This local content should be <em>automatically processed</em> with reader view.</p>
		<blockquote>This is an important quote that should be preserved.</blockquote>
	</div>
	<div class="footer">Footer content</div>
</body>
</html>`

	err := os.WriteFile(testFile, []byte(htmlContent), 0644)
	require.NoError(t, err, "Should create test file")

	// Run sz command with file path (no subcommand) - should default to reader view
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", testFile)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Default sz command with file should succeed")

	outputStr := string(output)
	assert.Contains(t, outputStr, "# Local Article Title", "Should extract title as markdown H1")
	assert.Contains(t, outputStr, "*automatically processed*", "Should format italic as markdown")
	assert.Contains(t, outputStr, "> This is an important quote", "Should format blockquotes as markdown")
	assert.NotContains(t, outputStr, "Site Navigation", "Should filter header")
	assert.NotContains(t, outputStr, "Footer content", "Should filter footer")
}

func TestRawFlagBypassesReaderViewSpec(t *testing.T) {
	t.Log("SPEC: Raw Flag Bypasses Reader View")
	t.Log("GIVEN a URL with HTML content")
	t.Log("WHEN the user runs `sz --raw https://example.com`")
	t.Log("THEN the output should show raw HTML without reader view processing")

	htmlContent := `<html><body><div class="nav">Navigation</div><p>Content</p></body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	// Run sz command with --raw flag
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "--raw", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Raw flag command should succeed")

	outputStr := string(output)
	// Should contain raw HTML including navigation
	assert.Contains(t, outputStr, "Navigation", "Should preserve all content with --raw flag")
	assert.Contains(t, outputStr, "Content", "Should preserve main content")
	// Should contain HTML tags when using --raw flag
	containsHTML := strings.Contains(outputStr, "<") && strings.Contains(outputStr, ">")
	assert.True(t, containsHTML, "Should preserve HTML structure with --raw flag")
}

func TestNoArgsShowsHelpSpec(t *testing.T) {
	t.Log("SPEC: No Arguments Shows Help")
	t.Log("GIVEN the sz command line tool is available")
	t.Log("WHEN the user runs `sz` without any arguments")
	t.Log("THEN the output should display help information")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)
	assert.Contains(t, outputStr, "sz is a CLI web browser", "Should contain tool description")
	assert.Contains(t, outputStr, "Usage:", "Should show usage information")
	assert.Contains(t, outputStr, "Available Commands:", "Should list available commands")
}

func TestFetchCommandBackwardCompatibilitySpec(t *testing.T) {
	t.Log("SPEC: Fetch Command Backward Compatibility")
	t.Log("GIVEN the existing fetch command")
	t.Log("WHEN the user runs `sz fetch https://example.com` (explicit fetch)")
	t.Log("THEN the fetch command should work as before (raw HTML by default)")

	htmlContent := `<html><body><div class="nav">Nav</div><p>Content</p></body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	// Test explicit fetch command (should maintain current behavior)
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch command should still work")

	outputStr := string(output)
	// Should contain raw HTML including navigation (backward compatibility)
	assert.Contains(t, outputStr, "Nav", "Should preserve all content with fetch command")
	assert.Contains(t, outputStr, "Content", "Should preserve main content")
	containsHTML := strings.Contains(outputStr, "<") && strings.Contains(outputStr, ">")
	assert.True(t, containsHTML, "Should preserve HTML structure with fetch command")
}

func TestFetchWithReaderViewFlagSpec(t *testing.T) {
	t.Log("SPEC: Fetch Command with Reader View Flag")
	t.Log("GIVEN the fetch command with --reader-view flag")
	t.Log("WHEN the user runs `sz fetch --reader-view https://example.com`")
	t.Log("THEN the output should show reader view processed content")

	htmlContent := `<!DOCTYPE html>
<html>
<body>
	<nav>Navigation Menu</nav>
	<article>
		<h1>Test Article</h1>
		<p>Main content here.</p>
	</article>
	<footer>Footer</footer>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	// Test fetch command with reader view flag
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "--reader-view", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch with reader view should work")

	outputStr := string(output)
	assert.Contains(t, outputStr, "# Test Article", "Should format as markdown")
	assert.Contains(t, outputStr, "Main content here", "Should extract main content")
	assert.NotContains(t, outputStr, "Navigation Menu", "Should filter navigation")
	assert.NotContains(t, outputStr, "Footer", "Should filter footer")
}
