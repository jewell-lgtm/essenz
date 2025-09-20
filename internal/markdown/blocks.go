package markdown

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// BlockRenderer defines an interface for rendering block-level elements
type BlockRenderer interface {
	CanRender(node *tree.TextNode) bool
	Render(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error)
	Priority() int
}

// HeadingRenderer handles heading elements (h1-h6)
type HeadingRenderer struct{}

// NewHeadingRenderer creates a new HeadingRenderer
func NewHeadingRenderer() *HeadingRenderer {
	return &HeadingRenderer{}
}

// CanRender checks if this renderer can handle the node
func (hr *HeadingRenderer) CanRender(node *tree.TextNode) bool {
	tag := strings.ToLower(node.Tag)
	return tag == "h1" || tag == "h2" || tag == "h3" || tag == "h4" || tag == "h5" || tag == "h6"
}

// Render renders a heading element
func (hr *HeadingRenderer) Render(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	level := hr.getHeadingLevel(node.Tag)
	content := hr.extractTextContent(node)

	if content == "" {
		return "", nil
	}

	// Generate ATX-style heading
	prefix := strings.Repeat("#", level)
	return fmt.Sprintf("\n%s %s\n\n", prefix, content), nil
}

// Priority returns the priority of this renderer
func (hr *HeadingRenderer) Priority() int {
	return 100
}

// getHeadingLevel extracts numeric level from heading tag
func (hr *HeadingRenderer) getHeadingLevel(tag string) int {
	tag = strings.ToLower(tag)
	if len(tag) == 2 && tag[0] == 'h' {
		if level, err := strconv.Atoi(string(tag[1])); err == nil && level >= 1 && level <= 6 {
			return level
		}
	}
	return 1
}

// extractTextContent recursively extracts text from a node
func (hr *HeadingRenderer) extractTextContent(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var parts []string
	for _, child := range node.Children {
		if text := hr.extractTextContent(child); text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, " ")
}

// ParagraphRenderer handles paragraph elements
type ParagraphRenderer struct{}

// NewParagraphRenderer creates a new ParagraphRenderer
func NewParagraphRenderer() *ParagraphRenderer {
	return &ParagraphRenderer{}
}

// CanRender checks if this renderer can handle the node
func (pr *ParagraphRenderer) CanRender(node *tree.TextNode) bool {
	return strings.ToLower(node.Tag) == "p"
}

// Render renders a paragraph element
func (pr *ParagraphRenderer) Render(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	content, err := pr.renderParagraphContent(node, state, renderer)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", nil
	}

	return content + "\n\n", nil
}

// Priority returns the priority of this renderer
func (pr *ParagraphRenderer) Priority() int {
	return 50
}

// renderParagraphContent renders paragraph content with inline formatting
func (pr *ParagraphRenderer) renderParagraphContent(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	var result strings.Builder

	for _, child := range node.Children {
		if child.Tag == "#text" {
			result.WriteString(renderer.renderTextContent(child.Text, state))
		} else {
			// Handle inline elements
			inline, err := pr.renderInlineElement(child, state, renderer)
			if err != nil {
				return "", err
			}
			result.WriteString(inline)
		}
	}

	return strings.TrimSpace(result.String()), nil
}

// renderInlineElement renders inline elements within paragraphs
func (pr *ParagraphRenderer) renderInlineElement(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	tag := strings.ToLower(node.Tag)

	switch tag {
	case "strong", "b":
		content := pr.extractTextContent(node)
		return renderer.style.FormatStrong(content), nil
	case "em", "i":
		content := pr.extractTextContent(node)
		return renderer.style.FormatEmphasis(content), nil
	case "code":
		content := pr.extractTextContent(node)
		return renderer.style.FormatInlineCode(content), nil
	case "a":
		return pr.renderLink(node, renderer), nil
	default:
		// For other inline elements, just extract text
		return pr.extractTextContent(node), nil
	}
}

// renderLink renders link elements
func (pr *ParagraphRenderer) renderLink(node *tree.TextNode, renderer *TreeRenderer) string {
	href := node.Attributes["href"]
	text := pr.extractTextContent(node)

	if href == "" {
		return text
	}

	return fmt.Sprintf("[%s](%s)", text, href)
}

// extractTextContent recursively extracts text from a node
func (pr *ParagraphRenderer) extractTextContent(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var parts []string
	for _, child := range node.Children {
		if text := pr.extractTextContent(child); text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, " ")
}

// ListRenderer handles list elements (ul, ol, li)
type ListRenderer struct{}

// NewListRenderer creates a new ListRenderer
func NewListRenderer() *ListRenderer {
	return &ListRenderer{}
}

// CanRender checks if this renderer can handle the node
func (lr *ListRenderer) CanRender(node *tree.TextNode) bool {
	tag := strings.ToLower(node.Tag)
	return tag == "ul" || tag == "ol"
}

// Render renders a list element
func (lr *ListRenderer) Render(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	tag := strings.ToLower(node.Tag)
	isOrdered := tag == "ol"

	var result strings.Builder
	counter := 1

	for _, child := range node.Children {
		if strings.ToLower(child.Tag) == "li" {
			item, err := lr.renderListItem(child, state, renderer, isOrdered, counter)
			if err != nil {
				return "", err
			}
			if item != "" {
				result.WriteString(item)
				counter++
			}
		}
	}

	return result.String() + "\n", nil
}

// Priority returns the priority of this renderer
func (lr *ListRenderer) Priority() int {
	return 80
}

// renderListItem renders a single list item
func (lr *ListRenderer) renderListItem(node *tree.TextNode, state *RenderState, renderer *TreeRenderer, isOrdered bool, counter int) (string, error) {
	var marker string
	indent := strings.Repeat(" ", state.CurrentDepth*renderer.config.ListStyle.IndentSize)

	if isOrdered {
		marker = fmt.Sprintf("%d. ", counter)
	} else {
		marker = renderer.config.ListStyle.UnorderedMarker + " "
	}

	content, err := lr.renderItemContent(node, state, renderer)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", nil
	}

	return fmt.Sprintf("%s%s%s\n", indent, marker, content), nil
}

// renderItemContent renders the content of a list item
func (lr *ListRenderer) renderItemContent(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	var result strings.Builder

	// Increase depth for nested elements
	state.CurrentDepth++
	defer func() { state.CurrentDepth-- }()

	for _, child := range node.Children {
		if child.Tag == "#text" {
			text := strings.TrimSpace(child.Text)
			if text != "" {
				result.WriteString(text)
			}
		} else if strings.ToLower(child.Tag) == "ul" || strings.ToLower(child.Tag) == "ol" {
			// Handle nested lists
			nested, err := lr.Render(child, state, renderer)
			if err != nil {
				return "", err
			}
			if nested != "" {
				result.WriteString("\n" + nested)
			}
		} else {
			// Handle other inline elements
			content, err := renderer.renderNode(context.Background(), child, state)
			if err != nil {
				return "", err
			}
			result.WriteString(content)
		}
	}

	return strings.TrimSpace(result.String()), nil
}

// BlockquoteRenderer handles blockquote elements
type BlockquoteRenderer struct{}

// NewBlockquoteRenderer creates a new BlockquoteRenderer
func NewBlockquoteRenderer() *BlockquoteRenderer {
	return &BlockquoteRenderer{}
}

// CanRender checks if this renderer can handle the node
func (br *BlockquoteRenderer) CanRender(node *tree.TextNode) bool {
	return strings.ToLower(node.Tag) == "blockquote"
}

// Render renders a blockquote element
func (br *BlockquoteRenderer) Render(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	content, err := br.extractBlockquoteContent(node, state, renderer)
	if err != nil {
		return "", err
	}

	if content == "" {
		return "", nil
	}

	// Format as blockquote with > prefix
	lines := strings.Split(content, "\n")
	var quotedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			quotedLines = append(quotedLines, "> "+line)
		}
	}

	return strings.Join(quotedLines, "\n") + "\n\n", nil
}

// Priority returns the priority of this renderer
func (br *BlockquoteRenderer) Priority() int {
	return 70
}

// extractBlockquoteContent extracts content from blockquote
func (br *BlockquoteRenderer) extractBlockquoteContent(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	var result strings.Builder

	for _, child := range node.Children {
		if child.Tag == "#text" {
			text := strings.TrimSpace(child.Text)
			if text != "" {
				result.WriteString(text + " ")
			}
		} else if strings.ToLower(child.Tag) == "p" {
			// Render paragraph content without extra newlines
			content, err := br.renderParagraphContent(child, state, renderer)
			if err != nil {
				return "", err
			}
			result.WriteString(content + " ")
		} else {
			// Render other elements
			content, err := renderer.renderNode(context.Background(), child, state)
			if err != nil {
				return "", err
			}
			result.WriteString(content)
		}
	}

	return strings.TrimSpace(result.String()), nil
}

// renderParagraphContent renders paragraph content for blockquotes
func (br *BlockquoteRenderer) renderParagraphContent(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	var result strings.Builder

	for _, child := range node.Children {
		if child.Tag == "#text" {
			result.WriteString(renderer.renderTextContent(child.Text, state))
		} else {
			// Handle inline elements
			tag := strings.ToLower(child.Tag)
			content := br.extractTextContent(child)

			switch tag {
			case "strong", "b":
				result.WriteString(renderer.style.FormatStrong(content))
			case "em", "i":
				result.WriteString(renderer.style.FormatEmphasis(content))
			case "code":
				result.WriteString(renderer.style.FormatInlineCode(content))
			case "a":
				href := child.Attributes["href"]
				if href != "" {
					result.WriteString(fmt.Sprintf("[%s](%s)", content, href))
				} else {
					result.WriteString(content)
				}
			default:
				result.WriteString(content)
			}
		}
	}

	return strings.TrimSpace(result.String()), nil
}

// extractTextContent recursively extracts text from a node
func (br *BlockquoteRenderer) extractTextContent(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var parts []string
	for _, child := range node.Children {
		if text := br.extractTextContent(child); text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, " ")
}

// CodeBlockRenderer handles pre/code elements
type CodeBlockRenderer struct{}

// NewCodeBlockRenderer creates a new CodeBlockRenderer
func NewCodeBlockRenderer() *CodeBlockRenderer {
	return &CodeBlockRenderer{}
}

// CanRender checks if this renderer can handle the node
func (cbr *CodeBlockRenderer) CanRender(node *tree.TextNode) bool {
	return strings.ToLower(node.Tag) == "pre"
}

// Render renders a code block element
func (cbr *CodeBlockRenderer) Render(node *tree.TextNode, state *RenderState, renderer *TreeRenderer) (string, error) {
	// Look for code element inside pre
	var codeNode *tree.TextNode
	var language string

	for _, child := range node.Children {
		if strings.ToLower(child.Tag) == "code" {
			codeNode = child
			// Extract language from class attribute
			if class, exists := child.Attributes["class"]; exists {
				if strings.HasPrefix(class, "language-") {
					language = strings.TrimPrefix(class, "language-")
				}
			}
			break
		}
	}

	var content string
	if codeNode != nil {
		content = cbr.extractCodeContent(codeNode)
	} else {
		content = cbr.extractCodeContent(node)
	}

	if content == "" {
		return "", nil
	}

	// Generate fenced code block
	if language != "" {
		return fmt.Sprintf("\n```%s\n%s\n```\n\n", language, content), nil
	}
	return fmt.Sprintf("\n```\n%s\n```\n\n", content), nil
}

// Priority returns the priority of this renderer
func (cbr *CodeBlockRenderer) Priority() int {
	return 90
}

// extractCodeContent extracts code content preserving formatting
func (cbr *CodeBlockRenderer) extractCodeContent(node *tree.TextNode) string {
	if node == nil {
		return ""
	}

	if node.Tag == "#text" {
		return node.Text
	}

	var result strings.Builder
	for _, child := range node.Children {
		result.WriteString(cbr.extractCodeContent(child))
	}

	return result.String()
}
