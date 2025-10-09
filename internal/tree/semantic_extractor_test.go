package tree

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemanticExtractor_SingleH1(t *testing.T) {
	// Test case: Simple document with single h1
	htmlContent, err := os.ReadFile("../../testdata/semantic_hierarchy_simple.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should find the h1 and its content first
	assert.Contains(t, result, "Main Article Title")
	assert.Contains(t, result, "This is the first paragraph")
	assert.Contains(t, result, "This is the second paragraph")

	// Should filter out navigation
	assert.NotContains(t, result, "Home")
	assert.NotContains(t, result, "About")

	// Should filter out footer
	assert.NotContains(t, result, "Copyright 2024")

	// h1 should appear near the beginning
	lines := strings.Split(result, "\n")
	h1Found := false
	h1Line := -1
	for i, line := range lines {
		if strings.Contains(line, "Main Article Title") {
			h1Found = true
			h1Line = i
			break
		}
	}
	assert.True(t, h1Found, "H1 should be found in output")
	assert.Less(t, h1Line, 5, "H1 should appear near the beginning")
}

func TestSemanticExtractor_MultipleH1Siblings(t *testing.T) {
	// Test case: Multiple h1s that are siblings
	htmlContent, err := os.ReadFile("../../testdata/semantic_multiple_h1_siblings.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should contain all h1s and their content
	assert.Contains(t, result, "First Article")
	assert.Contains(t, result, "Content for first article")
	assert.Contains(t, result, "Second Article")
	assert.Contains(t, result, "Content for second article")
	assert.Contains(t, result, "Third Article")
	assert.Contains(t, result, "Content for third article")

	// Should maintain document order for siblings
	firstPos := strings.Index(result, "First Article")
	secondPos := strings.Index(result, "Second Article")
	thirdPos := strings.Index(result, "Third Article")

	assert.Less(t, firstPos, secondPos, "First article should come before second")
	assert.Less(t, secondPos, thirdPos, "Second article should come before third")

	// Should filter out navigation
	assert.NotContains(t, result, "Home")
	assert.NotContains(t, result, "About")
}

func TestSemanticExtractor_NestedH1s(t *testing.T) {
	// Test case: h1s at different nesting levels
	htmlContent, err := os.ReadFile("../../testdata/semantic_nested_h1s.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should find all h1s
	assert.Contains(t, result, "Top Level Title")
	assert.Contains(t, result, "Article Title")
	assert.Contains(t, result, "Section Title")

	// Should include their content
	assert.Contains(t, result, "Top level content")
	assert.Contains(t, result, "Article content paragraph 1")
	assert.Contains(t, result, "Section content")

	// First h1 in document order should come first
	topPos := strings.Index(result, "Top Level Title")
	articlePos := strings.Index(result, "Article Title")
	sectionPos := strings.Index(result, "Section Title")

	assert.Less(t, topPos, articlePos, "Top level h1 should come before article h1")
	assert.Less(t, articlePos, sectionPos, "Article h1 should come before section h1")
}

func TestSemanticExtractor_H1AndH2Hierarchy(t *testing.T) {
	// Test case: h1s followed by h2s
	htmlContent, err := os.ReadFile("../../testdata/semantic_h1_h2_hierarchy.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should contain h1 and all h2s
	assert.Contains(t, result, "Main Article Title")
	assert.Contains(t, result, "First Section")
	assert.Contains(t, result, "Second Section")
	assert.Contains(t, result, "Third Section")

	// h1 should come before any h2
	h1Pos := strings.Index(result, "Main Article Title")
	h2FirstPos := strings.Index(result, "First Section")

	assert.Less(t, h1Pos, h2FirstPos, "H1 should come before any H2")

	// h2s should maintain document order
	firstH2Pos := strings.Index(result, "First Section")
	secondH2Pos := strings.Index(result, "Second Section")
	thirdH2Pos := strings.Index(result, "Third Section")

	assert.Less(t, firstH2Pos, secondH2Pos, "First H2 should come before second H2")
	assert.Less(t, secondH2Pos, thirdH2Pos, "Second H2 should come before third H2")

	// Should filter out navigation
	assert.NotContains(t, result, "Home")
	assert.NotContains(t, result, "Blog")
}

func TestSemanticExtractor_NoH1OnlyH2(t *testing.T) {
	// Test case: No h1, only h2s - should extract h2s
	htmlContent, err := os.ReadFile("../../testdata/semantic_no_h1_only_h2.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should contain all h2s and their content
	assert.Contains(t, result, "First Major Section")
	assert.Contains(t, result, "Second Major Section")
	assert.Contains(t, result, "Third Major Section")

	// Should include content
	assert.Contains(t, result, "Content for first section")
	assert.Contains(t, result, "Content for second section")

	// Should maintain document order
	firstPos := strings.Index(result, "First Major Section")
	secondPos := strings.Index(result, "Second Major Section")
	thirdPos := strings.Index(result, "Third Major Section")

	assert.Less(t, firstPos, secondPos, "First H2 should come before second")
	assert.Less(t, secondPos, thirdPos, "Second H2 should come before third")

	// Should include introduction paragraph
	assert.Contains(t, result, "Introduction paragraph with no heading")
}

func TestSemanticExtractor_H1InNavigation(t *testing.T) {
	// Test case: h1 inside navigation should be filtered
	htmlContent, err := os.ReadFile("../../testdata/semantic_h1_in_navigation.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should contain the real article h1
	assert.Contains(t, result, "Real Article Title")
	assert.Contains(t, result, "This is the main article content")

	// Should NOT contain navigation h1s
	assert.NotContains(t, result, "Site Navigation")
	assert.NotContains(t, result, "Sidebar Navigation")

	// Should filter out navigation links
	assert.NotContains(t, result, "Home")
	assert.NotContains(t, result, "About")
	assert.NotContains(t, result, "Category 1")
}

func TestSemanticExtractor_MalformedHTML(t *testing.T) {
	// Test case: Malformed HTML should still work (Chrome fixes it)
	htmlContent, err := os.ReadFile("../../testdata/semantic_malformed_html.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	// Chrome should fix malformed HTML when building tree
	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should still extract h1s even from malformed HTML
	assert.Contains(t, result, "Main Title")
	assert.Contains(t, result, "Another Title")

	// Should handle h2s properly
	assert.Contains(t, result, "Section Title")
	assert.Contains(t, result, "Heading inside paragraph")

	// Should include content
	assert.Contains(t, result, "Paragraph after unclosed h1")
	assert.Contains(t, result, "Section content")
}

func TestSemanticExtractor_NoHeadings(t *testing.T) {
	// Test case: No headings at all - should extract remaining content
	htmlContent, err := os.ReadFile("../../testdata/semantic_no_headings.html")
	require.NoError(t, err)

	extractor := NewSemanticExtractor()
	tb := NewTreeBuilder().WithPreserveAttributes(true)
	ctx := context.Background()

	root, err := tb.BuildTree(ctx, string(htmlContent))
	require.NoError(t, err)

	result := extractor.Extract(root)

	// Should contain main content
	assert.Contains(t, result, "Bold introduction text")
	assert.Contains(t, result, "First paragraph of content")
	assert.Contains(t, result, "Second paragraph of content")
	assert.Contains(t, result, "A quote that should be included")
	assert.Contains(t, result, "List item 1")
	assert.Contains(t, result, "Final paragraph of content")

	// Should filter out navigation
	assert.NotContains(t, result, "Home")
	assert.NotContains(t, result, "About")

	// May or may not include sidebar and footer - depends on implementation
}

func TestSemanticExtractor_ExtractionOrder(t *testing.T) {
	// Test the specific extraction order algorithm
	extractor := NewSemanticExtractor()

	// This will be a unit test for the extraction order logic
	// We'll need to implement the algorithm first, then test it
	assert.NotNil(t, extractor, "Extractor should be created")
}
