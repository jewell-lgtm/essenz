// Package extractor provides content extraction functionality for reader view.
package extractor

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Extractor handles content extraction from HTML documents.
type Extractor struct {
	// Configuration options
	minContentLength   int
	preserveFormatting bool
}

// New creates a new content extractor with default settings.
func New() *Extractor {
	return &Extractor{
		minContentLength:   100,
		preserveFormatting: true,
	}
}

// ExtractContent extracts the main content from HTML and converts it to markdown.
func (e *Extractor) ExtractContent(htmlContent string) (string, error) {
	// Parse HTML
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find the main content
	contentNode := e.findMainContent(doc)
	if contentNode == nil {
		// Fallback: use the body or entire document
		if bodyNode := e.findNode(doc, "body"); bodyNode != nil {
			contentNode = bodyNode
		} else {
			contentNode = doc
		}
	}

	// Extract text and convert to markdown
	markdown := e.nodeToMarkdown(contentNode)

	// Clean up the output
	markdown = e.cleanMarkdown(markdown)

	return markdown, nil
}

// findMainContent attempts to identify the main content area of the page.
func (e *Extractor) findMainContent(n *html.Node) *html.Node {
	// Look for semantic HTML5 elements first
	if mainNode := e.findNode(n, "main"); mainNode != nil {
		return mainNode
	}
	if articleNode := e.findNode(n, "article"); articleNode != nil {
		return articleNode
	}

	// Look for common content containers by class/id
	contentSelectors := []string{
		"content", "main-content", "article", "post", "entry",
		"story", "text", "body-content", "primary",
	}

	bestNode := e.findBestContentNode(n, contentSelectors)
	if bestNode != nil {
		return bestNode
	}

	return nil
}

// findBestContentNode finds the node with the highest content score.
func (e *Extractor) findBestContentNode(n *html.Node, selectors []string) *html.Node {
	var bestNode *html.Node
	bestScore := 0

	e.walkNodes(n, func(node *html.Node) {
		if node.Type != html.ElementNode {
			return
		}

		score := e.scoreNode(node, selectors)
		if score > bestScore {
			bestScore = score
			bestNode = node
		}
	})

	return bestNode
}

// scoreNode calculates a content score for a node.
func (e *Extractor) scoreNode(n *html.Node, selectors []string) int {
	if n.Type != html.ElementNode {
		return 0
	}

	score := 0

	// Positive scores for content-like elements
	switch n.Data {
	case "article", "main", "section":
		score += 25
	case "div", "p":
		score += 5
	case "h1", "h2", "h3", "h4", "h5", "h6":
		score += 10
	}

	// Check attributes for content indicators
	for _, attr := range n.Attr {
		if attr.Key == "class" || attr.Key == "id" {
			value := strings.ToLower(attr.Val)
			for _, selector := range selectors {
				if strings.Contains(value, selector) {
					score += 15
				}
			}

			// Negative scores for non-content elements
			if containsAny(value, []string{"nav", "menu", "sidebar", "footer", "header", "ad", "social", "comment"}) {
				score -= 25
			}
		}
	}

	// Add text content length bonus
	textLength := len(strings.TrimSpace(e.getTextContent(n)))
	if textLength > e.minContentLength {
		score += textLength / 25 // Bonus points for longer content
	}

	return score
}

// nodeToMarkdown converts an HTML node tree to markdown.
func (e *Extractor) nodeToMarkdown(n *html.Node) string {
	if n == nil {
		return ""
	}

	var result strings.Builder
	e.convertNode(n, &result, 0)
	return result.String()
}

// convertNode recursively converts HTML nodes to markdown.
func (e *Extractor) convertNode(n *html.Node, result *strings.Builder, depth int) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			result.WriteString(text)
		}
		return
	}

	if n.Type != html.ElementNode {
		return
	}

	// Skip unwanted elements
	if e.shouldSkipElement(n) {
		return
	}

	// Handle opening tags
	e.writeOpeningTag(n, result)

	// Process children
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		e.convertNode(child, result, depth+1)
	}

	// Handle closing tags
	e.writeClosingTag(n, result)
}

// shouldSkipElement determines if an element should be skipped entirely.
func (e *Extractor) shouldSkipElement(n *html.Node) bool {
	switch n.Data {
	case "nav", "footer", "header", "aside", "script", "style", "noscript":
		return true
	}

	// Check for unwanted classes/IDs
	for _, attr := range n.Attr {
		if attr.Key == "class" || attr.Key == "id" {
			value := strings.ToLower(attr.Val)
			if containsAny(value, []string{"nav", "menu", "sidebar", "footer", "header", "ad", "social", "comment"}) {
				return true
			}
		}
	}

	return false
}

// writeOpeningTag handles opening markdown syntax.
func (e *Extractor) writeOpeningTag(n *html.Node, result *strings.Builder) {
	switch n.Data {
	case "h1":
		result.WriteString("# ")
	case "h2":
		result.WriteString("## ")
	case "h3":
		result.WriteString("### ")
	case "h4":
		result.WriteString("#### ")
	case "h5":
		result.WriteString("##### ")
	case "h6":
		result.WriteString("###### ")
	case "strong", "b":
		result.WriteString("**")
	case "em", "i":
		result.WriteString("*")
	case "blockquote":
		result.WriteString("> ")
	case "li":
		result.WriteString("- ")
	case "a":
		result.WriteString("[")
	}
}

// writeClosingTag handles closing markdown syntax.
func (e *Extractor) writeClosingTag(n *html.Node, result *strings.Builder) {
	switch n.Data {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		result.WriteString("\n\n")
	case "p", "div":
		if e.hasTextContent(n) {
			result.WriteString("\n\n")
		}
	case "br":
		result.WriteString("\n")
	case "strong", "b":
		result.WriteString("**")
	case "em", "i":
		result.WriteString("*")
	case "blockquote":
		result.WriteString("\n\n")
	case "li":
		result.WriteString("\n")
	case "ul", "ol":
		result.WriteString("\n")
	case "a":
		// Extract link URL
		href := ""
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}
		if href != "" {
			result.WriteString(fmt.Sprintf("](%s)", href))
		} else {
			result.WriteString("]")
		}
	}
}

// Helper functions

func (e *Extractor) findNode(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if found := e.findNode(child, tagName); found != nil {
			return found
		}
	}
	return nil
}

func (e *Extractor) walkNodes(n *html.Node, fn func(*html.Node)) {
	fn(n)
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		e.walkNodes(child, fn)
	}
}

func (e *Extractor) getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text strings.Builder
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		text.WriteString(e.getTextContent(child))
	}
	return text.String()
}

func (e *Extractor) hasTextContent(n *html.Node) bool {
	return len(strings.TrimSpace(e.getTextContent(n))) > 0
}

func (e *Extractor) cleanMarkdown(markdown string) string {
	// Remove excessive newlines
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")

	// Trim leading/trailing whitespace
	markdown = strings.TrimSpace(markdown)

	return markdown
}

func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
