package shared

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte(getEnv("JWT_SECRET", "platforma-sp-secret-key-2026"))

type User struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Email              string    `json:"email" gorm:"uniqueIndex"`
	Password           string    `json:"password,omitempty"`
	Phone              string    `json:"phone"`
	RepresentativeName string    `json:"representativeName"`
	CompanyName        string    `json:"companyName"`
	Role               string    `json:"role"` // business, entrepreneur, admin
	ConsentGiven       bool      `json:"consentGiven"`
	PortfolioURL       string    `json:"portfolioUrl"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Budget      string    `json:"budget"`
	Deadline    string    `json:"deadline"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	CreatorID   uint      `json:"creatorId"`
	Status      string    `json:"status"` // active, closed
	CreatedAt   time.Time `json:"createdAt"`
}

type AIChatMessage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId"`
	Role      string    `json:"role"` // user, assistant
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Claims struct {
	UserID uint   `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
