package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/tfiling/ai-chef/internal/pkg/llm"
)

type RecipeGenerationController struct {
	recipeGenerator llm.IRecipeGenerator
}

func NewRecipeGenerationController(recipeGenerator llm.IRecipeGenerator) *RecipeGenerationController {
	return &RecipeGenerationController{
		recipeGenerator: recipeGenerator,
	}
}

func (c *RecipeGenerationController) RegisterRoutes(router fiber.Router) error {
	router.Post("/recipes/generate", c.GenerateRecipe)
	return nil
}

func (c *RecipeGenerationController) GenerateRecipe(ctx *fiber.Ctx) error {
	req := llm.RecipeRequest{}
	if err := ctx.BodyParser(&req); err != nil {
		err = errors.Wrap(err, "could not parse request")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := validator.New().Struct(req); err != nil {
		err = errors.Wrap(err, "received invalid request")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	recipe, err := c.recipeGenerator.GenerateRecipe(ctx.Context(), req)
	if err != nil {
		err = errors.Wrap(err, "recipe generation failed")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(recipe)
}
