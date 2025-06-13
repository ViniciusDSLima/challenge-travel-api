package entity

import (
	"challenge-travel-api/internal/domain/enums"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validation(t *testing.T) {
	t.Run("should create valid user", func(t *testing.T) {
		// Arrange
		id := uuid.New()
		now := time.Now()
		user := &User{
			Id:        id,
			Name:      "João Silva",
			Email:     "joao.silva@example.com",
			Password:  "senha123",
			CreatedAt: now,
			IsActive:  true,
			Role:      enums.UserTypeAdmin,
		}

		// Assert
		assert.NotNil(t, user)
		assert.Equal(t, id, user.Id)
		assert.Equal(t, "João Silva", user.Name)
		assert.Equal(t, "joao.silva@example.com", user.Email)
		assert.Equal(t, "senha123", user.Password)
		assert.Equal(t, now, user.CreatedAt)
		assert.True(t, user.IsActive)
		assert.Equal(t, enums.UserTypeAdmin, user.Role)
		assert.Nil(t, user.UpdatedAt)
		assert.Nil(t, user.DeletedAt)
	})

	t.Run("should handle user deactivation", func(t *testing.T) {
		// Arrange
		user := &User{
			Id:        uuid.New(),
			Name:      "Maria Santos",
			Email:     "maria.santos@example.com",
			Password:  "senha456",
			CreatedAt: time.Now(),
			IsActive:  true,
			Role:      enums.UserTypeCommon,
		}

		// Act
		user.IsActive = false
		now := time.Now()
		user.DeletedAt = &now

		// Assert
		assert.False(t, user.IsActive)
		assert.NotNil(t, user.DeletedAt)
		assert.Equal(t, now, *user.DeletedAt)
	})

	t.Run("should handle user update", func(t *testing.T) {
		// Arrange
		user := &User{
			Id:        uuid.New(),
			Name:      "Pedro Oliveira",
			Email:     "pedro.oliveira@example.com",
			Password:  "senha789",
			CreatedAt: time.Now(),
			IsActive:  true,
			Role:      enums.UserTypeCommon,
		}

		// Act
		user.Name = "Pedro Silva Oliveira"
		now := time.Now()
		user.UpdatedAt = &now

		// Assert
		assert.Equal(t, "Pedro Silva Oliveira", user.Name)
		assert.NotNil(t, user.UpdatedAt)
		assert.Equal(t, now, *user.UpdatedAt)
	})
}
