package llm

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type IRecipeParser interface {
	ParseRecipeFromCompletion(completion string) (Recipe, error)
}

type RecipeParser struct {
}

func NewRecipeParser() *RecipeParser {
	return &RecipeParser{}
}

func (p *RecipeParser) ParseRecipeFromCompletion(completion string) (Recipe, error) {
	var recipe Recipe
	err := json.Unmarshal([]byte(completion), &recipe)
	if err != nil {
		return Recipe{}, errors.Wrap(err, "failed to parse recipe from completion")
	}

	err = validator.New().Struct(recipe)
	if err != nil {
		return Recipe{}, errors.Wrap(err, "invalid recipe format")
	}

	return recipe, nil
}
