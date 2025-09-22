# F9 Interactive TUI with Bubble Tea Specification

## Overview

This specification defines the implementation of an interactive Terminal User Interface (TUI) using the Bubble Tea framework. The TUI displays web content in a scrollable viewport with vim-like navigation controls and applies a consistent pastel color scheme throughout.

## Requirements

### TUI Mode Activation
- Add `--interactive` CLI flag to launch TUI mode
- TUI mode replaces standard stdout output
- Graceful fallback to regular output if TUI unavailable

### Bubble Tea Integration
- Use Bubble Tea framework for TUI implementation
- Implement proper Model-View-Update (MVU) pattern
- Handle terminal resizing and cleanup gracefully

### Content Display
- Display processed markdown content in scrollable viewport
- Apply pastel color scheme to all content
- Show footer with statistics and links (if --footer enabled)
- Eszett (ß) character branding in status areas

### Vim-like Navigation
- `j` / `↓` - Scroll down one line
- `k` / `↑` - Scroll up one line
- `d` / `Ctrl+D` - Scroll down half page
- `u` / `Ctrl+U` - Scroll up half page
- `g` / `Home` - Go to top
- `G` / `End` - Go to bottom
- `q` / `Ctrl+C` - Quit application

### Status Bar
- Bottom status bar with eszett (ß) branding
- Show current position (line x of y)
- Display scroll indicators when applicable
- Use pink color for eszett character

### Color Scheme Application
- Content text: Light gray (#F5F5F5)
- Headers: Soft blue (#87CEEB)
- Links: Light blue (#ADD8E6)
- Emphasized text: Pale yellow (#FFFFE0)
- Background: Dark blue (#1E1E2E) for contrast
- Status bar: Muted purple (#9370DB)

## Test Cases

### TUI Mode Launch
```
GIVEN a URL or file to process
WHEN sz is run with --interactive flag
THEN TUI interface launches
AND content is displayed in scrollable viewport
AND vim navigation keys work
```

### Navigation Controls
```
GIVEN TUI mode is active with long content
WHEN vim navigation keys are pressed
THEN viewport scrolls appropriately
AND position indicators update
AND scrolling respects content boundaries
```

### Status Bar Display
```
GIVEN TUI mode is active
WHEN content is displayed
THEN status bar shows current position
AND eszett (ß) character appears in pink
AND scroll indicators show when applicable
```

### Terminal Resizing
```
GIVEN TUI mode is active
WHEN terminal is resized
THEN content reflows appropriately
AND viewport adjusts to new dimensions
AND navigation continues to work
```

### Graceful Exit
```
GIVEN TUI mode is active
WHEN quit keys are pressed (q, Ctrl+C)
THEN TUI exits cleanly
AND terminal state is restored
AND no artifacts remain
```

### Footer Integration
```
GIVEN TUI mode with --footer flag
WHEN content is displayed
THEN footer appears at bottom of content
AND footer statistics are visible
AND footer links (F8) are displayed if available
```

### Color Scheme Consistency
```
GIVEN TUI mode is active
WHEN content is displayed
THEN pastel color scheme is applied throughout
AND eszett character appears in pink
AND colors provide good readability
```

## Implementation Notes

- Follow TDD approach as outlined in CLAUDE.md
- Use Bubble Tea v0.24+ for modern TUI features
- Integrate with existing F7 footer and F8 links infrastructure
- Handle various terminal capabilities gracefully
- Ensure proper cleanup on exit

## TUI Architecture

### Model Structure
```go
type TUIModel struct {
    content     []string    // Processed content lines
    viewport    int         // Current scroll position
    height      int         // Terminal height
    width       int         // Terminal width
    footer      *Footer     // F7 footer data (optional)
    links       *Links      // F8 links data (optional)
    quitting    bool        // Exit flag
}
```

### Key Bindings
- Navigation: Vim-style movement keys
- Quit: `q`, `Ctrl+C`, `Esc` (when not in link mode)
- Help: `?` or `h` (future enhancement)

### Color Definitions
```go
var PastelColors = struct {
    Text        lipgloss.Color = "#F5F5F5"
    Headers     lipgloss.Color = "#87CEEB"
    Links       lipgloss.Color = "#ADD8E6"
    Emphasis    lipgloss.Color = "#FFFFE0"
    Background  lipgloss.Color = "#1E1E2E"
    StatusBar   lipgloss.Color = "#9370DB"
    Eszett      lipgloss.Color = "#FF69B4"  // Pink for ß
}
```

## Dependencies

- Bubble Tea framework (`github.com/charmbracelet/bubbletea`)
- Lipgloss for styling (`github.com/charmbracelet/lipgloss`)
- F7 footer implementation (optional integration)
- F8 links implementation (optional integration)

## Acceptance Criteria

1. `--interactive` flag launches TUI mode successfully
2. Bubble Tea framework properly integrated
3. Vim-like navigation keys work as specified
4. Content displays in scrollable viewport
5. Pastel color scheme applied consistently
6. Status bar shows position and eszett branding
7. Terminal resizing handled gracefully
8. Clean exit on quit commands
9. Footer and links integrated when available
10. Graceful fallback if TUI unavailable

## Future Enhancements (Out of Scope)

- Search functionality (`/` key)
- Help overlay (`?` key)
- Bookmarking support
- Multiple tab support
- Configuration file for color themes

---

**Note**: This specification must be implemented following the mandatory Git workflow described in CLAUDE.md, including writing failing tests first and using TDD approach. This feature builds on F7 and F8 for complete footer integration.