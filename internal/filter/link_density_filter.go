package filter

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// LinkDensityFilter removes sections with high link-to-text ratios.
type LinkDensityFilter struct {
	maxDensity float64
	minWords   int
}

// NewLinkDensityFilter creates a new LinkDensityFilter.
func NewLinkDensityFilter(maxDensity float64, minWords int) *LinkDensityFilter {
	return &LinkDensityFilter{
		maxDensity: maxDensity,
		minWords:   minWords,
	}
}

// ShouldExclude determines if a node should be excluded based on link density.
func (f *LinkDensityFilter) ShouldExclude(node *tree.TextNode, _ *FilterContext) bool {
	if node == nil || node.Tag == "#text" {
		return false
	}

	// Never filter structural elements
	structuralTags := map[string]bool{
		"document": true,
		"html":     true,
		"head":     true,
		"body":     true,
		"main":     true,
		"article":  true,
		"section":  true,
	}
	if structuralTags[strings.ToLower(node.Tag)] {
		return false
	}

	// Calculate link density for this node and its immediate children
	linkChars, totalChars, wordCount := f.calculateNodeStats(node)

	// Don't filter if there's insufficient content to analyze
	if wordCount < f.minWords || totalChars < 50 {
		return false
	}

	// Calculate link density
	var density float64
	if totalChars > 0 {
		density = float64(linkChars) / float64(totalChars)
	}

	// Exclude if link density is too high
	return density > f.maxDensity
}

// calculateNodeStats calculates link characters, total characters, and word count.
func (f *LinkDensityFilter) calculateNodeStats(node *tree.TextNode) (linkChars, totalChars, wordCount int) {
	f.collectNodeStats(node, &linkChars, &totalChars, &wordCount, false)
	return
}

// collectNodeStats recursively collects statistics from a node and its children.
func (f *LinkDensityFilter) collectNodeStats(node *tree.TextNode, linkChars, totalChars, wordCount *int, inLink bool) {
	if node == nil {
		return
	}

	if node.Tag == "#text" {
		text := strings.TrimSpace(node.Text)
		textLen := len(text)
		words := len(strings.Fields(text))

		*totalChars += textLen
		*wordCount += words

		if inLink {
			*linkChars += textLen
		}
		return
	}

	// Check if this is a link element
	isLink := strings.ToLower(node.Tag) == "a"

	// Process children
	for _, child := range node.Children {
		f.collectNodeStats(child, linkChars, totalChars, wordCount, inLink || isLink)
	}
}

// Priority returns the priority of this filter rule.
func (f *LinkDensityFilter) Priority() int {
	return 60 // Medium priority - content analysis is important but not definitive
}

// Name returns the name of this filter rule.
func (f *LinkDensityFilter) Name() string {
	return "LinkDensityFilter"
}
