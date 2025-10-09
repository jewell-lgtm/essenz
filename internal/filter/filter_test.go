package filter

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/jewell-lgtm/essenz/internal/tree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContentFilterWithSimpleGitHubHTML(t *testing.T) {
	// Read the simple test HTML file
	htmlContent, err := os.ReadFile("../../testdata/simple_github_test.html")
	require.NoError(t, err)

	// Build the tree
	tb := tree.NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()
	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)
	require.NotNil(t, root)

	// Apply content filtering
	filter := NewContentFilter().WithDebugMode(true)

	t.Run("FilteredContentAnalysis", func(t *testing.T) {
		filteredRoot, err := filter.FilterTree(ctx, root)
		require.NoError(t, err)
		require.NotNil(t, filteredRoot)

		// Check what content remains after filtering
		allText := getAllTextContent(filteredRoot)
		t.Logf("Filtered content length: %d characters", len(allText))
		t.Logf("Filtered content preview: %s", truncateString(allText, 200))

		// Should preserve important content
		expectedContent := []string{
			"essenz",
			"sharp Unix tool",
			"Installation",
			"Usage",
			"Features",
		}

		for _, expected := range expectedContent {
			assert.Contains(t, allText, expected,
				"Filtered content should contain: %s", expected)
		}

		// Should filter out navigation
		navigationContent := []string{
			"Sign in",
			"Sign up",
		}

		for _, nav := range navigationContent {
			assert.NotContains(t, allText, nav,
				"Filtered content should not contain navigation: %s", nav)
		}
	})

	t.Run("ContentRegionFilter", func(t *testing.T) {
		// Test just the content region filter
		contentRegionFilter := NewContentRegionFilter()

		// Find content regions before filtering
		contentRegions := findContentRegions(root)
		t.Logf("Found %d content regions before filtering", len(contentRegions))

		for _, region := range contentRegions {
			text := getTextContent(region)
			if len(text) > 50 {
				text = text[:50] + "..."
			}
			t.Logf("Content region: %s (words: %d) - %s",
				region.Tag, region.WordCount, text)
		}

		// Test filter decisions on specific nodes
		ctx := &FilterContext{
			DocumentRoot: root,
		}

		// Test on navigation
		navNodes := findNodesByTag(root, []string{"nav"})
		for _, nav := range navNodes {
			shouldExclude := contentRegionFilter.ShouldExclude(nav, ctx)
			t.Logf("Navigation node should exclude: %v (semantic weight: %.2f)",
				shouldExclude, nav.SemanticWeight)
		}

		// Test on article content
		articles := findNodesByTag(root, []string{"article"})
		for _, article := range articles {
			shouldExclude := contentRegionFilter.ShouldExclude(article, ctx)
			t.Logf("Article node should exclude: %v (semantic weight: %.2f, content region: %v)",
				shouldExclude, article.SemanticWeight, article.IsContentRegion)
		}
	})

	t.Run("FilterRuleAnalysis", func(t *testing.T) {
		// Test individual filter rules
		filters := []FilterRule{
			NewSemanticTagFilter(),
			NewClassNameFilter(),
			NewLinkDensityFilter(0.3, 5),
			NewLengthFilter(15),
		}

		ctx := &FilterContext{
			DocumentRoot: root,
		}

		// Test filters on different node types
		testNodes := []struct {
			name  string
			nodes []*tree.TextNode
		}{
			{"navigation", findNodesByTag(root, []string{"nav"})},
			{"headings", findNodesByTag(root, []string{"h1", "h2"})},
			{"articles", findNodesByTag(root, []string{"article"})},
			{"paragraphs", findNodesByTag(root, []string{"p"})},
		}

		for _, testCase := range testNodes {
			t.Logf("\nTesting filters on %s nodes:", testCase.name)
			for _, node := range testCase.nodes {
				text := getTextContent(node)
				if len(text) > 30 {
					text = text[:30] + "..."
				}
				t.Logf("  Node: %s - %s", node.Tag, text)

				for _, filter := range filters {
					shouldExclude := filter.ShouldExclude(node, ctx)
					t.Logf("    %s: exclude=%v", filter.Name(), shouldExclude)
				}
			}
		}
	})
}

func TestContentFilterComparison(t *testing.T) {
	// Compare filter behavior between simple case and real GitHub case
	testCases := []struct {
		name     string
		filePath string
	}{
		{"simple", "../../testdata/simple_github_test.html"},
		{"real", "../../testdata/github_essenz_raw.html"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Real test is now working with content region preservation
			// No skip needed anymore!
			htmlContent, err := os.ReadFile(tc.filePath)
			if err != nil {
				if tc.name == "real" {
					t.Skip("Real GitHub HTML not available")
					return
				}
				require.NoError(t, err)
			}

			// Build tree
			tb := tree.NewTreeBuilder().WithPreserveAttributes(true)
			ctx := context.Background()
			root, err := tb.BuildTree(ctx, string(htmlContent))
			require.NoError(t, err)

			// Get unfiltered content stats
			unfilteredText := getAllTextContent(root)
			unfilteredWords := len(strings.Fields(unfilteredText))

			// Apply filtering
			filter := NewContentFilter()
			filteredRoot, err := filter.FilterTree(ctx, root)
			require.NoError(t, err)

			filteredText := getAllTextContent(filteredRoot)
			filteredWords := len(strings.Fields(filteredText))

			t.Logf("Unfiltered: %d words", unfilteredWords)
			t.Logf("Filtered: %d words", filteredWords)
			t.Logf("Retention rate: %.1f%%", float64(filteredWords)/float64(unfilteredWords)*100)

			// Check for key content retention
			keyContent := []string{"essenz", "Installation", "Usage", "sz https://example.com"}
			for _, key := range keyContent {
				if strings.Contains(unfilteredText, key) {
					assert.Contains(t, filteredText, key,
						"Should retain key content: %s", key)
				}
			}
		})
	}
}

// Helper functions
func findContentRegions(node *tree.TextNode) []*tree.TextNode {
	var regions []*tree.TextNode
	if node.IsContentRegion {
		regions = append(regions, node)
	}
	for _, child := range node.Children {
		regions = append(regions, findContentRegions(child)...)
	}
	return regions
}

func findNodesByTag(node *tree.TextNode, tags []string) []*tree.TextNode {
	var nodes []*tree.TextNode
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

func getAllTextContent(node *tree.TextNode) string {
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

func getTextContent(node *tree.TextNode) string {
	return getAllTextContent(node)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
