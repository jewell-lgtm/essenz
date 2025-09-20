package summarizer

import (
	"fmt"
	"os"
	"time"
)

// Provider represents different AI service providers.
type Provider string

const (
	// OpenAIProvider represents OpenAI as the AI service provider
	OpenAIProvider Provider = "openai"
)

// New creates a new summarizer instance based on the provider and config.
func New(provider Provider, config Config) (Summarizer, error) {
	switch provider {
	case OpenAIProvider:
		return NewOpenAISummarizer(config)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// NewFromEnvironment creates a summarizer using environment variables.
func NewFromEnvironment() (Summarizer, error) {
	config := DefaultConfig()

	// Check for API key in environment
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	// Allow override of default model
	if model := os.Getenv("ESSENZ_DEFAULT_MODEL"); model != "" {
		config.DefaultModel = model
	}

	// Allow override of base URL (useful for proxies/testing)
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		config.BaseURL = baseURL
	}

	// Allow override of timeout
	if timeout := os.Getenv("ESSENZ_TIMEOUT"); timeout != "" {
		config.Timeout = timeout
	}

	return New(OpenAIProvider, config)
}

// ConfigFromFlags creates a config from command-line flags and environment.
func ConfigFromFlags(apiKey, model, baseURL, timeout string) Config {
	config := DefaultConfig()

	// Priority: command-line flags > environment variables > defaults
	if apiKey != "" {
		config.APIKey = apiKey
	} else if envKey := os.Getenv("OPENAI_API_KEY"); envKey != "" {
		config.APIKey = envKey
	}

	if model != "" {
		config.DefaultModel = model
	} else if envModel := os.Getenv("ESSENZ_DEFAULT_MODEL"); envModel != "" {
		config.DefaultModel = envModel
	}

	if baseURL != "" {
		config.BaseURL = baseURL
	} else if envURL := os.Getenv("OPENAI_BASE_URL"); envURL != "" {
		config.BaseURL = envURL
	}

	if timeout != "" {
		config.Timeout = timeout
	} else if envTimeout := os.Getenv("ESSENZ_TIMEOUT"); envTimeout != "" {
		config.Timeout = envTimeout
	}

	return config
}

// ParseSummaryLength converts a string to a SummaryLength type.
func ParseSummaryLength(s string) (SummaryLength, error) {
	switch s {
	case "short", "s":
		return ShortSummary, nil
	case "medium", "m", "":
		return MediumSummary, nil
	case "long", "l":
		return LongSummary, nil
	default:
		return "", fmt.Errorf("invalid summary length: %s (valid options: short, medium, long)", s)
	}
}

// GetDefaultTimeout returns a reasonable default timeout duration.
func GetDefaultTimeout() time.Duration {
	return 60 * time.Second
}
