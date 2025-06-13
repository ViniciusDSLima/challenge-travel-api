package gateway

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/utils"
	"context"

	"github.com/google/uuid"
)

type TravelRequestGateway interface {
	Create(ctx context.Context, travelRequest *entity.TravelRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.TravelRequest, error)
	Update(ctx context.Context, travelRequest *entity.TravelRequest) error
	List(ctx context.Context, filters utils.TravelRequestFilters) ([]entity.TravelRequest, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, filters utils.TravelRequestFilters) ([]entity.TravelRequest, error)
}
