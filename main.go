package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"platforma-sp/internal/shared"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for demo if DATABASE_URL is missing
		log.Println("DATABASE_URL not set. Please set it in .env")
		return
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migration
	db.AutoMigrate(&shared.User{}, &shared.Task{}, &shared.AIChatMessage{})
	log.Println("Database migrated successfully")
}

// --- Middleware ---

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

// --- Handlers ---

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		Phone        string `json:"phone"`
		Name         string `json:"name"`
		Role         string `json:"role"`
		ConsentGiven bool   `json:"consentGiven"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Validation
	if !govalidator.IsEmail(req.Email) {
		http.Error(w, "Invalid email format", 400)
		return
	}
	if req.Phone == "" {
		http.Error(w, "Phone number is required", 400)
		return
	}
	if !req.ConsentGiven {
		http.Error(w, "Consent is required", 400)
		return
	}

	var existing shared.User
	if err := db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		http.Error(w, "User already exists", 400)
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	user := shared.User{
		Email:        req.Email,
		Password:     string(hashed),
		Phone:        req.Phone,
		Name:         req.Name,
		Role:         req.Role,
		ConsentGiven: req.ConsentGiven,
	}

	if err := db.Create(&user).Error; err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	token, _ := createToken(user)
	user.Password = ""
	json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": user})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	json.NewDecoder(r.Body).Decode(&req)

	var user shared.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	token, _ := createToken(user)
	user.Password = ""
	json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": user})
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	var user shared.User
	userID := r.Header.Get("User-ID")
	if err := db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", 404)
		return
	}
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []shared.Task
	query := db.Model(&shared.Task{})

	// Search & Filters
	search := r.URL.Query().Get("search")
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	category := r.URL.Query().Get("category")
	if category != "" && category != "Все" {
		query = query.Where("category = ?", category)
	}
	location := r.URL.Query().Get("location")
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}

	query.Order("created_at desc").Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task shared.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	userID := r.Header.Get("User-ID")
	fmt.Sscanf(userID, "%d", &task.CreatorID)
	task.Status = "active"

	if err := db.Create(&task).Error; err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(task)
}

// --- AI Handlers ---

func aiAskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Prompt string `json:"prompt"` }
	json.NewDecoder(r.Body).Decode(&req)
	
	userID := r.Header.Get("User-ID")
	var uid uint
	fmt.Sscanf(userID, "%d", &uid)

	// Save User Message
	db.Create(&shared.AIChatMessage{UserID: uid, Role: "user", Content: req.Prompt})

	// Get History for Context (last 10)
	var history []shared.AIChatMessage
	db.Where("user_id = ?", uid).Order("created_at desc").Limit(10).Find(&history)

	// Call AI
	aiResponse, err := callAI(req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Save Assistant Message
	db.Create(&shared.AIChatMessage{UserID: uid, Role: "assistant", Content: aiResponse})

	json.NewEncoder(w).Encode(map[string]string{"response": aiResponse})
}

func generateProposalHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title, Description, Category, Budget string
	}
	json.NewDecoder(r.Body).Decode(&req)
	prompt := fmt.Sprintf(`Сформируй профессиональное предложение от социального предпринимателя для задачи: %s. Описание: %s. Категория: %s. Бюджет: %s.`, req.Title, req.Description, req.Category, req.Budget)
	
	content, err := callAI(prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"proposal": content})
}

func classifyTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Description string `json:"description"` }
	json.NewDecoder(r.Body).Decode(&req)

	prompt := fmt.Sprintf(`Классифицируй следующую бизнес-задачу в одну из категорий: Экология, Образование, Социальное жилье, Обучение ИТ, Помощь пожилым, Инклюзивность. Ответь только одним названием категории.
Задача: %s`, req.Description)

	category, err := callAI(prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"category": strings.TrimSpace(category)})
}

// --- Shared AI Logic ---

func callAI(prompt string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	provider := os.Getenv("AI_PROVIDER") // deepseek, gemini
	baseURL := os.Getenv("AI_BASE_URL")
	model := os.Getenv("AI_MODEL")

	if provider == "deepseek" {
		if baseURL == "" { baseURL = "https://api.deepseek.com" }
		if model == "" { model = "deepseek-chat" }
	} else {
		if baseURL == "" { baseURL = "https://generativelanguage.googleapis.com/v1beta/openai" }
		if model == "" { model = "gpt-4o" }
	}

	url := fmt.Sprintf("%s/v1/chat/completions", baseURL)
	body, _ := json.Marshal(map[string]interface{}{
		"model": model,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	})
	aiReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	aiReq.Header.Set("Content-Type", "application/json")
	aiReq.Header.Set("Authorization", "Bearer "+apiKey)
	
	resp, err := http.DefaultClient.Do(aiReq)
	if err != nil { return "", err }
	defer resp.Body.Close()

	var resData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil { return "", err }
	if len(resData.Choices) > 0 {
		return resData.Choices[0].Message.Content, nil
	}
	return "No AI Response", nil
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
	initDB()

	r := mux.NewRouter()

	// Auth
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/auth/me", authMiddleware(meHandler)).Methods("GET")

	// Tasks
	r.HandleFunc("/api/tasks", getTasksHandler).Methods("GET")
	r.HandleFunc("/api/tasks", authMiddleware(createTaskHandler)).Methods("POST")

	// AI
	r.HandleFunc("/api/ai/history", authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var messages []shared.AIChatMessage
		userID := r.Header.Get("User-ID")
		db.Where("user_id = ?", userID).Order("created_at asc").Find(&messages)
		json.NewEncoder(w).Encode(messages)
	})).Methods("GET")
	r.HandleFunc("/api/ai/ask", authMiddleware(aiAskHandler)).Methods("POST")
	r.HandleFunc("/api/ai/generate-proposal", authMiddleware(generateProposalHandler)).Methods("POST")
	r.HandleFunc("/api/ai/classify", authMiddleware(classifyTaskHandler)).Methods("POST")

	// Admin (Placeholder)
	r.HandleFunc("/api/admin/portfolios", authMiddleware(adminMiddleware(func(w http.ResponseWriter, r *http.Request) {
		var users []shared.User
		db.Where("role = ?", "entrepreneur").Find(&users)
		json.NewEncoder(w).Encode(users)
	}))).Methods("GET")

	// Static files
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
	fmt.Printf("Platforma SP (GORM/Postgres) running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
