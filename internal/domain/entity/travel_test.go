package entity

import (
	"challenge-travel-api/internal/domain/enums"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTravelRequest_UpdateTravelRequest(t *testing.T) {
	id := uuid.New()
	userId := uuid.New()
	now := time.Now()
	travelRequest := &TravelRequest{
		Id:              id,
		TravelerName:    "João Silva",
		UserId:          userId,
		DestinationName: "Paris",
		DepartureDate:   now,
		Status:          enums.TravelRequestStatusSolicited,
		CreatedAt:       now,
	}

	t.Run("should update all fields when all parameters are provided", func(t *testing.T) {
		newTravelerName := "Maria Santos"
		newDestinationName := "Londres"
		newDepartureDate := now.AddDate(0, 1, 0)
		newReturnDate := now.AddDate(0, 2, 0)
		newStatus := enums.TravelRequestStatusApproved
		approvedBy := uuid.New()

		// Act
		travelRequest.UpdateTravelRequest(
			&newDestinationName,
			&newTravelerName,
			&newDepartureDate,
			&newReturnDate,
			&newStatus,
			nil,
			&approvedBy,
		)

		// Assert
		assert.Equal(t, newTravelerName, travelRequest.TravelerName)
		assert.Equal(t, newDestinationName, travelRequest.DestinationName)
		assert.Equal(t, newDepartureDate, travelRequest.DepartureDate)
		assert.Equal(t, newReturnDate, *travelRequest.ReturnDate)
		assert.Equal(t, newStatus, travelRequest.Status)
		assert.Equal(t, approvedBy, *travelRequest.ApprovedBy)
		assert.NotNil(t, travelRequest.ApprovedAt)
		assert.NotNil(t, travelRequest.UpdatedAt)
	})

	t.Run("should update only provided fields", func(t *testing.T) {
		// Arrange
		originalTravelerName := travelRequest.TravelerName
		originalDepartureDate := travelRequest.DepartureDate
		newDestinationName := "Roma"

		// Act
		travelRequest.UpdateTravelRequest(
			&newDestinationName,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		// Assert
		assert.Equal(t, originalTravelerName, travelRequest.TravelerName)
		assert.Equal(t, newDestinationName, travelRequest.DestinationName)
		assert.Equal(t, originalDepartureDate, travelRequest.DepartureDate)
		assert.NotNil(t, travelRequest.UpdatedAt)
	})

	t.Run("should handle cancellation correctly", func(t *testing.T) {
		// Arrange
		canceledBy := uuid.New()
		newStatus := enums.TravelRequestStatusCanceled

		// Act
		travelRequest.UpdateTravelRequest(
			nil,
			nil,
			nil,
			nil,
			&newStatus,
			&canceledBy,
			nil,
		)

		// Assert
		assert.Equal(t, newStatus, travelRequest.Status)
		assert.Equal(t, canceledBy, *travelRequest.CanceledBy)
		assert.NotNil(t, travelRequest.CanceledAt)
		assert.NotNil(t, travelRequest.UpdatedAt)
	})
}
