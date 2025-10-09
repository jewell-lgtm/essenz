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

	// Content analysis fields
	ContentDensity  float64 `json:"content_density"`   // Ratio of text content to total content
	SemanticWeight  float64 `json:"semantic_weight"`   // Importance based on semantic markup
	WordCount       int     `json:"word_count"`        // Number of words in this node
	LinkDensity     float64 `json:"link_density"`      // Ratio of link text to total text
	IsContentRegion bool    `json:"is_content_region"` // Heuristic: likely main content
	ContentScore    float64 `json:"content_score"`     // BFS-calculated content importance score
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

	// Perform content analysis on the completed tree
	tb.analyzeContent(root)

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
		return tb.processElementNode(ctx, node, parent, depth, currentIndex)
	case html.TextNode:
		return tb.processTextNode(node, parent, depth, currentIndex)
	}

	return currentIndex
}

// processElementNode processes an HTML element node.
func (tb *TreeBuilder) processElementNode(ctx context.Context, node *html.Node, parent *TextNode, depth, index int) int {
	tagName := strings.ToLower(node.Data)

	// Always skip script, style, and noscript tags
	if tagName == "script" || tagName == "style" || tagName == "noscript" {
		return index
	}

	// Skip navigation elements if filtering is enabled
	if tb.filterNavigation && tb.navigationTags[tagName] {
		return index
	}

	// Check for hidden content when filtering is enabled
	if tb.filterNavigation && tb.isHiddenElement(node) {
		return index
	}

	// Create element node
	elementNode := &TextNode{
		Tag:        node.Data,
		Attributes: make(map[string]string),
		Children:   make([]*TextNode, 0),
		Parent:     parent,
		Depth:      depth,
		Index:      index,
	}

	// Preserve attributes if enabled
	if tb.preserveAttributes {
		for _, attr := range node.Attr {
			elementNode.Attributes[attr.Key] = attr.Val
		}
	}

	parent.Children = append(parent.Children, elementNode)
	currentIndex := index + 1

	// Process child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		currentIndex = tb.traverseNode(ctx, child, elementNode, depth+1, currentIndex)
	}

	return currentIndex
}

// processTextNode processes an HTML text node.
func (tb *TreeBuilder) processTextNode(node *html.Node, parent *TextNode, depth, index int) int {
	text := strings.TrimSpace(node.Data)

	// Skip empty text nodes unless whitespace is explicitly included
	if text == "" && !tb.includeWhitespace {
		return index
	}

	// Create text node
	textNode := &TextNode{
		Text:       node.Data, // Keep original text including whitespace
		Tag:        "#text",
		Attributes: make(map[string]string),
		Children:   make([]*TextNode, 0),
		Parent:     parent,
		Depth:      depth,
		Index:      index,
	}

	parent.Children = append(parent.Children, textNode)
	return index + 1
}

// isHiddenElement checks if an element should be considered hidden.
func (tb *TreeBuilder) isHiddenElement(node *html.Node) bool {
	for _, attr := range node.Attr {
		// Skip elements with hidden/invisible classes
		if attr.Key == "class" && (strings.Contains(attr.Val, "hidden") ||
			strings.Contains(attr.Val, "invisible") ||
			strings.Contains(attr.Val, "sr-only")) {
			return true
		}
		// Skip elements with display:none or visibility:hidden
		if attr.Key == "style" && (strings.Contains(attr.Val, "display:none") ||
			strings.Contains(attr.Val, "display: none") ||
			strings.Contains(attr.Val, "visibility:hidden") ||
			strings.Contains(attr.Val, "visibility: hidden")) {
			return true
		}
	}
	return false
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

// analyzeContent performs content density and semantic analysis on the tree.
func (tb *TreeBuilder) analyzeContent(root *TextNode) {
	// First pass: calculate basic metrics
	tb.calculateNodeMetrics(root)

	// Second pass: identify content regions using heuristics
	tb.identifyContentRegions(root)
}

// calculateNodeMetrics calculates word count, link density, and content density for each node.
func (tb *TreeBuilder) calculateNodeMetrics(node *TextNode) {
	if node == nil {
		return
	}

	// Calculate for this node
	if node.Tag == "#text" {
		node.WordCount = len(strings.Fields(strings.TrimSpace(node.Text)))
	} else {
		// For element nodes, calculate aggregate metrics from children
		totalWords := 0
		totalLinkWords := 0
		totalChildren := 0

		for _, child := range node.Children {
			tb.calculateNodeMetrics(child) // Recursive call first
			totalWords += child.WordCount
			totalChildren++

			// Count link words if this child is within a link
			if tb.isWithinLink(child) {
				totalLinkWords += child.WordCount
			}
		}

		node.WordCount = totalWords
		if totalWords > 0 {
			node.LinkDensity = float64(totalLinkWords) / float64(totalWords)
		}

		// Content density: ratio of text content to structural elements
		if totalChildren > 0 {
			node.ContentDensity = float64(totalWords) / float64(totalChildren+1)
		}
	}

	// Calculate semantic weight based on tag and attributes
	node.SemanticWeight = tb.calculateSemanticWeight(node)
}

// identifyContentRegions uses heuristics to identify likely content regions.
func (tb *TreeBuilder) identifyContentRegions(node *TextNode) {
	if node == nil {
		return
	}

	// Apply content region heuristics
	node.IsContentRegion = tb.isLikelyContentRegion(node)

	// Recurse to children
	for _, child := range node.Children {
		tb.identifyContentRegions(child)
	}
}

// isLikelyContentRegion determines if a node represents a content-heavy region.
func (tb *TreeBuilder) isLikelyContentRegion(node *TextNode) bool {
	// Skip text nodes - only evaluate element containers
	if node.Tag == "#text" {
		return false
	}

	// High word count suggests content
	if node.WordCount >= 50 {
		// Low link density suggests article content rather than navigation
		if node.LinkDensity < 0.3 {
			return true
		}
	}

	// Semantic tags that commonly contain main content
	contentTags := map[string]bool{
		"article": true, "main": true, "section": true,
		"div": true, // div is common but needs other signals
	}

	if contentTags[node.Tag] {
		// For semantic tags, lower word count threshold
		if node.WordCount >= 20 && node.LinkDensity < 0.4 {
			return true
		}
	}

	// Check for content-indicating class names or IDs
	if tb.hasContentIndicators(node) {
		if node.WordCount >= 15 && node.LinkDensity < 0.5 {
			return true
		}
	}

	// High content density can indicate content regions
	if node.ContentDensity > 10.0 && node.WordCount >= 30 {
		return true
	}

	return false
}

// calculateSemanticWeight assigns importance based on semantic markup.
func (tb *TreeBuilder) calculateSemanticWeight(node *TextNode) float64 {
	weight := 1.0

	// Semantic HTML5 tags get higher weight
	semanticWeights := map[string]float64{
		"h1": 3.0, "h2": 2.5, "h3": 2.0, "h4": 1.5, "h5": 1.2, "h6": 1.1,
		"article": 2.5, "main": 3.0, "section": 2.0,
		"p": 1.0, "blockquote": 1.8, "pre": 1.5, "code": 1.3,
		"strong": 1.4, "em": 1.2, "mark": 1.6,
		"figcaption": 1.3, "caption": 1.3,
		"nav": 0.3, "aside": 0.4, "footer": 0.2, "header": 0.5,
	}

	if w, exists := semanticWeights[node.Tag]; exists {
		weight = w
	}

	// Boost weight for content-indicating attributes
	if tb.hasContentIndicators(node) {
		weight *= 1.5
	}

	// Reduce weight for navigation-like attributes
	if tb.hasNavigationIndicators(node) {
		weight *= 0.3
	}

	return weight
}

// hasContentIndicators checks for class/id names that suggest main content.
func (tb *TreeBuilder) hasContentIndicators(node *TextNode) bool {
	contentIndicators := []string{
		"content", "article", "post", "entry", "main", "primary",
		"readme", "documentation", "description", "summary",
		"body", "text", "prose", "story", "news",
	}

	return tb.hasAttributePatterns(node, contentIndicators)
}

// hasNavigationIndicators checks for class/id names that suggest navigation.
func (tb *TreeBuilder) hasNavigationIndicators(node *TextNode) bool {
	navIndicators := []string{
		"nav", "menu", "sidebar", "aside", "header", "footer",
		"breadcrumb", "pagination", "social", "share", "related",
		"comments", "meta", "tags", "categories", "toolbar",
	}

	return tb.hasAttributePatterns(node, navIndicators)
}

// hasAttributePatterns checks if node attributes contain any of the given patterns.
func (tb *TreeBuilder) hasAttributePatterns(node *TextNode, patterns []string) bool {
	for _, attr := range []string{"class", "id", "role"} {
		if value, exists := node.Attributes[attr]; exists {
			lowerValue := strings.ToLower(value)
			for _, pattern := range patterns {
				if strings.Contains(lowerValue, pattern) {
					return true
				}
			}
		}
	}
	return false
}

// isWithinLink checks if a node is within a link element by traversing up the tree.
func (tb *TreeBuilder) isWithinLink(node *TextNode) bool {
	current := node.Parent
	for current != nil {
		if current.Tag == "a" {
			return true
		}
		current = current.Parent
	}
	return false
}
