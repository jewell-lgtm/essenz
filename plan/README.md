# DOM Event-Driven Content Extraction Implementation Plan

This directory contains detailed specifications for implementing a new DOM event-driven content extraction system for the `sz` tool. The new approach moves away from heuristic-based content detection to a more reliable text-node-based tree building system.

## Implementation Order

The features must be implemented in the following order, as each builds on the previous:

### [F1: DOM Ready Event System](01-dom-ready-events.md)
**Branch**: `feature/dom-ready-events`

Implements proper DOM event waiting to ensure all dynamic content is loaded before extraction begins. This is the foundation for reliable content extraction across modern JavaScript frameworks.

**Key Components**:
- DOM readiness detection (DOMContentLoaded, load events)
- JavaScript framework initialization waiting
- Configurable timeouts and retry logic
- Chrome integration enhancement

### [F2: Text Node Tree Builder](02-text-node-tree-builder.md)
**Branch**: `feature/text-node-tree-builder`

Builds a content tree starting from actual text nodes in the DOM, creating a bottom-up representation that preserves document structure while focusing on readable content.

**Key Components**:
- Text node discovery and analysis
- Hierarchical tree structure building
- Semantic context preservation
- Chrome integration for dynamic content

### [F3: Content Filter System](03-content-filter-system.md)
**Branch**: `feature/content-filter-system`

Implements sophisticated filtering to remove navigation, headers, footers, ads, and other non-content elements while preserving main article content.

**Key Components**:
- Rule-based filtering engine
- Semantic tag filtering
- CSS class pattern analysis
- Content density evaluation

### [F4: Image and Media Handler](04-image-media-handler.md)
**Branch**: `feature/image-media-handler`

Replaces images, videos, and other media with meaningful markdown equivalents that preserve semantic meaning and context.

**Key Components**:
- Media detection and analysis
- Alt text and caption extraction
- Context-aware description generation
- Multiple media type support

### [F5: Markdown Tree Renderer](05-markdown-tree-renderer.md)
**Branch**: `feature/markdown-tree-renderer`

Converts the processed content tree into clean, well-formatted markdown following best practices and maintaining document structure.

**Key Components**:
- Hierarchical markdown generation
- Style configuration options
- Clean formatting and spacing
- Quality assurance and validation

### [F6: Integration and Testing](06-integration-testing.md)
**Branch**: `feature/integration-testing`

Integrates all components into a cohesive system and implements comprehensive testing across real-world websites.

**Key Components**:
- Pipeline integration architecture
- Performance optimization
- Comprehensive test suite
- Real-world website testing

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

## Expected Outcomes

Upon completion, the system will provide:
- **Reliable extraction** across diverse website architectures
- **High-quality markdown** output with proper formatting
- **Better performance** through optimized processing pipeline
- **Improved maintainability** with clean component separation
- **Comprehensive testing** ensuring production readiness

## Getting Started

To begin implementation:

1. Read the [overview](00-overview.md) for architectural context
2. Start with [F1: DOM Ready Events](01-dom-ready-events.md)
3. Follow the TDD workflow for each feature
4. Ensure all specs pass before moving to the next feature

Each feature specification includes detailed technical requirements, acceptance criteria, test scenarios, and integration points to guide implementation.