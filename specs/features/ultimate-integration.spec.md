# Ultimate Integration: SPA + Overlays + Dynamic Content + Paywall

Goal

- Verify sz can extract main article content when a site:
  - Renders via a SPA framework (client-side routing and hydration)
  - Shows a cookie consent flyover/overlay
  - Populates the article via `fetch`/XHR after initial paint
  - Presents a paywall overlay (DOM element) that still leaves content in the DOM

Assumptions

- If the content lands in the DOM (even offscreen or behind overlays), sz can locate and extract it.
- Headless Chrome is available; sz can wait for page readiness via selectors and/or network-idle.

Test Setup

- Local test page served via `httptest` (or fixture) that:
  - Loads a minimal SPA shell (e.g., `<div id="app"></div>`)
  - Immediately injects a cookie consent overlay (dismissible via a button)
  - Uses `fetch` to load article JSON after a short delay and appends it under `.main-content`
  - Applies a "paywall" overlay element above the page (higher z-index); the article content remains present in the DOM underneath

Acceptance Criteria

- Given the SPA page with overlays and delayed content load,
  - When running `sz --wait-for=".main-content .article-loaded" <url>`
  - Then the output markdown contains:
    - H1: the article title
    - Article body paragraphs
    - List elements, quotes, and inline formatting if present
  - And does not contain:
    - Cookie consent text/buttons
    - Paywall prompt text/buttons
    - Navigation/menu scaffolding unrelated to the article

Pseudocode Spec

- Arrange
  - Start `httptest` server serving an HTML page with:
    - `<div id="app"></div>`
    - Inline script to:
      - Render cookie overlay: `<div class="cookie-consent">…<button class="accept">Accept</button></div>`
      - After `setTimeout`, `fetch('/article.json')`, then insert
        `<main class="main-content"><article class="article-loaded">…</article></main>`
      - Render a `.paywall` overlay: positioned fixed with high z-index
  - `/article.json` returns JSON with title and body blocks

- Act
  - Run `sz` with options that are realistic for SPAs, e.g.:
    - `--wait-for=".article-loaded"` to ensure content is present
    - Optionally `--wait-idle=1s` if supported
    - `--markdown-renderer` to produce final markdown

- Assert
  - Markdown contains `# <title>`
  - Markdown contains at least two body paragraphs
  - Markdown excludes strings unique to overlays: "We use cookies", "Subscribe to continue"

Edge Cases

- Paywall overlay toggles `overflow: hidden` on `body`; ensure extraction isn’t affected.
- Cookie overlay injects into shadow DOM; verify we still skip it if it appears in text.
- Article loads in multiple chunks (title first, body later): use a more generous wait/readiness.

Future Enhancements

- Add a real executable spec alongside this document once daemon/network are available in CI.
- Provide a small fixture SPA (static assets) under `specs/fixtures/spa/` to keep the behavior predictable.

