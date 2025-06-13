package router

import (
	_ "challenge-travel-api/docs"
	"challenge-travel-api/internal/interface/controller"
	"challenge-travel-api/internal/interface/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(
	authController *controller.AuthController,
	travelController *controller.TravelController,
) *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseRoute := router.Group("/api/v1")
	{
		auth := baseRoute.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}
	}

	baseRoute.Use(middleware.AuthMiddleware())
	{
		travels := baseRoute.Group("/travels")
		{
			travels.POST("", travelController.CreateTravelRequest)
			travels.GET("", travelController.ListTravelRequests)
			travels.GET("/:id", travelController.GetTravelRequest)
			travels.PUT("/:id", travelController.UpdateTravelRequest)
			travels.PATCH("/:id/status", travelController.UpdateStatusTravelRequest)
		}
	}

	return router
}
