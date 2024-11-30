package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tfiling/ai-chef/internal/pkg/models"
	"github.com/tfiling/ai-chef/internal/pkg/store"
)

type UserController struct {
	UserStore store.IUserStore
}

func (c *UserController) RegisterRoutes(router fiber.Router) error {
	router.Post("/users", c.CreateUser)
	router.Get("/users", c.GetUsers)
	router.Get("/users/:id", c.GetUser)
	router.Put("/users/:id", c.UpdateUser)
	router.Delete("/users/:id", c.DeleteUser)
	return nil
}

func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdUser, err := c.UserStore.Create(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdUser)
}

func (c *UserController) GetUsers(ctx *fiber.Ctx) error {
	users, err := c.UserStore.GetAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(users)
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.UserStore.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.JSON(user)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user.ID = id
	updatedUser, err := c.UserStore.Update(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(updatedUser)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.UserStore.Delete(ctx.Context(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
