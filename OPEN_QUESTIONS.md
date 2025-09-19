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

- [ ] **What's our target user?**
  - Developers who want to extract web content to markdown?
  - CLI enthusiasts learning Go?
  - People who need readable versions of web articles?

## 🏗️ Architecture & Design Questions

- [x] **Project structure - how should we organize the code?** ✅ **DECIDED**
  - **Approach**: Standard Go layout with thin CLI layer deferring to internal modules
  - **Structure**: `cmd/` (thin CLI) + `internal/` (business logic) + `pkg/` (public APIs)
  - **Modules**: `internal/browser/`, `internal/extractor/`, `internal/renderer/`, `internal/daemon/`
  - **Principle**: CLI commands orchestrate, modules do the work

- [ ] **What external dependencies do we want?**
  - The go.mod mentions chromedp, go-readability, cobra - do we use these?
  - Should we minimize dependencies to focus on Go fundamentals?
  - Are we building a "real" tool or a learning exercise?

- [ ] **Configuration approach?**
  - Should we implement the YAML config mentioned in README?
  - Environment variables? Command flags only?
  - How complex should our configuration be?

## 🧪 Testing & Quality Questions

- [ ] **Testing strategy depth?**
  - Are the current unit tests + executable specs sufficient?
  - Do we need integration tests with real web scraping?
  - Should we mock external dependencies or test against real sites?

- [ ] **Error handling patterns?**
  - What Go error handling patterns should we establish?
  - How should we handle network failures, parsing errors, etc.?
  - Should we use structured logging (logrus, zap) or keep it simple?

## 🚀 Development & Deployment Questions

- [ ] **Release strategy?**
  - Should we set up automated releases with GitHub Actions?
  - Do we need cross-platform builds (Windows, Linux, macOS)?
  - Should we publish to package managers (Homebrew, etc.)?

- [ ] **Documentation completeness?**
  - The README has extensive docs for features we don't have - update or implement?
  - Should we focus on godoc comments and keep README minimal?
  - Do we need contributor guidelines beyond what's there?

## 🎓 Learning & Scope Questions

- [ ] **What Go concepts should we prioritize learning?**
  - Interfaces and composition?
  - Goroutines and concurrency?
  - HTTP clients and web scraping?
  - File I/O and data processing?

- [x] **Feature scope for next iteration?** ✅ **DECIDED**
  - **Next feature**: `sz fetch` command supporting both local files and HTTP(S) URLs
  - **Scope**: Simple fetch and print content (no processing/parsing yet)
  - **Approach**: Start with executable specs, then implement incrementally
  - **Examples**: `sz fetch https://example.com` and `sz fetch /path/to/file.html`

- [ ] **Performance and optimization?**
  - Should we worry about performance from the start?
  - Are there specific Go performance patterns we should learn?
  - Should we add benchmarking and profiling?

## 🤝 Collaboration Questions

- [x] **How do you prefer to learn/work?** ✅ **DECIDED**
  - **Chosen approach**: Small, incremental steps with agent-driven implementation
  - Agent performs small, atomic pieces of work with granular planning
  - Human guides direction, answers questions, and tweaks plans before implementation
  - Focus on one concept/change at a time with lots of commits

- [ ] **What's your comfort level with complexity?**
  - Are you ready for more advanced Go concepts?
  - Should we stick to basics for now?
  - How much external library usage vs. building from scratch?

---

## 🚦 Recommended Starting Points

**High Priority** - Answer these first:

- [ ] **Feature scope for next iteration** - What's the next smallest meaningful feature?
- [ ] **How do you prefer to learn/work** - Approach to development and learning
- [ ] **What should `sz` actually do** - Should we build toward the README vision or pivot?

These will guide all our other decisions and help us maintain momentum with focused, conventional commits!