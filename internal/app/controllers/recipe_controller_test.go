package controllers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tfiling/ai-chef/internal/app/controllers"
	"github.com/tfiling/ai-chef/internal/pkg/models"
	"github.com/tfiling/ai-chef/internal/pkg/store"
)

func setupTestApp(t *testing.T, mockStore *store.MockRecipeStore) *fiber.App {
	app := fiber.New()
	controller := &controllers.RecipeController{RecipeStore: mockStore}
	err := controller.RegisterRoutes(app)
	assert.NoError(t, err)
	return app
}

func TestRecipeController_CreateRecipe_SuccessfullyCreateRecipe(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	recipe := &models.Recipe{Name: "Test Recipe"}
	expectedRecipe := &models.Recipe{ID: "123", Name: "Test Recipe"}

	mockStore.On("Create", mock.Anything, mock.AnythingOfType("*models.Recipe")).Return(expectedRecipe, nil)

	reqBody, _ := json.Marshal(recipe)
	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	var resultRecipe models.Recipe
	err = json.Unmarshal(body, &resultRecipe)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedRecipe.ID, resultRecipe.ID)
	assert.Equal(t, expectedRecipe.Name, resultRecipe.Name)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_CreateRecipe_InvalidJSON(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	invalidJSON := []byte(`{"name": Invalid JSON}`)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(invalidJSON)))

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRecipeController_CreateRecipe_StoreCreationFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	recipe := &models.Recipe{Name: "Test Recipe"}
	mockStore.On("Create", mock.Anything, mock.AnythingOfType("*models.Recipe")).Return(nil, fiber.NewError(fiber.StatusInternalServerError))

	reqBody, _ := json.Marshal(recipe)
	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_GetRecipes_SuccessfullyRetrieveAllRecipes(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	expectedRecipes := []models.Recipe{{ID: "1", Name: "Recipe 1"}, {ID: "2", Name: "Recipe 2"}}
	mockStore.On("GetAll", mock.Anything).Return(expectedRecipes, nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recipes", nil))
	body, _ := io.ReadAll(resp.Body)
	var resultRecipes []models.Recipe
	err := json.Unmarshal(body, &resultRecipes)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedRecipes, resultRecipes)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_GetRecipes_EmptyRecipeList(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	expectedRecipes := []models.Recipe{}
	mockStore.On("GetAll", mock.Anything).Return(expectedRecipes, nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recipes", nil))
	body, _ := io.ReadAll(resp.Body)
	var resultRecipes []models.Recipe
	err := json.Unmarshal(body, &resultRecipes)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Empty(t, resultRecipes)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_GetRecipes_StoreRetrievalFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	mockStore.On("GetAll", mock.Anything).Return(nil, fiber.NewError(fiber.StatusInternalServerError))

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recipes", nil))

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_GetRecipe_SuccessfullyRetrieveRecipe(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	expectedRecipe := &models.Recipe{ID: "123", Name: "Test Recipe"}
	mockStore.On("GetByID", mock.Anything, "123").Return(expectedRecipe, nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recipes/123", nil))
	body, _ := io.ReadAll(resp.Body)
	var resultRecipe models.Recipe
	err := json.Unmarshal(body, &resultRecipe)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedRecipe.ID, resultRecipe.ID)
	assert.Equal(t, expectedRecipe.Name, resultRecipe.Name)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_GetRecipe_RecipeNotFound(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	mockStore.On("GetByID", mock.Anything, "123").Return(nil, fiber.NewError(fiber.StatusNotFound))

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/recipes/123", nil))

	// Assert
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_UpdateRecipe_SuccessfullyUpdateRecipe(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	recipe := &models.Recipe{Name: "Updated Recipe"}
	expectedRecipe := &models.Recipe{ID: "123", Name: "Updated Recipe"}

	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*models.Recipe")).Return(expectedRecipe, nil)

	reqBody, _ := json.Marshal(recipe)
	req := httptest.NewRequest(http.MethodPut, "/recipes/123", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	var resultRecipe models.Recipe
	err = json.Unmarshal(body, &resultRecipe)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedRecipe.ID, resultRecipe.ID)
	assert.Equal(t, expectedRecipe.Name, resultRecipe.Name)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_UpdateRecipe_InvalidJSON(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	invalidJSON := []byte(`{"name": Invalid JSON}`)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodPut, "/recipes/123", bytes.NewBuffer(invalidJSON)))

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestRecipeController_UpdateRecipe_StoreUpdateFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	recipe := &models.Recipe{Name: "Updated Recipe"}
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*models.Recipe")).Return(nil, fiber.NewError(fiber.StatusInternalServerError))

	reqBody, _ := json.Marshal(recipe)
	req := httptest.NewRequest(http.MethodPut, "/recipes/123", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_DeleteRecipe_SuccessfullyDeleteRecipe(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	mockStore.On("Delete", mock.Anything, "123").Return(nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/recipes/123", nil))

	// Assert
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestRecipeController_DeleteRecipe_StoreDeletionFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockRecipeStore)
	app := setupTestApp(t, mockStore)

	mockStore.On("Delete", mock.Anything, "123").Return(fiber.NewError(fiber.StatusInternalServerError))

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/recipes/123", nil))

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}
