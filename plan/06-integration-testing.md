# F6: Integration and Comprehensive Testing

**Feature Branch**: `feature/integration-testing`

## Objective

Integrate all previous components into a cohesive content extraction system and implement comprehensive testing across real-world websites, ensuring robust performance and high-quality markdown output.

## Technical Requirements

### System Integration
- **Pipeline Architecture**: Integrate F1-F5 into a seamless extraction pipeline
- **Error Handling**: Robust error handling and recovery at each stage
- **Configuration Management**: Unified configuration system across all components
- **Performance Optimization**: End-to-end performance tuning and optimization
- **Monitoring and Logging**: Comprehensive logging and metrics collection

### Integration Pipeline
```go
type ExtractionPipeline struct {
    readinessChecker *pageready.ReadinessChecker
    treeBuilder     *tree.TreeBuilder
    contentFilter   *filter.ContentFilter
    mediaHandler    *media.MediaHandler
    markdownRenderer *markdown.TreeRenderer
    config          *PipelineConfig
    logger          Logger
}

type PipelineConfig struct {
    ReadinessConfig  pageready.ReadinessConfig
    TreeConfig      tree.TreeConfig
    FilterConfig    filter.FilterConfig
    MediaConfig     media.MediaConfig
    RenderConfig    markdown.RenderConfig
    Timeout         time.Duration
    RetryAttempts   int
    DebugMode       bool
}

func (ep *ExtractionPipeline) Extract(ctx context.Context, url string) (*ExtractionResult, error)
func (ep *ExtractionPipeline) ExtractFromHTML(html string) (*ExtractionResult, error)
```

### Performance Metrics and Monitoring
```go
type ExtractionResult struct {
    Markdown        string
    Stats          ExtractionStats
    Warnings       []string
    ProcessingTime time.Duration
    Success        bool
}

type ExtractionStats struct {
    OriginalNodes    int
    FilteredNodes    int
    MediaElements    int
    TextLength       int
    WordCount        int
    HeadingCount     int
    ListCount        int
    LinkCount        int
    ProcessingStages map[string]time.Duration
}
```

## Implementation Components

### 1. Pipeline Integration
```go
// internal/pipeline/extractor.go
type Extractor struct {
    pipeline *ExtractionPipeline
    metrics  *MetricsCollector
}

func (e *Extractor) ProcessURL(ctx context.Context, url string) (*ExtractionResult, error) {
    // Stage 1: Wait for DOM readiness
    // Stage 2: Build content tree from DOM
    // Stage 3: Filter non-content elements
    // Stage 4: Process media elements
    // Stage 5: Render to markdown
    // Stage 6: Post-process and validate
}
```

### 2. Configuration Management
```go
// internal/config/pipeline.go
type Config struct {
    Pipeline PipelineConfig `yaml:"pipeline"`
    Chrome   ChromeConfig   `yaml:"chrome"`
    Output   OutputConfig   `yaml:"output"`
}

func LoadConfig(path string) (*Config, error)
func (c *Config) Validate() error
func (c *Config) ApplyDefaults()
```

### 3. Error Handling and Recovery
```go
type PipelineError struct {
    Stage   string
    Cause   error
    Context map[string]interface{}
    Retry   bool
}

func (pe *PipelineError) Error() string
func (pe *PipelineError) Unwrap() error
func (pe *PipelineError) ShouldRetry() bool
```

### 4. Comprehensive Testing Suite
- **Unit Tests**: Individual component testing with mocks
- **Integration Tests**: End-to-end pipeline testing
- **Performance Tests**: Benchmarking and load testing
- **Real-World Tests**: Testing against actual websites
- **Regression Tests**: Preventing quality degradation

## Acceptance Criteria

### Pipeline Integration
1. Successfully integrates all F1-F5 components into working system
2. Handles errors gracefully at each stage without crashing
3. Provides detailed metrics and timing information
4. Supports both URL fetching and direct HTML processing

### Performance Requirements
1. Processes typical article pages (2000-5000 words) in under 10 seconds
2. Handles large pages (10000+ words) without memory issues
3. Supports concurrent processing of multiple URLs
4. Gracefully handles slow or unresponsive websites

### Quality Assurance
1. Produces clean, readable markdown for 95% of tested websites
2. Preserves main content while filtering navigation/ads effectively
3. Handles edge cases without crashing or producing broken output
4. Maintains consistent output quality across different site types

### Real-World Website Testing
1. Successfully extracts content from major news websites
2. Handles modern JavaScript framework sites (React/Vue/Angular)
3. Processes blog platforms (WordPress, Medium, Ghost)
4. Works with documentation sites and technical content

## Test Website Categories

### News and Media Sites
- **Major News**: CNN, BBC, Reuters, Associated Press
- **Tech News**: TechCrunch, Ars Technica, The Verge
- **Blogs**: Medium articles, personal blogs, company blogs
- **Magazines**: The Atlantic, The New Yorker, Wired

### Technical Documentation
- **API Docs**: Stripe, GitHub, Anthropic documentation
- **Framework Docs**: React, Vue, Angular official docs
- **Language Docs**: Python, Go, JavaScript MDN
- **Tutorial Sites**: freeCodeCamp, tutorials, guides

### E-commerce and Business
- **Product Pages**: Amazon, tech product pages
- **Company Pages**: About pages, press releases
- **Help Centers**: Support documentation, FAQs
- **Landing Pages**: Marketing pages, service descriptions

### Modern Framework Sites
- **Next.js Sites**: Vercel, Next.js examples
- **React SPAs**: Complex single-page applications
- **Vue.js Sites**: Vue ecosystem websites
- **Static Sites**: Gatsby, Hugo, Jekyll sites

### Edge Cases and Challenges
- **Paywall Sites**: NYTimes, WSJ (public articles)
- **Heavy JavaScript**: Complex web applications
- **Poor HTML**: Malformed markup, legacy sites
- **International**: Non-English content, RTL languages

## Performance Benchmarks

### Target Metrics
```yaml
performance_targets:
  extraction_time:
    small_page: "< 2s"      # < 1000 words
    medium_page: "< 5s"     # 1000-5000 words
    large_page: "< 10s"     # > 5000 words

  memory_usage:
    max_heap: "100MB"       # Maximum memory per extraction
    gc_pressure: "low"      # Minimal garbage collection impact

  quality_metrics:
    content_recall: "> 95%" # Percentage of main content extracted
    noise_filtering: "> 90%" # Percentage of navigation/ads filtered
    markdown_validity: "100%" # All output must be valid markdown
```

### Load Testing
- Concurrent extractions (10, 50, 100 simultaneous requests)
- Memory leak detection over extended runs
- Performance degradation analysis
- Resource cleanup verification

## Integration Test Scenarios

### End-to-End Happy Path
1. Fetch complex modern website (e.g., Anthropic engineering blog)
2. Wait for full DOM readiness and framework initialization
3. Build complete content tree from all text nodes
4. Filter out navigation, headers, footers, ads
5. Process all images and media with proper descriptions
6. Render clean, well-formatted markdown
7. Validate output quality and completeness

### Error Recovery Testing
- Network timeouts during page loading
- JavaScript errors preventing full page load
- Malformed HTML causing parsing issues
- Missing or corrupted media elements
- Very large pages causing memory pressure

### Configuration Testing
- Different readiness timeout values
- Aggressive vs. conservative filtering settings
- Various markdown output formats and styles
- Custom CSS selector patterns for filtering
- Different image handling preferences

## Monitoring and Observability

### Metrics Collection
```go
type MetricsCollector struct {
    extractionCounter    prometheus.Counter
    extractionDuration   prometheus.Histogram
    errorCounter        prometheus.CounterVec
    contentQuality      prometheus.Gauge
}

func (mc *MetricsCollector) RecordExtraction(result *ExtractionResult)
func (mc *MetricsCollector) RecordError(stage string, err error)
func (mc *MetricsCollector) RecordQuality(stats ExtractionStats)
```

### Logging Strategy
- Structured logging with consistent fields
- Debug mode for detailed pipeline tracing
- Performance logging for optimization
- Error context preservation for debugging

### Health Checks
- Pipeline component health verification
- Chrome browser connectivity testing
- Memory usage monitoring
- Performance regression detection

## Documentation and Examples

### API Documentation
- Complete API reference for all components
- Configuration options and examples
- Best practices and common patterns
- Troubleshooting guides

### Usage Examples
- Command-line usage examples
- Programmatic API usage
- Configuration file examples
- Integration with other tools

### Performance Tuning Guide
- Optimization recommendations
- Resource usage guidelines
- Scaling considerations
- Monitoring setup instructions