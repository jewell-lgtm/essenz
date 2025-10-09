package specs

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawFlagSpec(t *testing.T) {
	t.Log("SPEC: Raw HTML Output Flag")
	t.Log("GIVEN the sz command line tool with a local HTML file")
	t.Log("WHEN the user runs `sz --raw <file>`")
	t.Log("THEN the output should be the raw HTML without any processing")

	// Run sz with --raw flag on a test HTML file
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "--raw", "../testdata/raw_flag_test.html")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)

	// Verify the output contains raw HTML elements that would normally be filtered out
	assert.Contains(t, outputStr, "<!DOCTYPE html>", "Should preserve DOCTYPE declaration")
	assert.Contains(t, outputStr, "<html", "Should preserve opening html tag")
	assert.Contains(t, outputStr, "<head>", "Should preserve head tag")
	assert.Contains(t, outputStr, "<meta charset=\"UTF-8\">", "Should preserve meta tags")
	assert.Contains(t, outputStr, "<title>Raw Flag Test Page</title>", "Should preserve title tag")
	assert.Contains(t, outputStr, "</head>", "Should preserve closing head tag")
	assert.Contains(t, outputStr, "<body>", "Should preserve body tag")
	assert.Contains(t, outputStr, "<nav>", "Should preserve navigation elements")
	assert.Contains(t, outputStr, "<footer>", "Should preserve footer elements")
	assert.Contains(t, outputStr, "</html>", "Should preserve closing html tag")

	// Verify the output is complete HTML, not processed/extracted content
	assert.True(t, strings.HasPrefix(strings.TrimSpace(outputStr), "<!DOCTYPE html>"),
		"Output should start with DOCTYPE")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(outputStr), "</html>"),
		"Output should end with closing html tag")
}

func TestRawFlagPreservesAllContentSpec(t *testing.T) {
	t.Log("SPEC: Raw Flag Preserves All Content")
	t.Log("GIVEN the sz command line tool with a local HTML file")
	t.Log("WHEN the user runs `sz --raw <file>`")
	t.Log("THEN all HTML content should be preserved without extraction or filtering")

	// Run sz with --raw flag
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "--raw", "../testdata/raw_flag_test.html")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)

	// Check that navigation and footer elements are present (these would be filtered in normal mode)
	assert.Contains(t, outputStr, "<li><a href=\"/\">Home</a></li>",
		"Should preserve navigation links")
	assert.Contains(t, outputStr, "<li><a href=\"/about\">About</a></li>",
		"Should preserve all navigation items")
	assert.Contains(t, outputStr, "&copy; 2024 Test Site",
		"Should preserve footer content")
}

func TestDefaultModeWithoutRawFlagSpec(t *testing.T) {
	t.Log("SPEC: Default Mode Without Raw Flag")
	t.Log("GIVEN the sz command line tool with a local HTML file")
	t.Log("WHEN the user runs `sz <file>` without --raw flag")
	t.Log("THEN the output should be processed and NOT include raw HTML structure")

	// Run sz WITHOUT --raw flag (default behavior with semantic extraction)
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "../testdata/raw_flag_test.html")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)

	// In default mode with semantic extraction, we should NOT see the full HTML structure
	// The semantic extractor should extract only the main content
	assert.NotContains(t, outputStr, "<!DOCTYPE html>",
		"Default mode should NOT preserve DOCTYPE")
	assert.NotContains(t, outputStr, "<head>",
		"Default mode should NOT preserve head tag")
	assert.NotContains(t, outputStr, "<meta charset",
		"Default mode should NOT preserve meta tags")

	// The semantic extractor should still preserve the main content
	assert.Contains(t, outputStr, "Test Article",
		"Default mode should extract main heading")
	assert.Contains(t, outputStr, "This is a test article",
		"Default mode should extract main content")
}

func TestRawFlagWithRealWebPageSpec(t *testing.T) {
	t.Log("SPEC: Raw Flag With HTML File Containing Multiple Sections")
	t.Log("GIVEN a more complex HTML file with multiple content regions")
	t.Log("WHEN the user runs `sz --raw <file>`")
	t.Log("THEN all sections should be preserved in raw HTML format")

	// Use an existing test file with more complex structure
	cmd := exec.Command("go", "run", "../cmd/essenz/main.go", "--raw", "../testdata/blog_post.html")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command should execute successfully")

	outputStr := string(output)

	// Verify complete HTML structure is preserved
	assert.Contains(t, outputStr, "<!DOCTYPE html>", "Should preserve DOCTYPE for complex page")
	assert.Contains(t, outputStr, "<html", "Should preserve html tag for complex page")
	assert.Contains(t, outputStr, "</html>", "Should preserve closing html tag for complex page")

	// Check that the HTML structure is complete (starts and ends properly)
	trimmed := strings.TrimSpace(outputStr)
	assert.True(t, strings.HasPrefix(trimmed, "<!DOCTYPE html>") || strings.HasPrefix(trimmed, "<html"),
		"Complex page should start with DOCTYPE or html tag")
	assert.True(t, strings.HasSuffix(trimmed, "</html>"),
		"Complex page should end with closing html tag")
}
