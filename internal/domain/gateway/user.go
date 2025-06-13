package gateway

import (
	"challenge-travel-api/internal/domain/entity"
	"context"
	"github.com/google/uuid"
)

type UserGateway interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
}
