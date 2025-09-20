package markdown

import (
	"fmt"
	"strings"
)

// StyleManager handles formatting and style management for markdown output
type StyleManager struct {
	config RenderConfig
}

// NewStyleManager creates a new StyleManager with the given configuration
func NewStyleManager(config RenderConfig) *StyleManager {
	return &StyleManager{
		config: config,
	}
}

// FormatHeading formats a heading with the configured style
func (sm *StyleManager) FormatHeading(level int, text string) string {
	if text == "" {
		return ""
	}

	switch sm.config.HeadingStyle {
	case ATXHeading:
		prefix := strings.Repeat("#", level)
		return fmt.Sprintf("%s %s", prefix, text)
	case SetextHeading:
		if level == 1 {
			underline := strings.Repeat("=", len(text))
			return fmt.Sprintf("%s\n%s", text, underline)
		} else if level == 2 {
			underline := strings.Repeat("-", len(text))
			return fmt.Sprintf("%s\n%s", text, underline)
		} else {
			// Fallback to ATX for levels 3-6
			prefix := strings.Repeat("#", level)
			return fmt.Sprintf("%s %s", prefix, text)
		}
	default:
		prefix := strings.Repeat("#", level)
		return fmt.Sprintf("%s %s", prefix, text)
	}
}

// FormatEmphasis formats emphasis text with the configured style
func (sm *StyleManager) FormatEmphasis(text string) string {
	if text == "" {
		return ""
	}
	return sm.config.EmphasisStyle.Emphasis + text + sm.config.EmphasisStyle.Emphasis
}

// FormatStrong formats strong text with the configured style
func (sm *StyleManager) FormatStrong(text string) string {
	if text == "" {
		return ""
	}
	return sm.config.EmphasisStyle.Strong + text + sm.config.EmphasisStyle.Strong
}

// FormatInlineCode formats inline code with backticks
func (sm *StyleManager) FormatInlineCode(text string) string {
	if text == "" {
		return ""
	}
	return "`" + text + "`"
}

// FormatList formats a list with the configured style
func (sm *StyleManager) FormatList(items []string, ordered bool, level int) string {
	if len(items) == 0 {
		return ""
	}

	var result []string
	indent := strings.Repeat(" ", level*sm.config.ListStyle.IndentSize)

	for i, item := range items {
		var marker string
		if ordered {
			marker = fmt.Sprintf("%d. ", i+1)
		} else {
			marker = sm.config.ListStyle.UnorderedMarker + " "
		}
		result = append(result, indent+marker+item)
	}

	return strings.Join(result, "\n")
}

// FormatBlockquote formats a blockquote with > prefix
func (sm *StyleManager) FormatBlockquote(content string) string {
	if content == "" {
		return ""
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

// FormatCodeBlock formats a code block with the configured style
func (sm *StyleManager) FormatCodeBlock(content, language string) string {
	if content == "" {
		return ""
	}

	switch sm.config.CodeBlockStyle {
	case FencedCodeBlock:
		if language != "" {
			return fmt.Sprintf("```%s\n%s\n```", language, content)
		}
		return fmt.Sprintf("```\n%s\n```", content)
	case IndentedCodeBlock:
		lines := strings.Split(content, "\n")
		var indentedLines []string
		for _, line := range lines {
			indentedLines = append(indentedLines, "    "+line)
		}
		return strings.Join(indentedLines, "\n")
	default:
		if language != "" {
			return fmt.Sprintf("```%s\n%s\n```", language, content)
		}
		return fmt.Sprintf("```\n%s\n```", content)
	}
}

// FormatLink formats a link with the configured style
func (sm *StyleManager) FormatLink(text, url string) string {
	if url == "" {
		return text
	}
	if text == "" {
		text = url
	}
	return fmt.Sprintf("[%s](%s)", text, url)
}

// WrapText wraps text to the configured line width
func (sm *StyleManager) WrapText(text string, width int) string {
	if width <= 0 || len(text) <= width {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine []string
	currentLength := 0

	for _, word := range words {
		wordLength := len(word)

		// If adding this word would exceed the width, start a new line
		if currentLength > 0 && currentLength+1+wordLength > width {
			lines = append(lines, strings.Join(currentLine, " "))
			currentLine = []string{word}
			currentLength = wordLength
		} else {
			currentLine = append(currentLine, word)
			if currentLength > 0 {
				currentLength += 1 // Space
			}
			currentLength += wordLength
		}
	}

	// Add the last line
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return strings.Join(lines, "\n")
}

// EnsureProperSpacing ensures proper spacing between markdown elements
func (sm *StyleManager) EnsureProperSpacing(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for i, line := range lines {
		result = append(result, line)

		// Add spacing after headings
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) != "" {
				result = append(result, "")
			}
		}

		// Add spacing after paragraphs, lists, blockquotes
		if i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			currentLine := strings.TrimSpace(line)

			if currentLine != "" && nextLine != "" {
				// Check if we need spacing between different element types
				if sm.needsSpacing(currentLine, nextLine) {
					result = append(result, "")
				}
			}
		}
	}

	return strings.Join(result, "\n")
}

// needsSpacing determines if spacing is needed between two lines
func (sm *StyleManager) needsSpacing(current, next string) bool {
	// Add spacing before headings
	if strings.HasPrefix(next, "#") {
		return true
	}

	// Add spacing before blockquotes
	if strings.HasPrefix(next, ">") {
		return true
	}

	// Add spacing before code blocks
	if strings.HasPrefix(next, "```") {
		return true
	}

	// Add spacing before lists (but not between list items)
	if strings.HasPrefix(next, "-") || strings.HasPrefix(next, "*") || strings.HasPrefix(next, "+") {
		if !(strings.HasPrefix(current, "-") || strings.HasPrefix(current, "*") || strings.HasPrefix(current, "+")) {
			return true
		}
	}

	return false
}
