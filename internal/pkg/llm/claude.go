package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
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

type RecipeRequest struct {
	Style       string
	ServingSize int
	TimeLimit   int
}

type Recipe struct {
	ID           string
	Name         string
	Description  string
	Ingredients  []Ingredient
	Instructions string
	PrepTime     int
	CookTime     int
	TotalTime    int
	ServingSize  int
	Difficulty   string
	CreatedAt    time.Time
}

type Ingredient struct {
	Name   string
	Amount float64
	Unit   string
}

type IClaudeClient interface {
	CreateCompletion(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
}

type IRecipeGenerator interface {
	GenerateRecipe(ctx context.Context, req RecipeRequest) (Recipe, error)
}

type IRecipeParser interface {
	ParseRecipeFromCompletion(completion string) (Recipe, error)
}

type ClaudeClient struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type RecipeGenerator struct {
	claudeClient IClaudeClient
}

func NewRecipeGenerator(claudeClient IClaudeClient) *RecipeGenerator {
	return &RecipeGenerator{
		claudeClient: claudeClient,
	}
}

type RecipeParser struct {
}

func NewRecipeParser() *RecipeParser {
	return &RecipeParser{}
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

func NewClaudeClient(apiKey string) *ClaudeClient {
	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

func (c *ClaudeClient) CreateCompletion(ctx context.Context, req CompletionRequest) (CompletionResponse, error) {
	claudeReq := claudeCompletionRequest{
		Model:     modelName,
		Messages:  req.Messages,
		System:    req.SystemPrompt,
		MaxTokens: 4096,
	}

	reqBody, err := json.Marshal(claudeReq)
	if err != nil {
		return CompletionResponse{}, errors.Wrap(err, "failed to marshal request")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/messages", bytes.NewBuffer(reqBody))
	if err != nil {
		return CompletionResponse{}, errors.Wrap(err, "failed to create HTTP request")
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return CompletionResponse{}, errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return CompletionResponse{}, errors.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var claudeResp claudeCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return CompletionResponse{}, errors.Wrap(err, "failed to decode response")
	}

	return CompletionResponse{
		ID:               claudeResp.ID,
		Content:          claudeResp.Content.Text,
		Model:            claudeResp.Model,
		PromptTokens:     claudeResp.Usage.InputTokens,
		CompletionTokens: claudeResp.Usage.OutputTokens,
		TotalTokens:      claudeResp.Usage.TotalTokens,
		FinishReason:     claudeResp.StopReason,
	}, nil
}

func (c *ClaudeClient) setHeaders(httpReq *http.Request) {
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
}
