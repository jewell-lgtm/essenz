```
███████╗███████╗███████╗███████╗███╗   ██╗███████╗
██╔════╝██╔════╝╚══███╔╝██╔════╝████╗  ██║╚══███╔╝
█████╗  ███████╗  ███╔╝ █████╗  ██╔██╗ ██║  ███╔╝
██╔══╝  ╚════██║ ███╔╝  ██╔══╝  ██║╚██╗██║ ███╔╝
███████╗███████║███████╗███████╗██║ ╚████║███████╗
╚══════╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═══╝╚══════╝
                    ß - eszett
```

<div align="center">
  <h3>The web is bloated. Your reading shouldn't be.</h3>
  <p><strong>sz</strong> cuts through the noise, extracting what matters.</p>
  <br/>
</div>

---

## 🚀 See it in action

```bash
# Transform any article into clean, semantic markdown
$ sz https://longform.aeon.co/essays/why-efficiency-is-dangerous
```
```markdown
# Why Efficiency Is Dangerous and Slowing Down Makes Life Better

*By Barry Schwartz — Published April 2024*

> TL;DR: Our obsession with efficiency is backfiring. The very things we
> do to save time end up consuming more of it. Friction and inefficiency
> aren't bugs to eliminate—they're features that make life worth living.

[Article's most important insights appear first, less critical details follow...]
```

```bash
# Skip the ads, popups, and cookie banners - just get the content
$ sz https://news.ycombinator.com | head -20
```
```markdown
# Hacker News

## FTC bans hidden junk fees in hotel and ticket prices (337 comments)
The Federal Trade Commission finalized a rule requiring businesses to display
the total price upfront, ending deceptive pricing practices...

## Show HN: I built a local-first markdown editor with CRDT sync (89 comments)
After struggling with cloud-based note apps, I created an editor that works
offline-first with peer-to-peer sync when you're ready...
```

## ⚡ Why developers love sz

**🎯 It just works** — Handles SPAs, paywalls, and JavaScript-heavy sites that break other scrapers
**📖 Readable output** — Not just extracted, but intelligently reordered by importance
**⚙️ Unix philosophy** — Pipes beautifully with your existing workflow
**🔒 Privacy first** — Runs locally, no data leaves your machine
**🎨 Terminal native** — Built for developers who live in the terminal

## Installation

```bash
# Quick install with Go
go install github.com/jewell-lgtm/essenz/cmd/essenz@latest

# Or grab the binary
curl -L https://github.com/jewell-lgtm/essenz/releases/latest/download/sz-$(uname -s)-$(uname -m) -o sz
chmod +x sz && sudo mv sz /usr/local/bin/
```

**Requirements:** Chrome/Chromium installed, Go 1.21+ (for source builds)

## Core Features

### 🎯 Intelligent Extraction
```bash
# Basic extraction
sz https://example.com/article > article.md

# Wait for dynamic content
sz --wait-for=".article-content" https://spa-site.com

# Skip JavaScript for static sites (blazing fast)
sz --no-js https://static-site.com

# Interactive terminal UI mode
sz --tui
```

### 🧠 How It Works

1. **Fetch** → Headless Chrome renders the full page (including JavaScript)
2. **Extract** → Smart algorithms identify the main content
3. **Score** → Each block gets ranked by semantic importance
4. **Reorder** → Most important information surfaces to the top
5. **Render** → Clean, readable markdown output

### 📝 Output Example

```markdown
# The Rise of Local-First Software

*By Martin Kleppmann — November 2024*

> TL;DR: Cloud apps betrayed us. Local-first brings back ownership,
> privacy, and speed while keeping the collaboration we need.

## Key Insight: You Don't Need The Cloud
Local-first apps prove that the best features of cloud software can
exist without sacrificing user control or requiring constant connectivity...

[Rest of article, importance-ordered...]
```

## Configuration

Create `~/.config/essenz/config.yaml`:

```yaml
# Browser settings
browser:
  timeout: 30s
  viewport:
    width: 1920
    height: 1080

# Content extraction
extraction:
  top_blocks: 20
  summarize: true

# Scoring weights
scoring:
  tag_weights:
    h1: 3.0
    h2: 2.0
    p: 1.0
    nav: 0.2
```

## Advanced Usage

### TUI Mode

Launch the interactive terminal browser:

```bash
sz --tui
```

Key bindings:
- `Ctrl+L`: Focus URL bar
- `Enter`: Load URL
- `Tab`: Navigate links
- `Space`: Follow selected link
- `Backspace`: Go back
- `Ctrl+D`: Bookmark page
- `q`: Quit

### Wait Strategies

Handle different types of JavaScript rendering:

```bash
# Wait for specific selector
sz --wait-for="#content" https://example.com

# Wait for network idle
sz --wait-idle=2s https://example.com

# Custom timeout
sz --timeout=60s https://slow-site.com

# Combine strategies
sz --wait-for=".article" --wait-idle=1s https://example.com
```

By default, sz uses LCP (Largest Contentful Paint) as the primary readiness signal. It injects a PerformanceObserver and proceeds once the first LCP is observed (with a timeout fallback). To disable this behavior:

```bash
sz --no-lcp-wait https://example.com
```

### Output Formats

```bash
# Markdown (default)
sz https://example.com

# JSON structure
sz --format=json https://example.com

# Plain text
sz --format=text https://example.com

# HTML (cleaned)
sz --format=html https://example.com
```

## Development

### Prerequisites

- [asdf](https://asdf-vm.com/) version manager
- Go 1.21+ (managed via asdf)

### Building from Source

```bash
# Clone repository
git clone https://github.com/jewell-lgtm/essenz
cd essenz

# Install correct tool versions
asdf install

# Setup pre-commit hooks (runs checks before every commit)
make setup-pre-commit

# Install dependencies
go mod download

# Run all checks (includes tool version verification)
make check

# Build binary
make build

# Install locally
make install
```

### Writing Executable Specs

Essenz uses executable specifications for development. Create a spec file:

```markdown
# specs/features/my-feature.spec.md

## SPEC: Feature Description

GIVEN initial conditions
WHEN action occurs
THEN expected outcome

### Test Case
```spec
url: https://example.com
expected_contains:
  - "Expected text"
```
```

Run specs:

```bash
go test ./specs/...
```

Sandboxed environments

If your environment restricts network/exec/ports, skip env-dependent specs and use local caches:

```bash
make test-sandbox
```

### Project Structure

```
essenz/
├── cmd/essenz/        # CLI entry point
├── internal/
│   ├── browser/       # Chrome integration
│   ├── extractor/     # Content extraction
│   ├── renderer/      # Markdown rendering
│   └── tui/          # Terminal UI
├── specs/            # Executable specifications
└── test/             # Test fixtures
```

## Contributing

We welcome contributions! Please follow these guidelines:

1. **Write specs first**: Create executable specifications before implementing features
2. **Setup pre-commit**: Run `make setup-pre-commit` to install quality checks that run automatically
3. **Test thoroughly**: Use `make check` to run all quality checks and tests
4. **Follow conventions**: Pre-commit hooks enforce Go formatting, linting, and conventional commits
5. **Document changes**: Update specs and README as needed

### Development Workflow

Essenz follows a Test-Driven Development (TDD) workflow with feature branches:

#### 1. Start New Feature

```bash
# Fork and clone (first time only)
git clone https://github.com/YOUR_USERNAME/essenz
cd essenz

# Install tool versions
asdf install

# Setup pre-commit hooks (enforces quality automatically)
make setup-pre-commit

# Create feature branch from main
git checkout main
git pull origin main
git checkout -b feature/amazing-feature
```

#### 2. Write Executable Specification

Write the complete feature specification that should fail initially:

```bash
# Create spec file in specs/features/
# Write executable spec covering entire feature
# Run spec to confirm it fails
go test ./specs/features/amazing-feature.spec.go
```

#### 3. Initial Commit (Skip Hooks)

```bash
# Commit failing spec to establish feature scope
git add specs/features/amazing-feature.spec.go
SKIP=go-test git commit -m "feat: add executable spec for amazing feature

Comprehensive test specification covering:
- Core functionality requirements
- Error handling scenarios
- Expected behavior documentation

Tests define expected behavior before implementation (TDD approach).
Tests currently fail as expected - feature not yet implemented.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

#### 4. TDD Implementation Cycle

Make small, focused commits until the high-level spec passes:

```bash
# Make minimal change to pass one part of spec
# Edit code
make check  # Run all quality checks

# Commit small change
git add .
git commit -m "feat: add basic structure for amazing feature"

# Repeat: edit -> test -> commit until spec passes
go test ./specs/features/amazing-feature.spec.go
```

#### 5. Refactor (If Needed)

```bash
# Clean up implementation while keeping tests green
# Commit refactoring separately
git commit -m "refactor: improve amazing feature implementation"
```

#### 6. Merge to Main

```bash
# Switch back to main
git checkout main
git pull origin main

# Squash merge feature branch (keeps history clean)
git merge --squash feature/amazing-feature
git commit -m "feat: implement amazing feature

Complete implementation of amazing feature including:
- Core functionality with full spec coverage
- Comprehensive error handling
- Documentation and examples

Closes #123

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"

# Push to GitHub
git push origin main

# Clean up feature branch
git branch -d feature/amazing-feature
```

#### Key Principles

- **Spec First**: Always write executable specifications before implementation
- **Fail Fast**: Initial spec commit should show failing tests (demonstrates TDD approach)
- **Small Commits**: Make incremental progress with focused commits on feature branch
- **Clean History**: Use squash merge to main for clean, readable project history
- **Skip Hooks Selectively**: Use `SKIP=go-test` only for initial failing specs
- **Quality Checks**: Run `make check` before every commit to ensure code quality

## Performance

### Benchmarks

```bash
# Run benchmarks
go test -bench=. ./...

# Profile CPU usage
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof
```

Typical performance:
- Static page extraction: <500ms
- JavaScript SPA: 1-3s
- Heavy dynamic site: 3-5s

### Optimization Tips

- Use `--no-js` for static sites
- Enable caching with `--cache-dir`
- Adjust `--timeout` based on site complexity
- Use `--top-k` to limit output size

## Troubleshooting

### Common Issues

**Chrome not found**
```bash
# Install Chrome/Chromium
# macOS
brew install chromium

# Ubuntu/Debian
sudo apt-get install chromium-browser

# Set custom Chrome path
export ESSENZ_CHROME_PATH=/path/to/chrome
```

**JavaScript not rendering**
```bash
# Increase timeout
sz --timeout=60s https://slow-site.com

# Check for specific element
sz --wait-for=".content-loaded" https://site.com

# Debug mode
sz --debug https://site.com 2> debug.log
```

**Empty output**
```bash
# Try without JavaScript first
sz --no-js https://site.com

# Check if site blocks automation
sz --debug https://site.com

# Try different wait strategy
sz --wait-idle=3s https://site.com
```

## API Usage

Essenz can also be used as a Go library:

```go
package main

import (
    "github.com/jewell-lgtm/essenz/pkg/essenz"
)

func main() {
    // Create extractor
    e := essenz.New()
    defer e.Close()

    // Extract content
    content, err := e.Extract("https://example.com")
    if err != nil {
        panic(err)
    }

    // Output markdown
    fmt.Println(content.Markdown())
}
```

## License

MIT License - See [LICENSE](LICENSE) file

## Acknowledgments

- [go-readability](https://github.com/go-shiori/go-readability) for content extraction algorithms
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the amazing TUI framework
- [chromedp](https://github.com/chromedp/chromedp) for Chrome automation
- Mozilla's Readability.js for inspiration

## Support

- 🐛 [Report bugs](https://github.com/jewell-lgtm/essenz/issues)
- 💡 [Request features](https://github.com/jewell-lgtm/essenz/issues)
- 💬 [Discussions](https://github.com/jewell-lgtm/essenz/discussions)
- 📖 [Documentation](https://essenz.dev/docs)

---

<div align="center">
  <h1>ß</h1>
  <p><strong>essenz</strong> • <em>distill the web</em></p>
  <br/>
  <p>
    <a href="https://github.com/jewell-lgtm/essenz/releases">
      <img src="https://img.shields.io/github/v/release/jewell-lgtm/essenz?style=flat-square" alt="Release"/>
    </a>
    <a href="https://github.com/jewell-lgtm/essenz/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/jewell-lgtm/essenz?style=flat-square" alt="License"/>
    </a>
    <a href="https://github.com/jewell-lgtm/essenz">
      <img src="https://img.shields.io/github/stars/jewell-lgtm/essenz?style=flat-square" alt="Stars"/>
    </a>
  </p>
</div>
