# F8 Links Extraction and Footer Display Specification

## Overview

This specification defines the extraction of important links from page content and their display in the footer with single or double character codes. Links are shown only when `--footer` flag is enabled and hidden in non-interactive mode.

## Requirements

### Link Extraction
- Extract all important links from processed content
- Prioritize links by importance (content links over navigation)
- Filter out irrelevant links (ads, tracking, social media widgets)
- Assign unique single or double character codes (a-z, aa-zz)

### Link Codes Assignment
- Single character codes (a-z) for first 26 links
- Double character codes (aa-zz) for links 27-702
- Codes assigned in order of link importance
- Consistent codes for same URL within session

### Footer Links Section
- Links section appears in footer when `--footer` flag used
- Each link shows: `[code] Link Text → URL`
- Links are color-coded using pastel scheme
- Section header uses eszett (ß) character branding

### Non-Interactive Mode
- Link codes hidden in regular CLI output
- Only link text and URLs shown in footer
- Maintains clean, pipe-friendly output

## Test Cases

### Basic Link Extraction
```
GIVEN a page with multiple links
WHEN processed with --footer flag
THEN important links are extracted
AND each link gets a unique character code
AND links appear in footer section
```

### Link Prioritization
```
GIVEN a page with content and navigation links
WHEN links are extracted
THEN content links get priority (a, b, c, etc.)
AND navigation links get later codes
AND irrelevant links are filtered out
```

### Code Assignment Logic
```
GIVEN 30 important links on a page
WHEN codes are assigned
THEN first 26 get single characters (a-z)
AND remaining 4 get double characters (aa-ad)
AND codes are consistent for duplicate URLs
```

### Footer Integration
```
GIVEN extracted links with codes
WHEN footer is rendered
THEN links section includes eszett (ß) branding
AND each link shows format: [code] Text → URL
AND pastel color scheme is applied
```

### Non-Interactive Display
```
GIVEN extracted links
WHEN displayed in non-interactive mode
THEN codes are hidden from output
AND only link text and URLs shown
AND output remains clean for piping
```

### Link Filtering
```
GIVEN page with various link types
WHEN links are extracted
THEN content links are included
AND navigation links are included but deprioritized
AND ad/tracking/social widget links are excluded
```

## Implementation Notes

- Build on existing F7 footer infrastructure
- Follow TDD approach as outlined in CLAUDE.md
- Link extraction should integrate with existing content processing
- Importance scoring should align with content filtering algorithms
- Color scheme should match F7 footer implementation

## Link Importance Criteria

### High Priority (Single Character Codes)
- Links within main content area
- Article references and citations
- Related articles and content
- Author profiles and bylines

### Medium Priority (Later Codes)
- Navigation menu items
- Category and tag links
- Search and archive links
- Contact and about pages

### Excluded Links
- Social media sharing buttons
- Advertisement links
- Tracking and analytics URLs
- Cookie policy and privacy links
- Login/logout links

## Color Scheme (Pastel)

- Link codes: Soft green (#90EE90)
- Link text: Light blue (#ADD8E6)
- URLs: Light gray (#D3D3D3)
- Section header with ß: Pink (#FF69B4)

## Acceptance Criteria

1. Important links are extracted from page content
2. Links are prioritized by relevance and importance
3. Unique character codes assigned (a-z, then aa-zz)
4. Links section appears in footer with eszett branding
5. Pastel color scheme applied consistently
6. Codes hidden in non-interactive mode
7. Irrelevant links properly filtered out
8. Works with both URLs and local HTML files
9. Consistent code assignment for duplicate URLs
10. Integration with existing F7 footer display

## Dependencies

- F7 footer infrastructure (must be implemented first)
- Link extraction utilities
- URL normalization and deduplication
- Color scheme from F7 implementation

---

**Note**: This specification must be implemented following the mandatory Git workflow described in CLAUDE.md, including writing failing tests first and using TDD approach. This feature builds on F7 and requires it to be implemented first.