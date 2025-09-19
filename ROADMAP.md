# essenz Feature Roadmap

This roadmap defines the complete development path for essenz, organized into testable features. Each feature corresponds to one high-level executable specification that defines the expected end-to-end behavior.

## Phase 1: Foundation (Core CLI + Basic Fetching)

### ✅ F1: Basic CLI Structure
**Status**: IMPLEMENTED
**Spec**: `specs/cli_spec_test.go`
- `sz` shows help by default
- `sz version` displays version information
- `sz help` shows usage information

### ✅ F2: Simple HTTP/File Fetching
**Status**: IMPLEMENTED
**Spec**: `specs/fetch_spec_test.go`
- `sz fetch https://example.com` retrieves and prints web content
- `sz fetch /path/to/file.html` reads and prints local file content
- Error handling for invalid URLs and missing files

## Phase 2: Chrome Integration (Browser Automation)

### 🔄 F3: Chrome Daemon Management
**Status**: NEXT UP
**Spec**: `specs/daemon_spec_test.go` *(to be created)*
- Chrome daemon starts on first browser operation
- Multiple sz commands reuse the same Chrome instance
- Daemon shuts down gracefully when not needed
- Daemon restarts automatically if crashed

### 🔄 F4: JavaScript-Heavy Page Rendering
**Status**: PLANNED
**Spec**: `specs/chrome_fetch_spec_test.go` *(to be created)*
- `sz fetch https://spa-site.com` renders JavaScript before extraction
- Wait strategies for dynamic content loading
- Fallback to simple HTTP for static sites
- Timeout handling for slow-loading pages

## Phase 3: Content Intelligence (Extraction + Processing)

### 🔄 F5: Smart Content Extraction
**Status**: PLANNED
**Spec**: `specs/extraction_spec_test.go` *(to be created)*
- Extract main article content from complex web pages
- Remove navigation, ads, and clutter automatically
- Preserve important text structure and headings
- Handle various content layouts and CMS patterns

### 🔄 F6: Content Importance Scoring
**Status**: PLANNED
**Spec**: `specs/importance_spec_test.go` *(to be created)*
- Reorder content blocks by semantic importance
- Key information appears first in output
- Supporting details follow main content
- Consistent scoring across different page types

## Phase 4: Output Formatting (Clean Presentation)

### 🔄 F7: Markdown Rendering
**Status**: PLANNED
**Spec**: `specs/markdown_spec_test.go` *(to be created)*
- Clean, readable markdown output by default
- Proper heading hierarchy preservation
- Link formatting and preservation
- Code block and formatting retention

### 🔄 F8: Multiple Output Formats
**Status**: PLANNED
**Spec**: `specs/output_formats_spec_test.go` *(to be created)*
- `sz fetch --format=json` outputs structured data
- `sz fetch --format=text` provides plain text
- `sz fetch --format=html` gives cleaned HTML
- Format-specific optimizations and metadata

## Phase 5: Advanced Features (Performance + UX)

### 🔄 F9: Content Summarization
**Status**: PLANNED
**Spec**: `specs/summarization_spec_test.go` *(to be created)*
- `sz fetch --summarize` provides article summary
- Key points extraction from long content
- TL;DR generation for quick consumption
- Configurable summary length

### 🔄 F10: Caching System
**Status**: PLANNED
**Spec**: `specs/caching_spec_test.go` *(to be created)*
- `sz fetch --cache-dir=/path` enables persistent caching
- Cached responses for repeated requests
- Cache invalidation and freshness checking
- Performance improvements for common use cases

### 🔄 F11: Wait Strategies
**Status**: PLANNED
**Spec**: `specs/wait_strategies_spec_test.go` *(to be created)*
- `sz fetch --wait-for=".content"` waits for specific selectors
- `sz fetch --wait-idle=2s` waits for network idle
- `sz fetch --timeout=60s` custom timeout handling
- Reliable handling of complex JavaScript applications

## Phase 6: User Experience (Polish + Usability)

### 🔄 F12: Configuration System
**Status**: PLANNED
**Spec**: `specs/config_spec_test.go` *(to be created)*
- YAML configuration file support (~/.config/essenz/config.yaml)
- Environment variable overrides (ESSENZ_CHROME_PATH, etc.)
- Command-line flag precedence over config
- Config validation and helpful error messages

### 🔄 F13: Interactive TUI Mode
**Status**: PLANNED
**Spec**: `specs/tui_spec_test.go` *(to be created)*
- `sz --tui` launches terminal browser interface
- URL bar, navigation, and content display
- Keyboard shortcuts for common operations
- Bookmark management and history

### 🔄 F14: Debug and Monitoring
**Status**: PLANNED
**Spec**: `specs/debug_spec_test.go` *(to be created)*
- `sz fetch --debug` provides detailed operation logs
- Performance metrics and timing information
- Chrome process monitoring and health checks
- Troubleshooting information for failed extractions

## Phase 7: Performance & Reliability (Production Ready)

### 🔄 F15: Performance Optimization
**Status**: PLANNED
**Spec**: `specs/performance_spec_test.go` *(to be created)*
- Sub-second response time for cached content
- Memory usage stays below defined limits
- Chrome process resource management
- Concurrent request handling

### 🔄 F16: Error Recovery
**Status**: PLANNED
**Spec**: `specs/error_recovery_spec_test.go` *(to be created)*
- Graceful degradation when Chrome is unavailable
- Network failure retry logic
- Helpful error messages for common issues
- System resource exhaustion handling

### 🔄 F17: Cross-Platform Support
**Status**: PLANNED
**Spec**: `specs/cross_platform_spec_test.go` *(to be created)*
- Full functionality on macOS, Linux, Windows
- Platform-specific Chrome detection
- Path handling across operating systems
- Installation and setup processes

## Implementation Strategy

### Development Principles
1. **One feature, one spec**: Each feature has exactly one executable specification
2. **TDD workflow**: Write failing spec first, implement until passing
3. **Small iterations**: Each feature is 1-3 weeks of focused development
4. **Quality gates**: All specs must pass before moving to next feature

### Feature Dependencies
- **F3** (Chrome Daemon) blocks F4, F5, F6, F9, F11, F13
- **F5** (Content Extraction) blocks F6, F7, F9
- **F7** (Markdown) blocks F8
- **F12** (Config) enables advanced options in later features

### Success Metrics
- All executable specs pass consistently
- Performance benchmarks meet targets
- User feedback validates feature utility
- Code quality maintains high standards

---

## Current Status: Phase 2 Ready

✅ **Foundation Complete**: Basic CLI and HTTP fetching working
🎯 **Next Feature**: F3 - Chrome Daemon Management
📋 **Spec to Write**: `specs/daemon_spec_test.go`

This roadmap provides clear direction while maintaining flexibility for adjustments based on learning and feedback.