# Code Quality Improvement Checklist

## Test Failures to Fix (5 specs failing)

### 1. ContentFilterSpec - 5 sub-tests failing
- [ ] link_density_filtering - "Brief text" not being filtered
- [ ] whitelist_preservation - Navigation links not preserved in whitelisted containers
- [ ] edge_case_handling - Resource lists being filtered too aggressively
- [ ] filter_configuration - Custom navigation not preserved

### 2. DOMReadyEventsSpec - Multiple failures
- [ ] Basic DOM ready not extracting "main content"
- [ ] React framework detection not working ("React Article" missing)
- [ ] Timeout handling issues
- [ ] Custom selector waiting problems

### 3. ImageMediaHandlerSpec
- [ ] Build failure - need to investigate

### 4. MarkdownTreeRendererSpec
- [ ] Build failure - need to investigate

### 5. TextNodeTreeBuilderSpec
- [ ] Dynamic content handling not capturing JavaScript-generated content

## Linting Issues to Fix

### High Priority - Cyclomatic Complexity
- [ ] internal/tree/builder.go:102 - traverseNode complexity 26 (> 15)
- [ ] internal/filter/length_filter.go:81 - hasImportantChildren complexity 16 (> 15)

### Medium Priority - Stuttering Type Names
- [ ] tree.TreeBuilder → tree.Builder
- [ ] filter.FilterConfig → filter.Config
- [ ] filter.FilterRule → filter.Rule
- [ ] filter.FilterContext → filter.Context
- [ ] filter.FilterStats → filter.Stats
- [ ] media.MediaHandler → media.Handler
- [ ] media.MediaConfig → media.Config
- [ ] media.MediaReplacement → media.Replacement
- [ ] media.MediaDetector → media.Detector
- [ ] media.MediaMarkdownGenerator → media.MarkdownGenerator
- [ ] media.MediaElement → media.Element
- [ ] media.MediaType → media.Type

### Low Priority - Formatting
- [ ] internal/markdown/renderer.go:59 - gofmt formatting issue

## Chrome Process Cleanup
- [ ] Kill any lingering Chrome processes
- [ ] Ensure proper daemon cleanup in tests

## Documentation
- [ ] Add package-level documentation for all internal packages
- [ ] Document public APIs and interfaces
- [ ] Add examples for main user-facing commands

## Final Verification
- [ ] All tests passing
- [ ] Zero linting issues
- [ ] Build succeeds without warnings
- [ ] Manual smoke test of all features