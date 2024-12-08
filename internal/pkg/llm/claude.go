package llm

import (
	"context"
	"net/http"
	"time"
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

func NewClaudeClient(apiKey string, baseURL string) *ClaudeClient {
	return &ClaudeClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
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
