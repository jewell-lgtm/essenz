# LCP-Based Readiness (Default)

Goal

- Make Largest Contentful Paint (LCP) the default DOM readiness event so SPA and dynamic pages settle before extraction.

Behavior

- By default, sz injects a PerformanceObserver for `largest-contentful-paint` and waits for the first LCP entry up to the configured timeout (default 5s).
- On timeout, it falls back to legacy readiness (DOMContentLoaded/WaitReady).
- Custom selectors or framework hints are still honored after readiness triggers.
- Users can disable LCP waiting with `--no-lcp-wait`.

CLI

- `--dom-ready-timeout=5s` (default) controls max wait.
- `--wait-for-selector` and `--wait-for-frameworks` add extra conditions.
- `--no-lcp-wait` disables LCP readiness.
- `--debug-readiness` logs which condition fired (LCP, selector, framework) and timings.

Acceptance Criteria

- A page that renders the main article slightly after initial paint causes LCP to fire before timeout; extraction includes the article content.
- If LCP does not fire, the tool still proceeds via legacy readiness.
- Disabling LCP restores previous behavior.

Notes

- LCP may be an image or text. The observer stores a minimal payload on `window.__essenzLCP__` for potential scoring/metadata.
- Overlays (cookie/prompt) may become LCP if they dominate the viewport; content filtering should remove them even though readiness is satisfied.
