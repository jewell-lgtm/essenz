// Package markdown provides sophisticated markdown generation from content trees,
// converting structured content into clean, well-formatted markdown output.
package markdown

import (
	"context"
	"fmt"
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// TreeRenderer converts content trees to clean, well-formatted markdown
type TreeRenderer struct {
	config RenderConfig
	blocks []BlockRenderer
	inline []InlineRenderer
	style  *StyleManager
}

// RenderConfig configures markdown rendering behavior
type RenderConfig struct {
	HeadingStyle       HeadingStyle   // ATX (#) or Setext (===)
	ListStyle          ListStyle      // Ordered/unordered preferences
	EmphasisStyle      EmphasisStyle  // * or _ for emphasis
	CodeBlockStyle     CodeBlockStyle // ``` or indented
	LineWidth          int            // Max line width for wrapping
	PreserveLineBreaks bool           // Maintain original line breaks
}

// HeadingStyle controls how headings are rendered
type HeadingStyle string

const (
	ATXHeading    HeadingStyle = "atx"    // # ## ###
	SetextHeading HeadingStyle = "setext" // === ---
)

// ListStyle controls list formatting
type ListStyle struct {
	UnorderedMarker string // "-", "*", or "+"
	OrderedFormat   string // "1." or "1)"
	IndentSize      int    // Spaces for nested lists
}

// EmphasisStyle controls emphasis formatting
type EmphasisStyle struct {
	Emphasis string // "*" or "_"
	Strong   string // "**" or "__"
}

// CodeBlockStyle controls code block formatting
type CodeBlockStyle string

const (
	FencedCodeBlock   CodeBlockStyle = "fenced"   // ```
	IndentedCodeBlock CodeBlockStyle = "indented" // 4-space indent
)

// RenderState tracks rendering context
type RenderState struct {
	CurrentDepth int
	ListStack    []ListContext
	HeadingCount map[int]int
	WithinCode   bool
	LineBuffer   strings.Builder
}

// ListContext tracks nested list state
type ListContext struct {
	Type    string // "ul" or "ol"
	Level   int
	Counter int
	Marker  string
}

// NewTreeRenderer creates a new TreeRenderer with default configuration
func NewTreeRenderer() *TreeRenderer {
	renderer := &TreeRenderer{
		config: RenderConfig{
			HeadingStyle: ATXHeading,
			ListStyle: ListStyle{
				UnorderedMarker: "-",
				OrderedFormat:   "1.",
				IndentSize:      2,
			},
			EmphasisStyle: EmphasisStyle{
				Emphasis: "*",
				Strong:   "**",
			},
			CodeBlockStyle:     FencedCodeBlock,
			LineWidth:          80,
			PreserveLineBreaks: false,
		},
		blocks: make([]BlockRenderer, 0),
		inline: make([]InlineRenderer, 0),
	}

	// Add default block renderers
	renderer.AddBlockRenderer(NewHeadingRenderer())
	renderer.AddBlockRenderer(NewParagraphRenderer())
	renderer.AddBlockRenderer(NewListRenderer())
	renderer.AddBlockRenderer(NewBlockquoteRenderer())
	renderer.AddBlockRenderer(NewCodeBlockRenderer())

	// Add default inline renderers
	renderer.AddInlineRenderer(NewEmphasisRenderer())
	renderer.AddInlineRenderer(NewStrongRenderer())
	renderer.AddInlineRenderer(NewLinkRenderer())
	renderer.AddInlineRenderer(NewCodeSpanRenderer())

	// Create style manager
	renderer.style = NewStyleManager(renderer.config)

	return renderer
}

// WithConfig sets the render configuration
func (tr *TreeRenderer) WithConfig(config RenderConfig) *TreeRenderer {
	tr.config = config
	tr.style = NewStyleManager(config)
	return tr
}

// WithEmphasisStyle sets the emphasis style
func (tr *TreeRenderer) WithEmphasisStyle(style string) *TreeRenderer {
	switch style {
	case "underscore":
		tr.config.EmphasisStyle.Emphasis = "_"
		tr.config.EmphasisStyle.Strong = "__"
	case "asterisk":
		tr.config.EmphasisStyle.Emphasis = "*"
		tr.config.EmphasisStyle.Strong = "**"
	}
	tr.style = NewStyleManager(tr.config)
	return tr
}

// WithListStyle sets the list style
func (tr *TreeRenderer) WithListStyle(style string) *TreeRenderer {
	switch style {
	case "asterisk":
		tr.config.ListStyle.UnorderedMarker = "*"
	case "plus":
		tr.config.ListStyle.UnorderedMarker = "+"
	case "dash":
		tr.config.ListStyle.UnorderedMarker = "-"
	}
	tr.style = NewStyleManager(tr.config)
	return tr
}

// AddBlockRenderer adds a block-level renderer
func (tr *TreeRenderer) AddBlockRenderer(renderer BlockRenderer) {
	tr.blocks = append(tr.blocks, renderer)
}

// AddInlineRenderer adds an inline renderer
func (tr *TreeRenderer) AddInlineRenderer(renderer InlineRenderer) {
	tr.inline = append(tr.inline, renderer)
}

// RenderTree converts a content tree to markdown
func (tr *TreeRenderer) RenderTree(ctx context.Context, root *tree.TextNode) (string, error) {
	if root == nil {
		return "", nil
	}

	state := &RenderState{
		CurrentDepth: 0,
		ListStack:    make([]ListContext, 0),
		HeadingCount: make(map[int]int),
		WithinCode:   false,
	}

	result, err := tr.renderNode(ctx, root, state)
	if err != nil {
		return "", fmt.Errorf("failed to render tree: %w", err)
	}

	// Post-process the markdown
	return tr.postProcess(result), nil
}

// renderNode recursively renders a node and its children
func (tr *TreeRenderer) renderNode(ctx context.Context, node *tree.TextNode, state *RenderState) (string, error) {
	if node == nil {
		return "", nil
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Handle text nodes directly
	if node.Tag == "#text" {
		return tr.renderTextContent(node.Text, state), nil
	}

	// Try block renderers first
	for _, renderer := range tr.blocks {
		if renderer.CanRender(node) {
			return renderer.Render(node, state, tr)
		}
	}

	// If no block renderer handles it, render children
	var result strings.Builder
	for _, child := range node.Children {
		childResult, err := tr.renderNode(ctx, child, state)
		if err != nil {
			return "", err
		}
		if childResult != "" {
			result.WriteString(childResult)
		}
	}

	return result.String(), nil
}

// renderTextContent renders text content with inline formatting
func (tr *TreeRenderer) renderTextContent(text string, state *RenderState) string {
	if state.WithinCode {
		return text
	}

	// Clean up whitespace
	cleaned := strings.TrimSpace(text)
	if cleaned == "" {
		return ""
	}

	return cleaned
}

// postProcess cleans up the generated markdown
func (tr *TreeRenderer) postProcess(markdown string) string {
	lines := strings.Split(markdown, "\n")
	var result []string

	previousWasEmpty := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Remove excessive blank lines
		if trimmed == "" {
			if !previousWasEmpty {
				result = append(result, "")
				previousWasEmpty = true
			}
		} else {
			result = append(result, line)
			previousWasEmpty = false
		}
	}

	// Ensure document ends with single newline
	final := strings.Join(result, "\n")
	return strings.TrimSpace(final) + "\n"
}
