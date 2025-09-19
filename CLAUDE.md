# Claude Code Assistant Context

## Project Overview

**essenz (sz)** - A CLI web browser that extracts the essence of web pages, reordering content by importance rather than DOM structure. Built with Go and powered by headless Chrome for handling modern JavaScript-heavy sites while producing clean, readable markdown.

- **Primary Binary**: `sz` (built from `./cmd/essenz`)
- **Module Path**: `github.com/essenz/essenz`
- **Go Version**: 1.23.3 (managed via asdf)
- **Current Version**: 0.1.0

## Current Implementation Status

This is an early-stage project with basic CLI structure in place:

### ✅ Implemented Features
- Basic Cobra CLI structure with `sz` command
- `sz version` - Version information
- `sz fetch` - Basic HTTP/HTTPS URL fetching and local file reading
- Executable specification framework for TDD
- Pre-commit hooks for code quality
- Tool version management with asdf

### 🚧 In Development
- Core web scraping with headless Chrome
- Content extraction and semantic reordering
- Markdown rendering
- Interactive TUI mode
- Advanced wait strategies and JavaScript handling

## Project Structure

```
essenz/
├── cmd/essenz/           # CLI entry point (main.go)
├── specs/               # Executable specifications (TDD tests)
│   ├── cli_spec_test.go
│   └── fetch_spec_test.go
├── internal/            # Internal packages (future)
│   ├── browser/         # Chrome integration (planned)
│   ├── extractor/       # Content extraction (planned)
│   ├── renderer/        # Markdown rendering (planned)
│   └── tui/            # Terminal UI (planned)
├── .tool-versions       # asdf tool versions
├── Makefile            # Build automation
├── go.mod              # Go module definition
└── README.md           # Documentation
```

## Development Commands

### Essential Commands
```bash
# Check tool versions match .tool-versions
make check-tools

# Run all quality checks (format, vet, lint, test)
make check

# Build binary
make build

# Run tests
make test

# Setup pre-commit hooks
make setup-pre-commit

# Format code
make fmt

# Run linter
make lint
```

### Testing Commands
```bash
# Run all tests
go test ./...

# Run specific spec file
go test ./specs/cli_spec_test.go
go test ./specs/fetch_spec_test.go

# Run tests with verbose output
go test -v ./...
```

### Binary Usage
```bash
# Show help
./sz help

# Show version
./sz version

# Fetch URL content
./sz fetch https://example.com

# Read local file
./sz fetch /path/to/file.html
```

## Git Workflow

This project follows a TDD workflow with feature branches:

1. **Create feature branch** from main
2. **Write executable spec** covering entire feature (should fail)
3. **Initial commit** with `SKIP=go-test` to bypass failing tests
4. **TDD implementation** with small commits until spec passes
5. **Refactor** if needed
6. **Squash merge** to main for clean history
7. **Push to GitHub**

### Commit Message Format
Uses conventional commits enforced by pre-commit hooks:
- `feat:` - New features
- `fix:` - Bug fixes
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `docs:` - Documentation changes

All commits should end with:
```
🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Dependencies

### Core Dependencies
- `github.com/spf13/cobra` v1.8.0 - CLI framework
- `github.com/stretchr/testify` v1.11.1 - Testing assertions

### Development Tools
- Go 1.23.3
- golangci-lint (for linting)
- goimports (for import formatting)
- pre-commit 4.0.1 (for git hooks)

## Key Implementation Details

### Current CLI Structure
- **Root command** (`sz`) - Shows help by default
- **Version command** (`sz version`) - Shows version info
- **Fetch command** (`sz fetch <url|path>`) - Fetches web content or reads files

### Fetch Command Implementation
- Simple URL detection via `http://` or `https://` prefix
- HTTP client with 30s timeout and insecure TLS for testing
- File reading with proper error handling
- Content output to stdout, errors to stderr

### Testing Strategy
- **Executable specifications** in `specs/` directory
- Test-driven development approach
- Specifications written before implementation
- Tests should fail initially to demonstrate TDD

## Future Architecture Plans

Based on README documentation, the project will include:

- **Browser integration** - Headless Chrome via chromedp
- **Content extraction** - Readability algorithms for main content
- **Semantic scoring** - Importance-based content reordering
- **Multiple renderers** - Markdown, JSON, HTML, plain text
- **Interactive TUI** - Terminal UI with Bubble Tea
- **Configuration** - YAML config files
- **Performance optimizations** - Caching, connection pooling

## Quality Standards

- **Pre-commit hooks** enforce formatting, linting, and conventional commits
- **Tool version consistency** managed via asdf
- **Comprehensive testing** with executable specifications
- **Clean git history** via squash merging
- **Code quality** enforced by golangci-lint

## Notes for Claude Code

- Always run `make check` before committing
- Use `SKIP=go-test` only for initial failing spec commits
- Follow the documented git workflow for feature development
- Prefer editing existing files over creating new ones
- Binary is built as `sz` not `essenz`
- Current implementation is minimal - core features are planned but not implemented
- When implementing new features, start with executable specs in `specs/` directory