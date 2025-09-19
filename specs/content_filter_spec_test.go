package specs

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentFilterSpec(t *testing.T) {
	t.Run("semantic_tag_filtering", func(t *testing.T) {
		t.Log("SPEC: Semantic Tag Filtering")
		t.Log("GIVEN an HTML document with semantic navigation elements")
		t.Log("WHEN sz processes the document with content filtering")
		t.Log("THEN it should remove nav, header, footer, aside elements")

		binary := buildContentFilterBinary(t)

		// Create HTML with semantic navigation structure
		semanticHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Semantic Structure Test</title>
</head>
<body>
    <header>
        <h1>Site Header</h1>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="/about">About</a></li>
                <li><a href="/contact">Contact</a></li>
            </ul>
        </nav>
    </header>
    <main>
        <article>
            <h1>Main Article Title</h1>
            <p>This is the main article content that should be preserved.</p>
            <p>Second paragraph with important information.</p>
        </article>
    </main>
    <aside>
        <h2>Related Links</h2>
        <ul>
            <li><a href="/related1">Related Article 1</a></li>
            <li><a href="/related2">Related Article 2</a></li>
        </ul>
    </aside>
    <footer>
        <p>&copy; 2024 Test Site. All rights reserved.</p>
        <nav>
            <a href="/privacy">Privacy</a>
            <a href="/terms">Terms</a>
        </nav>
    </footer>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "semantic-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(semanticHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test content filtering
		cmd := exec.Command(binary, "--content-filter", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve main article content
		assert.Contains(t, outputStr, "Main Article Title", "Should preserve article title")
		assert.Contains(t, outputStr, "main article content", "Should preserve article content")
		assert.Contains(t, outputStr, "Second paragraph", "Should preserve all article paragraphs")

		// Should remove navigation elements
		assert.NotContains(t, outputStr, "Site Header", "Should remove header content")
		assert.NotContains(t, outputStr, "Home", "Should remove navigation links")
		assert.NotContains(t, outputStr, "About", "Should remove navigation links")
		assert.NotContains(t, outputStr, "Related Links", "Should remove aside content")
		assert.NotContains(t, outputStr, "All rights reserved", "Should remove footer content")
		assert.NotContains(t, outputStr, "Privacy", "Should remove footer navigation")
	})

	t.Run("css_class_filtering", func(t *testing.T) {
		t.Log("SPEC: CSS Class-Based Filtering")
		t.Log("GIVEN an HTML document with CSS class patterns indicating navigation")
		t.Log("WHEN sz applies content filtering")
		t.Log("THEN it should filter elements based on common CSS patterns")

		binary := buildContentFilterBinary(t)

		// Create HTML with CSS class patterns
		cssClassHTML := `<!DOCTYPE html>
<html>
<head>
    <title>CSS Class Pattern Test</title>
</head>
<body>
    <div class="navbar">
        <div class="nav-menu">
            <a href="/">Home</a>
            <a href="/products">Products</a>
        </div>
    </div>
    <div class="sidebar">
        <div class="advertisement">
            <h3>Special Offer!</h3>
            <p>Buy now and save 50%</p>
        </div>
        <div class="social-share">
            <button>Share on Facebook</button>
            <button>Tweet this</button>
        </div>
    </div>
    <div class="main-content">
        <h1>Product Review</h1>
        <p>This is a detailed review of the product.</p>
        <div class="related-posts">
            <h3>You might also like</h3>
            <a href="/post1">Similar Post 1</a>
            <a href="/post2">Similar Post 2</a>
        </div>
    </div>
    <div class="comments-section">
        <h3>Comments</h3>
        <div class="comment">User comment here</div>
    </div>
    <div class="breadcrumb">
        <a href="/">Home</a> > <a href="/category">Category</a> > Product
    </div>
    <div class="pagination">
        <a href="/prev">Previous</a>
        <a href="/next">Next</a>
    </div>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "css-class-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(cssClassHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test CSS class filtering
		cmd := exec.Command(binary, "--content-filter", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve main content
		assert.Contains(t, outputStr, "Product Review", "Should preserve main content title")
		assert.Contains(t, outputStr, "detailed review", "Should preserve main content")

		// Should remove navigation and non-content elements
		assert.NotContains(t, outputStr, "Home", "Should remove navbar links")
		assert.NotContains(t, outputStr, "Products", "Should remove navbar links")
		assert.NotContains(t, outputStr, "Special Offer", "Should remove advertisements")
		assert.NotContains(t, outputStr, "Share on Facebook", "Should remove social sharing")
		assert.NotContains(t, outputStr, "You might also like", "Should remove related posts")
		assert.NotContains(t, outputStr, "Comments", "Should remove comments section")
		assert.NotContains(t, outputStr, "Previous", "Should remove pagination")
		assert.NotContains(t, outputStr, "Category", "Should remove breadcrumb")
	})

	t.Run("link_density_filtering", func(t *testing.T) {
		t.Log("SPEC: Link Density Content Filtering")
		t.Log("GIVEN content sections with varying link-to-text ratios")
		t.Log("WHEN sz analyzes content density")
		t.Log("THEN it should filter high-link-density sections")

		binary := buildContentFilterBinary(t)

		// Create HTML with varying link densities
		linkDensityHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Link Density Test</title>
</head>
<body>
    <div class="high-link-density">
        <a href="/link1">Link 1</a>
        <a href="/link2">Link 2</a>
        <a href="/link3">Link 3</a>
        <span>Some text</span>
        <a href="/link4">Link 4</a>
        <a href="/link5">Link 5</a>
    </div>
    <div class="main-article">
        <h1>Article Title</h1>
        <p>This is a substantial paragraph with meaningful content that should be preserved.
        It contains valuable information about the topic at hand and provides context.</p>
        <p>Another paragraph with <a href="/relevant">one relevant link</a> but mostly text content
        that adds value to the reader's understanding of the subject matter.</p>
    </div>
    <div class="navigation-menu">
        <a href="/home">Home</a> |
        <a href="/about">About</a> |
        <a href="/services">Services</a> |
        <a href="/contact">Contact</a> |
        <a href="/blog">Blog</a>
    </div>
    <div class="short-content">
        <p>Brief text.</p>
    </div>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "link-density-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(linkDensityHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test link density filtering
		cmd := exec.Command(binary, "--content-filter", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve content with appropriate link density
		assert.Contains(t, outputStr, "Article Title", "Should preserve main article title")
		assert.Contains(t, outputStr, "substantial paragraph", "Should preserve main content")
		assert.Contains(t, outputStr, "Another paragraph", "Should preserve content with low link density")

		// Should remove high-link-density sections
		assert.NotContains(t, outputStr, "Link 1", "Should remove high-link-density section")
		assert.NotContains(t, outputStr, "Link 2", "Should remove high-link-density section")
		assert.NotContains(t, outputStr, "Home", "Should remove navigation menu")
		assert.NotContains(t, outputStr, "Services", "Should remove navigation menu")

		// Should remove very short content
		assert.NotContains(t, outputStr, "Brief text", "Should remove very short content blocks")
	})

	t.Run("whitelist_preservation", func(t *testing.T) {
		t.Log("SPEC: Whitelist Content Preservation")
		t.Log("GIVEN HTML with content in whitelisted containers")
		t.Log("WHEN sz applies content filtering")
		t.Log("THEN it should preserve whitelisted content regardless of other rules")

		binary := buildContentFilterBinary(t)

		// Create HTML with whitelisted and non-whitelisted content
		whitelistHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Whitelist Test</title>
</head>
<body>
    <nav class="main-nav">
        <a href="/">Home</a>
        <a href="/about">About</a>
    </nav>
    <main>
        <article>
            <h1>Protected Article</h1>
            <p>This content is in a main article tag and should be preserved.</p>
            <nav class="article-nav">
                <a href="/prev">Previous Article</a>
                <a href="/next">Next Article</a>
            </nav>
        </article>
    </main>
    <div class="content">
        <h2>Content Section</h2>
        <p>This is in a content div and should be preserved.</p>
        <div class="social">
            <a href="/share">Share</a>
        </div>
    </div>
    <div class="sidebar">
        <div class="ads">
            <p>Advertisement content that should be removed.</p>
        </div>
    </div>
    <div class="post">
        <h3>Blog Post</h3>
        <p>Content in a post container should be preserved.</p>
    </div>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "whitelist-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(whitelistHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test whitelist preservation
		cmd := exec.Command(binary, "--content-filter", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve whitelisted content containers
		assert.Contains(t, outputStr, "Protected Article", "Should preserve main/article content")
		assert.Contains(t, outputStr, "main article tag", "Should preserve main/article content")
		assert.Contains(t, outputStr, "Content Section", "Should preserve .content div")
		assert.Contains(t, outputStr, "content div", "Should preserve .content div")
		assert.Contains(t, outputStr, "Blog Post", "Should preserve .post div")
		assert.Contains(t, outputStr, "post container", "Should preserve .post div")

		// Should preserve navigation within whitelisted containers
		assert.Contains(t, outputStr, "Previous Article", "Should preserve article navigation")
		assert.Contains(t, outputStr, "Next Article", "Should preserve article navigation")

		// Should still remove non-whitelisted content
		assert.NotContains(t, outputStr, "Home", "Should remove main navigation")
		assert.NotContains(t, outputStr, "About", "Should remove main navigation")
		assert.NotContains(t, outputStr, "Advertisement content", "Should remove ads")
		assert.NotContains(t, outputStr, "Share", "Should remove social sharing outside whitelist")
	})

	t.Run("modern_framework_filtering", func(t *testing.T) {
		t.Log("SPEC: Modern Framework Content Filtering")
		t.Log("GIVEN HTML with modern framework CSS class patterns")
		t.Log("WHEN sz processes React/Vue/Next.js generated content")
		t.Log("THEN it should filter based on component patterns")

		binary := buildContentFilterBinary(t)

		// Create HTML with modern framework patterns
		frameworkHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Framework Pattern Test</title>
</head>
<body>
    <div class="Header_container__xyz123">
        <div class="Navigation_nav__abc456">
            <a href="/">Home</a>
            <a href="/products">Products</a>
        </div>
    </div>
    <div class="Layout_main__def789">
        <div class="Article_content__ghi012">
            <h1>Framework-Generated Article</h1>
            <p class="Body_text__jkl345">This is the main article content generated by a modern framework.</p>
            <p class="Paragraph_content__mno678">Second paragraph with substantial content that should be preserved.</p>
        </div>
    </div>
    <div class="Sidebar_aside__pqr901">
        <div class="Widget_ad__stu234">
            <h3>Sponsored Content</h3>
            <p>This is an advertisement.</p>
        </div>
        <div class="Social_share__vwx567">
            <button>Share</button>
            <button>Like</button>
        </div>
    </div>
    <div class="Footer_footer__yza890">
        <div class="Links_nav__bcd123">
            <a href="/privacy">Privacy</a>
            <a href="/terms">Terms</a>
        </div>
    </div>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "framework-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(frameworkHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test modern framework filtering
		cmd := exec.Command(binary, "--content-filter", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve main content despite framework class names
		assert.Contains(t, outputStr, "Framework-Generated Article", "Should preserve article title")
		assert.Contains(t, outputStr, "main article content", "Should preserve article content")
		assert.Contains(t, outputStr, "Second paragraph", "Should preserve article paragraphs")

		// Should remove framework-generated navigation and sidebars
		assert.NotContains(t, outputStr, "Home", "Should remove navigation despite framework classes")
		assert.NotContains(t, outputStr, "Products", "Should remove navigation despite framework classes")
		assert.NotContains(t, outputStr, "Sponsored Content", "Should remove ads despite framework classes")
		assert.NotContains(t, outputStr, "Privacy", "Should remove footer navigation")
		assert.NotContains(t, outputStr, "Share", "Should remove social widgets")
	})

	t.Run("edge_case_handling", func(t *testing.T) {
		t.Log("SPEC: Edge Case Content Filtering")
		t.Log("GIVEN complex HTML with edge cases")
		t.Log("WHEN sz applies intelligent content filtering")
		t.Log("THEN it should handle ambiguous cases correctly")

		binary := buildContentFilterBinary(t)

		// Create HTML with edge cases
		edgeCaseHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Edge Case Test</title>
</head>
<body>
    <article>
        <h1>Reference List Article</h1>
        <p>This article contains many links but is legitimate content.</p>
        <h2>Useful Resources</h2>
        <ul>
            <li><a href="/resource1">Important Resource 1</a> - Description of resource</li>
            <li><a href="/resource2">Important Resource 2</a> - Another useful resource</li>
            <li><a href="/resource3">Important Resource 3</a> - More helpful information</li>
        </ul>
        <p>This reference list provides valuable information despite high link density.</p>
    </article>
    <div class="navigation-looking">
        <h3>Article Navigation</h3>
        <p>This looks like navigation but is actually article content explaining navigation concepts.</p>
        <p>Navigation design is important for user experience.</p>
    </div>
    <main>
        <div class="short-but-important">
            <p><strong>Note:</strong> Critical information.</p>
        </div>
    </main>
    <section class="ads">
        <div class="content">
            <h3>Advertisement</h3>
            <p>This is clearly an ad despite being in a content div.</p>
        </div>
    </section>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "edge-case-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(edgeCaseHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test edge case handling
		cmd := exec.Command(binary, "--content-filter", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve legitimate content even with high link density
		assert.Contains(t, outputStr, "Reference List Article", "Should preserve article title")
		assert.Contains(t, outputStr, "many links but is legitimate", "Should preserve article intro")
		assert.Contains(t, outputStr, "Useful Resources", "Should preserve reference section")
		assert.Contains(t, outputStr, "Important Resource 1", "Should preserve reference links")
		assert.Contains(t, outputStr, "valuable information", "Should preserve article conclusion")

		// Should preserve content that looks like navigation but is actually content
		assert.Contains(t, outputStr, "Article Navigation", "Should preserve content about navigation")
		assert.Contains(t, outputStr, "navigation concepts", "Should preserve navigation discussion")

		// Should preserve short but important content in main containers
		assert.Contains(t, outputStr, "Critical information", "Should preserve important short content")

		// Should remove actual advertisements even in content containers
		assert.NotContains(t, outputStr, "Advertisement", "Should remove ads even in content div")
		assert.NotContains(t, outputStr, "clearly an ad", "Should remove ad content")
	})

	t.Run("filter_configuration", func(t *testing.T) {
		t.Log("SPEC: Content Filter Configuration")
		t.Log("GIVEN different filtering configuration options")
		t.Log("WHEN sz applies configurable content filtering")
		t.Log("THEN it should respect configuration parameters")

		binary := buildContentFilterBinary(t)

		// Create HTML for configuration testing
		configHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Configuration Test</title>
</head>
<body>
    <div class="main-content">
        <h1>Main Article</h1>
        <p>This is the main article content.</p>
    </div>
    <div class="custom-nav">
        <p>Custom navigation content.</p>
    </div>
    <div class="short">
        <p>Short content that might be filtered based on length settings.</p>
    </div>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "config-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(configHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test aggressive filtering mode
		cmd := exec.Command(binary, "--content-filter", "--aggressive-filtering", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve main content
		assert.Contains(t, outputStr, "Main Article", "Should preserve main content")
		assert.Contains(t, outputStr, "main article content", "Should preserve main content")

		// Test with custom whitelist
		cmd = exec.Command(binary, "--content-filter", "--preserve-selector=.custom-nav", tmpFile.Name())
		output, err = cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr = string(output)

		// Should preserve custom whitelisted content
		assert.Contains(t, outputStr, "Custom navigation", "Should preserve custom whitelisted content")
	})
}

// buildContentFilterBinary builds the sz binary for testing content filter functionality
func buildContentFilterBinary(t *testing.T) string {
	// Build the sz binary from project root
	cmd := exec.Command("go", "build", "-o", "/tmp/sz-content-filter-test", "./cmd/essenz")
	// Set working directory to project root
	cmd.Dir = ".."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build binary for content filter testing")

	return "/tmp/sz-content-filter-test"
}
