package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"platforma-sp/internal/ai"
	"platforma-sp/internal/auth"
	"platforma-sp/internal/db"
	"platforma-sp/internal/shared"
	"platforma-sp/internal/tasks"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", 401)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &shared.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return shared.JWT_SECRET, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", 401)
			return
		}
		r.Header.Set("User-ID", fmt.Sprintf("%d", claims.UserID))
		r.Header.Set("User-Role", claims.Role)
		next.ServeHTTP(w, r)
	}
}

func adminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("User-Role")
		if role != "admin" {
			http.Error(w, "Forbidden: Admins only", 403)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func main() {
	db.InitDB()

	r := mux.NewRouter()

	// Auth
	r.HandleFunc("/api/auth/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/api/auth/me", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var user shared.User
		userID := r.Header.Get("User-ID")
		db.DB.First(&user, userID)
		user.Password = ""
		json.NewEncoder(w).Encode(user)
	})).Methods("GET")

	// Tasks
	r.HandleFunc("/api/tasks", tasks.GetTasksHandler).Methods("GET")
	r.HandleFunc("/api/tasks", authMiddleware(tasks.CreateTaskHandler)).Methods("POST")

	// AI
	r.HandleFunc("/api/ai/history", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var messages []shared.AIChatMessage
		userID := r.Header.Get("User-ID")
		db.DB.Where("user_id = ?", userID).Order("created_at asc").Find(&messages)
		json.NewEncoder(w).Encode(messages)
	})).Methods("GET")
	r.HandleFunc("/api/ai/ask", authMiddleware(ai.AskAIHandler)).Methods("POST")
	r.HandleFunc("/api/ai/generate-proposal", authMiddleware(ai.GenerateProposalHandler)).Methods("POST")
	r.HandleFunc("/api/ai/classify", authMiddleware(ai.ClassifyTaskHandler)).Methods("POST")

	// Admin / Entrepreneur List
	r.HandleFunc("/api/admin/portfolios", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var users []shared.User
		db.DB.Where("role = ?", "entrepreneur").Find(&users)
		json.NewEncoder(w).Encode(users)
	})).Methods("GET")

	// Client serving
	distPath := "./dist"
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(distPath, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) || r.URL.Path == "/" || !strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			return
		}
		http.FileServer(http.Dir(distPath)).ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" { port = "3000" }
	fmt.Printf("Platforma SP (Microservices Gate) on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
