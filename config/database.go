package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DatabaseConnection() {
	var err error

	dbURL := os.Getenv("DATABASE_DSN")

	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@postgres:5432/app_database?sslmode=disable"
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	config := &gorm.Config{
		Logger: newLogger,
	}

	DB, err = gorm.Open(postgres.Open(dbURL), config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to the database: %v", err))
	}

	log.Println("Database connection established successfully")

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to configure connection pool: %v", err))
	}

	sqlDB.SetMaxOpenConns(25)

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Pool of connections configured successfully")
}

func GetDB() *gorm.DB {
	return DB
}
