package llm

type IRecipeParser interface {
	ParseRecipeFromCompletion(completion string) (Recipe, error)
}

type RecipeParser struct {
}
