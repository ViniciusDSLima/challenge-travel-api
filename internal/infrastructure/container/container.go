package container

import (
	"challenge-travel-api/internal/infrastructure/repository"
	"challenge-travel-api/internal/interface/controller"
	"challenge-travel-api/internal/usecase"

	"gorm.io/gorm"
)

func Container(db *gorm.DB) (*controller.AuthController, *controller.TravelController) {
	userRepo := repository.NewUserRepository(db)
	travelRepo := repository.NewTravelRequestRepository(db)

	authUseCase := usecase.NewAUthUseCase(userRepo)
	notificationService := usecase.NewEmailNotificationService()
	travelUseCase := usecase.NewTravelRequestUseCase(travelRepo, userRepo, notificationService)

	authController := controller.NewAuthController(authUseCase)
	travelController := controller.NewTravelController(travelUseCase)

	return authController, travelController

}
