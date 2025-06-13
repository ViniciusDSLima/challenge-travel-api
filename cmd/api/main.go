package main

import (
	database "challenge-travel-api/config"
	"challenge-travel-api/internal/infrastructure/container"
	"challenge-travel-api/internal/interface/router"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           API de Solicitações de Viagem
// @version         1.0
// @description     API para gerenciamento de solicitações de viagem
// @termsOfService  http://swagger.io/terms/

// @contact.name   Suporte
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Digite "Bearer" seguido de um espaço e o token JWT.
func main() {

	if err := godotenv.Load(); err != nil {
		log.Printf(".env file not found %v", err)
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	database.DatabaseConnection()

	log.Println("Executando migrations...")
	cmd := exec.Command("go", "run", "cmd/migrate/migrate.go", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}
	log.Println("Migrations executadas com sucesso!")

	db := database.GetDB()

	authController, travelController := container.Container(db)

	r := router.SetupRouter(authController, travelController)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on port: %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erorr to ruunning server on port: %v", port)
	}
}
