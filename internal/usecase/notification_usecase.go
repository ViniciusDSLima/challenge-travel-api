package usecase

import (
	"challenge-travel-api/internal/domain/entity"
	"challenge-travel-api/internal/domain/enums"
	"fmt"
	"log"
)

type NotificationUseCae interface {
	NotifyStatusChange(travelRequest *entity.TravelRequest, previousStatus enums.TravelRequestStatus)
}

type EmailNotificationService struct {
}

func NewEmailNotificationService() NotificationUseCae {
	return &EmailNotificationService{}
}

func (s *EmailNotificationService) NotifyStatusChange(travelRequest *entity.TravelRequest, previousStatus enums.TravelRequestStatus) {
	if (travelRequest.Status == enums.TravelRequestStatusApproved || travelRequest.Status == enums.TravelRequestStatusCanceled) &&
		travelRequest.Status != previousStatus {

		user := travelRequest.User
		message := ""

		switch travelRequest.Status {
		case enums.TravelRequestStatusApproved:
			message = fmt.Sprintf(
				"Olá %s, seu pedido de viagem para %s foi APROVADO! Datas: %s a %s",
				user.Name,
				travelRequest.DestinationName,
				travelRequest.DepartureDate.Format("02/01/2006"),
				travelRequest.ReturnDate.Format("02/01/2006"),
			)
		case enums.TravelRequestStatusCanceled:
			message = fmt.Sprintf(
				"Olá %s, seu pedido de viagem para %s foi CANCELADO.",
				user.Name,
				travelRequest.DestinationName,
			)
		}

		log.Printf("[NOTIFICATION] E-mail enviado para %s: %s", user.Email, message)
	}
}
