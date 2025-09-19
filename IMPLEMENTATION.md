# Essenz Implementation Guide

## Overview

This guide provides a detailed implementation roadmap for essenz, focusing on the Bubble Tea TUI framework and executable specifications architecture.

## Phase 1: Core Foundation (Week 1)

### 1.1 Project Setup

```go
// go.mod
module github.com/essenz/essenz

go 1.21

require (
    github.com/chromedp/chromedp v0.9.3
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/go-shiori/go-readability v0.0.0-20230421032831-c66949dfc0ad
    github.com/PuerkitoBio/goquery v1.8.1
    github.com/yuin/goldmark v1.6.0
    github.com/stretchr/testify v1.8.4
)
```

### 1.2 Core Interfaces

```go
// internal/core/interfaces.go
package core

import (
    "context"
    "time"
)

type Browser interface {
    Navigate(ctx context.Context, url string) error
    WaitFor(ctx context.Context, strategy WaitStrategy) error
    GetHTML(ctx context.Context) (string, error)
    Close() error
}

type Extractor interface {
    Extract(html string) (*Content, error)
    Score(content *Content) error
}

type Renderer interface {
    ToMarkdown(content *Content) (string, error)
    ToJSON(content *Content) ([]byte, error)
}

type Content struct {
    URL       string
    Title     string
    Author    string
    Published time.Time
    Blocks    []Block
    Metadata  map[string]interface{}
}

type Block struct {
    ID          string
    Tag         string
    Text        string
    HTML        string
    Score       float64
    WordCount   int
    LinkDensity float64
    Position    int
}
```

## Phase 2: Chrome Integration (Week 1-2)

### 2.1 Browser Implementation

```go
// internal/browser/chrome.go
package browser

import (
    "context"
    "time"

    "github.com/chromedp/chromedp"
    "github.com/essenz/essenz/internal/core"
)

type ChromeBrowser struct {
    ctx    context.Context
    cancel context.CancelFunc
}

func New(opts ...Option) (*ChromeBrowser, error) {
    opts = append(chromedp.DefaultExecAllocatorOptions[:],
        chromedp.Flag("headless", true),
        chromedp.Flag("disable-gpu", true),
        chromedp.Flag("no-sandbox", true),
    )

    allocCtx, cancel := chromedp.NewExecAllocator(
        context.Background(),
        opts...,
    )

    ctx, _ := chromedp.NewContext(allocCtx)

    return &ChromeBrowser{
        ctx:    ctx,
        cancel: cancel,
    }, nil
}

func (b *ChromeBrowser) Navigate(ctx context.Context, url string) error {
    return chromedp.Run(b.ctx,
        chromedp.Navigate(url),
    )
}

func (b *ChromeBrowser) WaitFor(ctx context.Context, strategy core.WaitStrategy) error {
    return strategy.Execute(b.ctx)
}

func (b *ChromeBrowser) GetHTML(ctx context.Context) (string, error) {
    var html string
    err := chromedp.Run(b.ctx,
        chromedp.OuterHTML("html", &html),
    )
    return html, err
}
```

### 2.2 Wait Strategies

```go
// internal/browser/strategies.go
package browser

import (
    "context"
    "time"

    "github.com/chromedp/chromedp"
)

type WaitStrategy interface {
    Execute(ctx context.Context) error
}

type SelectorWait struct {
    Selector string
    Timeout  time.Duration
}

func (w SelectorWait) Execute(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, w.Timeout)
    defer cancel()

    return chromedp.Run(ctx,
        chromedp.WaitVisible(w.Selector),
    )
}

type NetworkIdleWait struct {
    Duration time.Duration
}

func (w NetworkIdleWait) Execute(ctx context.Context) error {
    // Monitor network activity
    return chromedp.Run(ctx,
        chromedp.ActionFunc(func(ctx context.Context) error {
            time.Sleep(w.Duration)
            return nil
        }),
    )
}

type DOMStableWait struct {
    Checks   int
    Interval time.Duration
}

func (w DOMStableWait) Execute(ctx context.Context) error {
    var lastHTML, currentHTML string
    stable := 0

    for stable < w.Checks {
        chromedp.Run(ctx,
            chromedp.OuterHTML("body", &currentHTML),
        )

        if currentHTML == lastHTML {
            stable++
        } else {
            stable = 0
        }

        lastHTML = currentHTML
        time.Sleep(w.Interval)
    }

    return nil
}
```

## Phase 3: Bubble Tea TUI (Week 2-3)

### 3.1 Main TUI Application

```go
// internal/tui/app.go
package tui

import (
    "fmt"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/essenz/essenz/internal/browser"
    "github.com/essenz/essenz/internal/extractor"
)

type Model struct {
    browser     *browser.ChromeBrowser
    extractor   *extractor.Extractor

    // UI State
    url         string
    loading     bool
    content     string
    error       error

    // Navigation
    history     []string
    historyIdx  int
    links       []string
    selectedLink int

    // Viewport
    width       int
    height      int
    scrollY     int
}

func New() Model {
    b, _ := browser.New()
    e := extractor.New()

    return Model{
        browser:   b,
        extractor: e,
        history:   []string{},
    }
}

func (m Model) Init() tea.Cmd {
    return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit

        case "ctrl+l":
            // Focus URL bar
            m.focusMode = "url"
            return m, nil

        case "enter":
            if m.focusMode == "url" {
                return m, m.loadURL(m.url)
            }
            if m.focusMode == "content" && m.selectedLink >= 0 {
                return m, m.loadURL(m.links[m.selectedLink])
            }

        case "tab":
            // Navigate links
            m.selectedLink = (m.selectedLink + 1) % len(m.links)
            return m, nil

        case "backspace":
            // Go back
            if m.historyIdx > 0 {
                m.historyIdx--
                return m, m.loadURL(m.history[m.historyIdx])
            }

        case "up", "k":
            m.scrollY = max(0, m.scrollY-1)

        case "down", "j":
            m.scrollY++
        }

    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    case loadCompleteMsg:
        m.loading = false
        m.content = msg.content
        m.links = msg.links
        m.error = msg.err
    }

    return m, nil
}

func (m Model) View() string {
    var s strings.Builder

    // URL Bar
    urlBar := m.renderURLBar()
    s.WriteString(urlBar)
    s.WriteString("\n")

    // Loading indicator
    if m.loading {
        s.WriteString(m.renderLoading())
        return s.String()
    }

    // Error display
    if m.error != nil {
        s.WriteString(m.renderError())
        return s.String()
    }

    // Content viewport
    content := m.renderContent()
    s.WriteString(content)

    // Status bar
    status := m.renderStatusBar()
    s.WriteString(status)

    return s.String()
}
```

### 3.2 TUI Components

```go
// internal/tui/components.go
package tui

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/lipgloss"
)

var (
    urlBarStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62")).
        Padding(0, 1)

    contentStyle = lipgloss.NewStyle().
        Padding(1, 2)

    selectedLinkStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("212")).
        Background(lipgloss.Color("63")).
        Bold(true)

    statusBarStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("235")).
        Foreground(lipgloss.Color("252")).
        Padding(0, 1)
)

func (m Model) renderURLBar() string {
    prefix := "🌐 "
    if m.loading {
        prefix = "⟳ "
    }

    url := m.url
    if m.focusMode == "url" {
        url += "█"
    }

    return urlBarStyle.Render(prefix + url)
}

func (m Model) renderLoading() string {
    frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
    frame := frames[m.loadingFrame%len(frames)]

    return lipgloss.NewStyle().
        Padding(2, 2).
        Render(fmt.Sprintf("%s Loading...", frame))
}

func (m Model) renderContent() string {
    lines := strings.Split(m.content, "\n")

    // Apply viewport scrolling
    start := m.scrollY
    end := min(start+m.height-4, len(lines))

    visibleLines := lines[start:end]

    // Highlight links
    for i, line := range visibleLines {
        if strings.Contains(line, "[") && strings.Contains(line, "]") {
            // Simple link detection
            for j, link := range m.links {
                if j == m.selectedLink && strings.Contains(line, link) {
                    visibleLines[i] = selectedLinkStyle.Render(line)
                    break
                }
            }
        }
    }

    return contentStyle.Render(strings.Join(visibleLines, "\n"))
}

func (m Model) renderStatusBar() string {
    left := fmt.Sprintf(" %d links | Line %d/%d",
        len(m.links),
        m.scrollY+1,
        m.contentLines,
    )

    right := "q: quit | ?: help "

    width := m.width - lipgloss.Width(left) - lipgloss.Width(right)
    middle := strings.Repeat(" ", max(0, width))

    return statusBarStyle.Render(left + middle + right)
}
```

### 3.3 Content Loading Commands

```go
// internal/tui/commands.go
package tui

import (
    "context"
    "time"

    tea "github.com/charmbracelet/bubbletea"
)

type loadCompleteMsg struct {
    content string
    links   []string
    err     error
}

func (m Model) loadURL(url string) tea.Cmd {
    return func() tea.Msg {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        // Navigate to URL
        if err := m.browser.Navigate(ctx, url); err != nil {
            return loadCompleteMsg{err: err}
        }

        // Wait for content
        strategy := browser.DOMStableWait{
            Checks:   5,
            Interval: 200 * time.Millisecond,
        }
        if err := m.browser.WaitFor(ctx, strategy); err != nil {
            return loadCompleteMsg{err: err}
        }

        // Get HTML
        html, err := m.browser.GetHTML(ctx)
        if err != nil {
            return loadCompleteMsg{err: err}
        }

        // Extract content
        content, err := m.extractor.Extract(html)
        if err != nil {
            return loadCompleteMsg{err: err}
        }

        // Render to markdown
        markdown, err := m.renderer.ToMarkdown(content)
        if err != nil {
            return loadCompleteMsg{err: err}
        }

        // Extract links
        links := extractLinks(content)

        return loadCompleteMsg{
            content: markdown,
            links:   links,
        }
    }
}
```

## Phase 4: Executable Specifications (Week 3-4)

### 4.1 Spec Parser

```go
// specs/runner/parser.go
package runner

import (
    "bufio"
    "fmt"
    "regexp"
    "strings"

    "gopkg.in/yaml.v3"
)

type Spec struct {
    Name        string
    Given       string
    When        string
    Then        string
    TestCases   []TestCase
}

type TestCase struct {
    URL              string                 `yaml:"url"`
    WaitFor          string                 `yaml:"wait_for"`
    Timeout          string                 `yaml:"timeout"`
    ExpectedContains []string               `yaml:"expected_contains"`
    ExpectedStructure []ExpectedBlock       `yaml:"expected_structure"`
}

type ExpectedBlock struct {
    Tag       string `yaml:"tag"`
    Content   string `yaml:"content"`
    MinLength int    `yaml:"min_length"`
}

func ParseSpecFile(path string) (*Spec, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    spec := &Spec{}

    // Parse markdown structure
    lines := strings.Split(string(content), "\n")
    inSpecBlock := false
    specContent := []string{}

    for _, line := range lines {
        if strings.HasPrefix(line, "## SPEC:") {
            spec.Name = strings.TrimPrefix(line, "## SPEC:")
            spec.Name = strings.TrimSpace(spec.Name)
        }

        if strings.HasPrefix(line, "GIVEN") {
            spec.Given = strings.TrimPrefix(line, "GIVEN ")
        }

        if strings.HasPrefix(line, "WHEN") {
            spec.When = strings.TrimPrefix(line, "WHEN ")
        }

        if strings.HasPrefix(line, "THEN") {
            spec.Then = strings.TrimPrefix(line, "THEN ")
        }

        if strings.Contains(line, "```spec") {
            inSpecBlock = true
            continue
        }

        if inSpecBlock {
            if strings.Contains(line, "```") {
                // Parse YAML spec block
                var testCase TestCase
                err := yaml.Unmarshal([]byte(strings.Join(specContent, "\n")), &testCase)
                if err != nil {
                    return nil, err
                }
                spec.TestCases = append(spec.TestCases, testCase)

                inSpecBlock = false
                specContent = []string{}
            } else {
                specContent = append(specContent, line)
            }
        }
    }

    return spec, nil
}
```

### 4.2 Spec Executor

```go
// specs/runner/executor.go
package runner

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/essenz/essenz/internal/browser"
    "github.com/essenz/essenz/internal/extractor"
    "github.com/stretchr/testify/assert"
)

type Executor struct {
    browser   *browser.ChromeBrowser
    extractor *extractor.Extractor
}

func NewExecutor() *Executor {
    b, _ := browser.New()
    e := extractor.New()

    return &Executor{
        browser:   b,
        extractor: e,
    }
}

func (e *Executor) Execute(spec *Spec) []TestResult {
    results := []TestResult{}

    for _, tc := range spec.TestCases {
        result := e.executeTestCase(spec, tc)
        results = append(results, result)
    }

    return results
}

func (e *Executor) executeTestCase(spec *Spec, tc TestCase) TestResult {
    result := TestResult{
        SpecName: spec.Name,
        TestCase: tc,
    }

    ctx := context.Background()
    if tc.Timeout != "" {
        d, _ := time.ParseDuration(tc.Timeout)
        ctx, _ = context.WithTimeout(ctx, d)
    }

    // Navigate to URL
    err := e.browser.Navigate(ctx, tc.URL)
    if err != nil {
        result.Error = err
        return result
    }

    // Apply wait strategy
    if tc.WaitFor != "" {
        strategy := browser.SelectorWait{
            Selector: tc.WaitFor,
            Timeout:  10 * time.Second,
        }
        err = e.browser.WaitFor(ctx, strategy)
        if err != nil {
            result.Error = err
            return result
        }
    }

    // Get HTML and extract
    html, err := e.browser.GetHTML(ctx)
    if err != nil {
        result.Error = err
        return result
    }

    content, err := e.extractor.Extract(html)
    if err != nil {
        result.Error = err
        return result
    }

    // Validate expectations
    result.Passed = e.validateExpectations(content, tc)

    return result
}

func (e *Executor) validateExpectations(content *Content, tc TestCase) bool {
    markdown := content.ToMarkdown()

    // Check expected content
    for _, expected := range tc.ExpectedContains {
        if !strings.Contains(markdown, expected) {
            return false
        }
    }

    // Check structure
    for _, expBlock := range tc.ExpectedStructure {
        found := false
        for _, block := range content.Blocks {
            if block.Tag == expBlock.Tag {
                if expBlock.Content != "" && !strings.Contains(block.Text, expBlock.Content) {
                    continue
                }
                if expBlock.MinLength > 0 && len(block.Text) < expBlock.MinLength {
                    continue
                }
                found = true
                break
            }
        }
        if !found {
            return false
        }
    }

    return true
}
```

## Phase 5: Content Extraction & Scoring (Week 4-5)

### 5.1 Content Extractor

```go
// internal/extractor/extractor.go
package extractor

import (
    "strings"

    "github.com/PuerkitoBio/goquery"
    "github.com/go-shiori/go-readability"
    "github.com/essenz/essenz/internal/core"
)

type Extractor struct {
    scorer *Scorer
}

func New() *Extractor {
    return &Extractor{
        scorer: NewScorer(),
    }
}

func (e *Extractor) Extract(html string) (*core.Content, error) {
    // Try readability first
    article, err := readability.FromReader(strings.NewReader(html), "")
    if err == nil && article.Content != "" {
        return e.extractFromReadability(article)
    }

    // Fallback to custom extraction
    return e.extractFromHTML(html)
}

func (e *Extractor) extractFromHTML(html string) (*core.Content, error) {
    doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
    if err != nil {
        return nil, err
    }

    content := &core.Content{
        Blocks: []core.Block{},
    }

    // Extract title
    content.Title = doc.Find("title").First().Text()

    // Extract blocks
    doc.Find("h1, h2, h3, p, li, blockquote, pre").Each(func(i int, s *goquery.Selection) {
        block := core.Block{
            ID:       fmt.Sprintf("block-%d", i),
            Tag:      goquery.NodeName(s),
            Text:     strings.TrimSpace(s.Text()),
            HTML:     getOuterHTML(s),
            Position: i,
        }

        // Calculate features
        block.WordCount = len(strings.Fields(block.Text))
        block.LinkDensity = calculateLinkDensity(s)

        content.Blocks = append(content.Blocks, block)
    })

    // Score blocks
    e.scorer.Score(content)

    // Sort by score
    sort.Slice(content.Blocks, func(i, j int) bool {
        return content.Blocks[i].Score > content.Blocks[j].Score
    })

    return content, nil
}
```

### 5.2 Scoring Engine

```go
// internal/extractor/scorer.go
package extractor

import (
    "math"
    "strings"

    "github.com/essenz/essenz/internal/core"
)

type Scorer struct {
    weights ScoreWeights
}

type ScoreWeights struct {
    TagWeights map[string]float64
    LengthBonus float64
    PositionBonus float64
    LinkPenalty float64
    EmphasisBonus float64
}

func NewScorer() *Scorer {
    return &Scorer{
        weights: ScoreWeights{
            TagWeights: map[string]float64{
                "h1": 3.0,
                "h2": 2.0,
                "h3": 1.5,
                "p": 1.0,
                "blockquote": 1.4,
                "pre": 1.5,
                "li": 0.9,
                "nav": 0.2,
                "footer": 0.1,
            },
            LengthBonus: 0.3,
            PositionBonus: 0.2,
            LinkPenalty: 0.5,
            EmphasisBonus: 0.1,
        },
    }
}

func (s *Scorer) Score(content *core.Content) {
    totalBlocks := len(content.Blocks)

    for i := range content.Blocks {
        block := &content.Blocks[i]

        // Base tag weight
        tagWeight := s.weights.TagWeights[block.Tag]
        if tagWeight == 0 {
            tagWeight = 0.5
        }

        // Length score (sigmoid function)
        lengthScore := s.calculateLengthScore(block.WordCount)

        // Position bonus (earlier = better)
        positionBonus := 1.0 - (float64(block.Position) / float64(totalBlocks)) * s.weights.PositionBonus

        // Link penalty
        linkPenalty := math.Min(block.LinkDensity * s.weights.LinkPenalty, 1.0)

        // Calculate final score
        block.Score = tagWeight * lengthScore * positionBonus * (1.0 - linkPenalty)
    }
}

func (s *Scorer) calculateLengthScore(wordCount int) float64 {
    // Sigmoid function centered around optimal length (100-300 words)
    optimal := 200.0
    x := float64(wordCount)

    if x < 10 {
        return 0.1
    }
    if x > 1000 {
        return 0.7
    }

    // Smooth curve peaking around optimal
    return 1.0 / (1.0 + math.Exp(-0.01*(x-optimal)))
}
```

## Phase 6: Testing & CI/CD (Week 5-6)

### 6.1 Test Structure

```go
// test/integration/extraction_test.go
package integration

import (
    "testing"

    "github.com/essenz/essenz/specs/runner"
    "github.com/stretchr/testify/assert"
)

func TestJavaScriptExtraction(t *testing.T) {
    spec, err := runner.ParseSpecFile("../../specs/features/javascript-rendering.spec.md")
    assert.NoError(t, err)

    executor := runner.NewExecutor()
    defer executor.Close()

    results := executor.Execute(spec)

    for _, result := range results {
        assert.True(t, result.Passed, "Test case failed: %v", result.TestCase.URL)
        assert.NoError(t, result.Error)
    }
}

func TestScoring(t *testing.T) {
    // Test scoring algorithm
    scorer := extractor.NewScorer()

    blocks := []core.Block{
        {Tag: "h1", Text: "Title", WordCount: 2},
        {Tag: "p", Text: strings.Repeat("word ", 100), WordCount: 100},
        {Tag: "nav", Text: "Navigation", WordCount: 1},
    }

    content := &core.Content{Blocks: blocks}
    scorer.Score(content)

    // H1 should score highest
    assert.Greater(t, content.Blocks[0].Score, content.Blocks[2].Score)
}
```

### 6.2 Benchmarks

```go
// test/benchmark/performance_test.go
package benchmark

import (
    "testing"

    "github.com/essenz/essenz/internal/browser"
    "github.com/essenz/essenz/internal/extractor"
)

func BenchmarkExtraction(b *testing.B) {
    browser := browser.New()
    defer browser.Close()

    extractor := extractor.New()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        ctx := context.Background()
        browser.Navigate(ctx, "https://example.com")
        html, _ := browser.GetHTML(ctx)
        extractor.Extract(html)
    }
}

func BenchmarkScoring(b *testing.B) {
    // Create sample content
    content := generateSampleContent(100) // 100 blocks

    scorer := extractor.NewScorer()

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        scorer.Score(content)
    }
}
```

## Deployment & Distribution

### GitHub Actions CI/CD

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install Chrome
      run: |
        wget -q -O - https://dl.google.com/linux/linux_signing_key.pub | sudo apt-key add -
        sudo sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list'
        sudo apt-get update
        sudo apt-get install -y google-chrome-stable

    - name: Run tests
      run: go test ./...

    - name: Run benchmarks
      run: go test -bench=. ./test/benchmark

    - name: Build
      run: go build -o sz cmd/essenz/main.go

  release:
    needs: test
    if: startsWith(github.ref, 'refs/tags/')
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Build releases
      run: |
        GOOS=linux GOARCH=amd64 go build -o sz-linux-amd64 cmd/essenz/main.go
        GOOS=darwin GOARCH=amd64 go build -o sz-darwin-amd64 cmd/essenz/main.go
        GOOS=darwin GOARCH=arm64 go build -o sz-darwin-arm64 cmd/essenz/main.go
        GOOS=windows GOARCH=amd64 go build -o sz-windows-amd64.exe cmd/essenz/main.go

    - name: Create Release
      uses: softprops/action-gh-release@v1
      with:
        files: sz-*
```

## Next Steps

1. **Implement core browser integration** with Chrome
2. **Build TUI with Bubble Tea** for interactive mode
3. **Create spec parser and executor** for executable specifications
4. **Develop extraction and scoring algorithms**
5. **Write comprehensive test suite**
6. **Set up CI/CD pipeline**
7. **Create documentation and examples**
8. **Release v0.1.0** with core features

This implementation guide provides a solid foundation for building essenz with executable specifications and a beautiful TUI experience.