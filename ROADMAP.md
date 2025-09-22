# sz (essenz) Feature Roadmap

This roadmap reflects the current implementation status and future development plans for sz, organized into completed and planned features.

## ✅ Phase 1: Core Content Extraction (COMPLETED)

### ✅ F1: DOM Ready Event System
**Status**: IMPLEMENTED ✅
**Spec**: `specs/dom_ready_events_spec_test.go`
- DOM readiness detection with JavaScript framework support
- LCP (Largest Contentful Paint) based waiting
- Configurable timeouts and fallback strategies
- Chrome automation integration

### ✅ F2: Text Node Tree Builder
**Status**: IMPLEMENTED ✅
**Spec**: `specs/text_node_tree_builder_spec_test.go`
- Bottom-up content tree from actual text nodes
- Hierarchical structure preservation
- Dynamic content handling
- Semantic context maintenance

### ✅ F3: Content Filter System
**Status**: IMPLEMENTED ✅
**Spec**: `specs/content_filter_spec_test.go`
- Rule-based filtering for navigation, ads, boilerplate
- CSS class pattern analysis
- Content density evaluation
- Whitelist preservation for important elements

### ✅ F4: Image and Media Handler
**Status**: IMPLEMENTED ✅
**Spec**: `specs/image_media_handler_spec_test.go`
- Media detection and meaningful replacement
- Alt text and caption extraction
- Context-aware description generation
- Multiple media type support

### ✅ F5: Markdown Tree Renderer
**Status**: IMPLEMENTED ✅
**Spec**: `specs/markdown_tree_renderer_spec_test.go`
- Clean, well-formatted markdown output
- Hierarchical rendering with proper spacing
- Configurable emphasis and list styles
- Quality assurance and validation

### ✅ F6: TL;DR Summarizer
**Status**: IMPLEMENTED ✅
**Spec**: `specs/tldr_summarizer_spec_test.go`
- LLM-powered content summarization
- OpenAI and custom API endpoint support
- Multiple summary length options
- Secure API key management

### ✅ Additional Core Features
- **Chrome Daemon Management**: Persistent Chrome process with connection pooling
- **CLI Structure**: Complete Cobra-based CLI with version, fetch, tldr commands
- **Configuration System**: Flag-based configuration with environment variable support

## 🚧 Phase 2: Interactive Browser Mode (SPECS READY)

### 🆕 F7: Footer with Page Statistics
**Status**: SPECS READY 📋
**Spec**: `specs/features/F7-footer-stats.spec.md`
- `--footer` flag for statistics display
- Page metrics (source, size, processing time)
- Eszett (ß) branding in pink color
- Pastel color scheme foundation

### 🆕 F8: Links Extraction and Display
**Status**: SPECS READY 📋
**Spec**: `specs/features/F8-links-extraction.spec.md`
- Smart link extraction and prioritization
- Character code assignment (a-z, aa-zz)
- Footer links section with visual hierarchy
- Link filtering and importance scoring

### 🆕 F9: Interactive TUI with Bubble Tea
**Status**: SPECS READY 📋
**Spec**: `specs/features/F9-interactive-tui.spec.md`
- `--interactive` flag for TUI mode
- Scrollable viewport with vim-like navigation
- Status bar with position indicators
- Consistent pastel color scheme application

### 🆕 F10: Link Navigation in TUI
**Status**: SPECS READY 📋
**Spec**: `specs/features/F10-link-navigation.spec.md`
- Links overlay with keyboard toggle
- Character code navigation system
- Navigation history and back button
- Loading states and error handling

## 🔮 Phase 3: Future Enhancements (PLANNED)

### 📋 F11: Advanced Configuration
- YAML configuration file support
- Custom color schemes and themes
- Advanced filtering rules
- User-defined shortcuts

### 📋 F12: Performance Optimization
- Content caching system
- Memory usage optimization
- Concurrent request handling
- Response time improvements

### 📋 F13: Search and Discovery
- In-content search functionality (`/` key in TUI)
- Bookmark management system
- History and favorites
- Content tagging and organization

### 📋 F14: Extended Output Formats
- JSON structured output
- Plain text formatting
- HTML cleaned output
- Custom format templates

### 📋 F15: Integration Features
- Plugin system for custom processors
- Webhook support for automation
- API mode for programmatic access
- Export to external services

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