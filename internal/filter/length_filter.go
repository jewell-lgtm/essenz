package filter

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// LengthFilter removes very short text blocks that are likely not meaningful content.
type LengthFilter struct {
	minLength int
}

// NewLengthFilter creates a new LengthFilter.
func NewLengthFilter(minLength int) *LengthFilter {
	return &LengthFilter{
		minLength: minLength,
	}
}

// ShouldExclude determines if a node should be excluded based on content length.
func (f *LengthFilter) ShouldExclude(node *tree.TextNode, _ *FilterContext) bool {
	if node == nil || node.Tag == "#text" {
		return false
	}

	// Don't filter structural elements that might contain important short content
	if f.isStructuralElement(node) {
		return false
	}

	// Don't filter text nodes directly - only container elements
	if node.Tag == "#text" {
		return false
	}

	// Calculate total text content length for this node
	totalText := f.extractAllText(node)
	totalLength := len(strings.TrimSpace(totalText))

	// Exclude if content is too short
	if totalLength < f.minLength {
		// Exception: preserve if it contains important structural elements
		if f.hasImportantChildren(node) {
			return false
		}
		return true
	}

	return false
}

// isStructuralElement checks if the node is a structural element that should be preserved.
func (f *LengthFilter) isStructuralElement(node *tree.TextNode) bool {
	if node == nil {
		return false
	}

	structuralTags := map[string]bool{
		"main":    true,
		"article": true,
		"section": true,
		"div":     true, // Don't filter divs directly, check their content instead
		"header":  true, // These might be filtered by other rules, but not by length
		"footer":  true,
		"nav":     true,
		"aside":   true,
		"h1":      true,
		"h2":      true,
		"h3":      true,
		"h4":      true,
		"h5":      true,
		"h6":      true,
	}

	tagName := strings.ToLower(node.Tag)
	return structuralTags[tagName]
}

// hasImportantChildren checks if a node has children that indicate importance.
func (f *LengthFilter) hasImportantChildren(node *tree.TextNode) bool {
	if node == nil {
		return false
	}

	for _, child := range node.Children {
		tagName := strings.ToLower(child.Tag)

		// Check for important child elements
		switch tagName {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			return true // Headings are important
		case "img":
			return true // Images might be important content
		case "table":
			return true // Tables contain structured data
		case "ul", "ol", "dl":
			return true // Lists contain structured information
		case "blockquote":
			return true // Quotes are usually important
		case "code", "pre":
			return true // Code blocks are important
		}

		// Check for strong semantic indicators in attributes
		if classValue, exists := child.Attributes["class"]; exists {
			classLower := strings.ToLower(classValue)
			if strings.Contains(classLower, "important") ||
				strings.Contains(classLower, "highlight") ||
				strings.Contains(classLower, "note") ||
				strings.Contains(classLower, "warning") ||
				strings.Contains(classLower, "alert") {
				return true
			}
		}

		// Recurse into children
		if f.hasImportantChildren(child) {
			return true
		}
	}

	return false
}

// extractAllText recursively extracts all text content from a node.
func (f *LengthFilter) extractAllText(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return node.Text
	}

	var textParts []string
	for _, child := range node.Children {
		childText := f.extractAllText(child)
		if childText != "" {
			textParts = append(textParts, childText)
		}
	}

	return strings.Join(textParts, " ")
}

// Priority returns the priority of this filter rule.
func (f *LengthFilter) Priority() int {
	return 40 // Lower priority - length filtering should be done after more specific rules
}

// Name returns the name of this filter rule.
func (f *LengthFilter) Name() string {
	return "LengthFilter"
}
