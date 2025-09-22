# Open Questions

Before continuing development, we should answer these key questions to guide our direction and decisions.

## 🔥 Critical Chrome/Browser Management Questions

- [x] **How do we manage the headless Chrome daemon process?** ✅ **DECIDED**
  - **Approach**: Persistent Chrome daemon process with command connections
  - **Rationale**: Faster subsequent commands, shared resource pool, better for repeated use
  - **Implementation**: Long-running Chrome process, sz commands connect as needed
  - **Lifecycle**: Start daemon on first use, graceful shutdown, crash recovery needed

- [x] **How do we access browser scope inside sz commands?** ✅ **DECIDED**
  - **Approach**: Browser client that connects to persistent Chrome daemon
  - **Interface**: Abstract browser interface for CLI commands to use
  - **Connection**: Commands get browser context via connection manager
  - **Design**: CLI layer → Browser Manager → Chrome Daemon connection

- [x] **How do we prevent memory leaks with Chrome automation?** ✅ **DECIDED**
  - **Strategy**: Strict resource cleanup per command with timeout enforcement
  - **Approach**: Each command gets fresh browser context, cleaned up after use
  - **Implementation**: Context pooling with automatic cleanup, page limits, periodic GC
  - **Safeguards**: Command timeouts, memory monitoring, daemon restart if needed

## 🎯 Core Product Questions

- [x] **What should `sz` actually do?** ✅ **DECIDED**
  - **Core value**: Provide human readable text versions of web pages from the command line
  - **Scope**: Web scraping + content extraction + markdown output (not a full browser)
  - **Focus**: Clean, readable text extraction - no complex browser features needed

- [x] **What's our target user?** ✅ **DECIDED**
  - **Primary**: Developers and technical users who need clean text from web content
  - **Use cases**: Documentation extraction, content analysis, CLI workflows
  - **Secondary**: Anyone wanting readable article text without browser clutter
  - **Focus**: Technical users comfortable with command-line tools

## 🏗️ Architecture & Design Questions

- [x] **Project structure - how should we organize the code?** ✅ **DECIDED**
  - **Approach**: Standard Go layout with thin CLI layer deferring to internal modules
  - **Structure**: `cmd/` (thin CLI) + `internal/` (business logic) + `pkg/` (public APIs)
  - **Modules**: `internal/browser/`, `internal/extractor/`, `internal/renderer/`, `internal/daemon/`
  - **Principle**: CLI commands orchestrate, modules do the work

- [x] **What external dependencies do we want?** ✅ **DECIDED**
  - **Essential**: chromedp (Chrome automation), cobra (CLI framework)
  - **Content extraction**: go-readability or custom implementation
  - **Approach**: Use proven libraries for complex tasks, implement core logic ourselves
  - **Goal**: Real tool with educational value - practical dependencies are fine

- [x] **Configuration approach?** ✅ **DECIDED**
  - **Primary**: Command flags for most options (simple, transparent)
  - **Secondary**: Environment variables for defaults (ESSENZ_CHROME_PATH, etc.)
  - **Config file**: Optional YAML for advanced users (~/.config/essenz/config.yaml)
  - **Keep simple**: Start with flags, add config file later if needed

## 🧪 Testing & Quality Questions

- [x] **Testing strategy depth?** ✅ **DECIDED**
  - **Primary**: Executable specs for end-to-end behavior (current approach)
  - **Unit tests**: For individual modules and complex logic
  - **Integration**: Test against real sites for core functionality, mocks for edge cases
  - **Approach**: TDD with failing specs first, comprehensive but not excessive

- [x] **Error handling patterns?** ✅ **DECIDED**
  - **Pattern**: Standard Go errors with context (fmt.Errorf with %w wrapping)
  - **Logging**: Simple log package initially, structured logging (slog) later if needed
  - **User errors**: Clear, actionable messages to stderr, technical details in debug mode
  - **Recovery**: Graceful degradation where possible (fallback to HTTP if Chrome fails)

## 🚀 Development & Deployment Questions

- [x] **Release strategy?** ✅ **DECIDED**
  - **Automation**: GitHub Actions for releases and cross-platform builds
  - **Platforms**: macOS, Linux, Windows (Go makes this easy)
  - **Distribution**: GitHub releases initially, Homebrew later if tool gains traction
  - **Versioning**: Semantic versioning with automated changelog

- [x] **Documentation completeness?** ✅ **DECIDED**
  - **README**: Keep aspirational docs, mark unimplemented features clearly
  - **Code docs**: Comprehensive godoc comments for all public APIs
  - **Contributing**: Current guidelines are sufficient, expand as project grows
  - **Approach**: Document as we build, README shows vision and current state

## 🎓 Learning & Scope Questions

- [x] **What Go concepts should we prioritize learning?** ✅ **DECIDED**
  - **Phase 1**: Interfaces, composition, error handling (for modular architecture)
  - **Phase 2**: HTTP clients, context management (for web scraping)
  - **Phase 3**: Goroutines and concurrency (for daemon and performance)
  - **Focus**: Learn concepts as needed for current features

- [x] **Feature scope for next iteration?** ✅ **DECIDED**
  - **Next feature**: `sz fetch` command supporting both local files and HTTP(S) URLs
  - **Scope**: Simple fetch and print content (no processing/parsing yet)
  - **Approach**: Start with executable specs, then implement incrementally
  - **Examples**: `sz fetch https://example.com` and `sz fetch /path/to/file.html`

- [x] **Performance and optimization?** ✅ **DECIDED**
  - **Approach**: Write correct code first, optimize later based on actual usage
  - **Monitoring**: Add basic benchmarks for core operations, profile when needed
  - **Focus areas**: Memory management with Chrome, response time for commands
  - **Principle**: Premature optimization is root of all evil - measure first

## 🤝 Collaboration Questions

- [x] **How do you prefer to learn/work?** ✅ **DECIDED**
  - **Chosen approach**: Small, incremental steps with agent-driven implementation
  - Agent performs small, atomic pieces of work with granular planning
  - Human guides direction, answers questions, and tweaks plans before implementation
  - Focus on one concept/change at a time with lots of commits

- [x] **What's your comfort level with complexity?** ✅ **DECIDED**
  - **Approach**: Incremental complexity - start simple, add sophistication gradually
  - **Libraries**: Use proven libraries for complex tasks (chromedp, cobra)
  - **Learning**: Introduce advanced Go concepts as needed for features
  - **Balance**: Real-world tool with educational progression

---

## 🎯 Current Status: Core Implementation Complete

**All critical questions have been answered and implemented successfully.**

### ✅ Completed Implementations
- ✅ Chrome daemon management with persistent processes
- ✅ Full content extraction pipeline (F1-F6)
- ✅ TDD workflow with executable specifications
- ✅ Modular architecture with clean separation
- ✅ Production-ready CLI with comprehensive features
- ✅ LLM integration for TL;DR functionality

### 🚧 Next Phase: Interactive Browser Mode (F7-F10)
The next iteration focuses on **Interactive Browser Mode** with complete specifications ready:
- **F7**: Footer with page statistics and eszett branding
- **F8**: Links extraction with character codes
- **F9**: Interactive TUI using Bubble Tea framework
- **F10**: Link navigation with history and overlays

All questions answered, architecture proven, foundation solid. Ready for next phase implementation!