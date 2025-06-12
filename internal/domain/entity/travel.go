package entity

import (
	"challenge-travel-api/internal/domain/enums"
	"github.com/google/uuid"
	"time"
)

type TravelRequest struct {
	Id              uuid.UUID                 `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TravelerName    string                    `json:"traveler_name" gorm:"type:varchar(255);not null"`
	UserId          uuid.UUID                 `json:"user_id" gorm:"type:uuid;not null"`
	DestinationName string                    `json:"destination_name" gorm:"type:varchar(255);not null"`
	DepartureDate   time.Time                 `json:"departure_date" gorm:"type:timestamp;not null"`
	ReturnDate      *time.Time                `json:"return_date" gorm:"type:timestamp;"`
	Status          enums.TravelRequestStatus `json:"status" gorm:"type:travel_request_status;not null"`
	CreatedAt       time.Time                 `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt       *time.Time                `json:"updated_at" gorm:"type:timestamp"`

	User User `json:"user" gorm:"foreignkey:user_id"`
}
