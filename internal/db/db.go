package db

import (
	"log"
	"os"
	"platforma-sp/internal/shared"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set.")
		return
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB.AutoMigrate(&shared.User{}, &shared.Task{}, &shared.AIChatMessage{})
	log.Println("PostgreSQL initialized (Microservice Data Layer)")
}
