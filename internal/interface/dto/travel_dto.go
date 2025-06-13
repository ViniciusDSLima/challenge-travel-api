package dto

import (
	"challenge-travel-api/internal/domain/enums"
	"time"
)

type CreateTravelRequestDTO struct {
	TravelerName    string     `json:"traveler_name" binding:"required"`
	DestinationName string     `json:"destination_name" binding:"required"`
	DepartureDate   time.Time  `json:"departure_date" binding:"required"`
	ReturnDate      *time.Time `json:"return_date,omitempty"`
}

type UpdateTravelRequestDTO struct {
	TravelerName    *string    `json:"traveler_name,omitempty"`
	DestinationName *string    `json:"destination_name,omitempty" `
	DepartureDate   *time.Time `json:"departure_date,omitempty"`
	ReturnDate      *time.Time `json:"return_date,omitempty"`
}

type UpdateStatusTravelRequestDTO struct {
	TravelRequestId string                    `json:"travel_request_id" binding:"required"`
	Status          enums.TravelRequestStatus `json:"status" binding;"required"`
}
