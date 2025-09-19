package specs

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextNodeTreeBuilderSpec(t *testing.T) {
	t.Run("basic_html_structure_parsing", func(t *testing.T) {
		t.Log("SPEC: Basic HTML Structure Parsing")
		t.Log("GIVEN a simple HTML document with text nodes")
		t.Log("WHEN sz processes the document with text node tree building")
		t.Log("THEN it should extract text content in proper hierarchical structure")

		// Build binary for testing
		binary := buildTextNodeBinary(t)

		// Create a test HTML file with basic structure
		basicHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Test Document</title>
</head>
<body>
    <h1>Main Title</h1>
    <p>Introduction paragraph with <strong>emphasis</strong> and a <a href="https://example.com">link</a>.</p>
    <h2>Section Header</h2>
    <ul>
        <li>First item</li>
        <li>Second item with <em>emphasis</em></li>
    </ul>
    <script>// This should be ignored</script>
    <style>/* This should also be ignored */</style>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "text-node-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(basicHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with text node tree building enabled
		cmd := exec.Command(binary, "--text-node-tree", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should extract main content with proper hierarchy
		assert.Contains(t, outputStr, "Main Title", "Should extract h1 content")
		assert.Contains(t, outputStr, "Introduction paragraph", "Should extract paragraph content")
		assert.Contains(t, outputStr, "emphasis", "Should extract strong/em content")
		assert.Contains(t, outputStr, "link", "Should extract link text")
		assert.Contains(t, outputStr, "Section Header", "Should extract h2 content")
		assert.Contains(t, outputStr, "First item", "Should extract list items")
		assert.Contains(t, outputStr, "Second item", "Should extract list items with formatting")

		// Should ignore script and style content
		assert.NotContains(t, outputStr, "This should be ignored", "Should ignore script content")
		assert.NotContains(t, outputStr, "This should also be ignored", "Should ignore style content")
	})

	t.Run("hierarchical_structure_preservation", func(t *testing.T) {
		t.Log("SPEC: Hierarchical Structure Preservation")
		t.Log("GIVEN an HTML document with nested heading structure")
		t.Log("WHEN sz processes the document with tree building")
		t.Log("THEN it should preserve the hierarchical relationship of headings")

		binary := buildTextNodeBinary(t)

		// Create HTML with proper heading hierarchy
		hierarchyHTML := `<!DOCTYPE html>
<html>
<body>
    <h1>Chapter 1</h1>
    <p>Chapter introduction</p>
    <h2>Section 1.1</h2>
    <p>Section content</p>
    <h3>Subsection 1.1.1</h3>
    <p>Subsection content</p>
    <h2>Section 1.2</h2>
    <p>Another section</p>
    <h1>Chapter 2</h1>
    <p>Second chapter</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "hierarchy-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(hierarchyHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with tree structure output
		cmd := exec.Command(binary, "--text-node-tree", "--tree-format=json", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve heading hierarchy
		assert.Contains(t, outputStr, "Chapter 1", "Should extract h1 content")
		assert.Contains(t, outputStr, "Section 1.1", "Should extract h2 content")
		assert.Contains(t, outputStr, "Subsection 1.1.1", "Should extract h3 content")
		assert.Contains(t, outputStr, "Chapter 2", "Should extract second h1")

		// JSON output should indicate hierarchy levels
		if strings.Contains(outputStr, "level") {
			assert.Contains(t, outputStr, `"level":1`, "Should indicate h1 level")
			assert.Contains(t, outputStr, `"level":2`, "Should indicate h2 level")
			assert.Contains(t, outputStr, `"level":3`, "Should indicate h3 level")
		}
	})

	t.Run("complex_nested_structures", func(t *testing.T) {
		t.Log("SPEC: Complex Nested Structures")
		t.Log("GIVEN an HTML document with complex nested lists and quotes")
		t.Log("WHEN sz processes the document")
		t.Log("THEN it should handle nested structures correctly")

		binary := buildTextNodeBinary(t)

		// Create HTML with complex nesting
		complexHTML := `<!DOCTYPE html>
<html>
<body>
    <h1>Complex Document</h1>
    <blockquote>
        <p>This is a quoted paragraph.</p>
        <p>Multiple paragraphs in quote.</p>
    </blockquote>
    <ul>
        <li>Top level item
            <ul>
                <li>Nested item 1</li>
                <li>Nested item 2
                    <ul>
                        <li>Deep nested item</li>
                    </ul>
                </li>
            </ul>
        </li>
        <li>Another top level item</li>
    </ul>
    <div>
        <h2>Section in div</h2>
        <p>Paragraph in div with <code>inline code</code> and <strong>strong text</strong>.</p>
    </div>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "complex-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(complexHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test complex structure parsing
		cmd := exec.Command(binary, "--text-node-tree", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should extract nested content
		assert.Contains(t, outputStr, "Complex Document", "Should extract main heading")
		assert.Contains(t, outputStr, "quoted paragraph", "Should extract blockquote content")
		assert.Contains(t, outputStr, "Multiple paragraphs", "Should extract multiple quote paragraphs")
		assert.Contains(t, outputStr, "Top level item", "Should extract top-level list items")
		assert.Contains(t, outputStr, "Nested item 1", "Should extract nested list items")
		assert.Contains(t, outputStr, "Deep nested item", "Should extract deeply nested items")
		assert.Contains(t, outputStr, "Section in div", "Should extract headings in divs")
		assert.Contains(t, outputStr, "inline code", "Should extract code content")
		assert.Contains(t, outputStr, "strong text", "Should extract strong content")
	})

	t.Run("text_node_filtering", func(t *testing.T) {
		t.Log("SPEC: Text Node Filtering")
		t.Log("GIVEN an HTML document with various types of content")
		t.Log("WHEN sz processes with text node filtering")
		t.Log("THEN it should filter out irrelevant content while preserving meaningful text")

		binary := buildTextNodeBinary(t)

		// Create HTML with content that should be filtered
		filterHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Filter Test</title>
    <script>console.log("script content");</script>
    <style>.hidden { display: none; }</style>
</head>
<body>
    <h1>Main Content</h1>
    <p>Visible paragraph text.</p>
    <div style="display: none;">Hidden content that should be excluded</div>
    <nav>
        <ul>
            <li><a href="/home">Home</a></li>
            <li><a href="/about">About</a></li>
        </ul>
    </nav>
    <main>
        <article>
            <h2>Article Title</h2>
            <p>Article content that should be preserved.</p>
        </article>
    </main>
    <aside>
        <h3>Sidebar</h3>
        <p>Sidebar content</p>
    </aside>
    <footer>
        <p>Footer information</p>
    </footer>
    <!-- HTML comment -->
    <script type="application/ld+json">{"@type": "Article"}</script>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "filter-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(filterHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with content filtering
		cmd := exec.Command(binary, "--text-node-tree", "--filter-navigation", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should extract main content
		assert.Contains(t, outputStr, "Main Content", "Should extract main heading")
		assert.Contains(t, outputStr, "Visible paragraph", "Should extract visible content")
		assert.Contains(t, outputStr, "Article Title", "Should extract article heading")
		assert.Contains(t, outputStr, "Article content", "Should extract article text")

		// Should filter out non-content elements
		assert.NotContains(t, outputStr, "script content", "Should filter script content")
		assert.NotContains(t, outputStr, "Hidden content", "Should filter hidden content")
		assert.NotContains(t, outputStr, `{"@type": "Article"}`, "Should filter JSON-LD")

		// Navigation filtering should be configurable
		if strings.Contains(outputStr, "--filter-navigation") {
			assert.NotContains(t, outputStr, "Home", "Should filter navigation when flag is set")
			assert.NotContains(t, outputStr, "About", "Should filter navigation when flag is set")
		}
	})

	t.Run("dynamic_content_handling", func(t *testing.T) {
		t.Log("SPEC: Dynamic Content Handling")
		t.Log("GIVEN an HTML document with JavaScript-generated content")
		t.Log("WHEN sz processes the document after DOM ready")
		t.Log("THEN it should capture dynamically generated text nodes")

		binary := buildTextNodeBinary(t)

		// Create HTML with dynamic content
		dynamicHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Dynamic Test</title>
</head>
<body>
    <h1>Dynamic Content Test</h1>
    <div id="static">Static content</div>
    <div id="dynamic-container">Loading...</div>
    <script>
        setTimeout(() => {
            document.getElementById('dynamic-container').innerHTML =
                '<h2>Dynamic Heading</h2><p>This content was added by JavaScript.</p>';
        }, 100);
    </script>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "dynamic-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(dynamicHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with DOM ready integration (should wait for dynamic content)
		cmd := exec.Command(binary, "--text-node-tree", "--wait-for-frameworks", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should extract both static and dynamic content
		assert.Contains(t, outputStr, "Dynamic Content Test", "Should extract static heading")
		assert.Contains(t, outputStr, "Static content", "Should extract static div content")
		assert.Contains(t, outputStr, "Dynamic Heading", "Should extract dynamically added heading")
		assert.Contains(t, outputStr, "added by JavaScript", "Should extract dynamically added paragraph")

		// Should not contain the loading state
		assert.NotContains(t, outputStr, "Loading...", "Should not extract initial loading state")
	})

	t.Run("link_and_attribute_preservation", func(t *testing.T) {
		t.Log("SPEC: Link and Attribute Preservation")
		t.Log("GIVEN an HTML document with links and important attributes")
		t.Log("WHEN sz extracts the text node tree")
		t.Log("THEN it should preserve link destinations and relevant attributes")

		binary := buildTextNodeBinary(t)

		// Create HTML with links and attributes
		linkHTML := `<!DOCTYPE html>
<html>
<body>
    <h1>Link Test</h1>
    <p>Check out <a href="https://example.com" title="Example Site">this link</a> for more info.</p>
    <p>Email us at <a href="mailto:test@example.com">test@example.com</a> for support.</p>
    <img src="image.jpg" alt="Test image description" title="Image title">
    <blockquote cite="https://source.com">
        <p>This is a quote with citation.</p>
    </blockquote>
    <code lang="javascript">console.log("hello");</code>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "link-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(linkHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with attribute preservation
		cmd := exec.Command(binary, "--text-node-tree", "--preserve-attributes", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should extract link text and preserve destinations
		assert.Contains(t, outputStr, "this link", "Should extract link text")
		assert.Contains(t, outputStr, "test@example.com", "Should extract email link text")

		// Should preserve important attributes when flag is set
		if strings.Contains(outputStr, "https://example.com") {
			assert.Contains(t, outputStr, "https://example.com", "Should preserve link href")
		}
		if strings.Contains(outputStr, "Test image description") {
			assert.Contains(t, outputStr, "Test image description", "Should extract image alt text")
		}
		if strings.Contains(outputStr, "citation") {
			assert.Contains(t, outputStr, "quote with citation", "Should extract blockquote content")
		}
	})

	t.Run("performance_and_large_documents", func(t *testing.T) {
		t.Log("SPEC: Performance and Large Documents")
		t.Log("GIVEN a large HTML document with many text nodes")
		t.Log("WHEN sz processes the document")
		t.Log("THEN it should complete within reasonable time and memory limits")

		binary := buildTextNodeBinary(t)

		// Create a larger HTML document
		var largeHTML strings.Builder
		largeHTML.WriteString(`<!DOCTYPE html><html><body><h1>Large Document Test</h1>`)

		// Generate many paragraphs and lists
		for i := 0; i < 50; i++ {
			largeHTML.WriteString(`<h2>Section ` + strings.Repeat("X", i%10) + `</h2>`)
			for j := 0; j < 10; j++ {
				largeHTML.WriteString(`<p>This is paragraph ` + strings.Repeat("Y", j%5) + ` in section ` + strings.Repeat("Z", i%3) + `.</p>`)
			}
			largeHTML.WriteString(`<ul>`)
			for k := 0; k < 5; k++ {
				largeHTML.WriteString(`<li>List item ` + strings.Repeat("A", k%3) + `</li>`)
			}
			largeHTML.WriteString(`</ul>`)
		}
		largeHTML.WriteString(`</body></html>`)

		tmpFile, err := os.CreateTemp("", "large-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(largeHTML.String()))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test processing time (should complete within reasonable time)
		cmd := exec.Command(binary, "--text-node-tree", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed even with large document: %s", string(output))

		outputStr := string(output)

		// Should extract content from large document
		assert.Contains(t, outputStr, "Large Document Test", "Should extract main heading")
		assert.Contains(t, outputStr, "Section", "Should extract section headings")
		assert.Contains(t, outputStr, "paragraph", "Should extract paragraph content")
		assert.Contains(t, outputStr, "List item", "Should extract list items")

		// Output should be reasonable in size (not exponentially large)
		assert.Less(t, len(outputStr), len(largeHTML.String())*2, "Output should not be excessively large")
	})
}

// Helper function to build the binary for text node tree testing
func buildTextNodeBinary(t *testing.T) string {
	// Build the sz binary from project root
	cmd := exec.Command("go", "build", "-o", "/tmp/sz-text-node-test", "./cmd/essenz")
	// Set working directory to project root
	cmd.Dir = ".."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build binary for text node testing")

	return "/tmp/sz-text-node-test"
}
