# DOM Event-Driven Content Extraction Plan

## Overview

This plan implements a new approach to content extraction that waits for proper DOM readiness events and builds a tree structure from text nodes up, rather than trying to guess content areas through heuristics.

## Core Principles

1. **DOM Event-Driven**: Wait for the correct DOM events to indicate the page is ready for content extraction
2. **Text Node Foundation**: Start from actual text nodes and build the content tree from the bottom up
3. **Structural Awareness**: Maintain document structure while filtering out non-content elements
4. **Semantic Replacement**: Replace non-text elements (images, etc.) with meaningful markdown equivalents
5. **Filter by Context**: Exclude navigation, headers, footers, ads based on semantic tags and CSS classes

## Implementation Strategy

Each feature will be implemented as a separate branch following TDD methodology:

1. **F1: DOM Ready Event System** - Implement proper DOM event waiting
2. **F2: Text Node Tree Builder** - Build content tree from text nodes
3. **F3: Content Filter System** - Filter out non-content elements
4. **F4: Image and Media Handler** - Replace images/media with markdown equivalents
5. **F5: Markdown Tree Renderer** - Convert content tree to clean markdown
6. **F6: Integration and Testing** - Integrate all components and comprehensive testing

## Expected Benefits

- More reliable content extraction across different website architectures
- Better handling of modern JavaScript frameworks
- Cleaner separation of concerns
- More maintainable and testable codebase
- Improved markdown output quality

## Execution Order

Features should be implemented in numerical order, with each feature building on the previous ones. Each feature branch should include comprehensive specs that fail initially, then TDD implementation until specs pass.