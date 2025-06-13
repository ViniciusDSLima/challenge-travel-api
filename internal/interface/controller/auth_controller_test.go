package controller

import (
	"bytes"
	"challenge-travel-api/internal/interface/dto"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Register(ctx context.Context, input dto.RegisterRequestDTO) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockAuthUseCase) Login(ctx context.Context, input dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(dto.LoginResponseDTO), args.Error(1)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthController_Register(t *testing.T) {
	// Setup
	mockUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockUseCase)
	router := setupTestRouter()
	router.POST("/auth/register", controller.Register)

	t.Run("should register user successfully", func(t *testing.T) {
		// Arrange
		request := dto.RegisterRequestDTO{
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Password: "password123",
			Role:     "USER",
		}

		mockUseCase.On("Register", mock.Anything, request).Return(nil)

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return error for invalid request", func(t *testing.T) {
		// Arrange
		invalidRequest := map[string]interface{}{
			"invalid_field": "value",
		}

		// Act
		body, _ := json.Marshal(invalidRequest)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUseCase.AssertNotCalled(t, "Register")
	})

	t.Run("should return error for existing user", func(t *testing.T) {
		// Arrange
		request := dto.RegisterRequestDTO{
			Name:     "John Doe",
			Email:    "existing@example.com",
			Password: "password123",
			Role:     "USER",
		}

		mockUseCase.On("Register", mock.Anything, request).Return(errors.New("Usuário já existe no sistema"))

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Usuário já existe no sistema", response["error"])
		mockUseCase.AssertExpectations(t)
	})
}

func TestAuthController_Login(t *testing.T) {
	// Setup
	mockUseCase := new(MockAuthUseCase)
	controller := NewAuthController(mockUseCase)
	router := setupTestRouter()
	router.POST("/auth/login", controller.Login)

	t.Run("should login successfully", func(t *testing.T) {
		// Arrange
		request := dto.LoginRequestDTO{
			Email:    "john.doe@example.com",
			Password: "password123",
		}

		expectedResponse := dto.LoginResponseDTO{
			AccessToken: "valid-jwt-token",
		}

		mockUseCase.On("Login", mock.Anything, request).Return(expectedResponse, nil)

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.LoginResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.AccessToken, response.AccessToken)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return error for invalid credentials", func(t *testing.T) {
		// Arrange
		request := dto.LoginRequestDTO{
			Email:    "john.doe@example.com",
			Password: "wrongpassword",
		}

		mockUseCase.On("Login", mock.Anything, request).Return(dto.LoginResponseDTO{}, errors.New("Credenciais inválidas"))

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Credenciais inválidas", response["error"])
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return error for invalid request", func(t *testing.T) {
		// Arrange
		invalidRequest := map[string]interface{}{
			"invalid_field": "value",
		}

		// Act
		body, _ := json.Marshal(invalidRequest)
		req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUseCase.AssertNotCalled(t, "Login")
	})
}
