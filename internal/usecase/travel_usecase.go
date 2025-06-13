package usecase

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/enums"
	"challenge-travel-api/internal/domain/gateway"
	"challenge-travel-api/internal/interface/dto"
	"challenge-travel-api/internal/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidDates          = errors.New("data de ida deve ser anterior à data de volta")
	ErrFutureDatesOnly       = errors.New("as datas devem ser futuras")
	ErrInvalidDestination    = errors.New("destino é obrigatório")
	ErrUnauthorized          = errors.New("usuário não autorizado para esta operação")
	ErrTravelAlreadyApproved = errors.New("Viajem já aprovada")
)

type TravelUseCase interface {
	CreateTravelRequest(ctx context.Context, userID uuid.UUID, input dto.CreateTravelRequestDTO) (*entity.TravelRequest, error)
	UpdateTravelRequest(ctx context.Context, id uuid.UUID, userID uuid.UUID, input dto.UpdateTravelRequestDTO) (*entity.TravelRequest, error)
	UpdateStatusTravelRequest(ctx context.Context, userId string, input dto.UpdateStatusTravelRequestDTO) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.TravelRequest, error)
	ListTravelRequests(
		ctx context.Context,
		userID uuid.UUID,
		status *enums.TravelRequestStatus,
		startDate *time.Time,
		endDate *time.Time,
		destinationName *string,
		page int,
		pageSize int,
	) ([]entity.TravelRequest, error)
}

type TravelRequestUseCaseImpl struct {
	travelGateway       gateway.TravelRequestGateway
	userGateway         gateway.UserGateway
	notificationService NotificationUseCae
}

func NewTravelRequestUseCase(
	travelGateway gateway.TravelRequestGateway,
	userGateway gateway.UserGateway,
	notificationService NotificationUseCae,
) *TravelRequestUseCaseImpl {
	return &TravelRequestUseCaseImpl{
		travelGateway:       travelGateway,
		userGateway:         userGateway,
		notificationService: notificationService,
	}
}

func (uc *TravelRequestUseCaseImpl) CreateTravelRequest(
	ctx context.Context,
	userID uuid.UUID,
	input dto.CreateTravelRequestDTO,
) (*entity.TravelRequest, error) {
	if input.DestinationName == "" {
		return nil, ErrInvalidDestination
	}

	now := time.Now()
	if input.DepartureDate.Before(now) {
		return nil, ErrFutureDatesOnly
	}

	if input.ReturnDate != nil && input.DepartureDate.After(*input.ReturnDate) {
		return nil, ErrInvalidDates
	}

	user, err := uc.userGateway.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	travelRequest := &entity.TravelRequest{
		Id:              uuid.New(),
		TravelerName:    input.TravelerName,
		UserId:          userID,
		DestinationName: input.DestinationName,
		DepartureDate:   input.DepartureDate,
		ReturnDate:      input.ReturnDate,
		Status:          enums.TravelRequestStatusSolicited,
		CreatedAt:       now,
		UpdatedAt:       &now,
		User:            *user,
	}

	err = uc.travelGateway.Create(ctx, travelRequest)
	if err != nil {
		return nil, err
	}

	return travelRequest, nil
}

func (uc *TravelRequestUseCaseImpl) UpdateTravelRequest(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	input dto.UpdateTravelRequestDTO,
) (*entity.TravelRequest, error) {
	travelRequest, err := uc.travelGateway.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if travelRequest.UserId != userID {
		return nil, ErrUnauthorized
	}

	if travelRequest.Status != enums.TravelRequestStatusSolicited {
		return nil, errors.New("não é possível alterar um pedido que já foi aprovado ou cancelado")
	}

	if input.DestinationName == nil {
		return nil, ErrInvalidDestination
	}

	now := time.Now()
	if input.DepartureDate.Before(now) {
		return nil, ErrFutureDatesOnly
	}

	if input.ReturnDate != nil && input.DepartureDate.After(*input.ReturnDate) {
		return nil, ErrInvalidDates
	}

	travelRequest.UpdateTravelRequest(input.DestinationName, input.TravelerName, input.DepartureDate, input.ReturnDate, nil, nil, nil)

	err = uc.travelGateway.Update(ctx, travelRequest)
	if err != nil {
		return nil, err
	}

	return travelRequest, nil
}

func (uc *TravelRequestUseCaseImpl) UpdateStatusTravelRequest(ctx context.Context, userId string, input dto.UpdateStatusTravelRequestDTO) error {
	userIdUUID, err := uuid.Parse(userId)

	if err != nil {
		return err
	}

	travelRequestIdUUID, err := uuid.Parse(input.TravelRequestId)

	if err != nil {
		return err
	}

	user, err := uc.userGateway.FindByID(ctx, userIdUUID)

	if err != nil {
		return err
	}

	travel, err := uc.travelGateway.FindByID(ctx, travelRequestIdUUID)

	if err != nil {
		return err
	}

	if travel.UserId == user.Id {
		return ErrUnauthorized
	}

	if user.Role != enums.UserTypeAdmin {
		return ErrUnauthorized
	}

	previousStatus := travel.Status

	if previousStatus == enums.TravelRequestStatusApproved {
		return ErrTravelAlreadyApproved
	}

	var canceledBy *uuid.UUID
	if input.Status == enums.TravelRequestStatusCanceled {
		canceledBy = &user.Id
	}

	var approvedBy *uuid.UUID
	if input.Status == enums.TravelRequestStatusApproved {
		approvedBy = &user.Id
	}

	travel.UpdateTravelRequest(nil, nil, nil, nil, &input.Status, canceledBy, approvedBy)

	err = uc.travelGateway.Update(ctx, travel)

	if err != nil {
		return err
	}

	uc.notificationService.NotifyStatusChange(travel, previousStatus)

	return nil

}
func (uc *TravelRequestUseCaseImpl) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*entity.TravelRequest, error) {
	travelRequest, err := uc.travelGateway.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if travelRequest.UserId != userID {
		return nil, ErrUnauthorized
	}

	return travelRequest, nil
}

func (uc *TravelRequestUseCaseImpl) ListTravelRequests(
	ctx context.Context,
	userID uuid.UUID,
	status *enums.TravelRequestStatus,
	startDate *time.Time,
	endDate *time.Time,
	destinationName *string,
	page int,
	pageSize int,
) ([]entity.TravelRequest, error) {
	filters := utils.TravelRequestFilters{
		Status:          status,
		StartDate:       startDate,
		EndDate:         endDate,
		DestinationName: destinationName,
		Page:            page,
		PageSize:        pageSize,
	}

	return uc.travelGateway.ListByUserID(ctx, userID, filters)
}
