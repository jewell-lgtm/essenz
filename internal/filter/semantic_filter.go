package filter

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// SemanticTagFilter removes content based on HTML5 semantic tags.
type SemanticTagFilter struct {
	excludedTags map[string]bool
}

// NewSemanticTagFilter creates a new SemanticTagFilter.
func NewSemanticTagFilter() *SemanticTagFilter {
	return &SemanticTagFilter{
		excludedTags: map[string]bool{
			"nav":      true,
			"header":   true,
			"footer":   true,
			"aside":    true,
			"script":   true,
			"style":    true,
			"noscript": true,
		},
	}
}

// ShouldExclude determines if a node should be excluded based on semantic tags.
func (f *SemanticTagFilter) ShouldExclude(node *tree.TextNode, _ *FilterContext) bool {
	if node == nil {
		return false
	}

	tagName := strings.ToLower(node.Tag)
	return f.excludedTags[tagName]
}

// Priority returns the priority of this filter rule.
func (f *SemanticTagFilter) Priority() int {
	return 100 // High priority - semantic tags are clear indicators
}

// Name returns the name of this filter rule.
func (f *SemanticTagFilter) Name() string {
	return "SemanticTagFilter"
}
