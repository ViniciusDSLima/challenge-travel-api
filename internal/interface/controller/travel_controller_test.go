package controller

import (
	"bytes"
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/enums"
	"challenge-travel-api/internal/interface/dto"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTravelUseCase struct {
	mock.Mock
}

func (m *MockTravelUseCase) CreateTravelRequest(ctx context.Context, userID uuid.UUID, input dto.CreateTravelRequestDTO) (*entity.TravelRequest, error) {
	args := m.Called(ctx, userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TravelRequest), args.Error(1)
}

func (m *MockTravelUseCase) UpdateTravelRequest(ctx context.Context, id uuid.UUID, userID uuid.UUID, input dto.UpdateTravelRequestDTO) (*entity.TravelRequest, error) {
	args := m.Called(ctx, id, userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TravelRequest), args.Error(1)
}

func (m *MockTravelUseCase) UpdateStatusTravelRequest(ctx context.Context, userId string, input dto.UpdateStatusTravelRequestDTO) error {
	args := m.Called(ctx, userId, input)
	return args.Error(0)
}

func (m *MockTravelUseCase) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.TravelRequest, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TravelRequest), args.Error(1)
}

func (m *MockTravelUseCase) ListTravelRequests(ctx context.Context, userID uuid.UUID, status *enums.TravelRequestStatus, startDate *time.Time, endDate *time.Time, destinationName *string, page int, pageSize int) ([]entity.TravelRequest, error) {
	args := m.Called(ctx, userID, status, startDate, endDate, destinationName, page, pageSize)
	return args.Get(0).([]entity.TravelRequest), args.Error(1)
}

func TestTravelController_CreateTravelRequest(t *testing.T) {
	// Setup
	mockUseCase := new(MockTravelUseCase)
	controller := NewTravelController(mockUseCase)
	router := setupTestRouter()

	userID := uuid.New()
	router.POST("/travels", func(c *gin.Context) {
		c.Set("user_id", userID)
		controller.CreateTravelRequest(c)
	})

	t.Run("should create travel request successfully", func(t *testing.T) {
		// Arrange
		departureDate := time.Now().AddDate(0, 1, 0)
		returnDate := time.Now().AddDate(0, 2, 0)
		request := dto.CreateTravelRequestDTO{
			TravelerName:    "John Doe",
			DestinationName: "Paris",
			DepartureDate:   departureDate,
			ReturnDate:      &returnDate,
		}

		expectedTravel := &entity.TravelRequest{
			Id:              uuid.New(),
			TravelerName:    request.TravelerName,
			UserId:          userID,
			DestinationName: request.DestinationName,
			DepartureDate:   request.DepartureDate,
			ReturnDate:      request.ReturnDate,
			Status:          enums.TravelRequestStatusSolicited,
		}

		mockUseCase.On("CreateTravelRequest", mock.Anything, userID, request).Return(expectedTravel, nil)

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/travels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)
		var response entity.TravelRequest
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTravel.Id, response.Id)
		assert.Equal(t, expectedTravel.TravelerName, response.TravelerName)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return error for invalid request", func(t *testing.T) {
		// Arrange
		invalidRequest := map[string]interface{}{
			"invalid_field": "value",
		}

		// Act
		body, _ := json.Marshal(invalidRequest)
		req := httptest.NewRequest(http.MethodPost, "/travels", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUseCase.AssertNotCalled(t, "CreateTravelRequest")
	})
}

func TestTravelController_UpdateTravelRequest(t *testing.T) {
	// Setup
	mockUseCase := new(MockTravelUseCase)
	controller := NewTravelController(mockUseCase)
	router := setupTestRouter()

	userID := uuid.New()
	travelID := uuid.New()
	router.PUT("/travels/:id", func(c *gin.Context) {
		c.Set("user_id", userID)
		controller.UpdateTravelRequest(c)
	})

	t.Run("should update travel request successfully", func(t *testing.T) {
		// Arrange
		departureDate := time.Now().AddDate(0, 1, 0)
		returnDate := time.Now().AddDate(0, 2, 0)
		request := dto.UpdateTravelRequestDTO{
			TravelerName:    stringPtr("John Doe Updated"),
			DestinationName: stringPtr("London"),
			DepartureDate:   &departureDate,
			ReturnDate:      &returnDate,
		}

		expectedTravel := &entity.TravelRequest{
			Id:              travelID,
			TravelerName:    *request.TravelerName,
			UserId:          userID,
			DestinationName: *request.DestinationName,
			DepartureDate:   *request.DepartureDate,
			ReturnDate:      request.ReturnDate,
			Status:          enums.TravelRequestStatusSolicited,
		}

		mockUseCase.On("UpdateTravelRequest", mock.Anything, travelID, userID, request).Return(expectedTravel, nil)

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPut, "/travels/"+travelID.String(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response entity.TravelRequest
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedTravel.Id, response.Id)
		assert.Equal(t, expectedTravel.TravelerName, response.TravelerName)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("should return error for invalid ID", func(t *testing.T) {
		// Arrange
		request := dto.UpdateTravelRequestDTO{
			TravelerName:    stringPtr("John Doe"),
			DestinationName: stringPtr("Paris"),
		}

		// Act
		body, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPut, "/travels/invalid-id", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockUseCase.AssertNotCalled(t, "UpdateTravelRequest")
	})
}

func stringPtr(s string) *string {
	return &s
}
