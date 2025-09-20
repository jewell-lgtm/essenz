package media

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// ContextAnalyzer analyzes the context around media elements to enhance descriptions.
type ContextAnalyzer struct {
	radiusWords int
}

// NewContextAnalyzer creates a new ContextAnalyzer.
func NewContextAnalyzer(radiusWords int) *ContextAnalyzer {
	return &ContextAnalyzer{
		radiusWords: radiusWords,
	}
}

// ExtractContext extracts contextual information around a media node.
func (ca *ContextAnalyzer) ExtractContext(node *tree.TextNode) string {
	if node == nil || node.Parent == nil {
		return ""
	}

	// Get surrounding text from parent and siblings
	var contextParts []string

	// Add text from previous siblings
	prevText := ca.getPreviousSiblingText(node)
	if prevText != "" {
		contextParts = append(contextParts, prevText)
	}

	// Add text from next siblings
	nextText := ca.getNextSiblingText(node)
	if nextText != "" {
		contextParts = append(contextParts, nextText)
	}

	// Add text from parent context if needed
	if len(contextParts) == 0 {
		parentText := ca.getParentContext(node)
		if parentText != "" {
			contextParts = append(contextParts, parentText)
		}
	}

	context := strings.Join(contextParts, " ")
	return ca.limitWords(context, ca.radiusWords)
}

// FindAssociatedCaption finds captions associated with a media element.
func (ca *ContextAnalyzer) FindAssociatedCaption(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	// Check if this node is inside a figure with figcaption
	if caption := ca.findFigcaption(node); caption != "" {
		return caption
	}

	// Check for caption in parent or sibling elements
	if node.Parent != nil {
		// Look for figcaption siblings
		for _, sibling := range node.Parent.Children {
			if strings.ToLower(sibling.Tag) == "figcaption" {
				return ca.extractTextFromNode(sibling)
			}
		}

		// Check if parent is a figure
		if strings.ToLower(node.Parent.Tag) == "figure" {
			return ca.findFigcaption(node.Parent)
		}
	}

	// Look for title attributes
	if title := node.Attributes["title"]; title != "" {
		return title
	}

	return ""
}

// AnalyzeSurroundingText analyzes the text surrounding a media element.
func (ca *ContextAnalyzer) AnalyzeSurroundingText(node *tree.TextNode) string {
	context := ca.ExtractContext(node)
	return ca.extractRelevantKeywords(context)
}

// getPreviousSiblingText gets text from previous sibling nodes.
func (ca *ContextAnalyzer) getPreviousSiblingText(node *tree.TextNode) string {
	if node.Parent == nil {
		return ""
	}

	var textParts []string
	var foundNode bool

	// Traverse siblings in reverse order
	siblings := node.Parent.Children
	for i := len(siblings) - 1; i >= 0; i-- {
		sibling := siblings[i]

		if sibling == node {
			foundNode = true
			continue
		}

		if foundNode {
			// This is a previous sibling
			text := ca.extractTextFromNode(sibling)
			if text != "" {
				textParts = append([]string{text}, textParts...) // Prepend to maintain order
			}
		}
	}

	return strings.Join(textParts, " ")
}

// getNextSiblingText gets text from next sibling nodes.
func (ca *ContextAnalyzer) getNextSiblingText(node *tree.TextNode) string {
	if node.Parent == nil {
		return ""
	}

	var textParts []string
	var foundNode bool

	for _, sibling := range node.Parent.Children {
		if foundNode {
			// This is a next sibling
			text := ca.extractTextFromNode(sibling)
			if text != "" {
				textParts = append(textParts, text)
			}
		}

		if sibling == node {
			foundNode = true
		}
	}

	return strings.Join(textParts, " ")
}

// getParentContext gets contextual text from parent elements.
func (ca *ContextAnalyzer) getParentContext(node *tree.TextNode) string {
	current := node.Parent
	for current != nil {
		// Look for text in the parent that's not in child elements
		text := ca.getDirectTextContent(current)
		if text != "" {
			return text
		}
		current = current.Parent
	}
	return ""
}

// getDirectTextContent gets direct text content from a node (not from children).
func (ca *ContextAnalyzer) getDirectTextContent(node *tree.TextNode) string {
	var textParts []string

	for _, child := range node.Children {
		if child.Tag == "#text" {
			text := strings.TrimSpace(child.Text)
			if text != "" {
				textParts = append(textParts, text)
			}
		}
	}

	return strings.Join(textParts, " ")
}

// extractTextFromNode recursively extracts all text from a node.
func (ca *ContextAnalyzer) extractTextFromNode(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var textParts []string
	for _, child := range node.Children {
		text := ca.extractTextFromNode(child)
		if text != "" {
			textParts = append(textParts, text)
		}
	}

	return strings.Join(textParts, " ")
}

// findFigcaption finds figcaption elements within a node tree.
func (ca *ContextAnalyzer) findFigcaption(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	for _, child := range node.Children {
		if strings.ToLower(child.Tag) == "figcaption" {
			return ca.extractTextFromNode(child)
		}

		// Recursively search in children
		if caption := ca.findFigcaption(child); caption != "" {
			return caption
		}
	}

	return ""
}

// limitWords limits text to a specified number of words.
func (ca *ContextAnalyzer) limitWords(text string, maxWords int) string {
	if text == "" {
		return ""
	}

	words := strings.Fields(text)
	if len(words) <= maxWords {
		return text
	}

	return strings.Join(words[:maxWords], " ") + "..."
}

// extractRelevantKeywords extracts relevant keywords from context text.
func (ca *ContextAnalyzer) extractRelevantKeywords(context string) string {
	if context == "" {
		return ""
	}

	words := strings.Fields(strings.ToLower(context))
	var keywords []string

	// Keywords that are likely to be descriptive
	descriptiveWords := map[string]bool{
		"architecture": true, "building": true, "design": true,
		"chart": true, "graph": true, "diagram": true,
		"photo": true, "picture": true, "illustration": true,
		"screenshot": true, "logo": true, "icon": true,
		"modern": true, "vintage": true, "colorful": true,
		"large": true, "small": true, "beautiful": true,
		"office": true, "exterior": true, "interior": true,
		"glass": true, "steel": true, "concrete": true,
	}

	for _, word := range words {
		// Remove punctuation
		word = strings.Trim(word, ".,!?;:")
		if len(word) > 3 && descriptiveWords[word] {
			keywords = append(keywords, word)
		}
	}

	if len(keywords) > 0 {
		return strings.Join(keywords, " ")
	}

	return ""
}
