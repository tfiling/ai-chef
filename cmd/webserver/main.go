package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/tfiling/ai-chef/internal/app/controllers"
	"github.com/tfiling/ai-chef/internal/pkg/configs"
	"github.com/tfiling/ai-chef/internal/pkg/llm"
	"github.com/tfiling/ai-chef/internal/pkg/logging"
	"github.com/tfiling/ai-chef/internal/pkg/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initControllersDependencies() controllers.ControllersDependencies {
	db, err := newMongoDB()
	if err != nil {
		panic(errors.Wrap(err, "could not init mongo db client"))
	}
	recipeParser := llm.NewRecipeParser()
	claudeClient := llm.NewClaudeClient(configs.GetConfig().ClaudeAPIKey)
	recipeStore := store.NewMongoRecipeStore(db)
	userStore := store.NewMongoUserStore(db)
	recipeGenerator := llm.NewRecipeGenerator(claudeClient, recipeParser)
	return controllers.ControllersDependencies{
		RecipeStore:     recipeStore,
		UserStore:       userStore,
		RecipeGenerator: recipeGenerator,
	}
}

func newMongoDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://foo:bar@localhost:27017"))
	if err != nil {
		return nil, err
	}
	return client.Database("aichef"), nil
}

func main() {
	logging.Info("started WS", nil)
	app := fiber.New()
	apiGroup := app.Group(controllers.APIRouteBasePath)
	APIControllers, err := controllers.InitControllers(initControllersDependencies())
	if err != nil {
		logging.Panic(err, "error setting up controllers", nil)
	}
	if err := controllers.SetupRoutes(apiGroup, APIControllers); err != nil {
		logging.Panic(err, "error setting up routes", nil)
	}

	if err := app.Listen(":3000"); err != nil {
		logging.Error(err, "Error starting server", nil)
	}
}
