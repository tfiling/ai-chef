package store

import (
	"context"

	"github.com/tfiling/ai-chef/internal/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRecipeStore interface {
	Create(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error)
	GetAll(ctx context.Context) ([]models.Recipe, error)
	GetByID(ctx context.Context, id string) (*models.Recipe, error)
	Update(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error)
	Delete(ctx context.Context, id string) error
}

type MongoRecipeStore struct {
	collection *mongo.Collection
}

func NewMongoRecipeStore(db *mongo.Database) *MongoRecipeStore {
	return &MongoRecipeStore{
		collection: db.Collection("recipes"),
	}
}

func (s *MongoRecipeStore) Create(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	result, err := s.collection.InsertOne(ctx, recipe)
	if err != nil {
		return nil, err
	}

	recipe.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return recipe, nil
}

func (s *MongoRecipeStore) GetAll(ctx context.Context) ([]models.Recipe, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var recipes []models.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, err
	}

	return recipes, nil
}

func (s *MongoRecipeStore) GetByID(ctx context.Context, id string) (*models.Recipe, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var recipe models.Recipe
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&recipe)
	if err != nil {
		return nil, err
	}

	return &recipe, nil
}

func (s *MongoRecipeStore) Update(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	objID, err := primitive.ObjectIDFromHex(recipe.ID)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"name": recipe.Name,
		},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

func (s *MongoRecipeStore) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
