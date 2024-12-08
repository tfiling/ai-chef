package llm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tfiling/ai-chef/internal/pkg/llm"
)

func TestRecipeParser_ParseRecipeFromCompletion__ValidJSON(t *testing.T) {
	// Arrange
	parser := llm.NewRecipeParser()
	expected := llm.Recipe{
		ID:          "123",
		Name:        "Test Recipe",
		Description: "Test Description",
		Ingredients: []llm.Ingredient{
			{
				Name:   "Ingredient 1",
				Amount: 1.0,
				Unit:   "cup",
			},
		},
		Instructions: "Test Instructions",
		PrepTime:     10,
		CookTime:     20,
		TotalTime:    30,
		ServingSize:  4,
		Difficulty:   "medium",
		CreatedAt:    time.Now(),
	}
	json := `{
		"id": "123",
		"name": "Test Recipe",
		"description": "Test Description",
		"ingredients": [
			{
				"name": "Ingredient 1",
				"amount": 1.0,
				"unit": "cup"
			}
		],
		"instructions": "Test Instructions",
		"prepTime": 10,
		"cookTime": 20,
		"totalTime": 30,
		"servingSize": 4,
		"difficulty": "medium",
		"createdAt": "2024-01-01T00:00:00Z"
	}`

	// Act
	result, err := parser.ParseRecipeFromCompletion(json)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
	assert.Equal(t, expected.Description, result.Description)
	assert.Equal(t, expected.Instructions, result.Instructions)
	assert.Equal(t, expected.PrepTime, result.PrepTime)
	assert.Equal(t, expected.CookTime, result.CookTime)
	assert.Equal(t, expected.TotalTime, result.TotalTime)
	assert.Equal(t, expected.ServingSize, result.ServingSize)
	assert.Equal(t, expected.Difficulty, result.Difficulty)
	assert.Len(t, result.Ingredients, 1)
	assert.Equal(t, expected.Ingredients[0].Name, result.Ingredients[0].Name)
	assert.Equal(t, expected.Ingredients[0].Amount, result.Ingredients[0].Amount)
	assert.Equal(t, expected.Ingredients[0].Unit, result.Ingredients[0].Unit)
}

func TestRecipeParser_ParseRecipeFromCompletion__InvalidJSON(t *testing.T) {
	// Arrange
	parser := llm.NewRecipeParser()
	invalidJSON := `{"name": Invalid JSON`

	// Act
	result, err := parser.ParseRecipeFromCompletion(invalidJSON)

	// Assert
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestRecipeParser_ParseRecipeFromCompletion__MissingRequiredFields(t *testing.T) {
	// Arrange
	parser := llm.NewRecipeParser()
	testCases := []struct {
		name string
		json string
	}{
		{
			name: "missing name",
			json: `{
				"id": "123",
				"description": "Test Description",
				"ingredients": [{"name": "Ingredient 1", "amount": 1.0, "unit": "cup"}],
				"instructions": "Test Instructions",
				"prepTime": 10,
				"cookTime": 20,
				"totalTime": 30,
				"servingSize": 4,
				"difficulty": "medium",
				"createdAt": "2024-01-01T00:00:00Z"
			}`,
		},
		{
			name: "missing instructions",
			json: `{
				"id": "123",
				"name": "Test Recipe",
				"description": "Test Description",
				"ingredients": [{"name": "Ingredient 1", "amount": 1.0, "unit": "cup"}],
				"prepTime": 10,
				"cookTime": 20,
				"totalTime": 30,
				"servingSize": 4,
				"difficulty": "medium",
				"createdAt": "2024-01-01T00:00:00Z"
			}`,
		},
		{
			name: "empty ingredients",
			json: `{
				"id": "123",
				"name": "Test Recipe",
				"description": "Test Description",
				"ingredients": [],
				"instructions": "Test Instructions",
				"prepTime": 10,
				"cookTime": 20,
				"totalTime": 30,
				"servingSize": 4,
				"difficulty": "medium",
				"createdAt": "2024-01-01T00:00:00Z"
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := parser.ParseRecipeFromCompletion(tc.json)

			// Assert
			assert.Error(t, err)
			assert.Empty(t, result)
		})
	}
}

func TestRecipeParser_ParseRecipeFromCompletion__InvalidIngredients(t *testing.T) {
	// Arrange
	parser := llm.NewRecipeParser()
	testCases := []struct {
		name string
		json string
	}{
		{
			name: "missing ingredient name",
			json: `{
				"id": "123",
				"name": "Test Recipe",
				"description": "Test Description",
				"ingredients": [{"amount": 1.0, "unit": "cup"}],
				"instructions": "Test Instructions",
				"prepTime": 10,
				"cookTime": 20,
				"totalTime": 30,
				"servingSize": 4,
				"difficulty": "medium",
				"createdAt": "2024-01-01T00:00:00Z"
			}`,
		},
		{
			name: "zero amount",
			json: `{
				"id": "123",
				"name": "Test Recipe",
				"description": "Test Description",
				"ingredients": [{"name": "Ingredient 1", "amount": 0, "unit": "cup"}],
				"instructions": "Test Instructions",
				"prepTime": 10,
				"cookTime": 20,
				"totalTime": 30,
				"servingSize": 4,
				"difficulty": "medium",
				"createdAt": "2024-01-01T00:00:00Z"
			}`,
		},
		{
			name: "missing unit",
			json: `{
				"id": "123",
				"name": "Test Recipe",
				"description": "Test Description",
				"ingredients": [{"name": "Ingredient 1", "amount": 1.0}],
				"instructions": "Test Instructions",
				"prepTime": 10,
				"cookTime": 20,
				"totalTime": 30,
				"servingSize": 4,
				"difficulty": "medium",
				"createdAt": "2024-01-01T00:00:00Z"
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := parser.ParseRecipeFromCompletion(tc.json)

			// Assert
			assert.Error(t, err)
			assert.Empty(t, result)
		})
	}
}
