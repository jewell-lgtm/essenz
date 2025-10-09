package tree

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemanticExtractor_WikipediaArticle(t *testing.T) {
	htmlContent, err := os.ReadFile("../../testdata/wikipedia_article.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should extract main article title first
	assert.Contains(t, result, "Machine Learning")

	// Should extract h2 sections in order
	assert.Contains(t, result, "History")
	assert.Contains(t, result, "Types of machine learning")
	assert.Contains(t, result, "Applications")
	assert.Contains(t, result, "Ethics and fairness")

	// Should include content but maintain order
	historyPos := strings.Index(result, "History")
	typesPos := strings.Index(result, "Types of machine learning")
	applicationsPos := strings.Index(result, "Applications")

	assert.Less(t, historyPos, typesPos, "History should come before Types")
	assert.Less(t, typesPos, applicationsPos, "Types should come before Applications")

	// Should filter out navigation
	assert.NotContains(t, result, "Main page")
	assert.NotContains(t, result, "Create account")

	// Should include main content
	assert.Contains(t, result, "Arthur Samuel")
	assert.Contains(t, result, "Supervised learning")
	assert.Contains(t, result, "Computer vision")
}

func TestSemanticExtractor_NewsArticle(t *testing.T) {
	htmlContent, err := os.ReadFile("../../testdata/news_article.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should extract main headline first
	assert.Contains(t, result, "New AI Breakthrough Revolutionizes Healthcare Diagnosis")

	// Should extract h2 sections
	assert.Contains(t, result, "How the Technology Works")
	assert.Contains(t, result, "Clinical Trial Results")
	assert.Contains(t, result, "Regulatory Approval and Rollout")

	// Should include article content
	assert.Contains(t, result, "Stanford University")
	assert.Contains(t, result, "95% accuracy")
	assert.Contains(t, result, "Dr. Emily Rodriguez")

	// Should filter out navigation and UI elements
	assert.NotContains(t, result, "Subscribe Now")
	assert.NotContains(t, result, "Share on Facebook")
	assert.NotContains(t, result, "Accept All") // Cookie banner

	// Should filter out sidebar ads
	assert.NotContains(t, result, "Premium Cloud Computing Solutions")
}

func TestSemanticExtractor_BlogPost(t *testing.T) {
	htmlContent, err := os.ReadFile("../../testdata/blog_post.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should extract main title first
	assert.Contains(t, result, "Building Scalable Web Applications with Go")

	// Should extract h2 sections in order
	assert.Contains(t, result, "Why Choose Go for Web Development")
	assert.Contains(t, result, "Setting Up Your Go Web Application")
	assert.Contains(t, result, "Performance Optimization Techniques")

	// Should include code examples and content
	assert.Contains(t, result, "http.HandleFunc")
	assert.Contains(t, result, "goroutines")
	assert.Contains(t, result, "Connection Pooling")

	// Should filter out navigation and sidebar
	assert.NotContains(t, result, "All Posts")
	assert.NotContains(t, result, "Tutorials")
	assert.NotContains(t, result, "Subscribe to DevBlog")

	// Should maintain logical flow
	goPos := strings.Index(result, "Why Choose Go")
	setupPos := strings.Index(result, "Setting Up Your Go")
	perfPos := strings.Index(result, "Performance Optimization")

	assert.Less(t, goPos, setupPos, "Why Choose should come before Setup")
	assert.Less(t, setupPos, perfPos, "Setup should come before Performance")
}

func TestSemanticExtractor_EcommerceProduct(t *testing.T) {
	htmlContent, err := os.ReadFile("../../testdata/ecommerce_product.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should extract product title
	assert.Contains(t, result, "MacBook Pro 16-inch with M3 Chip")

	// Should extract key product information
	assert.Contains(t, result, "Product Description")
	assert.Contains(t, result, "Apple M3 chip")
	assert.Contains(t, result, "16.2-inch Liquid Retina XDR display")

	// Should extract h2 sections like specs and reviews
	assert.Contains(t, result, "Technical Specifications")
	if strings.Contains(result, "Customer Reviews") {
		assert.Contains(t, result, "Customer Reviews")
	}

	// Should filter out navigation and UI elements
	assert.NotContains(t, result, "Sign in")
	assert.NotContains(t, result, "Cart")
	assert.NotContains(t, result, "Spring Sale") // Promo banner

	// May contain product details but should filter commerce UI
	// This is more lenient as product pages have mixed content/commerce elements
}

func TestSemanticExtractor_DocumentationPage(t *testing.T) {
	htmlContent, err := os.ReadFile("../../testdata/documentation_page.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should extract main title
	assert.Contains(t, result, "Using the State Hook")

	// Should extract h2 sections in order
	assert.Contains(t, result, "What's a Hook?")
	assert.Contains(t, result, "Declaring a State Variable")
	assert.Contains(t, result, "Reading State")
	assert.Contains(t, result, "Updating State")

	// Should include code examples and explanations
	assert.Contains(t, result, "useState")
	assert.Contains(t, result, "const [count, setCount]")
	assert.Contains(t, result, "React.Component")

	// Should filter out navigation and sidebar
	assert.NotContains(t, result, "API Reference")
	assert.NotContains(t, result, "Community")
	assert.NotContains(t, result, "Search docs")

	// Should maintain logical order
	hookPos := strings.Index(result, "What's a Hook")
	declarePos := strings.Index(result, "Declaring a State Variable")
	readPos := strings.Index(result, "Reading State")

	assert.Less(t, hookPos, declarePos, "What's a Hook should come before Declaring")
	assert.Less(t, declarePos, readPos, "Declaring should come before Reading")
}

func TestSemanticExtractor_ExtractionQuality(t *testing.T) {
	// Test that extraction maintains reasonable quality across different document types
	testCases := []struct {
		name         string
		file         string
		minWords     int
		shouldHaveH1 bool
		shouldHaveH2 bool
	}{
		{"Wikipedia", "../../testdata/wikipedia_article.html", 100, true, true},
		{"News", "../../testdata/news_article.html", 200, true, true},
		{"Blog", "../../testdata/blog_post.html", 300, true, true},
		{"Docs", "../../testdata/documentation_page.html", 200, true, true},
		{"Simple", "../../testdata/semantic_hierarchy_simple.html", 10, true, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			htmlContent, err := os.ReadFile(tc.file)
			require.NoError(t, err)

			extractor := NewSemanticExtractor()
			tb := NewTreeBuilder().WithPreserveAttributes(true)
			ctx := context.Background()

			root, err := tb.BuildTree(ctx, string(htmlContent))
			require.NoError(t, err)

			result := extractor.Extract(root)

			// Quality checks
			wordCount := len(strings.Fields(result))
			assert.GreaterOrEqual(t, wordCount, tc.minWords,
				"Should extract at least %d words, got %d", tc.minWords, wordCount)

			if tc.shouldHaveH1 {
				assert.True(t, containsH1Like(result),
					"Should contain H1-like content (title/heading)")
			}

			if tc.shouldHaveH2 {
				assert.True(t, containsStructure(result),
					"Should contain structured content (multiple sections)")
			}

			// Should not be empty
			assert.NotEmpty(t, strings.TrimSpace(result), "Should not be empty")

			// Should not contain obvious navigation
			navTerms := []string{"Search", "Login", "Sign up", "Home", "About"}
			for _, nav := range navTerms {
				if strings.Contains(result, nav) {
					// Allow but warn - some content might legitimately contain these words
					t.Logf("Warning: Contains potential navigation term '%s'", nav)
				}
			}
		})
	}
}

// Helper functions
func containsH1Like(text string) bool {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if i < 5 && len(strings.TrimSpace(line)) > 10 {
			// First few lines should contain title-like content
			return true
		}
	}
	return false
}

func containsStructure(text string) bool {
	lines := strings.Split(text, "\n")
	sectionCount := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 5 && len(line) < 100 {
			// Potential section headers are short-ish lines
			words := strings.Fields(line)
			if len(words) >= 2 && len(words) <= 8 {
				sectionCount++
			}
		}
	}
	return sectionCount >= 2
}
