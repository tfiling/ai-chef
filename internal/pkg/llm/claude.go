package llm

import (
	"context"
	"net/http"
	"time"
)

const (
	RoleSystem    = "system"
	RoleAssistant = "assistant"
	RoleHuman     = "human"
)

const (
	baseURL   = "https://api.anthropic.com/v1"
	modelName = "claude-3-5-sonnet-latest"
)

type Message struct {
	Role    string
	Content string
}

type CompletionRequest struct {
	Messages     []Message
	SystemPrompt string
}

type CompletionResponse struct {
	ID               string
	Content          string
	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	FinishReason     string
	CreatedAt        time.Time
}

type IClaudeClient interface {
	CreateCompletion(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

type IRecipeGenerator interface {
	GenerateRecipe(ctx context.Context, req RecipeRequest) (Recipe, error)
}

type ClaudeClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type claudeCompletionRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
	MaxTokens int       `json:"max_tokens"`
}

type claudeCompletionResponse struct {
	ID           string  `json:"id"`
	Model        string  `json:"model"`
	Content      message `json:"content"`
	Usage        usage   `json:"usage"`
	StopReason   string  `json:"stop_reason"`
	StopSequence string  `json:"stop_sequence"`
}

type message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
