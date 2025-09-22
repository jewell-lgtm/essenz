# F7 Footer with Page Statistics Specification

## Overview

This specification defines the implementation of a footer feature that displays page statistics when the `--footer` flag is provided. The footer includes source information, size metrics, and branding using the eszett (ß) character in pink color.

## Requirements

### Footer Display
- Add `--footer` CLI flag to enable footer output
- Footer appears after main content with visual separator
- Uses eszett (ß) character as branding element in pink color
- Displays page statistics in pastel color scheme

### Statistics Included
- Source URL or file path
- Original page size (before processing)
- Output size (after markdown conversion)
- Processing time
- JavaScript enabled/disabled status

### Color Scheme
- Eszett (ß) character: Pink (#FF69B4)
- Statistics labels: Soft blue (#87CEEB)
- Values: Light gray (#D3D3D3)
- Separator line: Muted gray (#A9A9A9)

## Test Cases

### Basic Footer Display
```
GIVEN a web page URL
WHEN fetched with --footer flag
THEN output includes footer with page statistics
AND eszett character appears in pink
AND statistics are color-coded
```

### Local File Processing
```
GIVEN a local HTML file
WHEN processed with --footer flag
THEN footer shows file path as source
AND displays accurate size metrics
```

### Size Calculation Accuracy
```
GIVEN various content sizes
WHEN processed with --footer
THEN original size matches actual HTML content
AND output size matches generated markdown
```

### No Footer by Default
```
GIVEN any URL or file
WHEN processed without --footer flag
THEN no footer appears in output
AND only main content is displayed
```

## Implementation Notes

- Follow TDD approach as outlined in CLAUDE.md
- Use feature branch workflow from CLAUDE.md
- Footer should not interfere with existing output formats
- Color output should gracefully degrade for non-terminal environments
- Statistics calculation should be efficient and accurate

## Acceptance Criteria

1. `--footer` flag adds statistics footer to output
2. Eszett (ß) character appears in pink color
3. All statistics are accurately calculated and displayed
4. Pastel color scheme is applied consistently
5. Footer is visually separated from main content
6. No footer appears when flag is not provided
7. Works with both URLs and local files
8. Colors degrade gracefully in non-terminal environments

## Dependencies

- Terminal color support library (e.g., lipgloss, fatih/color)
- Size calculation utilities
- Time measurement for processing metrics

---

**Note**: This specification must be implemented following the mandatory Git workflow described in CLAUDE.md, including writing failing tests first and using TDD approach.