package entity

import (
	"challenge-travel-api/internal/domain/enums"
	"time"

	"github.com/google/uuid"
)

type TravelRequest struct {
	Id              uuid.UUID                 `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TravelerName    string                    `json:"traveler_name" gorm:"type:varchar(255);not null"`
	UserId          uuid.UUID                 `json:"user_id" gorm:"type:uuid;not null"`
	DestinationName string                    `json:"destination_name" gorm:"type:varchar(255);not null"`
	DepartureDate   time.Time                 `json:"departure_date" gorm:"type:timestamp;not null"`
	ReturnDate      *time.Time                `json:"return_date" gorm:"type:timestamp;"`
	Status          enums.TravelRequestStatus `json:"status" gorm:"type:travel_request_status;not null"`
	CanceledBy      *uuid.UUID                `json:"canceled_by" gorm:"type:uuid"`
	ApprovedBy      *uuid.UUID                `json:"approved_by" gorm:"type:uuid"`
	CreatedAt       time.Time                 `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time                `json:"updated_at" gorm:"type:timestamp"`
	CanceledAt      *time.Time                `json:"canceled_at" gorm:"type:timestamp"`
	ApprovedAt      *time.Time                `json:"approved_at" gorm:"type:timestamp"`

	User User `json:"user" gorm:"foreignkey:user_id"`
}

func (e *TravelRequest) UpdateTravelRequest(
	destinationName,
	travelerName *string,
	departureDate,
	returnDate *time.Time,
	status *enums.TravelRequestStatus,
	canceledBy,
	approvedBy *uuid.UUID) {
	if destinationName != nil {
		e.DestinationName = *destinationName
	}

	if travelerName != nil {
		e.TravelerName = *travelerName
	}

	if departureDate != nil {
		e.DepartureDate = *departureDate
	}

	if returnDate != nil {
		e.ReturnDate = returnDate
	}

	if status != nil {
		e.Status = *status
	}

	if canceledBy != nil {
		e.CanceledBy = canceledBy

		canceledAt := time.Now()
		e.CanceledAt = &canceledAt
	}

	if approvedBy != nil {
		e.ApprovedBy = approvedBy

		approvedAt := time.Now()
		e.ApprovedAt = &approvedAt
	}

	updatedDate := time.Now()
	e.UpdatedAt = &updatedDate
}
