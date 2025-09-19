# F4: Image and Media Handler

**Feature Branch**: `feature/image-media-handler`

## Objective

Replace images, videos, and other media elements with meaningful markdown equivalents that preserve the semantic meaning and context of the media within the content flow.

## Technical Requirements

### Media Detection and Analysis
- **Image Elements**: Handle `<img>`, `<picture>`, CSS background images
- **Video Content**: Process `<video>`, embedded players (YouTube, Vimeo)
- **Audio Elements**: Handle `<audio>` and podcast embeds
- **Interactive Media**: Deal with `<canvas>`, `<svg>`, charts and diagrams
- **Social Embeds**: Process Twitter, Instagram, LinkedIn embeds

### Content Replacement Strategy
- **Alt Text Utilization**: Use existing alt attributes as primary content
- **Caption Detection**: Find associated captions, figcaptions, or titles
- **Context Analysis**: Analyze surrounding text for media description
- **URL Analysis**: Extract meaningful information from image URLs
- **Fallback Generation**: Generate descriptive text when no alt text exists

### Media Replacement Formats
```go
type MediaReplacement struct {
    Type        MediaType
    Description string
    URL         string
    Caption     string
    Context     string
    Dimensions  *Dimensions
    Alternative string  // Fallback description
}

type MediaType int
const (
    IMAGE MediaType = iota
    VIDEO
    AUDIO
    CHART
    DIAGRAM
    SOCIAL_EMBED
    INTERACTIVE
)

type Dimensions struct {
    Width  int
    Height int
}
```

## Implementation Components

### 1. Internal Package: `internal/media`
```go
type MediaHandler struct {
    config MediaConfig
    detectors []MediaDetector
    replacers []MediaReplacer
}

type MediaConfig struct {
    IncludeImageURLs     bool
    IncludeVideoDuration bool
    GenerateDescriptions bool
    MaxDescriptionLength int
    PreferAltText       bool
    ContextRadius       int  // Words around media for context
}

func (mh *MediaHandler) ProcessMediaInTree(root *ContentNode) error
func (mh *MediaHandler) DetectMedia(node *ContentNode) []MediaElement
func (mh *MediaHandler) GenerateReplacement(media MediaElement) string
```

### 2. Media Detection System
```go
type MediaDetector interface {
    CanHandle(node *ContentNode) bool
    Extract(node *ContentNode) []MediaElement
    Priority() int
}

type ImageDetector struct{}
type VideoDetector struct{}
type AudioDetector struct{}
type SocialEmbedDetector struct{}
type InteractiveMediaDetector struct{}
```

### 3. Context Analysis
```go
type ContextAnalyzer struct {
    radiusWords int
}

func (ca *ContextAnalyzer) ExtractContext(node *ContentNode, media MediaElement) string
func (ca *ContextAnalyzer) FindAssociatedCaption(node *ContentNode) string
func (ca *ContextAnalyzer) AnalyzeSurroundingText(node *ContentNode) string
```

### 4. Markdown Generation
```go
type MediaMarkdownGenerator struct {
    config GeneratorConfig
}

type GeneratorConfig struct {
    ImageFormat      string  // "![alt](url)" or "An image: alt"
    VideoFormat      string
    AudioFormat      string
    IncludeURLs      bool
    UseDescriptiveText bool
}

func (mg *MediaMarkdownGenerator) GenerateMarkdown(replacement MediaReplacement) string
```

## Acceptance Criteria

### Image Handling
1. Converts `<img>` tags to descriptive text using alt attributes
2. Extracts captions from `<figcaption>` elements
3. Handles `<picture>` elements with multiple sources
4. Processes CSS background images when contextually important

### Video and Audio
1. Replaces video embeds with descriptive text and duration if available
2. Handles YouTube/Vimeo embeds with title extraction
3. Processes audio elements with title and duration
4. Maintains context of media within content flow

### Alternative Text Generation
1. Generates meaningful descriptions when alt text is missing
2. Uses surrounding context to enhance descriptions
3. Analyzes image URLs for descriptive information
4. Provides fallback text for complex interactive media

### Social Media Embeds
1. Extracts meaningful content from Twitter embeds
2. Handles Instagram post descriptions
3. Processes LinkedIn article embeds
4. Maintains social context and attribution

## Test Scenarios

### Basic Image Handling
```html
<p>This is an article about cats.</p>
<img src="cat.jpg" alt="A fluffy orange cat sitting in a sunny window">
<p>Cats love warm sunny spots.</p>
```

Expected output:
```
This is an article about cats.

An image: A fluffy orange cat sitting in a sunny window

Cats love warm sunny spots.
```

### Complex Media Structure
```html
<figure>
    <img src="chart.png" alt="Sales growth chart">
    <figcaption>Sales increased 25% over the past quarter</figcaption>
</figure>
```

Expected output:
```
An image: Sales growth chart
*Sales increased 25% over the past quarter*
```

### Video Content
```html
<video controls>
    <source src="tutorial.mp4" type="video/mp4">
    Your browser does not support the video tag.
</video>
<p>This tutorial shows the basic workflow.</p>
```

Expected output:
```
A video: tutorial (MP4 format)

This tutorial shows the basic workflow.
```

### Missing Alt Text
```html
<p>The new office building has impressive architecture.</p>
<img src="office-building-exterior-modern-glass.jpg">
<p>The design won several awards.</p>
```

Expected output:
```
The new office building has impressive architecture.

An image (office building exterior, modern glass architecture)

The design won several awards.
```

### Social Media Embeds
```html
<blockquote class="twitter-tweet">
    <p>Excited to announce our new product launch! #innovation</p>
    <a href="https://twitter.com/company/status/123">@company</a>
</blockquote>
```

Expected output:
```
> Excited to announce our new product launch! #innovation
>
> — @company on Twitter
```

## Configuration Options

### Output Formats
```yaml
image_format: "descriptive"  # "descriptive" or "markdown"
video_format: "descriptive"
audio_format: "descriptive"

descriptive_format:
  template: "An image: {alt_text}"
  include_url: false
  include_dimensions: false

markdown_format:
  template: "![{alt_text}]({url})"
  fallback: "An image: {description}"
```

### Context Analysis
```yaml
context:
  radius_words: 20
  analyze_captions: true
  analyze_titles: true
  use_filename_hints: true
  generate_missing_alt: true
```

### Media Type Handling
```yaml
media_types:
  images:
    enabled: true
    include_decorative: false
  videos:
    enabled: true
    extract_duration: true
  audio:
    enabled: true
    extract_metadata: true
  social_embeds:
    enabled: true
    preserve_formatting: true
```

## Integration Points

### With Content Tree (F2)
- Operates on content tree nodes containing media elements
- Preserves tree structure while replacing media nodes
- Maintains proper parent-child relationships

### With Content Filter (F3)
- Coordinates with filtering to handle decorative vs. content images
- Respects filter decisions about media-heavy sections
- Preserves media within whitelisted content areas

### With Markdown Renderer (F5)
- Provides properly formatted media replacements for markdown output
- Ensures media descriptions integrate naturally with text flow
- Supports different markdown dialects and formatting preferences

## Advanced Features

### AI-Powered Descriptions (Future Enhancement)
- Integration with image recognition APIs
- Automatic generation of meaningful alt text
- Context-aware description enhancement
- Content type-specific description patterns

### Media Metadata Extraction
- EXIF data analysis for images
- Video duration and quality information
- Audio metadata (title, artist, duration)
- Social media engagement metrics

### Accessibility Improvements
- Screen reader optimized descriptions
- Color and visual element descriptions
- Motion and animation descriptions
- Interactive element accessibility text