package filter

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// ClassNameFilter removes content based on CSS class and ID patterns.
type ClassNameFilter struct {
	excludePatterns []string
}

// NewClassNameFilter creates a new ClassNameFilter.
func NewClassNameFilter() *ClassNameFilter {
	return &ClassNameFilter{
		excludePatterns: []string{
			// Navigation patterns
			"nav", "menu", "navigation", "navbar", "nav-menu", "navigation-menu",

			// Layout patterns
			"sidebar", "aside", "header", "footer",

			// Advertisement patterns
			"ad", "ads", "advertisement", "sponsored", "promo",

			// Social patterns
			"social", "share", "sharing", "social-share", "social-media",

			// Comment patterns
			"comment", "comments", "comment-section", "disqus",

			// Related content patterns
			"related", "related-posts", "related-links", "you-might-like", "similar", "related-content",

			// Navigation aid patterns
			"breadcrumb", "breadcrumbs", "pagination", "pager",

			// Utility patterns
			"skip", "sr-only", "screen-reader", "hidden", "invisible",
		},
	}
}

// ShouldExclude determines if a node should be excluded based on class/ID patterns.
func (f *ClassNameFilter) ShouldExclude(node *tree.TextNode, _ *FilterContext) bool {
	if node == nil {
		return false
	}

	// Check class attribute
	if classValue, exists := node.Attributes["class"]; exists {
		if f.matchesPattern(strings.ToLower(classValue)) {
			return true
		}
	}

	// Check id attribute
	if idValue, exists := node.Attributes["id"]; exists {
		if f.matchesPattern(strings.ToLower(idValue)) {
			return true
		}
	}

	return false
}

// matchesPattern checks if a value matches any of the exclude patterns.
func (f *ClassNameFilter) matchesPattern(value string) bool {
	for _, pattern := range f.excludePatterns {
		// Check for exact word match or pattern within CSS class names
		if strings.Contains(value, pattern) {
			// Additional check: ensure it's a word boundary to avoid false positives
			// For example, "content" shouldn't match "advertisement"
			if f.isWordBoundaryMatch(value, pattern) {
				return true
			}
		}
	}
	return false
}

// isWordBoundaryMatch checks if pattern appears as a complete word or CSS class.
func (f *ClassNameFilter) isWordBoundaryMatch(value, pattern string) bool {
	// Exact match
	if value == pattern {
		return true
	}

	// Split by common CSS class separators and check each part
	separators := []string{" ", "-", "_", "."}
	for _, sep := range separators {
		parts := strings.Split(value, sep)
		for _, part := range parts {
			if part == pattern {
				return true
			}
		}
	}

	// For patterns like "related" matching "related-posts", be more flexible
	if strings.HasPrefix(value, pattern+"-") ||
		strings.HasPrefix(value, pattern+"_") ||
		strings.HasSuffix(value, "-"+pattern) ||
		strings.HasSuffix(value, "_"+pattern) {
		return true
	}

	return false
}

// Priority returns the priority of this filter rule.
func (f *ClassNameFilter) Priority() int {
	return 80 // High priority - CSS patterns are strong indicators
}

// Name returns the name of this filter rule.
func (f *ClassNameFilter) Name() string {
	return "ClassNameFilter"
}
