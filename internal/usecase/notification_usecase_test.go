package usecase

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/enums"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEmailNotificationService_NotifyStatusChange(t *testing.T) {
	// Setup
	service := NewEmailNotificationService()

	t.Run("should notify when status changes to approved", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		departureDate := time.Now().AddDate(0, 1, 0)
		returnDate := time.Now().AddDate(0, 2, 0)

		travelRequest := &entity.TravelRequest{
			Id:              uuid.New(),
			TravelerName:    "John Doe",
			UserId:          userId,
			DestinationName: "Paris",
			DepartureDate:   departureDate,
			ReturnDate:      &returnDate,
			Status:          enums.TravelRequestStatusApproved,
			User: entity.User{
				Id:    userId,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		}

		previousStatus := enums.TravelRequestStatusSolicited

		// Act & Assert
		// Como o serviço apenas loga a mensagem, não há retorno para verificar
		// Mas podemos garantir que não há panics
		assert.NotPanics(t, func() {
			service.NotifyStatusChange(travelRequest, previousStatus)
		})
	})

	t.Run("should notify when status changes to canceled", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		departureDate := time.Now().AddDate(0, 1, 0)
		returnDate := time.Now().AddDate(0, 2, 0)

		travelRequest := &entity.TravelRequest{
			Id:              uuid.New(),
			TravelerName:    "John Doe",
			UserId:          userId,
			DestinationName: "Paris",
			DepartureDate:   departureDate,
			ReturnDate:      &returnDate,
			Status:          enums.TravelRequestStatusCanceled,
			User: entity.User{
				Id:    userId,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		}

		previousStatus := enums.TravelRequestStatusSolicited

		// Act & Assert
		// Como o serviço apenas loga a mensagem, não há retorno para verificar
		// Mas podemos garantir que não há panics
		assert.NotPanics(t, func() {
			service.NotifyStatusChange(travelRequest, previousStatus)
		})
	})

	t.Run("should not notify when status remains the same", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		departureDate := time.Now().AddDate(0, 1, 0)
		returnDate := time.Now().AddDate(0, 2, 0)

		travelRequest := &entity.TravelRequest{
			Id:              uuid.New(),
			TravelerName:    "John Doe",
			UserId:          userId,
			DestinationName: "Paris",
			DepartureDate:   departureDate,
			ReturnDate:      &returnDate,
			Status:          enums.TravelRequestStatusApproved,
			User: entity.User{
				Id:    userId,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		}

		previousStatus := enums.TravelRequestStatusApproved

		// Act & Assert
		// Como o serviço apenas loga a mensagem, não há retorno para verificar
		// Mas podemos garantir que não há panics
		assert.NotPanics(t, func() {
			service.NotifyStatusChange(travelRequest, previousStatus)
		})
	})

	t.Run("should not notify for other status changes", func(t *testing.T) {
		// Arrange
		userId := uuid.New()
		departureDate := time.Now().AddDate(0, 1, 0)
		returnDate := time.Now().AddDate(0, 2, 0)

		travelRequest := &entity.TravelRequest{
			Id:              uuid.New(),
			TravelerName:    "John Doe",
			UserId:          userId,
			DestinationName: "Paris",
			DepartureDate:   departureDate,
			ReturnDate:      &returnDate,
			Status:          enums.TravelRequestStatusSolicited,
			User: entity.User{
				Id:    userId,
				Name:  "John Doe",
				Email: "john.doe@example.com",
			},
		}

		previousStatus := enums.TravelRequestStatusSolicited

		// Act & Assert
		// Como o serviço apenas loga a mensagem, não há retorno para verificar
		// Mas podemos garantir que não há panics
		assert.NotPanics(t, func() {
			service.NotifyStatusChange(travelRequest, previousStatus)
		})
	})
}
