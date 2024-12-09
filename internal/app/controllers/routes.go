package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/tfiling/ai-chef/internal/pkg/llm"
	"github.com/tfiling/ai-chef/internal/pkg/store"
)

const (
	APIRouteBasePath = "/api/v1"
)

type ControllersDependencies struct {
	RecipeStore     store.IRecipeStore
	UserStore       store.IUserStore
	RecipeGenerator llm.IRecipeGenerator
}

type Controller interface {
	RegisterRoutes(router fiber.Router) error
}

func SetupRoutes(v1Router fiber.Router, controllers []Controller) error {
	for _, controller := range controllers {
		if err := controller.RegisterRoutes(v1Router); err != nil {
			return errors.Wrap(err, "failed to register routes")
		}
	}
	return nil
}

func InitControllers(dependencies ControllersDependencies) ([]Controller, error) {
	return []Controller{
		&RecipeController{RecipeStore: dependencies.RecipeStore},
		&UserController{UserStore: dependencies.UserStore},
		NewRecipeGenerationController(dependencies.RecipeGenerator),
	}, nil
}
