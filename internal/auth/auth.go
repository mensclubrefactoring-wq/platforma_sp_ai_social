package auth

import (
	"encoding/json"
	// "fmt"  // временно закомментировать
	"net/http"
	"time"

	"platforma-sp/internal/db"
	"platforma-sp/internal/shared"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Phone              string `json:"phone"`
		RepresentativeName string `json:"representativeName"`
		CompanyName        string `json:"companyName"`
		Role               string `json:"role"`
		ConsentGiven       bool   `json:"consentGiven"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if !govalidator.IsEmail(req.Email) || req.Phone == "" || !req.ConsentGiven {
		http.Error(w, "Invalid input or missing consent", 400)
		return
	}

	var existing shared.User
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		http.Error(w, "User already exists", 400)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user := shared.User{
		Email:              req.Email,
		Password:           string(hashed),
		Phone:              req.Phone,
		RepresentativeName: req.RepresentativeName,
		CompanyName:        req.CompanyName,
		Role:               req.Role,
		ConsentGiven:       req.ConsentGiven,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	token, _ := CreateToken(user)
	user.Password = ""
	json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": user})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	json.NewDecoder(r.Body).Decode(&req)

	var user shared.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	token, _ := CreateToken(user)
	user.Password = ""
	json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": user})
}

func CreateToken(u shared.User) (string, error) {
	claims := &shared.Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(shared.JWT_SECRET)
}
