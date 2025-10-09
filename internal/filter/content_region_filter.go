package filter

import (
	"github.com/jewell-lgtm/essenz/internal/tree"
)

// ContentRegionFilter preserves nodes identified as content regions and filters out low-value content.
type ContentRegionFilter struct{}

// NewContentRegionFilter creates a new ContentRegionFilter.
func NewContentRegionFilter() *ContentRegionFilter {
	return &ContentRegionFilter{}
}

// ShouldExclude determines if a node should be excluded based on content region analysis.
func (f *ContentRegionFilter) ShouldExclude(node *tree.TextNode, ctx *FilterContext) bool {
	// Never exclude nodes that are identified as content regions
	if node.IsContentRegion {
		return false
	}

	// Never exclude nodes with high semantic weight (important markup)
	if node.SemanticWeight >= 2.0 {
		return false
	}

	// For container elements (not #text), check if they contain content regions
	if node.Tag != "#text" {
		// If this container has child content regions, preserve it
		if f.hasContentRegionChildren(node) {
			return false
		}

		// If it has substantial word count and good content density, preserve it
		if node.WordCount >= 20 && node.ContentDensity > 5.0 && node.LinkDensity < 0.4 {
			return false
		}
	}

	// For text nodes, check if parent or ancestors are content regions
	if node.Tag == "#text" {
		// Check if we're within a content region (check up to 3 levels up)
		if f.isWithinContentRegion(node, 3) {
			return false
		}

		// Check semantic weight of text content
		if node.SemanticWeight >= 1.5 {
			return false
		}
	}

	// Apply stricter filtering to nodes outside content regions
	// These are likely navigation, ads, or other non-content elements

	// Filter out very short content with high link density
	if node.Tag != "#text" && node.WordCount < 10 && node.LinkDensity > 0.5 {
		return true
	}

	// Filter out nodes with very low semantic weight and content density
	if node.SemanticWeight < 0.5 && node.ContentDensity < 2.0 && node.WordCount < 5 {
		return true
	}

	// Default: preserve (let other filters decide)
	return false
}

// Priority returns the priority of this filter (higher = runs first).
func (f *ContentRegionFilter) Priority() int {
	return 100 // Run first - content regions should override other filters
}

// Name returns the name of this filter.
func (f *ContentRegionFilter) Name() string {
	return "ContentRegionFilter"
}

// hasContentRegionChildren checks if any children are marked as content regions.
func (f *ContentRegionFilter) hasContentRegionChildren(node *tree.TextNode) bool {
	for _, child := range node.Children {
		if child.IsContentRegion {
			return true
		}
		// Check recursively (up to reasonable depth to avoid performance issues)
		if len(child.Children) > 0 && f.hasContentRegionChildren(child) {
			return true
		}
	}
	return false
}

// isWithinContentRegion checks if a node is within a content region by traversing up the tree.
func (f *ContentRegionFilter) isWithinContentRegion(node *tree.TextNode, maxLevels int) bool {
	current := node.Parent
	levels := 0

	for current != nil && levels < maxLevels {
		if current.IsContentRegion {
			return true
		}
		current = current.Parent
		levels++
	}

	return false
}
