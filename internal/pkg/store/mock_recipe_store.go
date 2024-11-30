package store

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tfiling/ai-chef/internal/pkg/models"
)

type MockRecipeStore struct {
	mock.Mock
}

func (m *MockRecipeStore) Create(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	args := m.Called(ctx, recipe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *MockRecipeStore) GetAll(ctx context.Context) ([]models.Recipe, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Recipe), args.Error(1)
}

func (m *MockRecipeStore) GetByID(ctx context.Context, id string) (*models.Recipe, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *MockRecipeStore) Update(ctx context.Context, recipe *models.Recipe) (*models.Recipe, error) {
	args := m.Called(ctx, recipe)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Recipe), args.Error(1)
}

func (m *MockRecipeStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
