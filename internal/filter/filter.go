// Package filter provides sophisticated content filtering to remove non-content elements.
package filter

import (
	"context"
	"fmt"
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// ContentFilter provides sophisticated filtering to remove non-content elements.
type ContentFilter struct {
	rules  []FilterRule
	config FilterConfig
}

// FilterConfig configures the content filtering behavior.
type FilterConfig struct {
	MaxLinkDensity    float64  // 0.3 = 30% links max
	MinContentLength  int      // Minimum characters for content blocks
	PreserveWhitelist []string // CSS selectors to always preserve
	AggressiveMode    bool     // More strict filtering
	DebugMode         bool     // Log filtering decisions
}

// FilterRule defines an interface for content filtering rules.
type FilterRule interface {
	ShouldExclude(node *tree.TextNode, context *FilterContext) bool
	Priority() int
	Name() string
}

// FilterContext provides context information for filtering decisions.
type FilterContext struct {
	DocumentRoot  *tree.TextNode
	CurrentDepth  int
	ParentNodes   []*tree.TextNode
	SiblingNodes  []*tree.TextNode
	DocumentStats *DocumentStats
}

// DocumentStats contains document-level statistics for filtering decisions.
type DocumentStats struct {
	TotalTextLength  int
	AverageWordCount float64
	LinkDensity      float64
	HeadingCount     int
	ParagraphCount   int
}

// FilterStats contains statistics about the filtering process.
type FilterStats struct {
	NodesProcessed int
	NodesRemoved   int
	RulesApplied   map[string]int
}

// NewContentFilter creates a new ContentFilter with default configuration.
func NewContentFilter() *ContentFilter {
	filter := &ContentFilter{
		rules: make([]FilterRule, 0),
		config: FilterConfig{
			MaxLinkDensity:    0.2, // More aggressive - 20% instead of 30%
			MinContentLength:  20,  // Reduce from 50 to 20 to be less aggressive
			PreserveWhitelist: []string{"main", "article", ".content", ".post", ".entry", ".main-article", ".main-content"},
			AggressiveMode:    false,
			DebugMode:         false,
		},
	}

	// Add default filter rules
	filter.AddRule(NewSemanticTagFilter())
	filter.AddRule(NewClassNameFilter())
	filter.AddRule(NewLinkDensityFilter(0.3, 5)) // Balanced: 30% max link density, 5 min words
	filter.AddRule(NewLengthFilter(10))          // Very low threshold but won't affect whitelist

	return filter
}

// WithConfig sets the filter configuration.
func (cf *ContentFilter) WithConfig(config FilterConfig) *ContentFilter {
	cf.config = config
	return cf
}

// WithAggressiveMode enables aggressive filtering.
func (cf *ContentFilter) WithAggressiveMode(aggressive bool) *ContentFilter {
	cf.config.AggressiveMode = aggressive
	return cf
}

// WithDebugMode enables debug logging.
func (cf *ContentFilter) WithDebugMode(debug bool) *ContentFilter {
	cf.config.DebugMode = debug
	return cf
}

// WithPreserveSelector adds a CSS selector to the whitelist.
func (cf *ContentFilter) WithPreserveSelector(selector string) *ContentFilter {
	cf.config.PreserveWhitelist = append(cf.config.PreserveWhitelist, selector)
	return cf
}

// AddRule adds a new filtering rule.
func (cf *ContentFilter) AddRule(rule FilterRule) {
	cf.rules = append(cf.rules, rule)
}

// FilterTree applies content filtering to a text node tree.
func (cf *ContentFilter) FilterTree(ctx context.Context, root *tree.TextNode) (*tree.TextNode, error) {
	if root == nil {
		return nil, fmt.Errorf("root node cannot be nil")
	}

	// Calculate document statistics
	stats := cf.calculateDocumentStats(root)

	// Create filter context
	filterCtx := &FilterContext{
		DocumentRoot:  root,
		CurrentDepth:  0,
		ParentNodes:   make([]*tree.TextNode, 0),
		SiblingNodes:  make([]*tree.TextNode, 0),
		DocumentStats: stats,
	}

	// Apply filtering recursively
	filtered := cf.filterNode(ctx, root, filterCtx)

	// Ensure we don't return a nil root
	if filtered == nil {
		// Return empty document root instead of nil
		filtered = &tree.TextNode{
			Tag:        "document",
			Text:       "",
			Attributes: make(map[string]string),
			Children:   make([]*tree.TextNode, 0),
			Parent:     nil,
			Index:      0,
		}
	}

	return filtered, nil
}

// filterNode recursively filters a node and its children.
func (cf *ContentFilter) filterNode(ctx context.Context, node *tree.TextNode, filterCtx *FilterContext) *tree.TextNode {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return node
	default:
	}

	if node == nil {
		return nil
	}

	// Check if node should be excluded by high-priority rules first (SemanticTagFilter, ClassNameFilter)
	// These rules override whitelist for strong negative indicators
	for _, rule := range cf.rules {
		if rule.Priority() >= 80 && rule.ShouldExclude(node, filterCtx) {
			if cf.config.DebugMode {
				fmt.Printf("DEBUG: Excluding node by high-priority rule %s: %s (class=%v)\n", rule.Name(), node.Tag, node.Attributes["class"])
			}
			return nil // Remove this node
		}
	}

	// Check whitelist protection for remaining rules
	isWhitelisted := cf.isWhitelisted(node)
	if !isWhitelisted {
		// Apply remaining lower-priority rules
		for _, rule := range cf.rules {
			if rule.Priority() < 80 && rule.ShouldExclude(node, filterCtx) {
				if cf.config.DebugMode {
					fmt.Printf("DEBUG: Excluding node by rule %s: %s (class=%v)\n", rule.Name(), node.Tag, node.Attributes["class"])
				}
				return nil // Remove this node
			}
		}
	} else {
		if cf.config.DebugMode {
			fmt.Printf("DEBUG: Preserving whitelisted node: %s\n", node.Tag)
		}
	}

	// Node passes all filters, process its children
	return cf.filterChildren(ctx, node, filterCtx)
}

// filterChildren filters the children of a node.
func (cf *ContentFilter) filterChildren(ctx context.Context, node *tree.TextNode, filterCtx *FilterContext) *tree.TextNode {
	if len(node.Children) == 0 {
		return node
	}

	// Create new context for children
	childCtx := &FilterContext{
		DocumentRoot:  filterCtx.DocumentRoot,
		CurrentDepth:  filterCtx.CurrentDepth + 1,
		ParentNodes:   append(filterCtx.ParentNodes, node),
		SiblingNodes:  node.Children, // Current children become siblings for the recursive call
		DocumentStats: filterCtx.DocumentStats,
	}

	// Filter children
	filteredChildren := make([]*tree.TextNode, 0)
	for _, child := range node.Children {
		filtered := cf.filterNode(ctx, child, childCtx)
		if filtered != nil {
			filtered.Parent = node // Maintain parent relationship
			filteredChildren = append(filteredChildren, filtered)
		}
	}

	// Update node's children
	node.Children = filteredChildren
	return node
}

// isWhitelisted checks if a node is in the whitelist.
func (cf *ContentFilter) isWhitelisted(node *tree.TextNode) bool {
	// Check tag-based whitelist
	for _, selector := range cf.config.PreserveWhitelist {
		if strings.HasPrefix(selector, ".") {
			// CSS class selector
			className := strings.TrimPrefix(selector, ".")
			if classValue, exists := node.Attributes["class"]; exists {
				if strings.Contains(classValue, className) {
					return true
				}
			}
		} else {
			// Tag selector
			if strings.EqualFold(node.Tag, selector) {
				return true
			}
		}
	}
	return false
}

// calculateDocumentStats calculates statistics about the document.
func (cf *ContentFilter) calculateDocumentStats(root *tree.TextNode) *DocumentStats {
	stats := &DocumentStats{}
	cf.collectStats(root, stats)

	// Calculate derived statistics
	if stats.ParagraphCount > 0 {
		stats.AverageWordCount = float64(stats.TotalTextLength) / float64(stats.ParagraphCount)
	}

	return stats
}

// collectStats recursively collects statistics from the tree.
func (cf *ContentFilter) collectStats(node *tree.TextNode, stats *DocumentStats) {
	if node == nil {
		return
	}

	if node.Tag == "#text" {
		stats.TotalTextLength += len(strings.TrimSpace(node.Text))
	} else {
		switch strings.ToLower(node.Tag) {
		case "h1", "h2", "h3", "h4", "h5", "h6":
			stats.HeadingCount++
		case "p":
			stats.ParagraphCount++
		case "a":
			// Count links for link density calculation
			// This is a simplified calculation - more sophisticated link density
			// calculation would be done per-section rather than document-wide
		}
	}

	// Recurse through children
	for _, child := range node.Children {
		cf.collectStats(child, stats)
	}
}

// GetFilterStats returns statistics about the last filtering operation.
func (cf *ContentFilter) GetFilterStats() *FilterStats {
	// This would be populated during filtering
	return &FilterStats{
		NodesProcessed: 0,
		NodesRemoved:   0,
		RulesApplied:   make(map[string]int),
	}
}
