package filter

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// ClassNameFilter removes content based on CSS class and ID patterns.
type ClassNameFilter struct {
	strongPatterns []string
	weakNav        []string
}

// NewClassNameFilter creates a new ClassNameFilter.
func NewClassNameFilter() *ClassNameFilter {
	return &ClassNameFilter{
		// Strong negative signals: remove regardless of link density
		strongPatterns: []string{
			"sidebar", "aside", "header", "footer",
			"ad", "ads", "advertisement", "sponsored", "promo",
			"social", "share", "sharing", "social-share", "social-media",
			"comment", "comments", "comment-section", "disqus",
			"related", "related-posts", "related-links", "you-might-like", "similar", "related-content",
			"breadcrumb", "breadcrumbs", "pagination", "pager",
			// Trivial/low-value content containers
			"short", "short-content",
			"skip", "sr-only", "screen-reader", "hidden", "invisible",
		},
		// Weak nav-indicator tokens: only remove if link-heavy
		weakNav: []string{"nav", "menu", "navigation", "navbar", "nav-menu", "navigation-menu"},
	}
}

// ShouldExclude determines if a node should be excluded based on class/ID patterns.
func (f *ClassNameFilter) ShouldExclude(node *tree.TextNode, _ *FilterContext) bool {
	if node == nil {
		return false
	}

	// Check class attribute
	if classValue, exists := node.Attributes["class"]; exists {
		lower := strings.ToLower(classValue)
		if f.matchesAny(lower, f.strongPatterns) {
			return true
		}
		if f.matchesAny(lower, f.weakNav) {
			// Only exclude nav-like classes when content is link-heavy or too short
			linkChars, totalChars := analyzeLinkVsText(node)
			if totalChars == 0 {
				return true
			}
			density := float64(linkChars) / float64(totalChars)
			if density > 0.35 || totalChars < 80 {
				return true
			}
		}
	}

	// Check id attribute
	if idValue, exists := node.Attributes["id"]; exists {
		lower := strings.ToLower(idValue)
		if f.matchesAny(lower, f.strongPatterns) || f.matchesAny(lower, f.weakNav) {
			return true
		}
	}

	return false
}

// matchesAny checks if a value matches any of the provided patterns with word boundaries.
func (f *ClassNameFilter) matchesAny(value string, patterns []string) bool {
	for _, pattern := range patterns {
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

// analyzeLinkVsText computes link and total text characters for a node subtree.
func analyzeLinkVsText(node *tree.TextNode) (linkChars, totalChars int) {
	var walk func(n *tree.TextNode, inLink bool)
	walk = func(n *tree.TextNode, inLink bool) {
		if n == nil {
			return
		}
		if n.Tag == "#text" {
			t := strings.TrimSpace(n.Text)
			l := len(t)
			totalChars += l
			if inLink {
				linkChars += l
			}
			return
		}
		isLink := strings.ToLower(n.Tag) == "a"
		for _, c := range n.Children {
			walk(c, inLink || isLink)
		}
	}
	walk(node, false)
	return
}
