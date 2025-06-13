package repository

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/gateway"
	"challenge-travel-api/internal/utils"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrTravelRequestNotFound = errors.New("solicitação de viagem não encontrada")
	ErrCannotCancelApproved  = errors.New("não é possível cancelar uma solicitação já aprovada após 24 horas")
	ErrUnauthorized          = errors.New("usuário não autorizado para esta operação")
)

type TravelRequestRepository struct {
	db *gorm.DB
}

func NewTravelRequestRepository(db *gorm.DB) gateway.TravelRequestGateway {
	return &TravelRequestRepository{
		db: db,
	}
}

func (r *TravelRequestRepository) Create(ctx context.Context, travelRequest *entity.TravelRequest) error {

	return r.db.WithContext(ctx).Create(travelRequest).Error
}

func (r *TravelRequestRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.TravelRequest, error) {
	var travelRequest entity.TravelRequest

	err := r.db.WithContext(ctx).First(&travelRequest, id).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, ErrTravelRequestNotFound
	}

	return &travelRequest, nil
}

func (r *TravelRequestRepository) Update(ctx context.Context, travelRequest *entity.TravelRequest) error {
	return r.db.WithContext(ctx).Updates(travelRequest).Where("id = ?", travelRequest.Id).Error
}

func (r *TravelRequestRepository) List(ctx context.Context, filters utils.TravelRequestFilters) ([]entity.TravelRequest, error) {
	var requests []entity.TravelRequest
	query := r.db.WithContext(ctx).Find(requests)

	if filters.Status != nil {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.StartDate != nil {
		query = query.Where("start_date >= ?", filters.StartDate)
	}

	if filters.EndDate != nil {
		query = query.Where("end_date <= ?", filters.EndDate)
	}

	if filters.DestinationName != nil {
		query = query.Where("destination ILIKE ?", "%"+*filters.DestinationName+"%")
	}

	if filters.UserId != nil {
		query = query.Where("user_id = ?", filters.UserId)
	}

	err := query.Error
	return requests, err
}

func (r *TravelRequestRepository) ListByUserID(ctx context.Context, userID uuid.UUID, filters utils.TravelRequestFilters) ([]entity.TravelRequest, error) {
	var requests []entity.TravelRequest

	err := r.db.WithContext(ctx).Find(requests).Where("user_id = ?", userID).Error

	if err != nil {
		return requests, err
	}

	return requests, nil
}
