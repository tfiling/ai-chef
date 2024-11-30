package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tfiling/ai-chef/internal/pkg/models"
	"github.com/tfiling/ai-chef/internal/pkg/store"
)

type RecipeController struct {
	RecipeStore store.IRecipeStore
}

func (c *RecipeController) RegisterRoutes(router fiber.Router) error {
	router.Post("/recipes", c.CreateRecipe)
	router.Get("/recipes", c.GetRecipes)
	router.Get("/recipes/:id", c.GetRecipe)
	router.Put("/recipes/:id", c.UpdateRecipe)
	router.Delete("/recipes/:id", c.DeleteRecipe)
	return nil
}

func (c *RecipeController) CreateRecipe(ctx *fiber.Ctx) error {
	recipe := new(models.Recipe)
	if err := ctx.BodyParser(recipe); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdRecipe, err := c.RecipeStore.Create(ctx.Context(), recipe)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdRecipe)
}

func (c *RecipeController) GetRecipes(ctx *fiber.Ctx) error {
	recipes, err := c.RecipeStore.GetAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(recipes)
}

func (c *RecipeController) GetRecipe(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	recipe, err := c.RecipeStore.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Recipe not found"})
	}

	return ctx.JSON(recipe)
}

func (c *RecipeController) UpdateRecipe(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	recipe := new(models.Recipe)
	if err := ctx.BodyParser(recipe); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	recipe.ID = id
	updatedRecipe, err := c.RecipeStore.Update(ctx.Context(), recipe)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(updatedRecipe)
}

func (c *RecipeController) DeleteRecipe(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.RecipeStore.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
