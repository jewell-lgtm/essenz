package markdown

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// InlineRenderer defines an interface for rendering inline elements
type InlineRenderer interface {
	CanRender(node *tree.TextNode) bool
	Render(content string, node *tree.TextNode) string
}

// EmphasisRenderer handles emphasis elements (em, i)
type EmphasisRenderer struct{}

// NewEmphasisRenderer creates a new EmphasisRenderer
func NewEmphasisRenderer() *EmphasisRenderer {
	return &EmphasisRenderer{}
}

// CanRender checks if this renderer can handle the node
func (er *EmphasisRenderer) CanRender(node *tree.TextNode) bool {
	tag := strings.ToLower(node.Tag)
	return tag == "em" || tag == "i"
}

// Render renders emphasis formatting
func (er *EmphasisRenderer) Render(content string, node *tree.TextNode) string {
	if content == "" {
		return ""
	}
	return "*" + content + "*"
}

// StrongRenderer handles strong elements (strong, b)
type StrongRenderer struct{}

// NewStrongRenderer creates a new StrongRenderer
func NewStrongRenderer() *StrongRenderer {
	return &StrongRenderer{}
}

// CanRender checks if this renderer can handle the node
func (sr *StrongRenderer) CanRender(node *tree.TextNode) bool {
	tag := strings.ToLower(node.Tag)
	return tag == "strong" || tag == "b"
}

// Render renders strong formatting
func (sr *StrongRenderer) Render(content string, node *tree.TextNode) string {
	if content == "" {
		return ""
	}
	return "**" + content + "**"
}

// LinkRenderer handles link elements (a)
type LinkRenderer struct{}

// NewLinkRenderer creates a new LinkRenderer
func NewLinkRenderer() *LinkRenderer {
	return &LinkRenderer{}
}

// CanRender checks if this renderer can handle the node
func (lr *LinkRenderer) CanRender(node *tree.TextNode) bool {
	return strings.ToLower(node.Tag) == "a"
}

// Render renders link formatting
func (lr *LinkRenderer) Render(content string, node *tree.TextNode) string {
	href := node.Attributes["href"]
	if href == "" {
		return content
	}

	if content == "" {
		content = href
	}

	return "[" + content + "](" + href + ")"
}

// CodeSpanRenderer handles inline code elements (code)
type CodeSpanRenderer struct{}

// NewCodeSpanRenderer creates a new CodeSpanRenderer
func NewCodeSpanRenderer() *CodeSpanRenderer {
	return &CodeSpanRenderer{}
}

// CanRender checks if this renderer can handle the node
func (csr *CodeSpanRenderer) CanRender(node *tree.TextNode) bool {
	return strings.ToLower(node.Tag) == "code"
}

// Render renders inline code formatting
func (csr *CodeSpanRenderer) Render(content string, node *tree.TextNode) string {
	if content == "" {
		return ""
	}
	return "`" + content + "`"
}
