# F2: Text Node Tree Builder

**Feature Branch**: `feature/text-node-tree-builder`

## Objective

Build a content tree structure starting from actual text nodes in the DOM, creating a bottom-up representation that preserves document structure while focusing on actual readable content.

## Technical Requirements

### Text Node Discovery
- Traverse the DOM to find all text nodes containing meaningful content
- Filter out whitespace-only text nodes and script/style content
- Identify text nodes within semantic containers (paragraphs, headers, lists)
- Preserve the hierarchical relationship between text nodes and their parent elements

### Tree Structure Design
```go
type ContentNode struct {
    Type        NodeType        // TEXT, HEADING, PARAGRAPH, LIST, etc.
    Content     string          // Actual text content
    Level       int            // For headings (1-6), list nesting, etc.
    Parent      *ContentNode   // Parent node in tree
    Children    []*ContentNode // Child nodes
    Attributes  map[string]string // Relevant attributes (href for links, alt for images)
    SourceTag   string         // Original HTML tag name
    Position    int            // Order within parent
}

type NodeType int
const (
    TEXT NodeType = iota
    HEADING
    PARAGRAPH
    LIST
    LIST_ITEM
    LINK
    EMPHASIS
    STRONG
    BLOCKQUOTE
    CODE
    IMAGE_PLACEHOLDER
)
```

### Content Tree Building Algorithm
1. **Text Node Collection**: Gather all meaningful text nodes from DOM
2. **Parent Element Analysis**: Analyze the semantic meaning of parent elements
3. **Tree Construction**: Build hierarchical structure preserving document flow
4. **Content Consolidation**: Merge adjacent text nodes within same semantic context
5. **Structure Validation**: Ensure tree represents logical document structure

## Implementation Components

### 1. Internal Package: `internal/tree`
```go
type TreeBuilder struct {
    filterRules []FilterRule
    config      TreeConfig
}

type TreeConfig struct {
    MinTextLength       int
    PreserveWhitespace  bool
    IncludeHiddenText   bool
    MaxDepth           int
}

func (tb *TreeBuilder) BuildFromDOM(doc *html.Node) (*ContentNode, error)
func (tb *TreeBuilder) BuildFromChrome(ctx context.Context, page chromedp.Page) (*ContentNode, error)
```

### 2. Chrome Integration
- Use chromedp to inject JavaScript for text node discovery
- Extract text nodes with their computed styles and visibility
- Capture semantic context from parent elements
- Handle dynamic content that appears after page load

### 3. Content Analysis
- Detect heading levels and hierarchy
- Identify list structures and nesting
- Recognize emphasis and formatting elements
- Handle inline links and maintain context

## Acceptance Criteria

### Text Node Discovery
1. Successfully finds all visible text nodes in a document
2. Filters out script, style, and comment nodes
3. Identifies text within hidden elements (display: none)
4. Handles dynamically generated content

### Tree Structure
1. Maintains proper parent-child relationships
2. Preserves document reading order
3. Correctly identifies heading hierarchy (h1-h6)
4. Handles nested list structures accurately

### Content Preservation
1. Preserves original text content without modification
2. Maintains formatting indicators (bold, italic, links)
3. Captures image alt text and link destinations
4. Handles special characters and Unicode correctly

## Test Scenarios

### Basic HTML Structure
```html
<!DOCTYPE html>
<html>
<body>
    <h1>Main Title</h1>
    <p>Introduction paragraph with <strong>emphasis</strong> and a <a href="link">link</a>.</p>
    <h2>Section Header</h2>
    <ul>
        <li>First item</li>
        <li>Second item with <em>emphasis</em></li>
    </ul>
</body>
</html>
```

Expected tree structure with proper hierarchy and content preservation.

### Complex Nested Structure
- Multiple heading levels with proper nesting
- Lists within lists (nested ul/ol)
- Blockquotes containing multiple paragraphs
- Tables with meaningful data structure

### Modern Framework Output
- React components with complex DOM structure
- Vue.js rendered content with dynamic elements
- Next.js hydrated content
- CSS-in-JS styled components

### Edge Cases
- Empty elements and whitespace handling
- Mixed content (text and inline elements)
- Malformed HTML structure
- Very deep nesting levels

## Content Type Handling

### Text Elements
- **Paragraphs**: Group text nodes within `<p>` tags
- **Headings**: Identify level and hierarchical position
- **Lists**: Preserve item structure and nesting
- **Inline Formatting**: Maintain emphasis, strong, code markers

### Structural Elements
- **Sections**: Group related content logically
- **Articles**: Identify main content boundaries
- **Blockquotes**: Preserve citation structure
- **Code Blocks**: Maintain formatting and language hints

### Interactive Elements
- **Links**: Capture destination URLs and anchor text
- **Buttons**: Include meaningful button text
- **Form Labels**: Associate labels with form context

## Integration Points

### With DOM Ready Events (F1)
- Waits for page readiness before tree building
- Ensures all dynamic content is loaded
- Handles framework-specific initialization

### With Content Filtering (F3)
- Provides structured input for filtering algorithms
- Enables semantic-aware content exclusion
- Supports whitelist/blacklist approaches

### With Future Components
- Foundation for markdown rendering (F5)
- Input for image/media handling (F4)
- Base structure for content analysis and scoring