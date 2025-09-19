// Package tree provides bottom-up text node tree construction for content analysis.
package tree

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

// TextNode represents a single text node in the hierarchical structure.
type TextNode struct {
	Text       string            `json:"text"`
	Tag        string            `json:"tag"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Parent     *TextNode         `json:"-"` // Exclude from JSON to avoid cycles
	Children   []*TextNode       `json:"children,omitempty"`
	Depth      int               `json:"depth"`
	Index      int               `json:"index"`
}

// TreeBuilder constructs hierarchical text node structures from HTML documents.
type TreeBuilder struct {
	filterNavigation   bool
	preserveAttributes bool
	includeWhitespace  bool
	maxDepth           int
	navigationTags     map[string]bool
}

// NewTreeBuilder creates a new TreeBuilder with default configuration.
func NewTreeBuilder() *TreeBuilder {
	return &TreeBuilder{
		filterNavigation:   false,
		preserveAttributes: false,
		includeWhitespace:  false,
		maxDepth:           100,
		navigationTags: map[string]bool{
			"nav":      true,
			"footer":   true,
			"header":   true,
			"aside":    true,
			"menu":     true,
			"script":   true,
			"style":    true,
			"noscript": true,
		},
	}
}

// WithFilterNavigation enables filtering of navigation elements.
func (tb *TreeBuilder) WithFilterNavigation(filter bool) *TreeBuilder {
	tb.filterNavigation = filter
	return tb
}

// WithPreserveAttributes enables preservation of element attributes.
func (tb *TreeBuilder) WithPreserveAttributes(preserve bool) *TreeBuilder {
	tb.preserveAttributes = preserve
	return tb
}

// WithIncludeWhitespace controls whether whitespace-only text nodes are included.
func (tb *TreeBuilder) WithIncludeWhitespace(include bool) *TreeBuilder {
	tb.includeWhitespace = include
	return tb
}

// WithMaxDepth sets the maximum depth for tree traversal.
func (tb *TreeBuilder) WithMaxDepth(depth int) *TreeBuilder {
	tb.maxDepth = depth
	return tb
}

// BuildTree constructs a text node tree from HTML content.
func (tb *TreeBuilder) BuildTree(ctx context.Context, htmlContent string) (*TextNode, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	root := &TextNode{
		Tag:        "document",
		Attributes: make(map[string]string),
		Children:   make([]*TextNode, 0),
		Depth:      0,
		Index:      0,
	}

	// Process all child nodes of the document
	currentIndex := 1
	for child := doc.FirstChild; child != nil; child = child.NextSibling {
		currentIndex = tb.traverseNode(ctx, child, root, 1, currentIndex)
	}

	return root, nil
}

// traverseNode recursively processes HTML nodes to build the text node tree.
func (tb *TreeBuilder) traverseNode(ctx context.Context, node *html.Node, parent *TextNode, depth, index int) int {
	if depth > tb.maxDepth {
		return index
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return index
	default:
	}

	currentIndex := index

	switch node.Type {
	case html.ElementNode:
		tagName := strings.ToLower(node.Data)

		// Always skip script, style, and noscript tags
		if tagName == "script" || tagName == "style" || tagName == "noscript" {
			return currentIndex
		}

		// Skip navigation elements if filtering is enabled
		if tb.filterNavigation && tb.navigationTags[tagName] {
			return currentIndex
		}

		// Check for hidden content when filtering is enabled
		if tb.filterNavigation {
			for _, attr := range node.Attr {
				// Skip elements with hidden/invisible classes
				if attr.Key == "class" && (strings.Contains(attr.Val, "hidden") ||
					strings.Contains(attr.Val, "invisible") ||
					strings.Contains(attr.Val, "sr-only")) {
					return currentIndex
				}
				// Skip elements with display:none or visibility:hidden
				if attr.Key == "style" && (strings.Contains(attr.Val, "display:none") ||
					strings.Contains(attr.Val, "display: none") ||
					strings.Contains(attr.Val, "visibility:hidden") ||
					strings.Contains(attr.Val, "visibility: hidden")) {
					return currentIndex
				}
			}
		}

		// Create element node
		elementNode := &TextNode{
			Tag:        node.Data,
			Attributes: make(map[string]string),
			Children:   make([]*TextNode, 0),
			Parent:     parent,
			Depth:      depth,
			Index:      currentIndex,
		}

		// Preserve attributes if enabled
		if tb.preserveAttributes {
			for _, attr := range node.Attr {
				elementNode.Attributes[attr.Key] = attr.Val
			}
		}

		parent.Children = append(parent.Children, elementNode)
		currentIndex++

		// Process child nodes
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			currentIndex = tb.traverseNode(ctx, child, elementNode, depth+1, currentIndex)
		}

	case html.TextNode:
		text := strings.TrimSpace(node.Data)

		// Skip empty text nodes unless whitespace is explicitly included
		if text == "" && !tb.includeWhitespace {
			return currentIndex
		}

		// Create text node
		textNode := &TextNode{
			Text:       node.Data, // Keep original text including whitespace
			Tag:        "#text",
			Attributes: make(map[string]string),
			Children:   make([]*TextNode, 0),
			Parent:     parent,
			Depth:      depth,
			Index:      currentIndex,
		}

		parent.Children = append(parent.Children, textNode)
		currentIndex++
	}

	return currentIndex
}

// GetTextNodes returns all text nodes from the tree structure.
func (tb *TreeBuilder) GetTextNodes(root *TextNode) []*TextNode {
	var textNodes []*TextNode
	tb.collectTextNodes(root, &textNodes)
	return textNodes
}

// collectTextNodes recursively collects all text nodes.
func (tb *TreeBuilder) collectTextNodes(node *TextNode, textNodes *[]*TextNode) {
	if node.Tag == "#text" && strings.TrimSpace(node.Text) != "" {
		*textNodes = append(*textNodes, node)
	}

	for _, child := range node.Children {
		tb.collectTextNodes(child, textNodes)
	}
}

// GetStats returns statistics about the tree structure.
func (tb *TreeBuilder) GetStats(root *TextNode) map[string]interface{} {
	stats := map[string]interface{}{
		"total_nodes":     0,
		"text_nodes":      0,
		"element_nodes":   0,
		"max_depth":       0,
		"text_characters": 0,
	}

	tb.calculateStats(root, stats, 0)
	return stats
}

// calculateStats recursively calculates tree statistics.
func (tb *TreeBuilder) calculateStats(node *TextNode, stats map[string]interface{}, depth int) {
	stats["total_nodes"] = stats["total_nodes"].(int) + 1

	if depth > stats["max_depth"].(int) {
		stats["max_depth"] = depth
	}

	if node.Tag == "#text" {
		stats["text_nodes"] = stats["text_nodes"].(int) + 1
		stats["text_characters"] = stats["text_characters"].(int) + len(strings.TrimSpace(node.Text))
	} else {
		stats["element_nodes"] = stats["element_nodes"].(int) + 1
	}

	for _, child := range node.Children {
		tb.calculateStats(child, stats, depth+1)
	}
}

// ToJSON converts the tree structure to JSON format.
func (tb *TreeBuilder) ToJSON(root *TextNode) (string, error) {
	data, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal tree to JSON: %w", err)
	}
	return string(data), nil
}

// ToText converts the tree structure to a readable text format.
func (tb *TreeBuilder) ToText(root *TextNode) string {
	var builder strings.Builder
	tb.writeTextNode(&builder, root, "")
	return builder.String()
}

// writeTextNode recursively writes nodes in text format.
func (tb *TreeBuilder) writeTextNode(builder *strings.Builder, node *TextNode, indent string) {
	if node == nil {
		return // Skip nil nodes gracefully
	}
	if node.Tag == "#text" {
		text := strings.TrimSpace(node.Text)
		if text != "" {
			builder.WriteString(fmt.Sprintf("%s[%d] %s: \"%s\"\n",
				indent, node.Index, node.Tag, text))
		}
	} else {
		attrs := ""
		if len(node.Attributes) > 0 && tb.preserveAttributes {
			var attrPairs []string
			for k, v := range node.Attributes {
				attrPairs = append(attrPairs, fmt.Sprintf("%s=\"%s\"", k, v))
			}
			attrs = fmt.Sprintf(" (%s)", strings.Join(attrPairs, ", "))
		}

		builder.WriteString(fmt.Sprintf("%s[%d] %s%s\n",
			indent, node.Index, node.Tag, attrs))
	}

	for _, child := range node.Children {
		tb.writeTextNode(builder, child, indent+"  ")
	}
}
