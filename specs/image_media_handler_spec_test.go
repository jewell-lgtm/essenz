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

// TestImageMediaHandlerSpec tests F4: Image and Media Handler functionality
func TestImageMediaHandlerSpec(t *testing.T) {
	// Build binary for testing
	binary := filepath.Join(os.TempDir(), "sz-media-test")
	err := exec.Command("go", "build", "-o", binary, "./cmd/essenz").Run()
	require.NoError(t, err, "Failed to build binary")
	defer func() { _ = os.Remove(binary) }()

	t.Run("basic_image_handling", func(t *testing.T) {
		t.Log("SPEC: Basic Image Handling")
		t.Log("GIVEN an HTML document with img tags containing alt text")
		t.Log("WHEN sz processes the document with media handling")
		t.Log("THEN it should replace images with descriptive text")

		// Create HTML with basic image
		imageHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Image Test</title>
</head>
<body>
    <p>This is an article about cats.</p>
    <img src="cat.jpg" alt="A fluffy orange cat sitting in a sunny window">
    <p>Cats love warm sunny spots.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "image-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(imageHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test image replacement
		cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve text content
		assert.Contains(t, outputStr, "This is an article about cats", "Should preserve article text")
		assert.Contains(t, outputStr, "Cats love warm sunny spots", "Should preserve following text")

		// Should replace image with descriptive text
		assert.Contains(t, outputStr, "An image: A fluffy orange cat sitting in a sunny window", "Should convert image to descriptive text")
		assert.NotContains(t, outputStr, "<img", "Should not contain img tags")
		assert.NotContains(t, outputStr, "cat.jpg", "Should not contain raw image URL")
	})

	t.Run("figure_with_caption", func(t *testing.T) {
		t.Log("SPEC: Figure Elements with Captions")
		t.Log("GIVEN an HTML document with figure and figcaption elements")
		t.Log("WHEN sz processes the document with media handling")
		t.Log("THEN it should combine image alt text with caption information")

		// Create HTML with figure and caption
		figureHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Figure Test</title>
</head>
<body>
    <p>Our quarterly results show strong growth.</p>
    <figure>
        <img src="chart.png" alt="Sales growth chart">
        <figcaption>Sales increased 25% over the past quarter</figcaption>
    </figure>
    <p>This trend is expected to continue.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "figure-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(figureHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test figure processing
		cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve surrounding text
		assert.Contains(t, outputStr, "Our quarterly results show strong growth", "Should preserve preceding text")
		assert.Contains(t, outputStr, "This trend is expected to continue", "Should preserve following text")

		// Should combine image and caption
		assert.Contains(t, outputStr, "An image: Sales growth chart", "Should include image alt text")
		assert.Contains(t, outputStr, "Sales increased 25% over the past quarter", "Should include caption text")

		// Should not contain raw HTML
		assert.NotContains(t, outputStr, "<figure>", "Should not contain figure tags")
		assert.NotContains(t, outputStr, "<figcaption>", "Should not contain figcaption tags")
	})

	t.Run("video_content_handling", func(t *testing.T) {
		t.Log("SPEC: Video Content Handling")
		t.Log("GIVEN an HTML document with video elements")
		t.Log("WHEN sz processes the document with media handling")
		t.Log("THEN it should replace videos with descriptive text")

		// Create HTML with video element
		videoHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Video Test</title>
</head>
<body>
    <p>Follow this tutorial to learn the basics.</p>
    <video controls>
        <source src="tutorial.mp4" type="video/mp4">
        <source src="tutorial.webm" type="video/webm">
        Your browser does not support the video tag.
    </video>
    <p>This tutorial shows the basic workflow.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "video-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(videoHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test video processing
		cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve text content
		assert.Contains(t, outputStr, "Follow this tutorial to learn the basics", "Should preserve preceding text")
		assert.Contains(t, outputStr, "This tutorial shows the basic workflow", "Should preserve following text")

		// Should replace video with descriptive text
		assert.Contains(t, outputStr, "A video: tutorial", "Should convert video to descriptive text")
		assert.NotContains(t, outputStr, "<video", "Should not contain video tags")
		assert.NotContains(t, outputStr, "<source", "Should not contain source tags")
	})

	t.Run("missing_alt_text_generation", func(t *testing.T) {
		t.Log("SPEC: Missing Alt Text Generation")
		t.Log("GIVEN an HTML document with images missing alt attributes")
		t.Log("WHEN sz processes the document with media handling")
		t.Log("THEN it should generate meaningful descriptions from context and URLs")

		// Create HTML with image missing alt text
		missingAltHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Missing Alt Test</title>
</head>
<body>
    <p>The new office building has impressive architecture.</p>
    <img src="office-building-exterior-modern-glass.jpg">
    <p>The design won several awards.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "missing-alt-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(missingAltHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test missing alt text handling
		cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve text content
		assert.Contains(t, outputStr, "The new office building has impressive architecture", "Should preserve preceding text")
		assert.Contains(t, outputStr, "The design won several awards", "Should preserve following text")

		// Should generate description from filename and context
		assert.Contains(t, outputStr, "An image", "Should indicate image presence")
		assert.Contains(t, outputStr, "office building", "Should extract context from filename")
		assert.Contains(t, outputStr, "modern glass", "Should extract architectural details from filename")
	})

	t.Run("social_media_embeds", func(t *testing.T) {
		t.Log("SPEC: Social Media Embed Handling")
		t.Log("GIVEN an HTML document with social media embeds")
		t.Log("WHEN sz processes the document with media handling")
		t.Log("THEN it should convert embeds to properly formatted quotes")

		// Create HTML with Twitter embed
		socialHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Social Media Test</title>
</head>
<body>
    <p>The company made an exciting announcement:</p>
    <blockquote class="twitter-tweet">
        <p>Excited to announce our new product launch! #innovation #tech</p>
        <a href="https://twitter.com/company/status/123">@company</a>
    </blockquote>
    <p>This generated significant engagement on social media.</p>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "social-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(socialHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test social media processing
		cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve surrounding text
		assert.Contains(t, outputStr, "The company made an exciting announcement", "Should preserve preceding text")
		assert.Contains(t, outputStr, "This generated significant engagement", "Should preserve following text")

		// Should convert to quote format
		assert.Contains(t, outputStr, "> Excited to announce our new product launch", "Should convert to blockquote format")
		assert.Contains(t, outputStr, "> — @company on Twitter", "Should include attribution")
		assert.NotContains(t, outputStr, "class=\"twitter-tweet\"", "Should not contain raw HTML attributes")
	})

	t.Run("media_handler_configuration", func(t *testing.T) {
		t.Log("SPEC: Media Handler Configuration")
		t.Log("GIVEN different media handler configuration options")
		t.Log("WHEN sz processes content with various media settings")
		t.Log("THEN it should respect configuration parameters")

		// Create HTML with multiple media types
		configHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Configuration Test</title>
</head>
<body>
    <h1>Media Configuration Test</h1>
    <img src="decorative-border.jpg" alt="Decorative border">
    <p>This is the main article content.</p>
    <img src="important-chart.jpg" alt="Important data visualization">
    <video controls>
        <source src="demo.mp4" type="video/mp4">
    </video>
</body>
</html>`

		tmpFile, err := os.CreateTemp("", "config-test*.html")
		require.NoError(t, err)
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		_, err = tmpFile.Write([]byte(configHTML))
		require.NoError(t, err)
		err = tmpFile.Close()
		require.NoError(t, err)

		// Test with different configurations
		t.Run("include_decorative_images", func(t *testing.T) {
			cmd := exec.Command(binary, "--media-handler", "--include-decorative", tmpFile.Name())
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Command should succeed: %s", string(output))

			outputStr := string(output)
			assert.Contains(t, outputStr, "Decorative border", "Should include decorative images when enabled")
		})

		t.Run("exclude_decorative_images", func(t *testing.T) {
			cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Command should succeed: %s", string(output))

			outputStr := string(output)
			assert.Contains(t, outputStr, "Important data visualization", "Should always include content images")
			// Decorative images should be excluded by default
		})
	})

	t.Run("integration_with_content_filter", func(t *testing.T) {
		t.Log("SPEC: Integration with Content Filter")
		t.Log("GIVEN HTML with media in both content and navigation areas")
		t.Log("WHEN sz processes with both content filtering and media handling")
		t.Log("THEN it should preserve media in content areas and filter media in navigation")

		// Create HTML with media in different contexts
		integrationHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Integration Test</title>
</head>
<body>
    <nav class="navigation">
        <img src="logo.jpg" alt="Company logo">
        <a href="/">Home</a>
    </nav>
    <main class="content">
        <h1>Article Title</h1>
        <p>This article discusses important topics.</p>
        <img src="content-image.jpg" alt="Relevant illustration">
        <p>The illustration shows key concepts.</p>
    </main>
    <footer>
        <img src="footer-logo.jpg" alt="Footer logo">
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

		// Test with both content filter and media handler
		cmd := exec.Command(binary, "--content-filter", "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed: %s", string(output))

		outputStr := string(output)

		// Should preserve content area media
		assert.Contains(t, outputStr, "Article Title", "Should preserve main content")
		assert.Contains(t, outputStr, "This article discusses important topics", "Should preserve article text")
		assert.Contains(t, outputStr, "An image: Relevant illustration", "Should preserve content images")

		// Should filter out navigation and footer
		assert.NotContains(t, outputStr, "Company logo", "Should filter navigation images")
		assert.NotContains(t, outputStr, "Footer logo", "Should filter footer images")
		assert.NotContains(t, outputStr, "Home", "Should filter navigation links")
	})

	t.Run("performance_with_large_documents", func(t *testing.T) {
		t.Log("SPEC: Performance with Large Documents")
		t.Log("GIVEN a large HTML document with many media elements")
		t.Log("WHEN sz processes the document with media handling")
		t.Log("THEN it should complete processing within reasonable time")

		// Create HTML with many media elements
		var htmlBuilder strings.Builder
		htmlBuilder.WriteString(`<!DOCTYPE html>
<html>
<head>
    <title>Performance Test</title>
</head>
<body>
    <h1>Large Document Test</h1>`)

		// Add many images to test performance
		for i := 0; i < 100; i++ {
			htmlBuilder.WriteString(`
    <p>Section ` + string(rune('A'+i%26)) + ` content with important information.</p>
    <img src="image` + string(rune('0'+i%10)) + `.jpg" alt="Test image ` + string(rune('0'+i%10)) + `">`)
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
		cmd := exec.Command(binary, "--media-handler", tmpFile.Name())
		output, err := cmd.CombinedOutput()
		duration := time.Since(start)

		require.NoError(t, err, "Command should succeed: %s", string(output))

		// Should complete within reasonable time (5 seconds)
		assert.Less(t, duration, 5*time.Second, "Should process large documents efficiently")

		outputStr := string(output)
		// Should process all images
		assert.Contains(t, outputStr, "An image: Test image", "Should process all images")
		assert.NotContains(t, outputStr, "<img", "Should not contain any img tags")
	})
}
