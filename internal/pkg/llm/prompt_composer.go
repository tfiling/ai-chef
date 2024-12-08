package llm

import "fmt"

// TODO - manually improve prompts
const (
	RecipeGenerationSystemPrompt = `You are an expert chef and recipe creator. Generate detailed recipes based on the given requirements.
Follow these rules:
1. Create recipes that are feasible and safe to make at home
2. Include exact measurements and cooking times
3. List ingredients with precise amounts and common units
4. Provide clear step-by-step instructions
5. Format the output as a JSON object with the following structure:
{
  "name": "Recipe Name",
  "description": "Brief description",
  "ingredients": [{"name": "ingredient", "amount": number, "unit": "unit"}],
  "instructions": "Numbered steps",
  "prepTime": number (minutes),
  "cookTime": number (minutes),
  "totalTime": number (minutes),
  "servingSize": number,
  "difficulty": "easy|medium|hard"
}`

	recipeRequestPromptTemplate = `Create a recipe with the following requirements:
Style: %s
Serving Size: %d people
Time Limit: %d minutes`
)

func composeRecipeRequest(req RecipeRequest) CompletionRequest {
	return CompletionRequest{
		SystemPrompt: RecipeGenerationSystemPrompt,
		Messages: []Message{
			{
				Role:    "user",
				Content: fmt.Sprintf(recipeRequestPromptTemplate, req.Style, req.ServingSize, req.TimeLimit),
			},
		},
	}
}
