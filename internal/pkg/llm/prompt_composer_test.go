package llm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tfiling/ai-chef/internal/pkg/llm"
)

func TestPromptComposer_ComposeCompletionRequest__ShouldComposeValidRequest(t *testing.T) {
	// Arrange
	composer := llm.NewPromptComposer()
	req := llm.RecipeRequest{
		Style:       "Italian",
		ServingSize: 4,
		TimeLimit:   60,
	}

	// Act
	result := composer.ComposeCompletionRequest(req)

	// Assert
	assert.Equal(t, llm.RecipeGenerationSystemPrompt, result.SystemPrompt)
	assert.Len(t, result.Messages, 1)
	assert.Equal(t, "user", result.Messages[0].Role)
	assert.Contains(t, result.Messages[0].Content, "Italian")
	assert.Contains(t, result.Messages[0].Content, "4 people")
	assert.Contains(t, result.Messages[0].Content, "60 minutes")
}
