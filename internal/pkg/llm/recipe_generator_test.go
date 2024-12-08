package llm_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tfiling/ai-chef/internal/pkg/llm"
)

type MockClaudeClient struct {
	mock.Mock
}

func (m *MockClaudeClient) CreateCompletion(ctx context.Context, req llm.CompletionRequest) (llm.CompletionResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(llm.CompletionResponse), args.Error(1)
}

type MockRecipeParser struct {
	mock.Mock
}

func (m *MockRecipeParser) ParseRecipeFromCompletion(completion string) (llm.Recipe, error) {
	args := m.Called(completion)
	return args.Get(0).(llm.Recipe), args.Error(1)
}

func TestRecipeGenerator_GenerateRecipe__SuccessfullyGenerateRecipe(t *testing.T) {
	// Arrange
	mockClient := new(MockClaudeClient)
	mockParser := new(MockRecipeParser)
	generator := llm.NewRecipeGenerator(mockClient, mockParser)

	req := llm.RecipeRequest{
		Style:       "Italian",
		ServingSize: 4,
		TimeLimit:   60,
	}

	expectedRecipe := llm.Recipe{
		ID:          "123",
		Name:        "Test Recipe",
		Description: "Description",
		CreatedAt:   time.Now(),
	}

	expectedCompletion := llm.CompletionResponse{
		Content: "recipe json",
	}

	mockClient.On("CreateCompletion", mock.Anything, mock.MatchedBy(func(req llm.CompletionRequest) bool {
		return req.SystemPrompt == llm.RecipeGenerationSystemPrompt
	})).Return(expectedCompletion, nil)

	mockParser.On("ParseRecipeFromCompletion", expectedCompletion.Content).Return(expectedRecipe, nil)

	// Act
	recipe, err := generator.GenerateRecipe(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedRecipe, recipe)
	mockClient.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}

func TestRecipeGenerator_GenerateRecipe__HandleGenerationFailure(t *testing.T) {
	// Arrange
	mockClient := new(MockClaudeClient)
	mockParser := new(MockRecipeParser)
	generator := llm.NewRecipeGenerator(mockClient, mockParser)

	req := llm.RecipeRequest{
		Style:       "Italian",
		ServingSize: 4,
		TimeLimit:   60,
	}

	expectedError := errors.New("api error")

	mockClient.On("CreateCompletion", mock.Anything, mock.Anything).Return(llm.CompletionResponse{}, expectedError)

	// Act
	recipe, err := generator.GenerateRecipe(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, recipe)
	mockClient.AssertExpectations(t)
	mockParser.AssertNotCalled(t, "ParseRecipeFromCompletion")
}

func TestRecipeGenerator_GenerateRecipe__HandleParserError(t *testing.T) {
	// Arrange
	mockClient := new(MockClaudeClient)
	mockParser := new(MockRecipeParser)
	generator := llm.NewRecipeGenerator(mockClient, mockParser)

	req := llm.RecipeRequest{
		Style:       "Italian",
		ServingSize: 4,
		TimeLimit:   60,
	}

	expectedCompletion := llm.CompletionResponse{
		Content: "invalid json",
	}
	expectedError := errors.New("parsing error")

	mockClient.On("CreateCompletion", mock.Anything, mock.Anything).Return(expectedCompletion, nil)
	mockParser.On("ParseRecipeFromCompletion", expectedCompletion.Content).Return(llm.Recipe{}, expectedError)

	// Act
	recipe, err := generator.GenerateRecipe(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Empty(t, recipe)
	mockClient.AssertExpectations(t)
	mockParser.AssertExpectations(t)
}
