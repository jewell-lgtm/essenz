# F5: Markdown Tree Renderer

**Feature Branch**: `feature/markdown-tree-renderer`

## Objective

Convert the filtered and media-processed content tree into clean, well-formatted markdown that preserves document structure, hierarchy, and readability while following markdown best practices.

## Technical Requirements

### Markdown Generation Strategy
- **Hierarchical Rendering**: Maintain proper heading hierarchy (H1-H6)
- **Structure Preservation**: Convert lists, blockquotes, and code blocks correctly
- **Inline Formatting**: Handle emphasis, strong, links, and inline code
- **Whitespace Management**: Ensure proper spacing between elements
- **Clean Output**: Generate readable, properly formatted markdown

### Content Tree to Markdown Mapping
```go
type MarkdownRenderer struct {
    config RenderConfig
    state  RenderState
}

type RenderConfig struct {
    HeadingStyle        HeadingStyle    // ATX (#) or Setext (===)
    ListStyle          ListStyle       // Ordered/unordered preferences
    EmphasisStyle      EmphasisStyle   // * or _ for emphasis
    CodeBlockStyle     CodeBlockStyle  // ``` or indented
    LinkStyle          LinkStyle       // Inline or reference
    LineWidth          int             // Max line width for wrapping
    PreserveLineBreaks bool           // Maintain original line breaks
}

type RenderState struct {
    currentDepth    int
    listStack      []ListContext
    headingCount   map[int]int
    linkReferences []LinkRef
    withinCode     bool
}
```

### Rendering Components
1. **Block Renderer**: Handles paragraphs, headings, lists, blockquotes
2. **Inline Renderer**: Processes emphasis, links, code spans
3. **List Renderer**: Manages nested lists and proper indentation
4. **Code Renderer**: Handles code blocks and inline code
5. **Link Renderer**: Generates inline or reference-style links

## Implementation Components

### 1. Internal Package: `internal/markdown`
```go
type TreeRenderer struct {
    config RenderConfig
    blocks []BlockRenderer
    inline []InlineRenderer
}

func (tr *TreeRenderer) RenderTree(root *ContentNode) (string, error)
func (tr *TreeRenderer) RenderNode(node *ContentNode, state *RenderState) (string, error)
func (tr *TreeRenderer) PostProcess(markdown string) string
```

### 2. Block-Level Rendering
```go
type BlockRenderer interface {
    CanRender(node *ContentNode) bool
    Render(node *ContentNode, state *RenderState) (string, error)
    Priority() int
}

type HeadingRenderer struct{}
type ParagraphRenderer struct{}
type ListRenderer struct{}
type BlockquoteRenderer struct{}
type CodeBlockRenderer struct{}
```

### 3. Inline Element Rendering
```go
type InlineRenderer interface {
    CanRender(node *ContentNode) bool
    Render(content string, node *ContentNode) string
}

type EmphasisRenderer struct{}
type StrongRenderer struct{}
type LinkRenderer struct{}
type CodeSpanRenderer struct{}
type ImageRenderer struct{}
```

### 4. Formatting and Style Management
```go
type StyleManager struct {
    config RenderConfig
}

func (sm *StyleManager) FormatHeading(level int, text string) string
func (sm *StyleManager) FormatList(items []string, ordered bool, level int) string
func (sm *StyleManager) FormatBlockquote(content string) string
func (sm *StyleManager) WrapText(text string, width int) string
```

## Acceptance Criteria

### Heading Hierarchy
1. Correctly maps content tree heading levels to markdown (H1-H6)
2. Maintains proper heading hierarchy without skipping levels
3. Handles multiple H1s appropriately
4. Generates clean ATX-style headings with proper spacing

### List Formatting
1. Renders unordered lists with consistent bullet style
2. Handles ordered lists with proper numbering
3. Manages nested lists with correct indentation (2 or 4 spaces)
4. Preserves list item content and inline formatting

### Text Formatting
1. Converts emphasis correctly (* or _ style)
2. Handles strong/bold formatting consistently
3. Processes inline code with backticks
4. Maintains link formatting with proper URL encoding

### Block Elements
1. Renders paragraphs with proper spacing (double newline)
2. Formats blockquotes with > prefix and proper indentation
3. Handles code blocks with triple backticks and language hints
4. Preserves table structure when present

### Clean Output
1. Generates consistent whitespace between elements
2. Avoids excessive blank lines or missing spacing
3. Produces readable, well-formatted markdown
4. Handles edge cases gracefully (empty elements, malformed content)

## Test Scenarios

### Basic Document Structure
Input content tree representing:
```html
<h1>Main Title</h1>
<p>Introduction paragraph with <strong>bold</strong> and <em>italic</em> text.</p>
<h2>Section Header</h2>
<ul>
    <li>First item</li>
    <li>Second item with <a href="link">link</a></li>
</ul>
```

Expected markdown output:
```markdown
# Main Title

Introduction paragraph with **bold** and *italic* text.

## Section Header

- First item
- Second item with [link](link)
```

### Complex Nested Structure
```html
<h1>Article Title</h1>
<h2>Introduction</h2>
<p>Some text with <code>inline code</code>.</p>
<blockquote>
    <p>This is a quote with <strong>emphasis</strong>.</p>
</blockquote>
<h2>Examples</h2>
<ol>
    <li>First example
        <ul>
            <li>Sub-item A</li>
            <li>Sub-item B</li>
        </ul>
    </li>
    <li>Second example</li>
</ol>
```

Expected markdown:
```markdown
# Article Title

## Introduction

Some text with `inline code`.

> This is a quote with **emphasis**.

## Examples

1. First example
   - Sub-item A
   - Sub-item B
2. Second example
```

### Code Block Handling
```html
<h2>Code Example</h2>
<pre><code class="language-javascript">
function example() {
    return "Hello, World!";
}
</code></pre>
<p>This function demonstrates basic syntax.</p>
```

Expected markdown:
```markdown
## Code Example

```javascript
function example() {
    return "Hello, World!";
}
```

This function demonstrates basic syntax.
```

### Link and Reference Management
- Inline links: `[text](url)`
- Reference-style links when configured
- Proper URL encoding and escaping
- Email address formatting

## Configuration Options

### Markdown Style Preferences
```yaml
heading_style: "atx"          # "atx" (#) or "setext" (===)
list_style:
  unordered: "-"              # "-", "*", or "+"
  ordered: "1."               # "1." or "1)"
  indent: 2                   # Spaces for nested lists

emphasis:
  style: "*"                  # "*" or "_"
  strong: "**"                # "**" or "__"

code:
  inline: "`"
  block: "```"                # "```" or 4-space indent
  language_hints: true

links:
  style: "inline"             # "inline" or "reference"
  encode_urls: true
```

### Formatting Options
```yaml
formatting:
  line_width: 80              # Maximum line width
  preserve_breaks: false      # Keep original line breaks
  trim_whitespace: true       # Remove excessive spacing
  normalize_headings: true    # Ensure proper hierarchy

spacing:
  paragraph_gap: 2            # Lines between paragraphs
  heading_gap: 2              # Lines around headings
  list_gap: 1                 # Lines around lists
  blockquote_gap: 2          # Lines around blockquotes
```

## Integration Points

### With Previous Components
- **Content Tree (F2)**: Receives structured content tree as input
- **Content Filter (F3)**: Works with pre-filtered, clean content
- **Media Handler (F4)**: Processes media replacements in tree

### Output Integration
- **CLI Interface**: Provides markdown output to stdout
- **File Output**: Supports writing to markdown files
- **API Integration**: Returns formatted markdown for external use

### Quality Assurance
- **Markdown Validation**: Ensures output is valid markdown
- **Lint Integration**: Supports markdown linting tools
- **Preview Generation**: Can generate HTML preview from markdown

## Performance Considerations

### Efficient Rendering
- Single-pass tree traversal for rendering
- Streaming output for large documents
- Memory-efficient string building
- Configurable rendering depth limits

### Output Quality
- Consistent formatting across all content types
- Proper escaping of special characters
- Handling of edge cases and malformed input
- Graceful degradation for unsupported elements

### Extensibility
- Plugin system for custom renderers
- Configurable output styles and formats
- Support for markdown extensions (tables, footnotes)
- Custom element handling for specific use cases