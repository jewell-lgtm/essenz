// Package extractor provides content extraction for reader view, powered by
// go-readability (the same algorithm as Firefox Reader View) for main-content
// isolation, with a lightweight HTML->markdown renderer on top.
package extractor

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	readability "github.com/go-shiori/go-readability"
	"golang.org/x/net/html"
)

// Extractor handles content extraction from HTML documents.
type Extractor struct {
	preserveFormatting bool
}

// New creates a new content extractor with default settings.
func New() *Extractor {
	return &Extractor{
		preserveFormatting: true,
	}
}

// ExtractContent isolates the main content from HTML and converts it to markdown.
// It uses go-readability to find the article body, falling back to the document
// body for pages readability can't classify (so trivial pages still pass through).
func (e *Extractor) ExtractContent(htmlContent string) (string, error) {
	if article, err := readability.FromReader(strings.NewReader(htmlContent), &url.URL{}); err == nil &&
		strings.TrimSpace(article.TextContent) != "" {
		if md := e.renderArticle(article); strings.TrimSpace(md) != "" {
			return md, nil
		}
	}
	return e.fallback(htmlContent)
}

// renderArticle converts a readability article to markdown, prepending the title
// as an H1 when readability split it out of the content.
func (e *Extractor) renderArticle(article readability.Article) string {
	doc, err := html.Parse(strings.NewReader(article.Content))
	if err != nil {
		return ""
	}
	md := e.cleanMarkdown(e.nodeToMarkdown(doc))
	if title := strings.TrimSpace(article.Title); title != "" && !strings.Contains(md, title) {
		md = "# " + title + "\n\n" + md
	}
	return md
}

// fallback renders the document body directly for pages readability rejects.
func (e *Extractor) fallback(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}
	node := e.findNode(doc, "body")
	if node == nil {
		node = doc
	}
	return e.cleanMarkdown(e.nodeToMarkdown(node)), nil
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
		if text := strings.TrimSpace(n.Data); text != "" {
			result.WriteString(text)
		}
		return
	}

	// Container nodes (document root) carry no markup themselves; recurse.
	if n.Type == html.DocumentNode {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			e.convertNode(child, result, depth)
		}
		return
	}

	if n.Type != html.ElementNode {
		return
	}

	if e.shouldSkipElement(n) {
		return
	}

	e.writeOpeningTag(n, result)
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		e.convertNode(child, result, depth+1)
	}
	e.writeClosingTag(n, result)
}

var (
	tokenSplitRe   = regexp.MustCompile(`[^a-z0-9]+`)
	boilerplateTok = map[string]bool{
		"nav": true, "menu": true, "sidebar": true, "footer": true,
		"header": true, "ad": true, "ads": true, "social": true, "comment": true,
	}
)

// shouldSkipElement drops non-content and non-rendering elements. readability
// already removes most boilerplate; this guards the fallback path too. Class/id
// matching is token-based to avoid false positives (e.g. "ad" inside
// "readability-page").
func (e *Extractor) shouldSkipElement(n *html.Node) bool {
	switch n.Data {
	case "nav", "footer", "header", "aside", "script", "style", "noscript":
		return true
	}
	for _, attr := range n.Attr {
		if attr.Key == "class" || attr.Key == "id" {
			for _, tok := range tokenSplitRe.Split(strings.ToLower(attr.Val), -1) {
				if boilerplateTok[tok] {
					return true
				}
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
	re := regexp.MustCompile(`\n{3,}`)
	markdown = re.ReplaceAllString(markdown, "\n\n")
	return strings.TrimSpace(markdown)
}
