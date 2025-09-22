# F10 Link Navigation in TUI Specification

## Overview

This specification defines the implementation of link navigation within the interactive TUI mode. Users can show/hide a links overlay, view extracted links with their character codes, and navigate to new pages using those codes.

## Requirements

### Links Overlay Activation
- Keyboard shortcut to toggle links overlay (default: `l`)
- Overlay shows extracted links with character codes
- Links displayed in scrollable list if many links
- ESC key closes the overlay

### Links Display
- Show format: `[code] Link Text → URL`
- Links sorted by importance (matching F8 extraction)
- Current selection highlighted
- Eszett (ß) character in overlay header
- Pastel color scheme applied

### Navigation Controls
- Type link code (a-z, aa-zz) to navigate to link
- Navigation loads new page in same TUI session
- Loading indicator while fetching new content
- Error handling for invalid URLs or failed fetches

### TUI Integration
- Links overlay appears over main content
- Main content dimmed when overlay active
- Smooth transitions between overlay states
- Maintains scroll position when returning

### Session Management
- New page replaces current content
- Previous page navigation (back button / `b` key)
- History stack for multiple page navigation
- Current URL shown in status bar

## Test Cases

### Links Overlay Toggle
```
GIVEN TUI mode is active with extracted links
WHEN links toggle key is pressed
THEN links overlay appears over content
AND extracted links are displayed with codes
AND main content is dimmed
```

### Link Code Navigation
```
GIVEN links overlay is active
WHEN a valid link code is typed (e.g., 'a')
THEN TUI navigates to that link
AND new content loads and displays
AND overlay closes automatically
```

### Invalid Link Code
```
GIVEN links overlay is active
WHEN an invalid code is typed
THEN error message briefly displays
AND overlay remains open
AND user can try again
```

### Overlay Closure
```
GIVEN links overlay is active
WHEN ESC key is pressed
THEN overlay closes
AND main content returns to normal view
AND previous scroll position restored
```

### Navigation History
```
GIVEN user has navigated to multiple pages
WHEN back navigation key is pressed
THEN previous page content loads
AND history stack is maintained
AND current URL updates in status bar
```

### Loading States
```
GIVEN user navigates to a new link
WHEN page is loading
THEN loading indicator displays
AND eszett (ß) character shows loading animation
AND user can cancel with ESC
```

### Error Handling
```
GIVEN user navigates to invalid/broken link
WHEN navigation fails
THEN error message displays
AND user remains on current page
AND can try different link or close overlay
```

### Long Link Lists
```
GIVEN page with many links (>20)
WHEN links overlay is shown
THEN links are scrollable within overlay
AND scroll position indicators shown
AND vim-like navigation works in overlay
```

## Implementation Notes

- Build on F8 links extraction and F9 TUI infrastructure
- Follow TDD approach as outlined in CLAUDE.md
- Use Bubble Tea's layered model for overlay
- Integrate with existing content fetching pipeline
- Handle loading states gracefully

## Overlay Design

### Layout
```
┌─────────────────────────────────────────┐
│ ß Links                          [l] ×  │
├─────────────────────────────────────────┤
│ [a] Article Title → example.com/article │
│ [b] Related Post → example.com/related  │
│ [c] Author Bio → example.com/author     │
│ [d] Comments → example.com/comments     │
│ ...                                     │
├─────────────────────────────────────────┤
│ Type code to navigate • ESC to close    │
└─────────────────────────────────────────┤
```

### Key Bindings in Overlay
- `l` - Toggle overlay (same key to close)
- `ESC` - Close overlay
- `a-z`, `aa-zz` - Navigate to link with that code
- `j/k` or `↑/↓` - Scroll through links (if many)
- `q` - Quit TUI entirely

### Color Scheme in Overlay
- Overlay background: Semi-transparent dark (#1E1E2E80)
- Border: Soft blue (#87CEEB)
- Selected link: Highlighted background (#9370DB40)
- Link codes: Soft green (#90EE90)
- Link text: Light blue (#ADD8E6)
- URLs: Light gray (#D3D3D3)
- Eszett header: Pink (#FF69B4)

## Navigation Flow

### Link Selection Process
1. User presses `l` to show links overlay
2. Links display with character codes
3. User types desired code (e.g., `a`)
4. TUI validates code and starts navigation
5. Loading indicator shows during fetch
6. New content replaces current content
7. Overlay closes, focus returns to content

### History Management
```go
type NavigationHistory struct {
    pages    []PageState
    current  int
    maxSize  int // e.g., 50 pages
}

type PageState struct {
    url      string
    content  []string
    position int  // scroll position
}
```

## Dependencies

- F8 links extraction (must be implemented first)
- F9 interactive TUI (must be implemented first)
- Bubble Tea layered models for overlay
- Existing content fetching pipeline
- URL validation utilities

## Acceptance Criteria

1. Links overlay toggles with keyboard shortcut
2. Extracted links display with character codes
3. Link codes navigate to corresponding URLs
4. Overlay closes on ESC or after navigation
5. Loading states handled gracefully
6. Error handling for invalid links/codes
7. Navigation history maintained
8. Back navigation works correctly
9. Status bar shows current URL
10. Pastel color scheme applied throughout
11. Eszett branding in overlay header
12. Long link lists scroll properly

## Future Enhancements (Out of Scope)

- Link previews on hover
- Bookmark management
- Multiple tab support
- Link filtering and search
- Custom key bindings

---

**Note**: This specification must be implemented following the mandatory Git workflow described in CLAUDE.md, including writing failing tests first and using TDD approach. This feature requires F7, F8, and F9 to be implemented first and builds on their infrastructure for complete interactive browser functionality.