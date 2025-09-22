# sz (essenz) Implementation Guide

## Overview

This guide describes the current implementation of sz, a CLI tool for extracting clean, readable content from web pages. The implementation is **COMPLETE** for core functionality and uses a modular architecture with comprehensive test coverage following TDD principles.

## Current Architecture

### Project Structure
```
essenz/
├── cmd/essenz/           # CLI entry point (sz binary)
│   └── main.go          # Cobra commands orchestrating internal modules
├── internal/            # Business logic modules (production-ready)
│   ├── browser/         # High-level browser operations
│   ├── daemon/          # Chrome process lifecycle management
│   ├── extractor/       # Legacy content extraction (replaced by newer modules)
│   ├── tree/            # DOM tree building and text node processing
│   ├── filter/          # Content filtering and cleaning
│   ├── markdown/        # Markdown rendering system
│   ├── media/           # Image and media handling
│   ├── pageready/       # DOM readiness detection (LCP-based)
│   └── summarizer/      # LLM-powered TL;DR functionality
├── specs/               # Executable specifications (TDD tests)
└── plan/                # Implementation documentation
```

### Core Dependencies
```go
// go.mod - Current production dependencies
module github.com/jewell-lgtm/essenz

go 1.23.3

require (
    github.com/chromedp/chromedp v0.11.2    // Chrome automation
    github.com/spf13/cobra v1.8.0           // CLI framework
    github.com/stretchr/testify v1.11.1     // Testing framework
    // Additional deps for content processing, HTTP clients, etc.
)
```

## Implemented Features (Production Ready)

### ✅ F1-F6: Core Content Extraction Pipeline

The following features are **fully implemented and production-ready**:

#### F1: DOM Ready Event System (`internal/pageready/`)
- LCP (Largest Contentful Paint) based readiness detection
- JavaScript framework initialization waiting
- Configurable timeouts with graceful fallbacks
- Chrome DevTools Protocol integration

#### F2: Text Node Tree Builder (`internal/tree/`)
- Bottom-up content tree from actual DOM text nodes
- Preserves semantic hierarchy and document structure
- Handles dynamic JavaScript-generated content
- Efficient tree traversal and manipulation

#### F3: Content Filter System (`internal/filter/`)
- Multi-stage filtering pipeline (class, length, navigation)
- Rule-based removal of ads, navigation, boilerplate
- Whitelist preservation for important content areas
- Configurable filtering aggressiveness

#### F4: Image and Media Handler (`internal/media/`)
- Intelligent media detection and replacement
- Alt text and caption extraction for context
- Support for images, videos, audio, and embeds
- Markdown-friendly media descriptions

#### F5: Markdown Tree Renderer (`internal/markdown/`)
- Clean, hierarchical markdown generation
- Configurable emphasis styles (`*` vs `_`)
- Proper spacing and formatting rules
- Quality assurance and validation

#### F6: TL;DR Summarizer (`internal/summarizer/`)
- OpenAI and custom API endpoint integration
- Multiple summary lengths (short/medium/long)
- Secure API key management
- Content-only summarization (no metadata leakage)

### Additional Production Features

#### Chrome Daemon Management (`internal/daemon/`)
- Persistent Chrome process with Unix socket communication
- Connection pooling and resource management
- Automatic startup, graceful shutdown, crash recovery
- Process lifecycle monitoring

#### CLI Interface (`cmd/essenz/`)
- Complete Cobra-based CLI with subcommands
- `sz version` - Version information
- `sz fetch <url|file>` - Content extraction
- `sz tldr <url|file>` - AI-powered summarization
- Comprehensive flag support and validation

## Current Usage

### Basic Content Extraction
```bash
# Extract content from URL
sz https://example.com/article

# Process local HTML file
sz /path/to/article.html

# Raw HTML output (bypass processing)
sz --raw https://example.com
```

### TL;DR Summarization
```bash
# Set API key and summarize
export OPENAI_API_KEY=sk-...
sz tldr https://example.com/article

# Control summary length
sz tldr --summary-length short https://example.com

# Custom model and endpoint
sz tldr --model gpt-4 --base-url https://api.openai.com/v1 https://example.com
```

### Advanced Options
```bash
# DOM readiness control
sz --wait-for-frameworks https://spa-site.com
sz --dom-ready-timeout 30s https://slow-site.com
sz --no-lcp-wait https://static-site.com

# Content filtering
sz --content-filter aggressive https://cluttered-site.com
sz --preserve-selector ".important-content" https://site.com

# Output formats
sz --text-node-tree https://site.com  # Debug tree structure
sz --markdown-renderer https://site.com  # Force markdown processing
```

## Next Phase: Interactive Browser Mode (F7-F10)

### 🚧 Ready for Implementation

The following features have **complete executable specifications** and are ready for implementation:

#### F7: Footer with Page Statistics
- **Spec**: `specs/features/F7-footer-stats.spec.md` ✅
- `--footer` flag for optional statistics display
- Page metrics (source, size, processing time)
- Eszett (ß) branding in pink color
- Foundation for pastel color scheme

#### F8: Links Extraction and Display
- **Spec**: `specs/features/F8-links-extraction.spec.md` ✅
- Smart link extraction from processed content
- Character code assignment (a-z, aa-zz)
- Footer display with visual hierarchy
- Link importance scoring and filtering

#### F9: Interactive TUI with Bubble Tea
- **Spec**: `specs/features/F9-interactive-tui.spec.md` ✅
- `--interactive` flag for Terminal User Interface
- Scrollable viewport with vim-like navigation
- Status bar with position indicators
- Complete pastel color scheme application

#### F10: Link Navigation in TUI
- **Spec**: `specs/features/F10-link-navigation.spec.md` ✅
- Links overlay with keyboard toggle
- Character code navigation system
- Navigation history and back functionality
- Loading states and comprehensive error handling

### Implementation Dependencies
```
F7 (Footer Stats) → F8 (Links) → F9 (TUI) → F10 (Navigation)
```
Each feature builds on the previous, following the established TDD workflow.

## Development Workflow

### TDD Approach (Proven)
The current implementation follows a strict Test-Driven Development approach:

1. **Write executable specifications** (markdown + test cases)
2. **Create failing tests** that define expected behavior
3. **Implement incrementally** until tests pass
4. **Refactor** while maintaining green tests
5. **Squash merge** completed features to main

### Quality Standards
- **Pre-commit hooks**: Automated formatting, linting, testing
- **Executable specs**: All features have comprehensive test coverage
- **Modular architecture**: Clean separation of concerns
- **Error handling**: Comprehensive error wrapping and recovery
- **Performance**: Optimized for real-world usage patterns

### Build and Development
```bash
# Install correct tool versions
asdf install

# Setup development environment
make setup-pre-commit

# Run all quality checks
make check

# Build binary
make build

# Run tests
make test
```

## Architectural Principles

### Modular Design
- **Thin CLI layer**: Commands orchestrate, don't implement
- **Internal packages**: Single responsibility, clear interfaces
- **Composability**: Features can be used independently
- **Testability**: Each module has comprehensive test coverage

### Error Handling
- **Context awareness**: All operations support cancellation
- **Graceful degradation**: Fallbacks when Chrome unavailable
- **User-friendly messages**: Clear error reporting
- **Comprehensive logging**: Debug information when needed

### Performance Considerations
- **Chrome daemon reuse**: Persistent process for multiple operations
- **Efficient filtering**: Multi-stage pipeline with early exits
- **Memory management**: Proper cleanup and resource management
- **Concurrent processing**: Where applicable and beneficial

## Summary

The sz tool has a **complete, production-ready core** (F1-F6) with:
- Full content extraction pipeline
- Chrome automation and daemon management
- LLM-powered summarization
- Comprehensive CLI interface
- Extensive test coverage

The **next phase** (F7-F10) has detailed specifications ready for implementation, adding:
- Interactive terminal browser mode
- Visual enhancements with eszett branding
- Link navigation with character codes
- Complete TUI experience using Bubble Tea

The architecture is proven, the workflow is established, and the foundation is solid for continued development.
