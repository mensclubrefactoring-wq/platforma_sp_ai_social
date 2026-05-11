package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"platforma-sp/internal/shared"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	Users []shared.User `json:"users"`
}

func getDB() DB {
	var db DB
	data, _ := os.ReadFile("db_auth.json")
	json.Unmarshal(data, &db)
	return db
}

func saveDB(db DB) {
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile("db_auth.json", data, 0644)
}

func createToken(u shared.User) (string, error) {
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

func main() {
	http.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		var req shared.User
		json.NewDecoder(r.Body).Decode(&req)
		
		db := getDB()
		req.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		req.Password = string(hashed)
		
		db.Users = append(db.Users, req)
		saveDB(db)
		
		token, _ := createToken(req)
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": req})
	})

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password string }
		json.NewDecoder(r.Body).Decode(&req)
		
		db := getDB()
		for _, u := range db.Users {
			if u.Email == req.Email {
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err == nil {
					token, _ := createToken(u)
					json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": u})
					return
				}
			}
		}
		http.Error(w, "Invalid credentials", 401)
	})

	fmt.Println("Auth Service running on :3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
