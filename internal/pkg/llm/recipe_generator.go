package llm

import (
	"time"
)

type RecipeRequest struct {
	Style       string        `validate:"required"`
	ServingSize int           `validate:"required,min=1"`
	TimeLimit   time.Duration `validate:"required"`
}

type Recipe struct {
	ID           string       `json:"id" validate:"required"`
	Name         string       `json:"name" validate:"required"`
	Description  string       `json:"description" validate:"required"`
	Ingredients  []Ingredient `json:"ingredients" validate:"required,min=1,dive"`
	Instructions string       `json:"instructions" validate:"required"`
	PrepTime     int          `json:"prepTime" validate:"required,gt=0"`
	CookTime     int          `json:"cookTime" validate:"required,gt=0"`
	TotalTime    int          `json:"totalTime" validate:"required,gt=0"`
	ServingSize  int          `json:"servingSize" validate:"required,gt=0"`
	Difficulty   string       `json:"difficulty" validate:"required,oneof=easy medium hard"`
	CreatedAt    time.Time    `json:"createdAt" validate:"required"`
}

type Ingredient struct {
	Name   string  `json:"name" validate:"required"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Unit   string  `json:"unit" validate:"required"`
}

type RecipeGenerator struct {
	claudeClient IClaudeClient
	parser       IRecipeParser
}
