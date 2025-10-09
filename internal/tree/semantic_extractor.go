package tree

import (
	"strings"
)

// SemanticExtractor implements semantic hierarchy-based content extraction
type SemanticExtractor struct {
	filterNavigation     bool
	filterPresentational bool
}

// NewSemanticExtractor creates a new semantic extractor with default settings
func NewSemanticExtractor() *SemanticExtractor {
	return &SemanticExtractor{
		filterNavigation:     true,
		filterPresentational: true,
	}
}

// Extract performs semantic hierarchy-based content extraction
func (se *SemanticExtractor) Extract(root *TextNode) string {
	if root == nil {
		return ""
	}

	// Step 1: Filter out presentational and navigational content
	filtered := se.filterTree(root)
	if filtered == nil {
		return ""
	}

	var extractedContent []string

	// Step 2: Extract content following semantic hierarchy
	// First, find and extract all h1 subtrees with their siblings
	remainingTree := filtered
	for {
		h1Subtree, siblings, newTree := se.extractNextH1WithSiblings(remainingTree)
		if h1Subtree == nil {
			break
		}

		// Add h1 subtree content
		if content := se.extractTextContent(h1Subtree); content != "" {
			extractedContent = append(extractedContent, content)
		}

		// Add sibling content in document order
		for _, sibling := range siblings {
			if content := se.extractTextContent(sibling); content != "" {
				extractedContent = append(extractedContent, content)
			}
		}

		remainingTree = newTree
	}

	// Step 3: Extract all h2 subtrees with their siblings from remaining content
	for {
		h2Subtree, siblings, newTree := se.extractNextH2WithSiblings(remainingTree)
		if h2Subtree == nil {
			break
		}

		// Add h2 subtree content
		if content := se.extractTextContent(h2Subtree); content != "" {
			extractedContent = append(extractedContent, content)
		}

		// Add sibling content in document order
		for _, sibling := range siblings {
			if content := se.extractTextContent(sibling); content != "" {
				extractedContent = append(extractedContent, content)
			}
		}

		remainingTree = newTree
	}

	// Step 4: Extract any remaining content in document order
	if remainingContent := se.extractRemainingContent(remainingTree); remainingContent != "" {
		extractedContent = append(extractedContent, remainingContent)
	}

	return strings.Join(extractedContent, "\n\n")
}

// filterTree removes presentational and navigational content
func (se *SemanticExtractor) filterTree(root *TextNode) *TextNode {
	if root == nil {
		return nil
	}

	// Check if this node should be filtered out
	if se.shouldFilterNode(root) {
		return nil
	}

	// Create a copy of the node to avoid modifying original
	filtered := &TextNode{
		Text:       root.Text,
		Tag:        root.Tag,
		Attributes: make(map[string]string),
		Depth:      root.Depth,
		Index:      root.Index,
	}

	// Copy attributes
	for k, v := range root.Attributes {
		filtered.Attributes[k] = v
	}

	// Filter and copy children
	var filteredChildren []*TextNode
	for _, child := range root.Children {
		if filteredChild := se.filterTree(child); filteredChild != nil {
			filteredChild.Parent = filtered
			filteredChildren = append(filteredChildren, filteredChild)
		}
	}

	filtered.Children = filteredChildren

	// Update indices
	for i, child := range filtered.Children {
		child.Index = i
	}

	return filtered
}

// shouldFilterNode determines if a node should be filtered out
func (se *SemanticExtractor) shouldFilterNode(node *TextNode) bool {
	tag := strings.ToLower(node.Tag)

	// Always filter these tags
	alwaysFilter := map[string]bool{
		"script": true, "style": true, "noscript": true,
		"meta": true, "link": true, "title": true,
	}
	if alwaysFilter[tag] {
		return true
	}

	// Filter navigation elements
	if se.filterNavigation {
		if se.isNavigationElement(node) {
			return true
		}
	}

	// Filter presentational elements
	if se.filterPresentational {
		if se.isPresentationalElement(node) {
			return true
		}
	}

	return false
}

// isNavigationElement checks if a node is a navigation element
func (se *SemanticExtractor) isNavigationElement(node *TextNode) bool {
	tag := strings.ToLower(node.Tag)

	// Direct navigation tags
	navTags := map[string]bool{
		"nav": true,
	}
	if navTags[tag] {
		return true
	}

	// Site-level header/footer (not article headers)
	if tag == "header" || tag == "footer" {
		// Only filter if it's likely a site header/footer, not article header
		if se.isSiteHeaderOrFooter(node) {
			return true
		}
	}

	// Check for navigation-related classes (be more specific)
	if class, exists := node.Attributes["class"]; exists {
		classLower := strings.ToLower(class)

		// Strong navigation indicators
		strongNavClasses := []string{"navigation", "main-nav", "site-nav", "breadcrumb", "pagination", "menu"}
		for _, navClass := range strongNavClasses {
			if strings.Contains(classLower, navClass) {
				return true
			}
		}

		// More specific header/footer filtering - only site-level, not article-level
		siteLevelIndicators := []string{"site-header", "site-footer", "main-header", "main-footer", "global-header", "global-footer"}
		for _, siteIndicator := range siteLevelIndicators {
			if strings.Contains(classLower, siteIndicator) {
				return true
			}
		}

		// Don't filter article/content headers
		contentHeaderIndicators := []string{"article-header", "post-header", "content-header", "entry-header"}
		for _, contentIndicator := range contentHeaderIndicators {
			if strings.Contains(classLower, contentIndicator) {
				return false // Explicitly preserve content headers
			}
		}

		// UI elements that should be filtered
		uiClasses := []string{"cookie-banner", "subscription-overlay", "social-share", "share-btn", "user-actions", "nav-menu"}
		for _, uiClass := range uiClasses {
			if strings.Contains(classLower, uiClass) {
				return true
			}
		}
	}

	// Check for navigation-related IDs
	if id, exists := node.Attributes["id"]; exists {
		idLower := strings.ToLower(id)
		navIds := []string{"navigation", "main-nav", "site-nav", "menu"}
		for _, navId := range navIds {
			if strings.Contains(idLower, navId) {
				return true
			}
		}
	}

	return false
}

// isSiteHeaderOrFooter checks if a header/footer is site-level vs article-level
func (se *SemanticExtractor) isSiteHeaderOrFooter(node *TextNode) bool {
	// Look for navigation-like content in header/footer
	if se.containsNavigationContent(node) {
		return true
	}

	// Check classes for site-level indicators
	if class, exists := node.Attributes["class"]; exists {
		classLower := strings.ToLower(class)
		siteLevelClasses := []string{"site-header", "site-footer", "main-header", "main-footer", "global-header", "global-footer"}
		for _, siteClass := range siteLevelClasses {
			if strings.Contains(classLower, siteClass) {
				return true
			}
		}
	}

	// Article headers usually contain h1-h6, while site headers contain nav elements
	if node.Tag == "header" {
		hasHeading := se.containsHeadings(node)
		hasNav := se.containsNavigation(node)
		// If it has navigation but no headings, likely a site header
		if hasNav && !hasHeading {
			return true
		}
	}

	// Check for footer with copyright/legal content
	if node.Tag == "footer" {
		if se.containsCopyrightContent(node) {
			return true
		}
	}

	return false
}

// containsNavigationContent checks if a node contains navigation-like content
func (se *SemanticExtractor) containsNavigationContent(node *TextNode) bool {
	if node == nil {
		return false
	}

	// Check for nav elements
	if node.Tag == "nav" {
		return true
	}

	// Check for navigation-like lists (ul/ol with multiple links)
	if node.Tag == "ul" || node.Tag == "ol" {
		linkCount := se.countLinks(node)
		if linkCount >= 3 { // Multiple links suggest navigation
			return true
		}
	}

	// Recursively check children
	for _, child := range node.Children {
		if se.containsNavigationContent(child) {
			return true
		}
	}

	return false
}

// containsHeadings checks if a node contains heading elements
func (se *SemanticExtractor) containsHeadings(node *TextNode) bool {
	if node == nil {
		return false
	}

	headings := map[string]bool{"h1": true, "h2": true, "h3": true, "h4": true, "h5": true, "h6": true}
	if headings[node.Tag] {
		return true
	}

	for _, child := range node.Children {
		if se.containsHeadings(child) {
			return true
		}
	}

	return false
}

// containsNavigation checks if a node contains nav elements
func (se *SemanticExtractor) containsNavigation(node *TextNode) bool {
	if node == nil {
		return false
	}

	if node.Tag == "nav" {
		return true
	}

	for _, child := range node.Children {
		if se.containsNavigation(child) {
			return true
		}
	}

	return false
}

// containsCopyrightContent checks if a node contains copyright or legal text
func (se *SemanticExtractor) containsCopyrightContent(node *TextNode) bool {
	if node == nil {
		return false
	}

	// Check text content for copyright indicators
	if node.Tag == "#text" {
		text := strings.ToLower(strings.TrimSpace(node.Text))
		copyrightIndicators := []string{"copyright", "©", "(c)", "all rights reserved", "terms of service", "privacy policy", "legal"}
		for _, indicator := range copyrightIndicators {
			if strings.Contains(text, indicator) {
				return true
			}
		}
	}

	// Recursively check children
	for _, child := range node.Children {
		if se.containsCopyrightContent(child) {
			return true
		}
	}

	return false
}

// countLinks counts the number of link elements in a node
func (se *SemanticExtractor) countLinks(node *TextNode) int {
	if node == nil {
		return 0
	}

	count := 0
	if node.Tag == "a" {
		count = 1
	}

	for _, child := range node.Children {
		count += se.countLinks(child)
	}

	return count
}

// isPresentationalElement checks if a node is purely presentational
func (se *SemanticExtractor) isPresentationalElement(node *TextNode) bool {
	tag := strings.ToLower(node.Tag)

	// Presentational tags
	presentationalTags := map[string]bool{
		"hr": true, "br": true, "wbr": true,
	}
	if presentationalTags[tag] {
		return true
	}

	// UI elements that should be filtered
	uiTags := map[string]bool{
		"button": true, "svg": true, "img": true,
	}
	if uiTags[tag] {
		return true
	}

	// Check for UI-related classes
	if class, exists := node.Attributes["class"]; exists {
		classLower := strings.ToLower(class)

		// Don't filter article/content headers (preserve content)
		contentHeaderIndicators := []string{"article-header", "post-header", "content-header", "entry-header"}
		for _, contentIndicator := range contentHeaderIndicators {
			if strings.Contains(classLower, contentIndicator) {
				return false // Explicitly preserve content headers
			}
		}

		// UI classes - be more specific to avoid false positives
		// Use word boundaries or specific patterns to avoid matching partial words
		if strings.Contains(classLower, " ad ") ||
			strings.HasPrefix(classLower, "ad ") ||
			strings.HasSuffix(classLower, " ad") ||
			classLower == "ad" {
			return true
		}

		uiClasses := []string{
			"btn", "button", "share", "social", "advertisement", "banner",
			"popup", "overlay", "modal", "tooltip", "dropdown", "widget",
			"sidebar", "newsletter", "signup", "subscribe", "cookie", "chat",
		}
		for _, uiClass := range uiClasses {
			if strings.Contains(classLower, uiClass) {
				return true
			}
		}
	}

	// Check for hidden elements
	if se.isHidden(node) {
		return true
	}

	return false
}

// isHidden checks if an element is visually hidden
func (se *SemanticExtractor) isHidden(node *TextNode) bool {
	// Check for hidden classes
	if class, exists := node.Attributes["class"]; exists {
		hiddenClasses := []string{"hidden", "invisible", "sr-only", "screen-reader", "visually-hidden"}
		classLower := strings.ToLower(class)
		for _, hiddenClass := range hiddenClasses {
			if strings.Contains(classLower, hiddenClass) {
				return true
			}
		}
	}

	// Check for display:none or visibility:hidden in style
	if style, exists := node.Attributes["style"]; exists {
		styleLower := strings.ToLower(style)
		if strings.Contains(styleLower, "display:none") ||
			strings.Contains(styleLower, "display: none") ||
			strings.Contains(styleLower, "visibility:hidden") ||
			strings.Contains(styleLower, "visibility: hidden") {
			return true
		}
	}

	return false
}

// extractNextH1WithSiblings finds the smallest subtree containing the next h1 and returns it with siblings
func (se *SemanticExtractor) extractNextH1WithSiblings(root *TextNode) (*TextNode, []*TextNode, *TextNode) {
	if root == nil {
		return nil, nil, nil
	}

	// Find the first h1 in document order
	h1Node := se.findFirstHeading(root, "h1")
	if h1Node == nil {
		return nil, nil, root // No h1 found
	}

	// Use the h1 node itself as the subtree (don't look for containing elements)
	subtree := h1Node

	// Get siblings of this h1
	siblings := se.getSiblings(subtree)

	// Remove this subtree and its siblings from the tree
	newTree := se.removeSubtreeAndSiblings(root, subtree, siblings)

	return subtree, siblings, newTree
}

// extractNextH2WithSiblings finds the smallest subtree containing the next h2 and returns it with siblings
func (se *SemanticExtractor) extractNextH2WithSiblings(root *TextNode) (*TextNode, []*TextNode, *TextNode) {
	if root == nil {
		return nil, nil, nil
	}

	// Find the first h2 in document order
	h2Node := se.findFirstHeading(root, "h2")
	if h2Node == nil {
		return nil, nil, root // No h2 found
	}

	// Use the h2 node itself as the subtree (don't look for containing elements)
	subtree := h2Node

	// Get siblings of this h2
	siblings := se.getSiblings(subtree)

	// Remove this subtree and its siblings from the tree
	newTree := se.removeSubtreeAndSiblings(root, subtree, siblings)

	return subtree, siblings, newTree
}

// findFirstHeading finds the first heading of the specified level in document order
func (se *SemanticExtractor) findFirstHeading(node *TextNode, headingTag string) *TextNode {
	if node == nil {
		return nil
	}

	// Check if this node is the heading we're looking for
	if strings.ToLower(node.Tag) == headingTag {
		return node
	}

	// Search children in document order
	for _, child := range node.Children {
		if found := se.findFirstHeading(child, headingTag); found != nil {
			return found
		}
	}

	return nil
}

// getSiblings returns the siblings of a node (including the node itself)
func (se *SemanticExtractor) getSiblings(node *TextNode) []*TextNode {
	if node == nil || node.Parent == nil {
		return []*TextNode{node}
	}

	// Return all children of the parent (which includes the node and its siblings)
	return node.Parent.Children
}

// removeSubtreeAndSiblings removes a subtree and its siblings from the tree
func (se *SemanticExtractor) removeSubtreeAndSiblings(root, subtree *TextNode, siblings []*TextNode) *TextNode {
	if root == nil {
		return nil
	}

	// Create a copy of the tree without the specified subtree and siblings
	return se.copyTreeExcluding(root, subtree, siblings)
}

// copyTreeExcluding creates a copy of the tree excluding specified nodes
func (se *SemanticExtractor) copyTreeExcluding(node *TextNode, excludeSubtree *TextNode, excludeSiblings []*TextNode) *TextNode {
	if node == nil {
		return nil
	}

	// Check if this node should be excluded
	if node == excludeSubtree {
		return nil
	}

	for _, sibling := range excludeSiblings {
		if node == sibling {
			return nil
		}
	}

	// Create a copy of this node
	copy := &TextNode{
		Text:       node.Text,
		Tag:        node.Tag,
		Attributes: make(map[string]string),
		Depth:      node.Depth,
		Index:      node.Index,
	}

	// Copy attributes
	for k, v := range node.Attributes {
		copy.Attributes[k] = v
	}

	// Copy children that aren't excluded
	var copiedChildren []*TextNode
	for _, child := range node.Children {
		if copiedChild := se.copyTreeExcluding(child, excludeSubtree, excludeSiblings); copiedChild != nil {
			copiedChild.Parent = copy
			copiedChildren = append(copiedChildren, copiedChild)
		}
	}

	copy.Children = copiedChildren

	// Update indices
	for i, child := range copy.Children {
		child.Index = i
	}

	return copy
}

// extractRemainingContent extracts all remaining content from the tree
func (se *SemanticExtractor) extractRemainingContent(root *TextNode) string {
	if root == nil {
		return ""
	}

	return se.extractTextContent(root)
}

// extractTextContent extracts readable text content from a node and its descendants
func (se *SemanticExtractor) extractTextContent(node *TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		text := strings.TrimSpace(node.Text)
		if text == "" {
			return ""
		}
		return text
	}

	var parts []string

	// Add spacing for block elements
	isBlockElement := se.isBlockElement(node.Tag)

	for _, child := range node.Children {
		content := se.extractTextContent(child)
		if content != "" {
			parts = append(parts, content)
		}
	}

	if len(parts) == 0 {
		return ""
	}

	if isBlockElement {
		return strings.Join(parts, "\n")
	} else {
		return strings.Join(parts, " ")
	}
}

// isBlockElement checks if a tag represents a block-level element
func (se *SemanticExtractor) isBlockElement(tag string) bool {
	blockElements := map[string]bool{
		"div": true, "p": true, "h1": true, "h2": true, "h3": true,
		"h4": true, "h5": true, "h6": true, "article": true, "section": true,
		"main": true, "header": true, "footer": true, "aside": true,
		"blockquote": true, "pre": true, "ul": true, "ol": true, "li": true,
		"dl": true, "dt": true, "dd": true, "table": true, "tr": true,
		"th": true, "td": true, "form": true, "fieldset": true,
	}

	return blockElements[strings.ToLower(tag)]
}
