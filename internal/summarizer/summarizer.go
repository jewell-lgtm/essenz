// Package summarizer provides AI-powered article summarization functionality.
package summarizer

import (
	"context"
	"fmt"
)

// SummaryLength represents the desired length of the summary.
type SummaryLength string

const (
	// ShortSummary generates a brief summary (1-2 sentences)
	ShortSummary SummaryLength = "short"
	// MediumSummary generates a moderate summary (3-5 sentences)
	MediumSummary SummaryLength = "medium"
	// LongSummary generates a detailed summary (6-10 sentences)
	LongSummary SummaryLength = "long"
)

// SummaryRequest contains the parameters for generating a summary.
type SummaryRequest struct {
	// Content is the article content to summarize (markdown format)
	Content string
	// Length specifies the desired summary length
	Length SummaryLength
	// Model specifies which AI model to use (e.g., "gpt-3.5-turbo", "gpt-4")
	Model string
	// Temperature controls creativity (0.0 = deterministic, 1.0 = creative)
	Temperature float32
}

// SummaryResponse contains the generated summary and metadata.
type SummaryResponse struct {
	// Summary is the generated summary text
	Summary string
	// Model is the model used for generation
	Model string
	// TokensUsed is the number of tokens consumed (if available)
	TokensUsed int
	// Duration is how long the summarization took
	Duration string
}

// Summarizer defines the interface for AI-powered summarization services.
type Summarizer interface {
	// Summarize generates a summary from the given content.
	Summarize(ctx context.Context, req SummaryRequest) (*SummaryResponse, error)

	// GetAvailableModels returns a list of available models.
	GetAvailableModels() []string

	// ValidateConfig checks if the summarizer is properly configured.
	ValidateConfig() error
}

// Config contains configuration for the summarizer.
type Config struct {
	// APIKey is the API key for the AI service
	APIKey string
	// DefaultModel is the default model to use if none specified
	DefaultModel string
	// DefaultLength is the default summary length if none specified
	DefaultLength SummaryLength
	// BaseURL allows overriding the API base URL (for testing/proxy)
	BaseURL string
	// Timeout is the request timeout duration
	Timeout string
}

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultModel:  "gpt-3.5-turbo",
		DefaultLength: MediumSummary,
		Timeout:       "60s",
	}
}

// Validate checks if the config is valid.
func (c Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if c.DefaultModel == "" {
		return fmt.Errorf("default model is required")
	}

	validLengths := map[SummaryLength]bool{
		ShortSummary:  true,
		MediumSummary: true,
		LongSummary:   true,
	}

	if !validLengths[c.DefaultLength] {
		return fmt.Errorf("invalid default length: %s", c.DefaultLength)
	}

	return nil
}

// GetPromptTemplate returns the prompt template for the given length.
func GetPromptTemplate(length SummaryLength) string {
	switch length {
	case ShortSummary:
		return `Please provide a concise summary of the following article in 1-2 sentences. Focus on the main point or conclusion:

%s

Summary:`

	case MediumSummary:
		return `Please provide a summary of the following article in 3-5 sentences. Include the main points and key details:

%s

Summary:`

	case LongSummary:
		return `Please provide a comprehensive summary of the following article in 6-10 sentences. Include the main points, supporting details, and conclusions:

%s

Summary:`

	default:
		return GetPromptTemplate(MediumSummary)
	}
}
