package llm

//TODO - test as a blackbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptComposer_ComposeCompletionRequest__ShouldComposeValidRequest(t *testing.T) {
	// Arrange
	req := RecipeRequest{
		Style:       "Italian",
		ServingSize: 4,
		TimeLimit:   60,
	}

	// Act
	result := composeRecipeRequest(req)

	// Assert
	assert.Equal(t, RecipeGenerationSystemPrompt, result.SystemPrompt)
	assert.Len(t, result.Messages, 1)
	assert.Equal(t, "user", result.Messages[0].Role)
	assert.Contains(t, result.Messages[0].Content, "Italian")
	assert.Contains(t, result.Messages[0].Content, "4 people")
	assert.Contains(t, result.Messages[0].Content, "60 minutes")
}
