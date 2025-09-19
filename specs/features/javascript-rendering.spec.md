# JavaScript Rendering Specification

## SPEC: Basic JavaScript Content Loading

GIVEN a page that renders content via JavaScript
WHEN essenz loads the page
THEN the JavaScript executes completely
AND the rendered content is extracted
AND no JavaScript code appears in the output

### Test Case
```spec
url: https://example.com/spa-article
wait_for: .article-body
timeout: 10s
expected_contains:
  - "Article Title"
  - "Main content paragraph"
expected_not_contains:
  - "function"
  - "document.getElementById"
  - "<script>"
```

## SPEC: React/Next.js Application Support

GIVEN a React-based single page application
WHEN the page uses client-side routing
THEN essenz waits for React to hydrate
AND extracts the fully rendered component tree
AND preserves semantic HTML structure

### Test Case
```spec
url: https://nextjs-blog.example.com/post/1
wait_strategy: network_idle
network_idle_duration: 2s
expected_structure:
  - tag: h1
    content: "Blog Post Title"
  - tag: article
    min_length: 500
  - tag: nav
    links_count: ">3"
```

## SPEC: Lazy-Loaded Content

GIVEN a page with lazy-loaded content
WHEN content loads as user scrolls
THEN essenz simulates scrolling
AND waits for new content to appear
AND captures all loaded sections

### Test Case
```spec
url: https://infinite-scroll.example.com
scroll_strategy: progressive
scroll_intervals:
  - distance: 1000px
    wait: 500ms
  - distance: 2000px
    wait: 1s
expected_blocks: ">=10"
max_scroll_time: 15s
```

## SPEC: Dynamic Content Updates

GIVEN a page that updates content dynamically
WHEN content changes without navigation
THEN essenz detects DOM mutations
AND waits for stability
AND captures the final state

### Test Case
```spec
url: https://live-updates.example.com
wait_strategy: dom_stable
stability_checks: 5
stability_interval: 200ms
expected_behavior:
  initial_content: "Loading..."
  final_content: "Article content loaded"
```

## SPEC: Authentication-Required Content

GIVEN a page behind authentication
WHEN cookies are provided
THEN essenz includes authentication cookies
AND accesses protected content
AND extracts authenticated view

### Test Case
```spec
url: https://members-only.example.com/article
cookies:
  - name: session_id
    value: "${SESSION_ID}"
    domain: members-only.example.com
expected_contains:
  - "Premium Content"
  - "Member-only article"
```

## SPEC: JavaScript Framework Detection

GIVEN various JavaScript frameworks
WHEN essenz analyzes the page
THEN it detects the framework in use
AND applies framework-specific wait strategies
AND reports framework in metadata

### Test Cases
```spec
frameworks:
  - url: https://angular-app.example.com
    expected_framework: angular
    wait_strategy: angular_stable
  - url: https://vue-app.example.com
    expected_framework: vue
    wait_strategy: vue_mounted
  - url: https://vanilla-js.example.com
    expected_framework: none
    wait_strategy: dom_stable
```

## SPEC: JavaScript Error Handling

GIVEN a page with JavaScript errors
WHEN errors occur during rendering
THEN essenz captures console errors
AND continues extraction attempt
AND reports errors in metadata

### Test Case
```spec
url: https://broken-js.example.com
continue_on_error: true
expected_metadata:
  js_errors:
    - contains: "TypeError"
    - contains: "undefined"
expected_output: partial_content
```

## SPEC: Progressive Enhancement

GIVEN a page with progressive enhancement
WHEN JavaScript is disabled
THEN essenz falls back to static HTML
AND extracts base content
AND notes degraded experience

### Test Case
```spec
url: https://progressive.example.com
javascript_enabled: false
expected_contains:
  - "Basic HTML content"
expected_metadata:
  javascript: disabled
  enhancement_level: basic
```

## SPEC: AJAX Content Loading

GIVEN a page that loads content via AJAX
WHEN content is fetched asynchronously
THEN essenz waits for XHR/fetch completion
AND extracts the loaded content
AND includes AJAX-loaded data

### Test Case
```spec
url: https://ajax-content.example.com
wait_for_xhr:
  - pattern: "/api/content"
  - pattern: "/api/comments"
timeout: 5s
expected_contains:
  - "Dynamically loaded content"
  - "User comments"
```

## SPEC: Web Components and Shadow DOM

GIVEN a page using Web Components
WHEN content is in Shadow DOM
THEN essenz pierces shadow boundaries
AND extracts encapsulated content
AND preserves component structure

### Test Case
```spec
url: https://web-components.example.com
pierce_shadow_dom: true
expected_structure:
  - web_component: "article-card"
    shadow_content:
      - tag: h2
        content: "Component Title"
      - tag: p
        content: "Shadow DOM content"
```