package specs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMarkdownTreeRendererSpec tests F5: Markdown Tree Renderer functionality
func TestMarkdownTreeRendererSpec(t *testing.T) {
	// Build binary for testing
	binary := filepath.Join(os.TempDir(), "sz-markdown-test")
	err := exec.Command("go", "build", "-o", binary, "./cmd/essenz").Run()
	require.NoError(t, err, "Failed to build binary")
	defer func() { _ = os.Remove(binary) }()

	t.Run("basic_document_structure", func(t *testing.T) {
		t.Log("SPEC: Basic Document Structure Rendering")
		t.Log("GIVEN an HTML document with headings, paragraphs, and basic formatting")
		t.Log("WHEN sz processes the document with markdown rendering")
		t.Log("THEN it should generate clean, well-formatted markdown")

		// Create HTML with basic structure
		basicHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Basic Document</title>
</head>
<body>
    <h1>Main Title</h1>
    <p>Introduction paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
    <h2>Section Header</h2>
    <p>Another paragraph with some content.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "basic-markdown-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(basicHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test markdown rendering
		cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should generate proper markdown headings
		assert.Contains(t, outputStr, "# Main Title", "Should generate H1 with # style")
		assert.Contains(t, outputStr, "## Section Header", "Should generate H2 with ## style")

		// Should handle inline formatting
		assert.Contains(t, outputStr, "**bold**", "Should convert strong to **bold**")
		assert.Contains(t, outputStr, "*italic*", "Should convert em to *italic*")

		// Should have proper paragraph structure
		assert.Contains(t, outputStr, "Introduction paragraph with **bold** and *italic* text.", "Should preserve paragraph content with formatting")
		assert.Contains(t, outputStr, "Another paragraph with some content.", "Should preserve second paragraph")

		// Should have proper spacing between elements
		lines := strings.Split(outputStr, "\n")
		found := false
		for i, line := range lines {
			if strings.Contains(line, "# Main Title") && i+1 < len(lines) {
				assert.Equal(t, "", lines[i+1], "Should have blank line after H1")
				found = true
				break
			}
		}
		assert.True(t, found, "Should find H1 with proper spacing")
	})

	t.Run("list_formatting", func(t *testing.T) {
		t.Log("SPEC: List Formatting")
		t.Log("GIVEN an HTML document with ordered and unordered lists")
		t.Log("WHEN sz processes the document with markdown rendering")
		t.Log("THEN it should convert lists to proper markdown format")

		// Create HTML with lists
		listHTML := `<!DOCTYPE html>
<html>
<head>
    <title>List Test</title>
</head>
<body>
    <h2>Shopping List</h2>
    <ul>
        <li>First item</li>
        <li>Second item with <a href="http://example.com">link</a></li>
        <li>Third item</li>
    </ul>
    <h2>Steps</h2>
    <ol>
        <li>Step one</li>
        <li>Step two</li>
        <li>Step three</li>
    </ol>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "list-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(listHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test list rendering
		cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should generate unordered list
		assert.Contains(t, outputStr, "- First item", "Should convert ul to - style")
		assert.Contains(t, outputStr, "- Second item with [link](http://example.com)", "Should preserve links in list items")
		assert.Contains(t, outputStr, "- Third item", "Should convert all list items")

		// Should generate ordered list
		assert.Contains(t, outputStr, "1. Step one", "Should convert ol to numbered style")
		assert.Contains(t, outputStr, "2. Step two", "Should properly number list items")
		assert.Contains(t, outputStr, "3. Step three", "Should continue numbering")

		// Should not contain HTML tags
		assert.NotContains(t, outputStr, "<ul>", "Should not contain ul tags")
		assert.NotContains(t, outputStr, "<ol>", "Should not contain ol tags")
		assert.NotContains(t, outputStr, "<li>", "Should not contain li tags")
	})

	t.Run("nested_list_structure", func(t *testing.T) {
		t.Log("SPEC: Nested List Structure")
		t.Log("GIVEN an HTML document with nested lists")
		t.Log("WHEN sz processes the document with markdown rendering")
		t.Log("THEN it should handle nested indentation correctly")

		// Create HTML with nested lists
		nestedHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Nested List Test</title>
</head>
<body>
    <h2>Project Structure</h2>
    <ol>
        <li>Frontend
            <ul>
                <li>React components</li>
                <li>CSS styles</li>
            </ul>
        </li>
        <li>Backend
            <ul>
                <li>API routes</li>
                <li>Database models</li>
            </ul>
        </li>
    </ol>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "nested-list-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(nestedHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test nested list rendering
		cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should handle nested structure
		assert.Contains(t, outputStr, "1. Frontend", "Should start with top-level numbered item")
		assert.Contains(t, outputStr, "   - React components", "Should indent nested list items")
		assert.Contains(t, outputStr, "   - CSS styles", "Should maintain consistent indentation")
		assert.Contains(t, outputStr, "2. Backend", "Should continue numbering for second item")
		assert.Contains(t, outputStr, "   - API routes", "Should indent second set of nested items")
		assert.Contains(t, outputStr, "   - Database models", "Should maintain nested structure")
	})

	t.Run("blockquote_and_code", func(t *testing.T) {
		t.Log("SPEC: Blockquote and Code Block Handling")
		t.Log("GIVEN an HTML document with blockquotes and code blocks")
		t.Log("WHEN sz processes the document with markdown rendering")
		t.Log("THEN it should format blockquotes and code correctly")

		// Create HTML with blockquotes and code
		blockHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Block Elements Test</title>
</head>
<body>
    <h2>Quote Example</h2>
    <blockquote>
        <p>This is a famous quote with <strong>emphasis</strong>.</p>
        <p>It spans multiple paragraphs.</p>
    </blockquote>
    <h2>Code Example</h2>
    <p>Here's some <code>inline code</code> in a paragraph.</p>
    <pre><code class="language-javascript">
function example() {
    return "Hello, World!";
}
</code></pre>
    <p>This function demonstrates basic syntax.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "block-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(blockHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test blockquote and code rendering
		cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should handle blockquotes
		assert.Contains(t, outputStr, "> This is a famous quote with **emphasis**.", "Should convert blockquote with > prefix")
		assert.Contains(t, outputStr, "> It spans multiple paragraphs.", "Should handle multi-paragraph blockquotes")

		// Should handle inline code
		assert.Contains(t, outputStr, "Here's some `inline code` in a paragraph.", "Should convert inline code with backticks")

		// Should handle code blocks
		assert.Contains(t, outputStr, "```javascript", "Should create fenced code block with language")
		assert.Contains(t, outputStr, "function example() {", "Should preserve code content")
		assert.Contains(t, outputStr, "return \"Hello, World!\";", "Should maintain code formatting")
		assert.Contains(t, outputStr, "```", "Should close code block")

		// Should not contain HTML
		assert.NotContains(t, outputStr, "<blockquote>", "Should not contain blockquote tags")
		assert.NotContains(t, outputStr, "<pre>", "Should not contain pre tags")
		assert.NotContains(t, outputStr, "<code>", "Should not contain code tags")
	})

	t.Run("link_formatting", func(t *testing.T) {
		t.Log("SPEC: Link Formatting")
		t.Log("GIVEN an HTML document with various types of links")
		t.Log("WHEN sz processes the document with markdown rendering")
		t.Log("THEN it should convert links to proper markdown format")

		// Create HTML with various links
		linkHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Link Test</title>
</head>
<body>
    <h2>External Links</h2>
    <p>Visit <a href="https://example.com">our website</a> for more information.</p>
    <p>You can also <a href="mailto:contact@example.com">email us</a> directly.</p>
    <h2>Internal Links</h2>
    <p>See the <a href="#section1">first section</a> for details.</p>
    <p>Complex link with <a href="https://github.com/user/repo" title="GitHub Repository">GitHub repo</a>.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "link-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(linkHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test link rendering
		cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should convert links to markdown format
		assert.Contains(t, outputStr, "[our website](https://example.com)", "Should convert external links")
		assert.Contains(t, outputStr, "[email us](mailto:contact@example.com)", "Should convert email links")
		assert.Contains(t, outputStr, "[first section](#section1)", "Should convert internal links")
		assert.Contains(t, outputStr, "[GitHub repo](https://github.com/user/repo)", "Should convert complex links")

		// Should preserve surrounding text
		assert.Contains(t, outputStr, "Visit [our website](https://example.com) for more information.", "Should preserve context")
		assert.Contains(t, outputStr, "You can also [email us](mailto:contact@example.com) directly.", "Should maintain paragraph flow")

		// Should not contain HTML
		assert.NotContains(t, outputStr, "<a href=", "Should not contain anchor tags")
		assert.NotContains(t, outputStr, "</a>", "Should not contain closing anchor tags")
	})

	t.Run("markdown_style_configuration", func(t *testing.T) {
		t.Log("SPEC: Markdown Style Configuration")
		t.Log("GIVEN different markdown style configuration options")
		t.Log("WHEN sz processes content with various style settings")
		t.Log("THEN it should respect configuration parameters")

		// Create HTML for style testing
		styleHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Style Test</title>
</head>
<body>
    <h1>Main Title</h1>
    <p>Text with <em>emphasis</em> and <strong>strong</strong> formatting.</p>
    <ul>
        <li>First item</li>
        <li>Second item</li>
    </ul>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "style-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(styleHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test default style
		t.Run("default_style", func(t *testing.T) {
			cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Command should succeed: %s", string(output))

			outputStr := string(output)
			assert.Contains(t, outputStr, "# Main Title", "Should use # for H1 by default")
			assert.Contains(t, outputStr, "*emphasis*", "Should use * for emphasis by default")
			assert.Contains(t, outputStr, "**strong**", "Should use ** for strong by default")
			assert.Contains(t, outputStr, "- First item", "Should use - for list items by default")
		})

		// Test alternative styles with configuration
		t.Run("underscore_emphasis", func(t *testing.T) {
			cmd := exec.Command(binary, "--markdown-renderer", "--emphasis-style=underscore", tmpFile.Name())
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Command should succeed: %s", string(output))

			outputStr := string(output)
			assert.Contains(t, outputStr, "_emphasis_", "Should use _ for emphasis when configured")
			assert.Contains(t, outputStr, "__strong__", "Should use __ for strong when configured")
		})

		t.Run("asterisk_lists", func(t *testing.T) {
			cmd := exec.Command(binary, "--markdown-renderer", "--list-style=asterisk", tmpFile.Name())
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Command should succeed: %s", string(output))

			outputStr := string(output)
			assert.Contains(t, outputStr, "* First item", "Should use * for list items when configured")
			assert.Contains(t, outputStr, "* Second item", "Should apply consistently")
		})
	})

	t.Run("performance_large_documents", func(t *testing.T) {
		t.Log("SPEC: Performance with Large Documents")
		t.Log("GIVEN a large HTML document with complex structure")
		t.Log("WHEN sz processes the document with markdown rendering")
		t.Log("THEN it should complete processing within reasonable time")

		// Create large HTML document
		var htmlBuilder strings.Builder
		htmlBuilder.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Large Document Test</title>
</head>
<body>
    <h1>Large Document Performance Test</h1>`)

		// Add many sections to test performance
		for i := 0; i < 50; i++ {
			htmlBuilder.WriteString(`
    <h2>Section ` + string(rune('A'+i%26)) + `</h2>
    <p>This is a paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
    <ul>
        <li>Item 1 with <a href="http://example.com">link</a></li>
        <li>Item 2 with <code>inline code</code></li>
        <li>Item 3 with more content</li>
    </ul>
    <blockquote>
        <p>This is a quote in section ` + string(rune('A'+i%26)) + `.</p>
    </blockquote>`)
		}

		htmlBuilder.WriteString(`
</body>
</html>`)

		tmpFile, err := os.CreateTemp("", "performance-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(htmlBuilder.String()))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test processing time
		start := time.Now()
		cmd := exec.Command(binary, "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)

		require.NoError(t, err, "Command should succeed: %s", string(output))

		// Should complete within reasonable time (3 seconds)
		assert.Less(t, duration, 3*time.Second, "Should process large documents efficiently")

		outputStr := string(output)
		// Should process all sections
		assert.Contains(t, outputStr, "# Large Document Performance Test", "Should process main heading")
		assert.Contains(t, outputStr, "## Section A", "Should process first section")
		assert.Contains(t, outputStr, "## Section Y", "Should process later sections")

		// Should maintain markdown quality
		assert.Contains(t, outputStr, "**bold**", "Should maintain formatting quality")
		assert.Contains(t, outputStr, "- Item 1 with [link](http://example.com)", "Should maintain complex structures")
		assert.Contains(t, outputStr, "> This is a quote", "Should maintain blockquotes")
	})

	t.Run("integration_with_content_processing", func(t *testing.T) {
		t.Log("SPEC: Integration with Content Processing")
		t.Log("GIVEN HTML with content filtering, media handling, and markdown rendering")
		t.Log("WHEN sz processes with all features enabled")
		t.Log("THEN it should apply processing pipeline and output clean markdown")

		// Create HTML with content, navigation, and media
		integrationHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Integration Test</title>
</head>
<body>
    <nav class="navigation">
        <ul>
            <li><a href="/">Home</a></li>
            <li><a href="/about">About</a></li>
        </ul>
    </nav>
    <main class="content">
        <h1>Article Title</h1>
        <p>This article discusses important topics.</p>
        <img src="chart.jpg" alt="Important chart">
        <h2>Key Points</h2>
        <ul>
            <li>First point</li>
            <li>Second point</li>
        </ul>
        <blockquote>
            <p>This is an important quote.</p>
        </blockquote>
    </main>
    <footer>
        <p>Copyright 2023</p>
    </footer>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "integration-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(integrationHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with full processing pipeline
		cmd := exec.Command(binary, "--content-filter", "--media-handler", "--markdown-renderer", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve main content as markdown
		assert.Contains(t, outputStr, "# Article Title", "Should preserve main heading as markdown")
		assert.Contains(t, outputStr, "This article discusses important topics.", "Should preserve article content")
		assert.Contains(t, outputStr, "## Key Points", "Should preserve subheading as markdown")

		// Should convert lists to markdown
		assert.Contains(t, outputStr, "- First point", "Should convert lists to markdown")
		assert.Contains(t, outputStr, "- Second point", "Should maintain list structure")

		// Should handle media elements
		assert.Contains(t, outputStr, "An image: Important chart", "Should process media with media handler")

		// Should convert blockquotes
		assert.Contains(t, outputStr, "> This is an important quote.", "Should convert blockquotes to markdown")

		// Should filter out navigation and footer
		assert.NotContains(t, outputStr, "Home", "Should filter navigation")
		assert.NotContains(t, outputStr, "About", "Should filter navigation links")
		assert.NotContains(t, outputStr, "Copyright 2023", "Should filter footer")

		// Should not contain HTML tags
		assert.NotContains(t, outputStr, "<h1>", "Should not contain HTML heading tags")
		assert.NotContains(t, outputStr, "<p>", "Should not contain HTML paragraph tags")
		assert.NotContains(t, outputStr, "<ul>", "Should not contain HTML list tags")
	})
}
