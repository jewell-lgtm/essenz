# F3: Content Filter System

**Feature Branch**: `feature/content-filter-system`

## Objective

Implement a sophisticated filtering system that removes non-content elements (navigation, headers, footers, ads, etc.) from the content tree while preserving the main article content and meaningful information.

## Technical Requirements

### Filtering Strategy
- **Semantic Tag Filtering**: Remove content based on HTML5 semantic tags
- **CSS Class/ID Analysis**: Filter based on common patterns in class names and IDs
- **Content Density Analysis**: Remove sections with high link-to-text ratios
- **Position-Based Filtering**: Remove header/footer content based on document position
- **Whitelist Approach**: Preserve explicitly identified content areas

### Filter Rule Engine
```go
type FilterRule interface {
    ShouldExclude(node *ContentNode, context *FilterContext) bool
    Priority() int
    Name() string
}

type FilterContext struct {
    DocumentRoot    *ContentNode
    CurrentDepth    int
    ParentNodes     []*ContentNode
    SiblingNodes    []*ContentNode
    DocumentStats   *DocumentStats
}

type DocumentStats struct {
    TotalTextLength int
    AverageWordCount float64
    LinkDensity     float64
    HeadingCount    int
    ParagraphCount  int
}
```

### Built-in Filter Rules
1. **SemanticTagFilter**: Removes `nav`, `header`, `footer`, `aside` content
2. **ClassNameFilter**: Filters based on common CSS patterns
3. **LinkDensityFilter**: Removes high-link-density sections
4. **LengthFilter**: Removes very short text blocks
5. **PositionFilter**: Removes top/bottom page elements
6. **DuplicateContentFilter**: Removes repeated navigation elements

## Implementation Components

### 1. Internal Package: `internal/filter`
```go
type ContentFilter struct {
    rules   []FilterRule
    config  FilterConfig
}

type FilterConfig struct {
    MaxLinkDensity      float64  // 0.3 = 30% links
    MinContentLength    int      // Minimum characters for content blocks
    PreserveWhitelist   []string // CSS selectors to always preserve
    AggressiveMode      bool     // More strict filtering
    DebugMode          bool     // Log filtering decisions
}

func (cf *ContentFilter) FilterTree(root *ContentNode) (*ContentNode, error)
func (cf *ContentFilter) AddRule(rule FilterRule)
func (cf *ContentFilter) GetFilterStats() *FilterStats
```

### 2. Semantic Tag Filtering
```go
type SemanticTagFilter struct {
    excludedTags []string
}

// Default excluded tags: nav, header, footer, aside, script, style, noscript
func (f *SemanticTagFilter) ShouldExclude(node *ContentNode, ctx *FilterContext) bool
```

### 3. CSS Pattern Filtering
```go
type ClassNameFilter struct {
    patterns []string
}

// Common patterns: nav, menu, sidebar, footer, header, ad, advertisement,
// social, share, comments, related, breadcrumb, pagination
func (f *ClassNameFilter) ShouldExclude(node *ContentNode, ctx *FilterContext) bool
```

### 4. Content Density Analysis
```go
type LinkDensityFilter struct {
    maxDensity float64
    minWords   int
}

func (f *LinkDensityFilter) calculateLinkDensity(node *ContentNode) float64
func (f *LinkDensityFilter) ShouldExclude(node *ContentNode, ctx *FilterContext) bool
```

## Acceptance Criteria

### Semantic Filtering
1. Removes all `<nav>` elements and their content
2. Filters out `<header>` and `<footer>` sections
3. Excludes `<aside>` sidebar content
4. Preserves `<main>` and `<article>` content

### CSS-Based Filtering
1. Removes elements with navigation-related class names
2. Filters out advertisement containers
3. Excludes social media sharing widgets
4. Removes comment sections and related content links

### Content Quality Filtering
1. Removes sections with >30% link density
2. Filters out very short text blocks (<50 characters)
3. Removes duplicate navigation menus
4. Preserves main article content regardless of container classes

### Preservation Rules
1. Never removes content from whitelisted selectors
2. Preserves content within `<main>` and `<article>` tags
3. Maintains content that passes length and density thresholds
4. Keeps structured content (lists, tables) even in filtered sections

## Test Scenarios

### Basic Site Structure
```html
<header>
    <nav>Site Navigation</nav>
</header>
<main>
    <article>
        <h1>Article Title</h1>
        <p>Main content paragraph.</p>
    </article>
</main>
<aside>
    <div class="ads">Advertisement</div>
    <div class="related">Related Links</div>
</aside>
<footer>
    <div class="social">Social Media</div>
</footer>
```

Expected: Only article title and content preserved.

### Modern Framework Structure
- React components with CSS modules
- Next.js layout with complex class names
- Vue.js components with scoped styles
- CSS-in-JS generated class names

### Complex Navigation Patterns
- Mega menus with multiple levels
- Breadcrumb navigation
- In-page navigation (table of contents)
- Pagination controls

### Edge Cases
- Content that looks like navigation but is actually article content
- Main content with high link density (like a reference list)
- Short but meaningful content blocks
- Mixed content areas

## Filter Rule Configuration

### Default Aggressive Rules
```yaml
semantic_tags:
  exclude: [nav, header, footer, aside, script, style, noscript]

css_patterns:
  exclude: [nav, menu, sidebar, footer, header, ad, advertisement,
           social, share, comments, related, breadcrumb, pagination]

content_quality:
  max_link_density: 0.3
  min_content_length: 50
  min_word_count: 10

position_rules:
  remove_top_percent: 5    # Remove top 5% of page
  remove_bottom_percent: 5 # Remove bottom 5% of page
```

### Whitelist Protection
```yaml
whitelist:
  selectors: [main, article, .content, .post, .entry]
  tags: [main, article]
  always_preserve: true
```

## Integration Points

### With Text Node Tree (F2)
- Operates on structured content tree from F2
- Preserves tree hierarchy while removing nodes
- Maintains parent-child relationships after filtering

### With Image Handling (F4)
- Coordinates with image filtering decisions
- Preserves images within main content areas
- Removes decorative images from filtered sections

### With Markdown Rendering (F5)
- Provides clean content tree for markdown conversion
- Ensures only meaningful content reaches renderer
- Maintains structure needed for proper markdown formatting

## Performance Considerations

### Efficient Filtering
- Single-pass tree traversal where possible
- Lazy evaluation of expensive filter rules
- Caching of filter decisions for similar nodes
- Configurable rule priority ordering

### Memory Management
- In-place tree modification to reduce memory usage
- Proper cleanup of filtered nodes
- Streaming processing for very large documents