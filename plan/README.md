# sz (essenz) Implementation Plan

This directory contains detailed specifications for implementing the `sz` tool's content extraction and interactive browser features. The implementation is organized into two phases: core content extraction (F1-F6) and interactive browser mode (F7-F10).

## Implementation Order

The features must be implemented in the following order, as each builds on the previous:

### [F1: DOM Ready Event System](01-dom-ready-events.md) ✅ **COMPLETED**
**Branch**: `feature/dom-ready-events`

Implements proper DOM event waiting to ensure all dynamic content is loaded before extraction begins. This is the foundation for reliable content extraction across modern JavaScript frameworks.

**Key Components**:
- DOM readiness detection (DOMContentLoaded, load events)
- JavaScript framework initialization waiting
- Configurable timeouts and retry logic
- Chrome integration enhancement

### [F2: Text Node Tree Builder](02-text-node-tree-builder.md) ✅ **COMPLETED**
**Branch**: `feature/text-node-tree-builder`

Builds a content tree starting from actual text nodes in the DOM, creating a bottom-up representation that preserves document structure while focusing on readable content.

**Key Components**:
- Text node discovery and analysis
- Hierarchical tree structure building
- Semantic context preservation
- Chrome integration for dynamic content

### [F3: Content Filter System](03-content-filter-system.md) ✅ **COMPLETED**
**Branch**: `feature/content-filter-system`

Implements sophisticated filtering to remove navigation, headers, footers, ads, and other non-content elements while preserving main article content.

**Key Components**:
- Rule-based filtering engine
- Semantic tag filtering
- CSS class pattern analysis
- Content density evaluation

### [F4: Image and Media Handler](04-image-media-handler.md) ✅ **COMPLETED**
**Branch**: `feature/image-media-handler`

Replaces images, videos, and other media with meaningful markdown equivalents that preserve semantic meaning and context.

**Key Components**:
- Media detection and analysis
- Alt text and caption extraction
- Context-aware description generation
- Multiple media type support

### [F5: Markdown Tree Renderer](05-markdown-tree-renderer.md) ✅ **COMPLETED**
**Branch**: `feature/markdown-tree-renderer`

Converts the processed content tree into clean, well-formatted markdown following best practices and maintaining document structure.

**Key Components**:
- Hierarchical markdown generation
- Style configuration options
- Clean formatting and spacing
- Quality assurance and validation

### [F6: TL;DR Summarizer](06-tldr-summarizer.md) ✅ **COMPLETED**
**Branch**: `feature/tldr-summarizer`

Integrates LLM-powered summarization to provide concise article summaries using OpenAI or compatible APIs.

**Key Components**:
- LLM integration (OpenAI, custom endpoints)
- Content summarization pipeline
- Multiple summary length options
- API key management and security

---

## Interactive Browser Mode Features (F7-F10)

The following features implement an interactive terminal browser mode with TUI capabilities:

### [F7: Footer with Page Statistics](../specs/features/F7-footer-stats.spec.md)
**Branch**: `feature/footer-stats`

Adds optional footer display with page statistics, metrics, and branding using the eszett (ß) character.

**Key Components**:
- `--footer` CLI flag for statistics display
- Page metrics (source, size, processing time)
- Eszett (ß) branding in pink color
- Pastel color scheme foundation

### [F8: Links Extraction and Display](../specs/features/F8-links-extraction.spec.md)
**Branch**: `feature/links-extraction`

Extracts important links from content and displays them in the footer with character codes for navigation.

**Key Components**:
- Smart link extraction and prioritization
- Character code assignment (a-z, aa-zz)
- Footer links section with visual hierarchy
- Link filtering and importance scoring

### [F9: Interactive TUI with Bubble Tea](../specs/features/F9-interactive-tui.spec.md)
**Branch**: `feature/interactive-tui`

Implements a full Terminal User Interface using Bubble Tea framework for interactive content browsing.

**Key Components**:
- `--interactive` flag for TUI mode
- Scrollable viewport with vim-like navigation
- Status bar with position indicators
- Consistent pastel color scheme application

### [F10: Link Navigation in TUI](../specs/features/F10-link-navigation.spec.md)
**Branch**: `feature/link-navigation`

Adds link navigation capabilities to the TUI mode with overlay interface and character code navigation.

**Key Components**:
- Links overlay with keyboard toggle
- Character code navigation system
- Navigation history and back button
- Loading states and error handling

## Development Workflow

Each feature follows the mandatory TDD workflow defined in the project:

1. **Create feature branch** from main
2. **Write executable specs** that initially fail
3. **Commit failing specs** with `SKIP=go-test`
4. **Implement incrementally** with small commits until specs pass
5. **Squash merge** to main and push

## Testing Strategy

### Spec-Driven Development
- Each feature includes comprehensive executable specifications
- Specs define expected behavior before implementation
- TDD approach ensures reliable, testable code

### Real-World Validation
- Testing against actual websites (news, blogs, documentation)
- Modern framework support (React, Vue, Next.js)
- Edge case handling and error recovery

### Performance Requirements
- Fast extraction (< 10 seconds for typical articles)
- Memory efficient processing
- Concurrent request support

## Configuration Philosophy

The new system emphasizes:
- **Configurability**: Extensive options for different use cases
- **Reasonable Defaults**: Works well out of the box
- **Extensibility**: Plugin system for custom requirements
- **Observability**: Detailed metrics and logging

## Migration Strategy

The new system will:
1. Coexist with the current extractor during development
2. Use feature flags for gradual rollout
3. Maintain backward compatibility for existing users
4. Provide migration path for custom configurations

## Implementation Status

### ✅ Core Content Extraction (F1-F6) - COMPLETED
The foundational content extraction system has been fully implemented and is production-ready:
- DOM event-driven content loading with JavaScript support
- Text node tree building for semantic content structure
- Sophisticated content filtering to remove boilerplate
- Image and media handling with markdown conversion
- Clean markdown rendering with configurable formatting
- LLM-powered TL;DR summarization

### 🚧 Interactive Browser Mode (F7-F10) - SPECS READY
Executable specifications have been created for the interactive browser features:
- Footer with page statistics and eszett branding
- Links extraction with character code navigation
- Full TUI implementation using Bubble Tea
- Interactive link navigation with history

## Expected Outcomes

Upon completion of all features, the system will provide:
- **Reliable extraction** across diverse website architectures ✅
- **High-quality markdown** output with proper formatting ✅
- **LLM integration** for content summarization ✅
- **Interactive terminal browser** with vim-like navigation 🚧
- **Link navigation system** with character codes 🚧
- **Comprehensive testing** ensuring production readiness ✅

## Getting Started

### For Core Features (F1-F6) - Already Implemented
The core content extraction features are complete and production-ready. See existing tests and implementation in:
- `internal/pageready/` - DOM readiness detection
- `internal/tree/` - Text node tree building
- `internal/filter/` - Content filtering
- `internal/media/` - Media handling
- `internal/markdown/` - Markdown rendering
- `internal/summarizer/` - TL;DR functionality

### For Interactive Browser Mode (F7-F10) - Ready for Implementation
To implement the interactive browser features:

1. Start with [F7: Footer with Page Statistics](../specs/features/F7-footer-stats.spec.md)
2. Continue with [F8: Links Extraction](../specs/features/F8-links-extraction.spec.md)
3. Implement [F9: Interactive TUI](../specs/features/F9-interactive-tui.spec.md)
4. Complete with [F10: Link Navigation](../specs/features/F10-link-navigation.spec.md)
5. Follow the TDD workflow for each feature
6. Ensure all specs pass before moving to the next feature

Each feature specification includes detailed technical requirements, acceptance criteria, test scenarios, and integration points to guide implementation.