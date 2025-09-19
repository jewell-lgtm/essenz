# Essenz (eßenz) Product Specification

## Overview

Essenz is a CLI web browser that distills web pages into their semantic essence, reordering content by importance and presenting it as clean Markdown. Built with Go and Bubble Tea, it leverages headless Chrome for JavaScript rendering and uses executable specifications to drive development.

## Core Philosophy

### Executable Specifications
Every feature in essenz is defined by executable specifications that:
- Serve as both documentation and tests
- Validate implementation against requirements
- Provide living examples of usage
- Ensure behavior matches intent

### Content Distillation
- Extract the semantic essence of web pages
- Reorder content by importance, not DOM order
- Support modern JavaScript-heavy sites
- Present clean, readable Markdown output

## Architecture

### Component Stack
```
┌─────────────────────────────────────────┐
│          Bubble Tea TUI                 │
├─────────────────────────────────────────┤
│         Command Interface               │
├─────────────────────────────────────────┤
│     Content Analysis Engine             │
├─────────────────────────────────────────┤
│   Headless Chrome + go-readability      │
├─────────────────────────────────────────┤
│        HTTP Client + Cache              │
└─────────────────────────────────────────┘
```

### Key Components

1. **Browser Engine**: Headless Chrome via CDP (Chrome DevTools Protocol)
2. **Content Extractor**: Modified readability algorithm with importance scoring
3. **TUI Interface**: Bubble Tea for interactive browsing
4. **Spec Runner**: Executable specification framework
5. **Markdown Renderer**: Semantic-preserving markdown conversion

## Executable Specifications

### Spec Format
Specifications are written in Markdown with embedded executable blocks:

```markdown
## SPEC: JavaScript Content Extraction

GIVEN a page with JavaScript-rendered content
WHEN essenz fetches the page
THEN it waits for content to render
AND extracts the fully rendered DOM
AND produces correct markdown output

### Example
```spec
url: https://example.com/spa
wait_for: .main-content
timeout: 5s
expected_output: |
  # Article Title
  Main content text...
```
```

### Spec Execution
Each spec block compiles to a test that:
1. Sets up the browser with specified conditions
2. Navigates to the URL
3. Waits for render completion
4. Extracts and analyzes content
5. Validates against expected output

## Feature Specifications

### Feature 1: JavaScript-Aware Content Extraction

#### Specification
```markdown
GIVEN a single-page application
WHEN the page is loaded
THEN essenz waits for JavaScript execution
AND detects when content is ready
AND extracts the rendered DOM
```

#### Implementation
- Use Chrome DevTools Protocol (CDP)
- Implement smart wait strategies:
  - Wait for specific selectors
  - Wait for network idle
  - Wait for DOM stability
  - Configurable timeout

### Feature 2: Importance-Based Content Reordering

#### Specification
```markdown
GIVEN extracted page content
WHEN scoring importance
THEN headlines score higher than body text
AND longer coherent text scores higher than fragments
AND main content scores higher than navigation
AND content is reordered by descending importance
```

#### Scoring Algorithm
```go
type ContentBlock struct {
    Text         string
    Tag          string
    WordCount    int
    LinkDensity  float64
    Position     int
    Depth        int
    Score        float64
}

// Scoring weights (configurable)
weights := ScoreWeights{
    TagWeight:    map[string]float64{
        "h1": 3.0, "h2": 2.0, "h3": 1.5,
        "p": 1.0, "nav": 0.2, "footer": 0.1,
    },
    LengthBonus:  0.3,
    PositionBonus: 0.2,
    LinkPenalty:  -0.5,
}
```

### Feature 3: Interactive TUI Browser

#### Specification
```markdown
GIVEN the TUI is running
WHEN user enters a URL
THEN content loads with progress indicator
AND rendered markdown is displayed
AND user can navigate links
AND history is maintained
```

#### TUI Components (Bubble Tea)
- URL input bar
- Loading spinner
- Content viewport (scrollable)
- Link highlighting and selection
- History navigation (back/forward)
- Bookmark management

### Feature 4: Smart Waiting Strategies

#### Specification
```markdown
GIVEN different types of web pages
WHEN determining render completion
THEN static pages return immediately
AND SPAs wait for framework signals
AND lazy-loaded content is detected
AND infinite scroll is handled gracefully
```

#### Wait Strategies
```go
type WaitStrategy interface {
    Wait(ctx context.Context, page *Page) error
}

type CompositeWait struct {
    Strategies []WaitStrategy
    Mode       WaitMode // All, Any, Smart
}

// Built-in strategies
- WaitForSelector(selector string)
- WaitForNetworkIdle(duration time.Duration)
- WaitForDOMStable(checks int, delay time.Duration)
- WaitForCustomJS(script string)
```

### Feature 5: Caching and Offline Support

#### Specification
```markdown
GIVEN previously visited pages
WHEN accessing with poor connectivity
THEN cached content is available
AND cache respects TTL headers
AND user can force refresh
```

## Technical Implementation

### Dependencies
```go
// go.mod
module github.com/essenz/essenz

require (
    github.com/chromedp/chromedp      // Chrome automation
    github.com/charmbracelet/bubbletea // TUI framework
    github.com/go-shiori/go-readability // Content extraction
    github.com/PuerkitoBio/goquery    // HTML parsing
    github.com/yuin/goldmark          // Markdown rendering
)
```

### Project Structure
```
essenz/
├── cmd/
│   └── essenz/
│       └── main.go
├── internal/
│   ├── browser/
│   │   ├── chrome.go      # Chrome CDP integration
│   │   ├── strategies.go  # Wait strategies
│   │   └── cache.go       # Content caching
│   ├── extractor/
│   │   ├── readability.go # Content extraction
│   │   ├── scorer.go      # Importance scoring
│   │   └── blocks.go      # Block segmentation
│   ├── renderer/
│   │   ├── markdown.go    # Markdown conversion
│   │   └── metadata.go    # Frontmatter generation
│   └── tui/
│       ├── app.go         # Main TUI application
│       ├── views/         # TUI views
│       └── components/    # Reusable components
├── specs/
│   ├── runner/
│   │   ├── parser.go      # Spec parser
│   │   ├── executor.go    # Spec executor
│   │   └── reporter.go    # Results reporter
│   └── features/          # Feature specs (*.spec.md)
├── test/
│   ├── fixtures/          # Test HTML pages
│   └── integration/       # Integration tests
└── README.md
```

### Chrome Integration

```go
// Browser initialization
type Browser struct {
    ctx    context.Context
    cancel context.CancelFunc
}

func NewBrowser(opts ...Option) (*Browser, error) {
    opts = append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
        chromedp.Flag("disable-gpu", true),
        chromedp.Flag("no-sandbox", true),
        chromedp.UserAgent("essenz/1.0"),
    )

    allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
    ctx, _ := chromedp.NewContext(allocCtx)

    return &Browser{ctx: ctx}, nil
}

// Content extraction
func (b *Browser) Extract(url string, wait WaitStrategy) (*Content, error) {
    var html string

    err := chromedp.Run(b.ctx,
        chromedp.Navigate(url),
        wait.Wait(b.ctx),
        chromedp.OuterHTML("html", &html),
    )

    return parseAndScore(html)
}
```

## CLI Interface

### Basic Usage
```bash
# Simple extraction
sz https://example.com

# With custom wait strategy
sz --wait-for=".article-content" https://example.com

# Interactive TUI mode
sz --tui

# Output to file
sz https://example.com > article.md

# With debugging
sz --debug https://example.com
```

### Command Flags
```
Usage: sz [flags] <url>

Flags:
  --tui                 Launch interactive TUI
  --wait-for string     CSS selector to wait for
  --wait-idle duration  Wait for network idle (default 2s)
  --timeout duration    Maximum wait time (default 30s)
  --no-js              Disable JavaScript execution
  --cache-dir string    Cache directory (default ~/.essenz/cache)
  --debug              Show debug information
  --format string      Output format: markdown, json, html (default markdown)
  --top-k int          Number of top blocks to output (default 20)
  --summarize          Generate summary (default true)
```

## Output Format

### Markdown Output Structure
```markdown
# [Page Title]

*By [Author] — [Date]*

> [TL;DR Summary - top extracted sentences]

## [Most Important Section]
[Content reordered by importance...]

## [Second Most Important Section]
[Content...]

---

### Metadata
```yaml
url: https://example.com/article
title: Article Title
author: Author Name
published: 2024-01-01T00:00:00Z
fetched: 2024-01-02T00:00:00Z
render_time: 1.23s
wait_strategy: dom_stable
javascript: enabled
links:
  - text: Related Article
    url: https://example.com/related
images:
  - url: https://example.com/image.jpg
    alt: Description
```
```

## Testing Strategy

### Executable Spec Tests
All features are tested via their specifications:

```bash
# Run all specs
go test ./specs/...

# Run specific feature specs
go test ./specs/features/javascript_test.go

# Generate spec report
go test ./specs/... -json | go run ./specs/reporter
```

### Integration Tests
Test against real websites and cached fixtures:

```go
func TestJavaScriptSite(t *testing.T) {
    spec := LoadSpec("specs/features/javascript.spec.md")
    result := spec.Execute()
    assert.NoError(t, result.Error)
    assert.Contains(t, result.Output, spec.Expected)
}
```

### Performance Benchmarks
```go
func BenchmarkExtraction(b *testing.B) {
    browser := NewBrowser()
    defer browser.Close()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        browser.Extract(testURL, WaitForDOMStable())
    }
}
```

## Development Workflow

### Setting Up
```bash
# Clone repository
git clone https://github.com/essenz/essenz
cd essenz

# Install Chrome/Chromium
# macOS: brew install chromium
# Linux: apt-get install chromium-browser

# Install Go dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o sz cmd/essenz/main.go
```

### Writing Specs
1. Create a new spec file in `specs/features/`
2. Write the specification in Markdown with embedded test blocks
3. Run `go generate` to create test scaffolding
4. Implement feature to make specs pass

### Contributing
- All features must have executable specifications
- Specs are written before implementation
- Tests must pass before merge
- Performance benchmarks for new extractors

## Performance Considerations

### Optimization Strategies
- Reuse Chrome instances for multiple extractions
- Implement connection pooling
- Cache DOM analysis results
- Lazy load content blocks
- Stream markdown output

### Resource Management
```go
// Chrome instance pool
type BrowserPool struct {
    instances chan *Browser
    max       int
}

func (p *BrowserPool) Get(ctx context.Context) (*Browser, error) {
    select {
    case browser := <-p.instances:
        return browser, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        return NewBrowser()
    }
}
```

## Future Enhancements

### Planned Features
1. **Multi-tab support** in TUI mode
2. **Reader mode themes** (dark, sepia, etc.)
3. **Export formats** (PDF, EPUB, plain text)
4. **Bookmarklet** for browser integration
5. **API server mode** for programmatic access
6. **Plugin system** for custom extractors
7. **ML-based importance scoring** (optional)

### Experimental Ideas
- Voice narration of extracted content
- Automatic translation
- Content diffing between visits
- Social sharing of distilled articles
- Collaborative annotations

## License

MIT License - See LICENSE file

## Brand Identity

### Visual Identity
- **Logo**: The ß character as the central brand element
- **Name**: essenz (executable: `sz`)
- **Tagline**: "Distill the web"
- **Colors**: Minimal, terminal-friendly

### ASCII Banner
```
   ____  ____  ____  ____  ____  ____
  | ___|| ___|| ___|| ___||  _ \|_  /
  | |_  | |_  | |_  | |_  | | | |/ /
  |  _| |___ ||___ ||___ || |_| / /_
  |____||____||____||____||____/____|
              ß
        distill the web
```