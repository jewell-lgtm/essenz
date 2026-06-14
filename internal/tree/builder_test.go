package tree

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTreeBuilderWithSimpleGitHubHTML(t *testing.T) {
	// Read the simple test HTML file
	htmlContent, err := os.ReadFile("../../testdata/simple_github_test.html")
	require.NoError(t, err)

	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)
	require.NotNil(t, root)

	// Test that content analysis identifies important sections
	t.Run("ContentRegionIdentification", func(t *testing.T) {
		contentRegions := findContentRegions(root)

		// Should find at least one content region
		assert.Greater(t, len(contentRegions), 0, "Should identify content regions")

		// Log content regions for debugging
		for _, region := range contentRegions {
			t.Logf("Content region: %s (words: %d, link density: %.2f, semantic weight: %.2f)",
				region.Tag, region.WordCount, region.LinkDensity, region.SemanticWeight)
		}
	})

	t.Run("SemanticWeighting", func(t *testing.T) {
		// Find headings and verify they have high semantic weight
		headings := findNodesByTag(root, []string{"h1", "h2"})
		assert.Greater(t, len(headings), 0, "Should find headings")

		for _, heading := range headings {
			assert.Greater(t, heading.SemanticWeight, 1.0,
				"Heading '%s' should have high semantic weight, got %.2f",
				getTextContent(heading), heading.SemanticWeight)
		}
	})

	t.Run("ContentDensityCalculation", func(t *testing.T) {
		// Find the main article content
		articles := findNodesByTag(root, []string{"article"})
		require.Greater(t, len(articles), 0, "Should find article element")

		article := articles[0]
		assert.Greater(t, article.WordCount, 20, "Article should have substantial word count")
		assert.Greater(t, article.ContentDensity, 5.0, "Article should have good content density")
		assert.Less(t, article.LinkDensity, 0.3, "Article should have low link density")
	})

	t.Run("NavigationDetection", func(t *testing.T) {
		// Find navigation elements
		navs := findNodesByTag(root, []string{"nav"})
		assert.Greater(t, len(navs), 0, "Should find navigation elements")

		for _, nav := range navs {
			assert.Less(t, nav.SemanticWeight, 1.0,
				"Navigation should have low semantic weight, got %.2f", nav.SemanticWeight)
		}
	})

	t.Run("ExpectedContentExtraction", func(t *testing.T) {
		// Verify specific content is captured
		allText := getAllTextContent(root)

		expectedContent := []string{
			"essenz",
			"sharp Unix tool",
			"Installation",
			"Usage",
			"Features",
			"sz https://example.com",
			"go install",
		}

		for _, expected := range expectedContent {
			assert.Contains(t, allText, expected,
				"Should capture important content: %s", expected)
		}

		// Navigation regions and their contents (e.g. the Sign in / Sign up /
		// About / Contact links) should be down-weighted, not the page containers
		// that merely happen to enclose them.
		navElements := findNodesByTag(root, []string{"nav"})
		require.Greater(t, len(navElements), 0, "fixture should contain nav elements")
		for _, nav := range navElements {
			assert.Less(t, nav.SemanticWeight, 1.0,
				"nav element should have low semantic weight, got %.2f", nav.SemanticWeight)
			for _, link := range findNodesByTag(nav, []string{"a"}) {
				assert.Less(t, link.SemanticWeight, 1.0,
					"link '%s' within nav should have low semantic weight, got %.2f",
					getTextContent(link), link.SemanticWeight)
			}
		}
	})
}

func TestContentAnalysisWithRealGitHubHTML(t *testing.T) {
	// Only run if we have the real GitHub HTML file
	htmlContent, err := os.ReadFile("../../testdata/github_essenz_raw.html")
	if err != nil {
		t.Skip("Real GitHub HTML not available, skipping test")
		return
	}

	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)
	require.NotNil(t, root)

	t.Run("ContentRegionIdentificationRealCase", func(t *testing.T) {
		contentRegions := findContentRegions(root)

		// Should find content regions in real GitHub page
		assert.Greater(t, len(contentRegions), 0, "Should identify content regions in real GitHub page")

		// Log the top content regions
		t.Logf("Found %d content regions in real GitHub HTML", len(contentRegions))
		for i, region := range contentRegions {
			if i < 5 { // Log first 5
				text := getTextContent(region)
				if len(text) > 100 {
					text = text[:100] + "..."
				}
				t.Logf("Content region %d: %s (words: %d, density: %.2f) - %s",
					i+1, region.Tag, region.WordCount, region.ContentDensity, text)
			}
		}
	})
}

// Helper functions

func findContentRegions(node *TextNode) []*TextNode {
	var regions []*TextNode
	if node.IsContentRegion {
		regions = append(regions, node)
	}
	for _, child := range node.Children {
		regions = append(regions, findContentRegions(child)...)
	}
	return regions
}

func findNodesByTag(node *TextNode, tags []string) []*TextNode {
	var nodes []*TextNode
	for _, tag := range tags {
		if strings.EqualFold(node.Tag, tag) {
			nodes = append(nodes, node)
		}
	}
	for _, child := range node.Children {
		nodes = append(nodes, findNodesByTag(child, tags)...)
	}
	return nodes
}

func getAllTextContent(node *TextNode) string {
	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}
	var parts []string
	for _, child := range node.Children {
		if text := getAllTextContent(child); text != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, " ")
}

func getTextContent(node *TextNode) string {
	return getAllTextContent(node)
}
