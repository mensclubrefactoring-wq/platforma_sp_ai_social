package shared

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte(getEnv("JWT_SECRET", "platforma-sp-secret-key-2026"))

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Budget      string    `json:"budget"`
	Deadline    string    `json:"deadline"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	CreatorID   string    `json:"creatorId"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Claims struct {
	UserID string `json:"id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
