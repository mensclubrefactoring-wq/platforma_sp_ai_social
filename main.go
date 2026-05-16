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

	"github.com/joho/godotenv"
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
			return shared.GetJWTSecret(), nil
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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, User-ID, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Загружаем .env
	godotenv.Load()

	// 1. Инициализация БД (PostgreSQL)
	db.InitDB()

	router := mux.NewRouter()

	// API endpoints (Реальная логика из внутренних пакетов)
	router.HandleFunc("/api/auth/register", auth.RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", auth.LoginHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/me", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var user shared.User
		userID := r.Header.Get("User-ID")
		if db.DB != nil {
			db.DB.First(&user, userID)
			user.Password = ""
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "DB not initialized", 500)
		}
	})).Methods("GET", "OPTIONS")

	// Tasks
	router.HandleFunc("/api/tasks", authMiddleware(tasks.GetTasksHandler)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tasks", authMiddleware(tasks.CreateTaskHandler)).Methods("POST", "OPTIONS")

	// AI эндпоинты (теперь напрямую)
	router.HandleFunc("/api/ai/history", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var messages []shared.AIChatMessage
		userIDStr := r.Header.Get("User-ID")
		var uid uint
		fmt.Sscanf(userIDStr, "%d", &uid)
		if db.DB != nil {
			db.DB.Where("user_id = ?", uid).Order("created_at asc").Find(&messages)
			json.NewEncoder(w).Encode(messages)
		}
	})).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/ai/ask", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Prompt string `json:"prompt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		userIDStr := r.Header.Get("User-ID")
		var uid uint
		fmt.Sscanf(userIDStr, "%d", &uid)

		// Сохраняем вопрос
		if db.DB != nil {
			db.DB.Create(&shared.AIChatMessage{UserID: uid, Role: "user", Content: req.Prompt})
		}

		resp, err := ai.CallGigaChat(req.Prompt)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Сохраняем ответ
		if db.DB != nil {
			db.DB.Create(&shared.AIChatMessage{UserID: uid, Role: "assistant", Content: resp})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": resp,
					},
				},
			},
		})
	})).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/ai/generate-proposal", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Budget      string `json:"budget"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		prompt := fmt.Sprintf(`Сформируй профессиональное предложение от социального предпринимателя для задачи бизнеса:
Название: %s
Описание: %s
Категория: %s
Бюджет: %s`, req.Title, req.Description, req.Category, req.Budget)

		resp, err := ai.CallGigaChat(prompt)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"proposal": resp})
	})).Methods("POST", "OPTIONS")

	// Admin
	router.HandleFunc("/api/admin/portfolios", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var users []shared.User
		if db.DB != nil {
			db.DB.Where("role = ?", "entrepreneur").Find(&users)
			json.NewEncoder(w).Encode(users)
		}
	})).Methods("GET", "OPTIONS")

	// Статические файлы фронтенда
	distPath := "./dist"
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(distPath, r.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) || r.URL.Path == "/" || !strings.Contains(r.URL.Path, ".") {
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			return
		}
		http.FileServer(http.Dir(distPath)).ServeHTTP(w, r)
	})

	router.Use(corsMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🌐 Platforma SP (Gateway) on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
}

