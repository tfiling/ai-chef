package models

type RecipeStyle string
type RecipeDifficulty string

const (
	StyleItalian       RecipeStyle = "italian"
	StyleFrench        RecipeStyle = "french"
	StyleAsian         RecipeStyle = "asian"
	StyleMediterranean RecipeStyle = "mediterranean"

	DifficultyEasy   RecipeDifficulty = "easy"
	DifficultyMedium RecipeDifficulty = "medium"
	DifficultyHard   RecipeDifficulty = "hard"
)

type Recipe struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" validate:"required"`
}
