package utils

import (
	"challenge-travel-api/internal/domain/enums"
	"time"

	"github.com/google/uuid"
)

type TravelRequestFilters struct {
	UserId          *uuid.UUID
	Status          *enums.TravelRequestStatus
	StartDate       *time.Time
	EndDate         *time.Time
	DestinationName *string
	Page            int
	PageSize        int
}
