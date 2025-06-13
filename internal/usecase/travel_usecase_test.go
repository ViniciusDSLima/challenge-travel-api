package usecase

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/enums"
	"challenge-travel-api/internal/interface/dto"
	"challenge-travel-api/internal/utils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTravelGateway struct {
	mock.Mock
}

func (m *MockTravelGateway) Create(ctx context.Context, travelRequest *entity.TravelRequest) error {
	args := m.Called(ctx, travelRequest)
	return args.Error(0)
}

func (m *MockTravelGateway) FindByID(ctx context.Context, id uuid.UUID) (*entity.TravelRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.TravelRequest), args.Error(1)
}

func (m *MockTravelGateway) Update(ctx context.Context, travelRequest *entity.TravelRequest) error {
	args := m.Called(ctx, travelRequest)
	return args.Error(0)
}

func (m *MockTravelGateway) List(ctx context.Context, filters utils.TravelRequestFilters) ([]entity.TravelRequest, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).([]entity.TravelRequest), args.Error(1)
}

func (m *MockTravelGateway) ListByUserID(ctx context.Context, userID uuid.UUID, filters utils.TravelRequestFilters) ([]entity.TravelRequest, error) {
	args := m.Called(ctx, userID, filters)
	return args.Get(0).([]entity.TravelRequest), args.Error(1)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) NotifyStatusChange(travel *entity.TravelRequest, previousStatus enums.TravelRequestStatus) {
	m.Called(travel, previousStatus)
}

func TestTravelRequestUseCase_CreateTravelRequest(t *testing.T) {
	// Setup
	mockTravelGateway := new(MockTravelGateway)
	mockUserGateway := new(MockUserGateway)
	mockNotificationService := new(MockNotificationService)
	useCase := NewTravelRequestUseCase(mockTravelGateway, mockUserGateway, mockNotificationService)

	ctx := context.Background()
	userID := uuid.New()
	now := time.Now()
	futureDate := now.AddDate(0, 1, 0)
	returnDate := now.AddDate(0, 2, 0)

	user := &entity.User{
		Id:   userID,
		Name: "Test User",
		Role: enums.UserTypeCommon,
	}

	t.Run("should create travel request successfully", func(t *testing.T) {
		// Arrange
		input := dto.CreateTravelRequestDTO{
			TravelerName:    "John Doe",
			DestinationName: "Paris",
			DepartureDate:   futureDate,
			ReturnDate:      &returnDate,
		}

		mockUserGateway.On("FindByID", ctx, userID).Return(user, nil)
		mockTravelGateway.On("Create", ctx, mock.AnythingOfType("*entity.TravelRequest")).Return(nil)

		// Act
		result, err := useCase.CreateTravelRequest(ctx, userID, input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, input.TravelerName, result.TravelerName)
		assert.Equal(t, input.DestinationName, result.DestinationName)
		assert.Equal(t, input.DepartureDate, result.DepartureDate)
		assert.Equal(t, input.ReturnDate, result.ReturnDate)
		assert.Equal(t, enums.TravelRequestStatusSolicited, result.Status)
		mockUserGateway.AssertExpectations(t)
		mockTravelGateway.AssertExpectations(t)
	})

	t.Run("should return error for invalid destination", func(t *testing.T) {
		// Arrange
		input := dto.CreateTravelRequestDTO{
			TravelerName:    "John Doe",
			DestinationName: "",
			DepartureDate:   futureDate,
			ReturnDate:      &returnDate,
		}

		// Act
		result, err := useCase.CreateTravelRequest(ctx, userID, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDestination, err)
		assert.Nil(t, result)
	})

	t.Run("should return error for past departure date", func(t *testing.T) {
		// Arrange
		pastDate := now.AddDate(0, -1, 0)
		input := dto.CreateTravelRequestDTO{
			TravelerName:    "John Doe",
			DestinationName: "Paris",
			DepartureDate:   pastDate,
			ReturnDate:      &returnDate,
		}

		// Act
		result, err := useCase.CreateTravelRequest(ctx, userID, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrFutureDatesOnly, err)
		assert.Nil(t, result)
	})

	t.Run("should return error for invalid dates", func(t *testing.T) {
		// Arrange
		invalidReturnDate := now.AddDate(0, -1, 0)
		input := dto.CreateTravelRequestDTO{
			TravelerName:    "John Doe",
			DestinationName: "Paris",
			DepartureDate:   futureDate,
			ReturnDate:      &invalidReturnDate,
		}

		// Act
		result, err := useCase.CreateTravelRequest(ctx, userID, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidDates, err)
		assert.Nil(t, result)
	})
}

func TestTravelRequestUseCase_UpdateStatusTravelRequest(t *testing.T) {
	// Setup
	mockTravelGateway := new(MockTravelGateway)
	mockUserGateway := new(MockUserGateway)
	mockNotificationService := new(MockNotificationService)
	useCase := NewTravelRequestUseCase(mockTravelGateway, mockUserGateway, mockNotificationService)

	ctx := context.Background()
	userID := uuid.New()
	adminID := uuid.New()
	travelID := uuid.New()

	admin := &entity.User{
		Id:   adminID,
		Name: "Admin User",
		Role: enums.UserTypeAdmin,
	}

	regularUser := &entity.User{
		Id:   userID,
		Name: "Regular User",
		Role: enums.UserTypeCommon,
	}

	travel := &entity.TravelRequest{
		Id:     travelID,
		UserId: userID,
		Status: enums.TravelRequestStatusSolicited,
	}

	t.Run("should update status successfully", func(t *testing.T) {
		// Arrange
		input := dto.UpdateStatusTravelRequestDTO{
			TravelRequestId: travelID.String(),
			Status:          enums.TravelRequestStatusApproved,
		}

		mockUserGateway.On("FindByID", ctx, adminID).Return(admin, nil)
		mockTravelGateway.On("FindByID", ctx, travelID).Return(travel, nil)
		mockTravelGateway.On("Update", ctx, mock.AnythingOfType("*entity.TravelRequest")).Return(nil)
		mockNotificationService.On("NotifyStatusChange", mock.AnythingOfType("*entity.TravelRequest"), enums.TravelRequestStatusSolicited).Return()

		// Act
		err := useCase.UpdateStatusTravelRequest(ctx, adminID.String(), input)

		// Assert
		assert.NoError(t, err)
		mockUserGateway.AssertExpectations(t)
		mockTravelGateway.AssertExpectations(t)
		mockNotificationService.AssertExpectations(t)
	})

	t.Run("should return error for unauthorized user", func(t *testing.T) {
		// Arrange
		input := dto.UpdateStatusTravelRequestDTO{
			TravelRequestId: travelID.String(),
			Status:          enums.TravelRequestStatusApproved,
		}

		mockUserGateway.On("FindByID", ctx, userID).Return(regularUser, nil)
		mockTravelGateway.On("FindByID", ctx, travelID).Return(travel, nil)

		// Act
		err := useCase.UpdateStatusTravelRequest(ctx, userID.String(), input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrUnauthorized, err)
		mockUserGateway.AssertExpectations(t)
		mockTravelGateway.AssertExpectations(t)
	})

	t.Run("should return error for already approved travel", func(t *testing.T) {
		// Arrange
		approvedTravel := &entity.TravelRequest{
			Id:     travelID,
			UserId: userID,
			Status: enums.TravelRequestStatusApproved,
		}

		input := dto.UpdateStatusTravelRequestDTO{
			TravelRequestId: travelID.String(),
			Status:          enums.TravelRequestStatusApproved,
		}

		mockUserGateway.On("FindByID", ctx, adminID).Return(admin, nil)
		mockTravelGateway.On("FindByID", ctx, travelID).Return(approvedTravel, nil)

		// Act
		err := useCase.UpdateStatusTravelRequest(ctx, adminID.String(), input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrTravelAlreadyApproved, err)
		mockUserGateway.AssertExpectations(t)
		mockTravelGateway.AssertExpectations(t)
	})
}
