# Claude Code Assistant Context

## 🚨 MANDATORY WORKFLOW - READ FIRST

**CRITICAL: Before implementing ANY feature, Claude MUST:**

1. ✅ Create feature branch from main
2. ✅ Write executable spec that FAILS initially
3. ✅ Commit failing spec with `SKIP=go-test`
4. ✅ Implement with small commits until spec passes
5. ✅ Squash merge to main and push

**NO EXCEPTIONS. See "MANDATORY Git Workflow" section below for details.**

## Project Overview

**essenz (sz)** - A CLI web browser that extracts the essence of web pages, reordering content by importance rather than DOM structure. Built with Go and powered by headless Chrome for handling modern JavaScript-heavy sites while producing clean, readable markdown.

- **Primary Binary**: `sz` (built from `./cmd/essenz`)
- **Module Path**: `github.com/jewell-lgtm/essenz`
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

**Architecture Principle**: Thin CLI layer orchestrating modular internal packages

```
essenz/
├── cmd/essenz/           # CLI entry point (thin orchestration layer)
│   └── main.go          # Cobra commands that delegate to internal modules
├── internal/            # Business logic modules (not importable externally)
│   ├── browser/         # Chrome daemon management and browser operations
│   ├── daemon/          # Chrome process lifecycle and connection management
│   ├── extractor/       # Content extraction and semantic analysis
│   ├── renderer/        # Markdown/text output formatting
│   └── config/          # Configuration management
├── pkg/                 # Public APIs (importable by other projects)
│   └── essenz/          # Core types and interfaces
├── specs/               # Executable specifications (TDD tests)
│   ├── cli_spec_test.go
│   └── fetch_spec_test.go
├── .tool-versions       # asdf tool versions
├── Makefile            # Build automation
├── go.mod              # Go module definition
└── README.md           # Documentation
```

### Module Responsibilities

- **`cmd/essenz`**: Argument parsing, command routing, minimal business logic
- **`internal/browser`**: Browser context management, page operations
- **`internal/daemon`**: Chrome process lifecycle, connection pooling
- **`internal/extractor`**: Content analysis, text extraction, importance scoring
- **`internal/renderer`**: Output formatting (markdown, JSON, text)
- **`internal/config`**: Configuration loading and validation
- **`pkg/essenz`**: Public types, interfaces for external integrations

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

## MANDATORY Git Workflow

**⚠️ CRITICAL: Claude MUST follow this workflow for ALL feature development. No exceptions.**

### Required Workflow Steps

**EVERY feature development MUST follow these exact steps:**

1. **Create feature branch** from main
   ```bash
   git checkout main
   git pull origin main
   git checkout -b feature/descriptive-name
   ```

2. **Write executable spec** covering entire feature (MUST fail initially)
   ```bash
   # Create spec file in specs/
   # Write comprehensive test covering full feature
   # Verify it fails: go test ./specs/feature-name.spec.go
   ```

3. **Initial commit** with hooks skipped (for failing spec)
   ```bash
   git add specs/feature-name.spec.go
   SKIP=go-test git commit -m "feat: add executable spec for [feature]"
   ```

4. **TDD implementation** with small, focused commits
   ```bash
   # Make minimal change
   make check  # MUST pass before commit
   git add .
   git commit -m "feat: implement basic [specific change]"
   # Repeat until spec passes
   ```

5. **Refactor** if needed (separate commits)
   ```bash
   git commit -m "refactor: improve [specific aspect]"
   ```

6. **Return to main and squash merge**
   ```bash
   git checkout main
   git merge --squash feature/branch-name
   git commit -m "feat: implement [complete feature]"
   ```

7. **Push to GitHub**
   ```bash
   git push origin main
   git branch -d feature/branch-name
   ```

### Workflow Enforcement Checklist

Before starting ANY feature work, Claude must verify:
- [ ] Currently on main branch
- [ ] Main branch is up to date
- [ ] Feature branch created with descriptive name
- [ ] Executable spec written and failing
- [ ] Initial commit made with SKIP=go-test

During development, Claude must verify:
- [ ] Each commit is small and focused
- [ ] `make check` passes before every commit
- [ ] Commit messages follow conventional format
- [ ] Tests are passing incrementally

Before completing feature, Claude must verify:
- [ ] All specs are passing
- [ ] Code is properly formatted and linted
- [ ] Returned to main branch
- [ ] Squash merge completed
- [ ] Changes pushed to GitHub
- [ ] Feature branch deleted

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

### WORKFLOW ENFORCEMENT
- **ALWAYS follow the mandatory workflow above - no shortcuts**
- **NEVER implement features directly on main branch**
- **NEVER skip writing failing specs first**
- **ALWAYS create feature branch before any code changes**
- **ALWAYS squash merge and push at completion**

### DEVELOPMENT STANDARDS
- Always run `make check` before committing (enforced by checklist)
- Use `SKIP=go-test` only for initial failing spec commits
- Prefer editing existing files over creating new ones
- Binary is built as `sz` not `essenz`
- Current implementation is minimal - core features are planned but not implemented
- When implementing new features, start with executable specs in `specs/` directory

### COMMIT MESSAGE TEMPLATES
All commits must end with:
```
🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
```