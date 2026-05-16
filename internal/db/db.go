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
		log.Println("WARNING: DATABASE_URL not set. Application will run in degraded mode (no database).")
		return
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("ERROR: Failed to connect to database: %v. Application will run in degraded mode.", err)
		return
	}

	DB.AutoMigrate(&shared.User{}, &shared.Task{}, &shared.AIChatMessage{})
	log.Println("PostgreSQL initialized (Microservice Data Layer)")
}
