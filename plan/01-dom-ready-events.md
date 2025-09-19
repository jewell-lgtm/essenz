# F1: DOM Ready Event System

**Feature Branch**: `feature/dom-ready-events`

## Objective

Implement a robust DOM event system that waits for the appropriate DOM readiness events before attempting content extraction, ensuring all dynamic content is loaded.

## Technical Requirements

### Chrome Integration Enhancement
- Enhance the existing Chrome browser integration to wait for proper DOM events
- Implement event listeners for `DOMContentLoaded`, `load`, and custom readiness indicators
- Add configurable timeouts for different types of content loading
- Support for waiting on JavaScript framework initialization events

### Event Detection Logic
- **DOMContentLoaded**: Basic DOM structure is ready
- **Load Event**: All resources including images, stylesheets loaded
- **Framework Events**: React/Vue/Angular app initialization completion
- **Custom Selectors**: Wait for specific CSS selectors to appear (configurable)
- **Timeout Handling**: Graceful fallback after maximum wait time

### Configuration Options
- Maximum wait time (default: 5 seconds)
- Framework detection patterns (React, Vue, Angular, Next.js)
- Custom readiness selectors
- Retry logic for failed page loads

## Implementation Components

### 1. Internal Package: `internal/pageready`
```go
type ReadinessChecker struct {
    MaxWaitTime     time.Duration
    FrameworkHints  []string
    CustomSelectors []string
}

type ReadinessResult struct {
    IsReady     bool
    EventType   string
    WaitTime    time.Duration
    Error       error
}

func (r *ReadinessChecker) WaitForReady(ctx context.Context, page chromedp.Page) (*ReadinessResult, error)
```

### 2. Chrome Integration Updates
- Enhance `internal/browser` package to use readiness checking
- Add JavaScript injection for framework detection
- Implement event waiting with chromedp

### 3. Configuration Integration
- Add readiness options to extractor configuration
- Environment variable support for wait times
- CLI flags for debugging readiness issues

## Acceptance Criteria

### Core Functionality
1. Successfully waits for DOMContentLoaded event before extraction
2. Detects and waits for common JavaScript framework initialization
3. Handles timeout scenarios gracefully without hanging
4. Returns detailed information about what event triggered readiness

### Framework Support
1. Detects React application readiness
2. Detects Vue.js application readiness
3. Detects Next.js hydration completion
4. Falls back to standard DOM events for static sites

### Error Handling
1. Graceful timeout handling with partial content extraction
2. Network error recovery with retry logic
3. Invalid page handling (404, 500 errors)
4. Logging of readiness detection process for debugging

## Test Scenarios

### Static HTML Sites
- Basic HTML page loads and triggers DOMContentLoaded
- Page with external CSS/JS loads and triggers load event
- Slow-loading page times out gracefully

### JavaScript Framework Sites
- React SPA completes hydration before extraction begins
- Next.js page waits for client-side rendering
- Vue.js app initialization detection
- Angular app bootstrap completion

### Edge Cases
- Pages that never trigger readiness events (timeout handling)
- Pages with broken JavaScript (fallback behavior)
- Very slow networks (extended timeout behavior)
- Pages with infinite loading states

## Integration Points

### With Existing Code
- Integrates with `internal/browser` Chrome automation
- Used by `internal/extractor` before content extraction begins
- Configurable via CLI flags and environment variables

### With Future Features
- Provides foundation for text node tree building (F2)
- Ensures page stability before content filtering (F3)
- Critical for reliable extraction across different site types