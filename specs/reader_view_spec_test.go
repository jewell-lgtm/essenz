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

func TestReaderViewExtractionSpec(t *testing.T) {
	t.Log("SPEC: Reader View Content Extraction")
	t.Log("GIVEN an HTML page with article content mixed with navigation and ads")
	t.Log("WHEN the user runs `sz fetch --reader-view https://example.com`")
	t.Log("THEN the output should contain only the main article content in clean markdown")

	// Create complex HTML with main content and noise
	htmlContent := `<!DOCTYPE html>
<html>
<head>
	<title>Important Article Title</title>
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
			<h1>Important Article Title</h1>
			<p class="byline">By John Doe, Published January 1, 2024</p>
			<p>This is the main article content that should be extracted. It contains important information that the user wants to read.</p>
			<p>This is another paragraph of the main content. It provides valuable insights and should be included in the reader view.</p>
			<h2>Section Header</h2>
			<p>More content under the section header. This should also be preserved in the clean output.</p>
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

	// Run fetch command with reader view flag
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "--reader-view", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch with reader view should succeed")

	outputStr := string(output)

	// Should contain main content
	assert.Contains(t, outputStr, "Important Article Title", "Should extract article title")
	assert.Contains(t, outputStr, "main article content", "Should extract main content")
	assert.Contains(t, outputStr, "Section Header", "Should extract section headers")
	assert.Contains(t, outputStr, "valuable insights", "Should extract all paragraphs")

	// Should NOT contain navigation, ads, or footer content
	assert.NotContains(t, outputStr, "Advertisement Content", "Should filter out ads")
	assert.NotContains(t, outputStr, "Related Links", "Should filter out sidebar")
	assert.NotContains(t, outputStr, "Home", "Should filter out navigation")
	assert.NotContains(t, outputStr, "Copyright 2024", "Should filter out footer")
	assert.NotContains(t, outputStr, "social media", "Should filter out social links")
}

func TestReaderViewMarkdownOutputSpec(t *testing.T) {
	t.Log("SPEC: Reader View Markdown Output Format")
	t.Log("GIVEN an HTML article with various formatting elements")
	t.Log("WHEN the user runs `sz fetch --reader-view https://example.com`")
	t.Log("THEN the output should be properly formatted markdown")

	htmlContent := `<!DOCTYPE html>
<html>
<body>
	<article>
		<h1>Main Title</h1>
		<h2>Subtitle</h2>
		<p>Regular paragraph with <strong>bold text</strong> and <em>italic text</em>.</p>
		<ul>
			<li>First bullet point</li>
			<li>Second bullet point</li>
		</ul>
		<blockquote>This is a quote that should be preserved.</blockquote>
		<p>Paragraph with <a href="https://example.com">a link</a> inside.</p>
	</article>
</body>
</html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "--reader-view", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch with reader view should succeed")

	outputStr := string(output)

	// Check markdown formatting
	assert.Contains(t, outputStr, "# Main Title", "Should format H1 as markdown")
	assert.Contains(t, outputStr, "## Subtitle", "Should format H2 as markdown")
	assert.Contains(t, outputStr, "**bold text**", "Should format bold as markdown")
	assert.Contains(t, outputStr, "*italic text*", "Should format italic as markdown")
	assert.Contains(t, outputStr, "- First bullet", "Should format lists as markdown")
	assert.Contains(t, outputStr, "> This is a quote", "Should format quotes as markdown")
	assert.Contains(t, outputStr, "[a link](https://example.com)", "Should format links as markdown")
}

func TestReaderViewLocalFileSpec(t *testing.T) {
	t.Log("SPEC: Reader View with Local HTML Files")
	t.Log("GIVEN a local HTML file with article content")
	t.Log("WHEN the user runs `sz fetch --reader-view /path/to/file.html`")
	t.Log("THEN the output should extract and format the content")

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "article.html")
	htmlContent := `<!DOCTYPE html>
<html>
<body>
	<div class="header">Site Navigation</div>
	<div class="content">
		<h1>Local Article Title</h1>
		<p>This is local content that should be extracted cleanly.</p>
	</div>
	<div class="footer">Footer content</div>
</body>
</html>`

	err := os.WriteFile(testFile, []byte(htmlContent), 0644)
	require.NoError(t, err, "Should create test file")

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "--reader-view", testFile)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Fetch local file with reader view should succeed")

	outputStr := string(output)
	assert.Contains(t, outputStr, "Local Article Title", "Should extract title")
	assert.Contains(t, outputStr, "local content", "Should extract main content")
	assert.NotContains(t, outputStr, "Site Navigation", "Should filter header")
	assert.NotContains(t, outputStr, "Footer content", "Should filter footer")
}

func TestReaderViewFallbackSpec(t *testing.T) {
	t.Log("SPEC: Reader View Graceful Fallback")
	t.Log("GIVEN content that cannot be processed by reader view")
	t.Log("WHEN the user runs `sz fetch --reader-view` on problematic content")
	t.Log("THEN the command should fall back to raw content with a warning")

	// Create minimal HTML that might be hard to parse
	htmlContent := `<html><body><div>Minimal content</div></body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", "--reader-view", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should succeed even with minimal content")

	outputStr := string(output)
	// Should contain the content, either extracted or as fallback
	assert.Contains(t, outputStr, "Minimal content", "Should show content in some form")
}

func TestReaderViewCompatibilitySpec(t *testing.T) {
	t.Log("SPEC: Reader View Backward Compatibility")
	t.Log("GIVEN the existing fetch command without reader view flag")
	t.Log("WHEN the user runs `sz fetch https://example.com` (no flag)")
	t.Log("THEN the output should remain unchanged (raw HTML)")

	htmlContent := `<html><body><div class="nav">Nav</div><p>Content</p></body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	// Test without reader view flag
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "fetch", server.URL)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Regular fetch should still work")

	outputStr := string(output)
	// Should contain raw HTML including navigation
	assert.Contains(t, outputStr, "Nav", "Should preserve all content without reader view")
	assert.Contains(t, outputStr, "Content", "Should preserve main content")
	// Should contain HTML tags when not using reader view
	containsHTML := strings.Contains(outputStr, "<") && strings.Contains(outputStr, ">")
	assert.True(t, containsHTML, "Should preserve HTML structure without reader view")
}
