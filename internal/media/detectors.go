package media

import (
	"strings"

	"github.com/jewell-lgtm/essenz/internal/tree"
)

// MediaDetector defines an interface for detecting and extracting media elements.
type MediaDetector interface {
	CanHandle(node *tree.TextNode) bool
	Extract(node *tree.TextNode) []MediaElement
	Priority() int
}

// MediaElement represents a detected media element.
type MediaElement struct {
	Type        MediaType
	Description string
	URL         string
	Alternative string
	Metadata    map[string]string
}

// ImageDetector handles image elements.
type ImageDetector struct{}

// NewImageDetector creates a new ImageDetector.
func NewImageDetector() *ImageDetector {
	return &ImageDetector{}
}

// CanHandle checks if this detector can handle the given node.
func (d *ImageDetector) CanHandle(node *tree.TextNode) bool {
	if node == nil {
		return false
	}
	tag := strings.ToLower(node.Tag)
	return tag == "img" || tag == "picture"
}

// Extract extracts image information from the node.
func (d *ImageDetector) Extract(node *tree.TextNode) []MediaElement {
	var elements []MediaElement

	tag := strings.ToLower(node.Tag)
	switch tag {
	case "img":
		element := MediaElement{
			Type: IMAGE,
			URL:  node.Attributes["src"],
		}

		// Prefer alt text for description
		if alt := node.Attributes["alt"]; alt != "" {
			element.Description = alt
		} else if title := node.Attributes["title"]; title != "" {
			element.Description = title
		}

		element.Alternative = element.Description
		if element.Alternative == "" {
			element.Alternative = "image"
		}

		elements = append(elements, element)

	case "picture":
		// For picture elements, look for img child or use first source
		for _, child := range node.Children {
			if strings.ToLower(child.Tag) == "img" {
				childElements := d.Extract(child)
				elements = append(elements, childElements...)
			}
		}
	}

	return elements
}

// Priority returns the priority of this detector.
func (d *ImageDetector) Priority() int {
	return 100
}

// VideoDetector handles video elements.
type VideoDetector struct{}

// NewVideoDetector creates a new VideoDetector.
func NewVideoDetector() *VideoDetector {
	return &VideoDetector{}
}

// CanHandle checks if this detector can handle the given node.
func (d *VideoDetector) CanHandle(node *tree.TextNode) bool {
	if node == nil {
		return false
	}
	return strings.ToLower(node.Tag) == "video"
}

// Extract extracts video information from the node.
func (d *VideoDetector) Extract(node *tree.TextNode) []MediaElement {
	element := MediaElement{
		Type: VIDEO,
	}

	// Try to find a source element
	var videoURL string
	var videoFormat string

	if src := node.Attributes["src"]; src != "" {
		videoURL = src
	} else {
		// Look for source children
		for _, child := range node.Children {
			if strings.ToLower(child.Tag) == "source" {
				if src := child.Attributes["src"]; src != "" {
					videoURL = src
					if srcType := child.Attributes["type"]; srcType != "" {
						videoFormat = srcType
					}
					break
				}
			}
		}
	}

	element.URL = videoURL

	// Generate description
	description := "video"
	if videoURL != "" {
		// Extract filename for description
		if lastSlash := strings.LastIndex(videoURL, "/"); lastSlash != -1 {
			filename := videoURL[lastSlash+1:]
			if lastDot := strings.LastIndex(filename, "."); lastDot != -1 {
				name := filename[:lastDot]
				description = strings.ReplaceAll(name, "-", " ")
				description = strings.ReplaceAll(description, "_", " ")
			}
		}
	}

	if videoFormat != "" {
		formatParts := strings.Split(videoFormat, "/")
		if len(formatParts) > 1 {
			description += " (" + strings.ToUpper(formatParts[1]) + " format)"
		}
	}

	element.Description = description
	element.Alternative = description

	return []MediaElement{element}
}

// Priority returns the priority of this detector.
func (d *VideoDetector) Priority() int {
	return 90
}

// AudioDetector handles audio elements.
type AudioDetector struct{}

// NewAudioDetector creates a new AudioDetector.
func NewAudioDetector() *AudioDetector {
	return &AudioDetector{}
}

// CanHandle checks if this detector can handle the given node.
func (d *AudioDetector) CanHandle(node *tree.TextNode) bool {
	if node == nil {
		return false
	}
	return strings.ToLower(node.Tag) == "audio"
}

// Extract extracts audio information from the node.
func (d *AudioDetector) Extract(node *tree.TextNode) []MediaElement {
	element := MediaElement{
		Type: AUDIO,
	}

	// Try to find a source element
	var audioURL string
	if src := node.Attributes["src"]; src != "" {
		audioURL = src
	} else {
		// Look for source children
		for _, child := range node.Children {
			if strings.ToLower(child.Tag) == "source" {
				if src := child.Attributes["src"]; src != "" {
					audioURL = src
					break
				}
			}
		}
	}

	element.URL = audioURL

	// Generate description
	description := "audio"
	if audioURL != "" {
		if lastSlash := strings.LastIndex(audioURL, "/"); lastSlash != -1 {
			filename := audioURL[lastSlash+1:]
			if lastDot := strings.LastIndex(filename, "."); lastDot != -1 {
				name := filename[:lastDot]
				description = strings.ReplaceAll(name, "-", " ")
				description = strings.ReplaceAll(description, "_", " ")
			}
		}
	}

	element.Description = description
	element.Alternative = description

	return []MediaElement{element}
}

// Priority returns the priority of this detector.
func (d *AudioDetector) Priority() int {
	return 80
}

// SocialEmbedDetector handles social media embeds.
type SocialEmbedDetector struct{}

// NewSocialEmbedDetector creates a new SocialEmbedDetector.
func NewSocialEmbedDetector() *SocialEmbedDetector {
	return &SocialEmbedDetector{}
}

// CanHandle checks if this detector can handle the given node.
func (d *SocialEmbedDetector) CanHandle(node *tree.TextNode) bool {
	if node == nil {
		return false
	}

	if strings.ToLower(node.Tag) == "blockquote" {
		if class := node.Attributes["class"]; class != "" {
			classLower := strings.ToLower(class)
			return strings.Contains(classLower, "twitter-tweet") ||
				strings.Contains(classLower, "instagram-media") ||
				strings.Contains(classLower, "linkedin-embed")
		}
	}

	return false
}

// Extract extracts social media embed information from the node.
func (d *SocialEmbedDetector) Extract(node *tree.TextNode) []MediaElement {
	element := MediaElement{
		Type: SOCIAL_EMBED,
	}

	class := strings.ToLower(node.Attributes["class"])
	var platform string

	if strings.Contains(class, "twitter-tweet") {
		platform = "Twitter"
	} else if strings.Contains(class, "instagram-media") {
		platform = "Instagram"
	} else if strings.Contains(class, "linkedin-embed") {
		platform = "LinkedIn"
	}

	// Extract content text
	content := d.extractTextContent(node)

	// Extract attribution
	attribution := d.extractAttribution(node, platform)

	element.Description = content
	element.Alternative = content
	if attribution != "" {
		element.Alternative += "\n\n— " + attribution + " on " + platform
	}

	return []MediaElement{element}
}

// extractTextContent extracts text content from social embed.
func (d *SocialEmbedDetector) extractTextContent(node *tree.TextNode) string {
	var textParts []string

	// Look for paragraph elements with content
	for _, child := range node.Children {
		if strings.ToLower(child.Tag) == "p" {
			text := d.getNodeText(child)
			if text != "" && !strings.HasPrefix(text, "http") {
				textParts = append(textParts, text)
			}
		}
	}

	return strings.Join(textParts, " ")
}

// extractAttribution extracts attribution information.
func (d *SocialEmbedDetector) extractAttribution(node *tree.TextNode, platform string) string {
	// Look for links that contain attribution
	for _, child := range node.Children {
		if strings.ToLower(child.Tag) == "a" {
			href := child.Attributes["href"]
			text := d.getNodeText(child)

			switch platform {
			case "Twitter":
				if strings.Contains(href, "twitter.com") && strings.HasPrefix(text, "@") {
					return text
				}
			case "Instagram":
				if strings.Contains(href, "instagram.com") {
					return text
				}
			case "LinkedIn":
				if strings.Contains(href, "linkedin.com") {
					return text
				}
			}
		}
	}

	return ""
}

// getNodeText recursively extracts text from a node.
func (d *SocialEmbedDetector) getNodeText(node *tree.TextNode) string {
	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var textParts []string
	for _, child := range node.Children {
		if text := d.getNodeText(child); text != "" {
			textParts = append(textParts, text)
		}
	}

	return strings.Join(textParts, " ")
}

// Priority returns the priority of this detector.
func (d *SocialEmbedDetector) Priority() int {
	return 70
}

// InteractiveMediaDetector handles interactive media elements.
type InteractiveMediaDetector struct{}

// NewInteractiveMediaDetector creates a new InteractiveMediaDetector.
func NewInteractiveMediaDetector() *InteractiveMediaDetector {
	return &InteractiveMediaDetector{}
}

// CanHandle checks if this detector can handle the given node.
func (d *InteractiveMediaDetector) CanHandle(node *tree.TextNode) bool {
	if node == nil {
		return false
	}
	tag := strings.ToLower(node.Tag)
	return tag == "canvas" || tag == "svg"
}

// Extract extracts interactive media information from the node.
func (d *InteractiveMediaDetector) Extract(node *tree.TextNode) []MediaElement {
	element := MediaElement{
		Type: INTERACTIVE,
	}

	tag := strings.ToLower(node.Tag)
	switch tag {
	case "canvas":
		element.Description = "interactive canvas element"
		if title := node.Attributes["title"]; title != "" {
			element.Description = title
		}
	case "svg":
		element.Description = "vector graphic"
		if title := node.Attributes["title"]; title != "" {
			element.Description = title
		}
		// Look for title element inside SVG
		for _, child := range node.Children {
			if strings.ToLower(child.Tag) == "title" {
				if titleText := d.getNodeText(child); titleText != "" {
					element.Description = titleText
					break
				}
			}
		}
	}

	element.Alternative = element.Description

	return []MediaElement{element}
}

// getNodeText recursively extracts text from a node.
func (d *InteractiveMediaDetector) getNodeText(node *tree.TextNode) string {
	if node.Tag == "#text" {
		return strings.TrimSpace(node.Text)
	}

	var textParts []string
	for _, child := range node.Children {
		if text := d.getNodeText(child); text != "" {
			textParts = append(textParts, text)
		}
	}

	return strings.Join(textParts, " ")
}

// Priority returns the priority of this detector.
func (d *InteractiveMediaDetector) Priority() int {
	return 60
}
