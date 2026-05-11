package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"platforma-sp/internal/shared"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// --- Shared Logic ---

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

// --- Auth Service Logic ---

type AuthDB struct {
	Users []shared.User `json:"users"`
}

func getAuthDB() AuthDB {
	var db AuthDB
	data, err := os.ReadFile("db_auth.json")
	if err != nil {
		return AuthDB{Users: []shared.User{}}
	}
	json.Unmarshal(data, &db)
	return db
}

func saveAuthDB(db AuthDB) {
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile("db_auth.json", data, 0644)
}

func startAuthService() {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		var req shared.User
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		db := getAuthDB()
		// Check for existing user
		for _, u := range db.Users {
			if u.Email == req.Email {
				http.Error(w, "User already exists", 400)
				return
			}
		}

		req.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		req.Password = string(hashed)
		db.Users = append(db.Users, req)
		saveAuthDB(db)
		
		token, _ := createToken(req)
		req.Password = "" // Hide password
		json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": req})
	})

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Email, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		db := getAuthDB()
		for _, u := range db.Users {
			if u.Email == req.Email {
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err == nil {
					token, _ := createToken(u)
					u.Password = "" // Hide password
					json.NewEncoder(w).Encode(map[string]interface{}{"token": token, "user": u})
					return
				}
			}
		}
		http.Error(w, "Invalid credentials", 401)
	})

	mux.HandleFunc("/auth/me", func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Unauthorized", 401)
			return
		}
		db := getAuthDB()
		for _, u := range db.Users {
			if u.ID == claims.UserID {
				u.Password = ""
				json.NewEncoder(w).Encode(u)
				return
			}
		}
		http.Error(w, "User not found", 404)
	})

	fmt.Println("Auth Service listening on :3001")
	log.Fatal(http.ListenAndServe(":3001", mux))
}

// --- Tasks Service Logic ---

type TasksDB struct {
	Tasks []shared.Task `json:"tasks"`
}

func getTasksDB() TasksDB {
	var db TasksDB
	data, err := os.ReadFile("db_tasks.json")
	if err != nil {
		return TasksDB{Tasks: []shared.Task{}}
	}
	json.Unmarshal(data, &db)
	return db
}

func saveTasksDB(db TasksDB) {
	data, _ := json.MarshalIndent(db, "", "  ")
	os.WriteFile("db_tasks.json", data, 0644)
}

func startTasksService() {
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode(getTasksDB().Tasks)
			return
		}
		if r.Method == "POST" {
			var task shared.Task
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			task.ID = fmt.Sprintf("%d", time.Now().UnixNano())
			task.CreatedAt = time.Now()
			task.Status = "active"
			db := getTasksDB()
			db.Tasks = append(db.Tasks, task)
			saveTasksDB(db)
			json.NewEncoder(w).Encode(task)
			return
		}
	})
	fmt.Println("Tasks Service listening on :3002")
	log.Fatal(http.ListenAndServe(":3002", mux))
}

// --- AI Service Logic ---

func startAIService() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ai/ask", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Prompt string `json:"prompt"` }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			http.Error(w, "GEMINI_API_KEY not set", 500)
			return
		}
		baseURL := os.Getenv("AI_BASE_URL")
		if baseURL == "" { baseURL = "https://generativelanguage.googleapis.com/v1beta/openai" }
		model := os.Getenv("AI_MODEL")
		if model == "" { model = "gpt-4o" }

		url := fmt.Sprintf("%s/v1/chat/completions", baseURL)
		body, _ := json.Marshal(map[string]interface{}{
			"model": model,
			"messages": []map[string]string{{"role": "user", "content": req.Prompt}},
		})
		aiReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		aiReq.Header.Set("Content-Type", "application/json")
		aiReq.Header.Set("Authorization", "Bearer "+apiKey)
		resp, err := http.DefaultClient.Do(aiReq)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer resp.Body.Close()
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		json.NewEncoder(w).Encode(result)
	})

	mux.HandleFunc("/ai/generate-proposal", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title, Description, Category, Budget string
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		prompt := fmt.Sprintf(`Сформируй профессиональное предложение от социального предпринимателя для задачи: %s. Описание: %s. Категория: %s. Бюджет: %s.`, req.Title, req.Description, req.Category, req.Budget)
		
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			http.Error(w, "GEMINI_API_KEY not set", 500)
			return
		}
		url := "https://generativelanguage.googleapis.com/v1beta/openai/v1/chat/completions"
		body, _ := json.Marshal(map[string]interface{}{
			"model": "gpt-4o",
			"messages": []map[string]string{{"role": "user", "content": prompt}},
		})
		aiReq, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		aiReq.Header.Set("Authorization", "Bearer "+apiKey)
		aiReq.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(aiReq)
		if err != nil { http.Error(w, err.Error(), 500); return }
		defer resp.Body.Close()
		var resData struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
			http.Error(w, "Failed to parse AI response: "+err.Error(), 500)
			return
		}
		content := ""
		if len(resData.Choices) > 0 {
			content = resData.Choices[0].Message.Content
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"proposal": content})
	})
	fmt.Println("AI Service listening on :3003")
	log.Fatal(http.ListenAndServe(":3003", mux))
}

// --- Gateway Logic ---

func proxyTo(targetURL string, prefixToStrip string) http.HandlerFunc {
	target, _ := url.Parse(targetURL)
	proxy := httputil.NewSingleHostReverseProxy(target)
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, prefixToStrip)
		r.Host = target.Host
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	// Start sub-services in goroutines
	go startAuthService()
	go startTasksService()
	go startAIService()

	// Wait a bit for services to bind
	time.Sleep(100 * time.Millisecond)

	r := mux.NewRouter()
	r.PathPrefix("/api/auth").HandlerFunc(proxyTo("http://localhost:3001", "/api"))
	r.PathPrefix("/api/tasks").HandlerFunc(proxyTo("http://localhost:3002", "/api"))
	r.PathPrefix("/api/ai").HandlerFunc(proxyTo("http://localhost:3003", "/api"))

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
	fmt.Printf("Gateway (Platforma SP) running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
