package summarizer

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenAISummarizer implements the Summarizer interface using OpenAI's API.
type OpenAISummarizer struct {
	client *openai.Client
	config Config
}

// NewOpenAISummarizer creates a new OpenAI-based summarizer.
func NewOpenAISummarizer(config Config) (*OpenAISummarizer, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAISummarizer{
		client: client,
		config: config,
	}, nil
}

// Summarize generates a summary using OpenAI's chat completion API.
func (s *OpenAISummarizer) Summarize(ctx context.Context, req SummaryRequest) (*SummaryResponse, error) {
	start := time.Now()

	// Use default values if not specified
	model := req.Model
	if model == "" {
		model = s.config.DefaultModel
	}

	length := req.Length
	if length == "" {
		length = s.config.DefaultLength
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.3 // Low temperature for more focused summaries
	}

	// Generate the prompt
	prompt := fmt.Sprintf(GetPromptTemplate(length), req.Content)

	// Create the chat completion request
	chatReq := openai.ChatCompletionRequest{
		Model:       model,
		Temperature: temperature,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that creates concise, accurate summaries of articles. Focus on the main points and avoid unnecessary details.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	// Make the API call
	resp, err := s.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned from API")
	}

	summary := resp.Choices[0].Message.Content
	if summary == "" {
		return nil, fmt.Errorf("empty summary returned from API")
	}

	duration := time.Since(start)

	return &SummaryResponse{
		Summary:    summary,
		Model:      model,
		TokensUsed: resp.Usage.TotalTokens,
		Duration:   duration.String(),
	}, nil
}

// GetAvailableModels returns a list of OpenAI models suitable for summarization.
func (s *OpenAISummarizer) GetAvailableModels() []string {
	return []string{
		"gpt-4",
		"gpt-4-turbo",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
	}
}

// ValidateConfig checks if the OpenAI summarizer is properly configured.
func (s *OpenAISummarizer) ValidateConfig() error {
	// Test the API key by making a simple request
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Make a minimal request to test the API key
	_, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo",
		MaxTokens: 1,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "test",
			},
		},
	})

	if err != nil {
		return fmt.Errorf("API key validation failed: %w", err)
	}

	return nil
}
