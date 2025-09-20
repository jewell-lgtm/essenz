package media

import (
	"fmt"
	"strings"
)

// MediaMarkdownGenerator generates markdown from media replacements.
type MediaMarkdownGenerator struct {
	config GeneratorConfig
}

// GeneratorConfig configures the markdown generation behavior.
type GeneratorConfig struct {
	ImageFormat        string // "descriptive" or "markdown"
	VideoFormat        string
	AudioFormat        string
	IncludeURLs        bool
	UseDescriptiveText bool
}

// NewMediaMarkdownGenerator creates a new MediaMarkdownGenerator.
func NewMediaMarkdownGenerator(config GeneratorConfig) *MediaMarkdownGenerator {
	return &MediaMarkdownGenerator{
		config: config,
	}
}

// GenerateMarkdown generates markdown text for a media replacement.
func (mg *MediaMarkdownGenerator) GenerateMarkdown(replacement MediaReplacement) string {
	switch replacement.Type {
	case IMAGE:
		return mg.generateImageMarkdown(replacement)
	case VIDEO:
		return mg.generateVideoMarkdown(replacement)
	case AUDIO:
		return mg.generateAudioMarkdown(replacement)
	case SOCIAL_EMBED:
		return mg.generateSocialMarkdown(replacement)
	case INTERACTIVE:
		return mg.generateInteractiveMarkdown(replacement)
	default:
		return mg.generateDefaultMarkdown(replacement)
	}
}

// generateImageMarkdown generates markdown for image elements.
func (mg *MediaMarkdownGenerator) generateImageMarkdown(replacement MediaReplacement) string {
	if mg.config.ImageFormat == "markdown" && replacement.URL != "" {
		// Standard markdown image format
		alt := replacement.Description
		if alt == "" {
			alt = replacement.Alternative
		}
		return fmt.Sprintf("![%s](%s)", alt, replacement.URL)
	}

	// Descriptive text format
	var parts []string

	// Main description
	description := replacement.Description
	if description == "" {
		description = replacement.Alternative
	}
	if description == "" {
		description = "image"
	}

	parts = append(parts, "An image: "+description)

	// Add caption if available
	if replacement.Caption != "" {
		parts = append(parts, "*"+replacement.Caption+"*")
	}

	return strings.Join(parts, "\n")
}

// generateVideoMarkdown generates markdown for video elements.
func (mg *MediaMarkdownGenerator) generateVideoMarkdown(replacement MediaReplacement) string {
	var parts []string

	// Main description
	description := replacement.Description
	if description == "" {
		description = replacement.Alternative
	}
	if description == "" {
		description = "video"
	}

	parts = append(parts, "A video: "+description)

	// Add caption if available
	if replacement.Caption != "" {
		parts = append(parts, "*"+replacement.Caption+"*")
	}

	return strings.Join(parts, "\n")
}

// generateAudioMarkdown generates markdown for audio elements.
func (mg *MediaMarkdownGenerator) generateAudioMarkdown(replacement MediaReplacement) string {
	var parts []string

	// Main description
	description := replacement.Description
	if description == "" {
		description = replacement.Alternative
	}
	if description == "" {
		description = "audio"
	}

	parts = append(parts, "An audio: "+description)

	// Add caption if available
	if replacement.Caption != "" {
		parts = append(parts, "*"+replacement.Caption+"*")
	}

	return strings.Join(parts, "\n")
}

// generateSocialMarkdown generates markdown for social media embeds.
func (mg *MediaMarkdownGenerator) generateSocialMarkdown(replacement MediaReplacement) string {
	// Format as blockquote
	content := replacement.Description
	if content == "" {
		content = replacement.Alternative
	}

	lines := strings.Split(content, "\n")
	var quotedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			quotedLines = append(quotedLines, "> "+line)
		}
	}

	return strings.Join(quotedLines, "\n")
}

// generateInteractiveMarkdown generates markdown for interactive media.
func (mg *MediaMarkdownGenerator) generateInteractiveMarkdown(replacement MediaReplacement) string {
	description := replacement.Description
	if description == "" {
		description = replacement.Alternative
	}
	if description == "" {
		description = "interactive element"
	}

	return "An interactive element: " + description
}

// generateDefaultMarkdown generates markdown for unknown media types.
func (mg *MediaMarkdownGenerator) generateDefaultMarkdown(replacement MediaReplacement) string {
	description := replacement.Description
	if description == "" {
		description = replacement.Alternative
	}
	if description == "" {
		description = "media element"
	}

	return "A media element: " + description
}
