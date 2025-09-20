// Package media provides sophisticated media element handling for converting images,
// videos, and other media into meaningful markdown equivalents.
package media

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// MediaHandler processes media elements in a content tree and replaces them with descriptive text.
type MediaHandler struct {
	config    MediaConfig
	detectors []MediaDetector
	generator *MediaMarkdownGenerator
	analyzer  *ContextAnalyzer
}

// MediaConfig configures the media handling behavior.
type MediaConfig struct {
	IncludeImageURLs        bool
	IncludeVideoDuration    bool
	GenerateDescriptions    bool
	MaxDescriptionLength    int
	PreferAltText           bool
	ContextRadius           int // Words around media for context
	IncludeDecorativeImages bool
}

// MediaReplacement represents the replacement information for a media element.
type MediaReplacement struct {
	Type        MediaType
	Description string
	URL         string
	Caption     string
	Context     string
	Dimensions  *Dimensions
	Alternative string // Fallback description
}

// MediaType represents the type of media element.
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

// Dimensions represents the dimensions of a media element.
type Dimensions struct {
	Width  int
	Height int
}

// NewMediaHandler creates a new MediaHandler with default configuration.
func NewMediaHandler() *MediaHandler {
	handler := &MediaHandler{
		config: MediaConfig{
			IncludeImageURLs:        false,
			IncludeVideoDuration:    true,
			GenerateDescriptions:    true,
			MaxDescriptionLength:    200,
			PreferAltText:           true,
			ContextRadius:           20,
			IncludeDecorativeImages: false,
		},
		detectors: make([]MediaDetector, 0),
		analyzer:  NewContextAnalyzer(20),
	}

	// Add default detectors
	handler.AddDetector(NewImageDetector())
	handler.AddDetector(NewVideoDetector())
	handler.AddDetector(NewAudioDetector())
	handler.AddDetector(NewSocialEmbedDetector())
	handler.AddDetector(NewInteractiveMediaDetector())

	// Create markdown generator
	handler.generator = NewMediaMarkdownGenerator(GeneratorConfig{
		ImageFormat:        "descriptive", // "An image: {description}"
		VideoFormat:        "descriptive",
		AudioFormat:        "descriptive",
		IncludeURLs:        handler.config.IncludeImageURLs,
		UseDescriptiveText: true,
	})

	return handler
}

// WithConfig sets the media handler configuration.
func (mh *MediaHandler) WithConfig(config MediaConfig) *MediaHandler {
	mh.config = config
	mh.analyzer = NewContextAnalyzer(config.ContextRadius)
	return mh
}

// WithIncludeDecorative enables or disables decorative image inclusion.
func (mh *MediaHandler) WithIncludeDecorative(include bool) *MediaHandler {
	mh.config.IncludeDecorativeImages = include
	return mh
}

// AddDetector adds a media detector to the handler.
func (mh *MediaHandler) AddDetector(detector MediaDetector) {
	mh.detectors = append(mh.detectors, detector)
}

// ProcessMediaInTree processes all media elements in a content tree.
func (mh *MediaHandler) ProcessMediaInTree(ctx context.Context, root *tree.TextNode) error {
	return mh.processNode(ctx, root)
}

// processNode recursively processes a node and its children.
func (mh *MediaHandler) processNode(ctx context.Context, node *tree.TextNode) error {
	if node == nil {
		return nil
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Process current node if it's a media element
	if mh.isMediaElement(node) {
		replacement, err := mh.generateReplacement(node)
		if err != nil {
			return fmt.Errorf("failed to generate media replacement: %w", err)
		}

		if replacement != "" {
			// Convert media element to text node
			mh.replaceWithText(node, replacement)
		}
	}

	// Process children
	for _, child := range node.Children {
		if err := mh.processNode(ctx, child); err != nil {
			return err
		}
	}

	return nil
}

// isMediaElement checks if a node represents a media element.
func (mh *MediaHandler) isMediaElement(node *tree.TextNode) bool {
	if node == nil {
		return false
	}

	tag := strings.ToLower(node.Tag)
	switch tag {
	case "img", "picture", "video", "audio", "canvas", "svg":
		return true
	case "blockquote":
		// Check for social media embeds
		if class, exists := node.Attributes["class"]; exists {
			return strings.Contains(strings.ToLower(class), "twitter-tweet") ||
				strings.Contains(strings.ToLower(class), "instagram-media") ||
				strings.Contains(strings.ToLower(class), "linkedin-embed")
		}
	case "figure":
		// Figure elements containing media
		return mh.containsMediaChild(node)
	}

	return false
}

// containsMediaChild checks if a node contains media child elements.
func (mh *MediaHandler) containsMediaChild(node *tree.TextNode) bool {
	for _, child := range node.Children {
		if mh.isMediaElement(child) {
			return true
		}
		if mh.containsMediaChild(child) {
			return true
		}
	}
	return false
}

// generateReplacement generates a replacement string for a media element.
func (mh *MediaHandler) generateReplacement(node *tree.TextNode) (string, error) {
	// Detect media type and extract information
	var replacement MediaReplacement
	var detected bool

	for _, detector := range mh.detectors {
		if detector.CanHandle(node) {
			elements := detector.Extract(node)
			if len(elements) > 0 {
				element := elements[0] // Use first detected element
				replacement = mh.createReplacement(element, node)
				detected = true
				break
			}
		}
	}

	if !detected {
		// Fallback for unknown media types
		replacement = mh.createFallbackReplacement(node)
	}

	// Generate markdown using the replacement
	return mh.generator.GenerateMarkdown(replacement), nil
}

// createReplacement creates a MediaReplacement from a detected media element.
func (mh *MediaHandler) createReplacement(element MediaElement, node *tree.TextNode) MediaReplacement {
	replacement := MediaReplacement{
		Type:        element.Type,
		Description: element.Description,
		URL:         element.URL,
		Alternative: element.Alternative,
	}

	// Add context analysis
	replacement.Context = mh.analyzer.ExtractContext(node)
	replacement.Caption = mh.analyzer.FindAssociatedCaption(node)

	// Enhance description if needed
	if replacement.Description == "" && mh.config.GenerateDescriptions {
		replacement.Description = mh.generateDescriptionFromContext(replacement.Context, replacement.URL)
	}

	return replacement
}

// createFallbackReplacement creates a fallback replacement for unknown media.
func (mh *MediaHandler) createFallbackReplacement(node *tree.TextNode) MediaReplacement {
	description := "Unknown media element"

	// Try to extract some meaningful information
	if alt, exists := node.Attributes["alt"]; exists && alt != "" {
		description = alt
	} else if title, exists := node.Attributes["title"]; exists && title != "" {
		description = title
	}

	return MediaReplacement{
		Type:        IMAGE, // Default to image
		Description: description,
		Context:     mh.analyzer.ExtractContext(node),
		Alternative: description,
	}
}

// generateDescriptionFromContext generates a description from context and URL.
func (mh *MediaHandler) generateDescriptionFromContext(context, url string) string {
	var parts []string

	// Extract meaningful words from URL
	if url != "" {
		urlParts := mh.extractDescriptiveWordsFromURL(url)
		if len(urlParts) > 0 {
			parts = append(parts, strings.Join(urlParts, " "))
		}
	}

	// Add relevant context words
	if context != "" {
		contextWords := mh.extractRelevantContextWords(context)
		if len(contextWords) > 0 {
			parts = append(parts, strings.Join(contextWords, " "))
		}
	}

	if len(parts) > 0 {
		description := strings.Join(parts, ", ")
		if len(description) > mh.config.MaxDescriptionLength {
			description = description[:mh.config.MaxDescriptionLength] + "..."
		}
		return description
	}

	return "image" // Fallback
}

// extractDescriptiveWordsFromURL extracts meaningful words from a URL.
func (mh *MediaHandler) extractDescriptiveWordsFromURL(url string) []string {
	// Remove file extension and path separators
	filename := url
	if lastSlash := strings.LastIndex(url, "/"); lastSlash != -1 {
		filename = url[lastSlash+1:]
	}
	if lastDot := strings.LastIndex(filename, "."); lastDot != -1 {
		filename = filename[:lastDot]
	}

	// Split on common separators and filter meaningful words
	separators := regexp.MustCompile(`[-_\s]+`)
	words := separators.Split(filename, -1)

	var meaningful []string
	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if len(word) > 2 && !isCommonWord(word) {
			meaningful = append(meaningful, word)
		}
	}

	return meaningful
}

// extractRelevantContextWords extracts relevant words from context.
func (mh *MediaHandler) extractRelevantContextWords(context string) []string {
	words := strings.Fields(strings.ToLower(context))
	var relevant []string

	// Look for descriptive words near the media
	descriptivePatterns := []string{
		"architecture", "building", "design", "chart", "graph", "diagram",
		"photo", "picture", "illustration", "screenshot", "logo", "icon",
		"modern", "vintage", "colorful", "large", "small", "beautiful",
	}

	for _, word := range words {
		for _, pattern := range descriptivePatterns {
			if strings.Contains(word, pattern) {
				relevant = append(relevant, word)
				break
			}
		}
	}

	return relevant
}

// isCommonWord checks if a word is too common to be descriptive.
func isCommonWord(word string) bool {
	common := map[string]bool{
		"the": true, "and": true, "for": true, "are": true, "but": true,
		"not": true, "you": true, "all": true, "can": true, "her": true,
		"was": true, "one": true, "our": true, "had": true, "day": true,
		"get": true, "use": true, "man": true, "new": true, "now": true,
		"way": true, "may": true, "say": true, "img": true, "src": true,
		"alt": true, "jpg": true, "png": true, "gif": true, "jpeg": true,
	}
	return common[word]
}

// replaceWithText replaces a media node with descriptive text.
func (mh *MediaHandler) replaceWithText(node *tree.TextNode, replacement string) {
	// Clear children and attributes
	node.Children = nil
	node.Attributes = make(map[string]string)

	// Set as text node
	node.Tag = "#text"
	node.Text = replacement
}
