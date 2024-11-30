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

func setupTestUserController(mockStore *store.MockUserStore) *fiber.App {
	app := fiber.New()
	controller := &controllers.UserController{UserStore: mockStore}
	controller.RegisterRoutes(app)
	return app
}

func TestUserController_CreateUser_SuccessfullyCreateUser(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	user := &models.User{Username: "TestUser"}
	expectedUser := &models.User{ID: "123", Username: "TestUser"}
	mockStore.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	var resultUser models.User
	json.Unmarshal(body, &resultUser)

	// Assert
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedUser.ID, resultUser.ID)
	assert.Equal(t, expectedUser.Username, resultUser.Username)
	mockStore.AssertExpectations(t)
}

func TestUserController_CreateUser_InvalidJSON(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	invalidJSON := []byte(`{"username": Invalid JSON}`)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(invalidJSON)))

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserController_CreateUser_StoreCreationFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	user := &models.User{Username: "TestUser"}
	mockStore.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil, fiber.NewError(fiber.StatusInternalServerError))

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestUserController_GetUsers_SuccessfullyRetrieveAllUsers(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	expectedUsers := []models.User{{ID: "1", Username: "User1"}, {ID: "2", Username: "User2"}}
	mockStore.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/users", nil))
	body, _ := io.ReadAll(resp.Body)
	var resultUsers []models.User
	json.Unmarshal(body, &resultUsers)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedUsers, resultUsers)
	mockStore.AssertExpectations(t)
}

func TestUserController_GetUsers_EmptyList(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	expectedUsers := []models.User{}
	mockStore.On("GetAll", mock.Anything).Return(expectedUsers, nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/users", nil))
	body, _ := io.ReadAll(resp.Body)
	var resultUsers []models.User
	json.Unmarshal(body, &resultUsers)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Empty(t, resultUsers)
	mockStore.AssertExpectations(t)
}

func TestUserController_GetUsers_StoreRetrievalFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	mockStore.On("GetAll", mock.Anything).Return(nil, fiber.NewError(fiber.StatusInternalServerError))

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/users", nil))

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestUserController_GetUser_SuccessfullyRetrieveUser(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	expectedUser := &models.User{ID: "123", Username: "TestUser"}
	mockStore.On("GetByID", mock.Anything, "123").Return(expectedUser, nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/users/123", nil))
	body, _ := io.ReadAll(resp.Body)
	var resultUser models.User
	json.Unmarshal(body, &resultUser)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedUser.ID, resultUser.ID)
	assert.Equal(t, expectedUser.Username, resultUser.Username)
	mockStore.AssertExpectations(t)
}

func TestUserController_GetUser_UserNotFound(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	mockStore.On("GetByID", mock.Anything, "123").Return(nil, fiber.NewError(fiber.StatusNotFound))

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/users/123", nil))

	// Assert
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestUserController_UpdateUser_SuccessfullyUpdateUser(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	user := &models.User{Username: "UpdatedUser"}
	expectedUser := &models.User{ID: "123", Username: "UpdatedUser"}
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(expectedUser, nil)

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPut, "/users/123", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	var resultUser models.User
	json.Unmarshal(body, &resultUser)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedUser.ID, resultUser.ID)
	assert.Equal(t, expectedUser.Username, resultUser.Username)
	mockStore.AssertExpectations(t)
}

func TestUserController_UpdateUser_InvalidJSON(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	invalidJSON := []byte(`{"username": Invalid JSON}`)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodPut, "/users/123", bytes.NewBuffer(invalidJSON)))

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserController_UpdateUser_StoreUpdateFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	user := &models.User{Username: "UpdatedUser"}
	mockStore.On("Update", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil, fiber.NewError(fiber.StatusInternalServerError))

	reqBody, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPut, "/users/123", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestUserController_DeleteUser_SuccessfullyDeleteUser(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	mockStore.On("Delete", mock.Anything, "123").Return(nil)

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/users/123", nil))

	// Assert
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	mockStore.AssertExpectations(t)
}

func TestUserController_DeleteUser_StoreDeletionFailure(t *testing.T) {
	// Arrange
	mockStore := new(store.MockUserStore)
	app := setupTestUserController(mockStore)

	mockStore.On("Delete", mock.Anything, "123").Return(fiber.NewError(fiber.StatusInternalServerError))

	// Act
	resp, _ := app.Test(httptest.NewRequest(http.MethodDelete, "/users/123", nil))

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockStore.AssertExpectations(t)
}
