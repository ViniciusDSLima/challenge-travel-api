package usecase

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/enums"
	"challenge-travel-api/internal/interface/dto"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockUserGateway struct {
	mock.Mock
}

func (m *MockUserGateway) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserGateway) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserGateway) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	fmt.Printf("FindByEmail called with email: %s\n", email) // Debug log

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserGateway) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func TestAuthUseCase_Register(t *testing.T) {
	// Setup
	ctx := context.Background()

	t.Run("should register user successfully", func(t *testing.T) {
		//setup
		mockUserGateway := new(MockUserGateway)
		useCase := NewAUthUseCase(mockUserGateway)

		// Arrange
		input := dto.RegisterRequestDTO{
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Password: "password123",
			Role:     enums.UserTypeCommon,
		}

		mockUserGateway.On("FindByEmail", ctx, input.Email).Return(nil, gorm.ErrRecordNotFound)
		mockUserGateway.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

		// Act
		err := useCase.Register(ctx, input)

		// Assert
		assert.NoError(t, err)
		mockUserGateway.AssertExpectations(t)
	})

	t.Run("should return error for existing user", func(t *testing.T) {
		//setup
		mockUserGateway := new(MockUserGateway)
		useCase := NewAUthUseCase(mockUserGateway)

		// Arrange
		input := dto.RegisterRequestDTO{
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Password: "password123",
			Role:     enums.UserTypeCommon,
		}

		existingUser := &entity.User{
			Id:    uuid.New(),
			Name:  "John Doe",
			Email: input.Email,
		}

		mockUserGateway.On("FindByEmail", ctx, input.Email).Return(existingUser, nil)

		// Act
		err := useCase.Register(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrUserAlreadyExists, err)
		mockUserGateway.AssertExpectations(t)
	})
}

func TestAuthUseCase_Login(t *testing.T) {
	// Setup
	mockUserGateway := new(MockUserGateway)
	useCase := NewAUthUseCase(mockUserGateway)
	ctx := context.Background()

	// Set JWT secret key for testing
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET_KEY")

	t.Run("should login successfully", func(t *testing.T) {
		// Arrange
		input := dto.LoginRequestDTO{
			Email:    "john.doe@example.com",
			Password: "password123",
		}

		hashedPassword, _ := useCase.hashPassword(input.Password)
		user := &entity.User{
			Id:       uuid.New(),
			Name:     "John Doe",
			Email:    input.Email,
			Password: hashedPassword,
			Role:     enums.UserTypeCommon,
		}

		mockUserGateway.On("FindByEmail", ctx, input.Email).Return(user, nil)

		// Act
		result, err := useCase.Login(ctx, input)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, result.AccessToken)
		mockUserGateway.AssertExpectations(t)
	})

	t.Run("should return error for invalid credentials", func(t *testing.T) {
		// Arrange
		input := dto.LoginRequestDTO{
			Email:    "john.doe@example.com",
			Password: "wrongpassword",
		}

		hashedPassword, _ := useCase.hashPassword("correctpassword")
		user := &entity.User{
			Id:       uuid.New(),
			Name:     "John Doe",
			Email:    input.Email,
			Password: hashedPassword,
			Role:     enums.UserTypeCommon,
		}

		mockUserGateway.On("FindByEmail", ctx, input.Email).Return(user, nil)

		// Act
		result, err := useCase.Login(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result.AccessToken)
		mockUserGateway.AssertExpectations(t)
	})

	t.Run("should return error for non-existent user", func(t *testing.T) {
		// Arrange
		input := dto.LoginRequestDTO{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		mockUserGateway.On("FindByEmail", ctx, input.Email).Return(nil, gorm.ErrRecordNotFound)

		// Act
		result, err := useCase.Login(ctx, input)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result.AccessToken)
		mockUserGateway.AssertExpectations(t)
	})
}
