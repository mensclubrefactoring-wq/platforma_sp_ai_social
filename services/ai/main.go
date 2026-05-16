package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// newUUID generates a random UUID v4 for RqUID
func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// Токен и время его истечения
var (
	tokenMutex  sync.Mutex
	cachedToken string
	tokenExpiry time.Time
)

func main() {
	// Загружаем .env для получения GIGACHAT_API_KEY
	godotenv.Load()

	// Эндпоинт для вопросов к AI
	http.HandleFunc("/ai/ask", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Prompt string `json:"prompt"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			sendJSONError(w, "Failed to read request", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		log.Printf("📝 AI Request: %s", req.Prompt)

		response, err := callGigaChat(req.Prompt)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			sendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"choices": []map[string]interface{}{
				{
					"message": map[string]string{
						"content": response,
					},
				},
			},
		})
	})

	// Эндпоинт для генерации предложений
	http.HandleFunc("/ai/generate-proposal", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Budget      string `json:"budget"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			sendJSONError(w, "Failed to read request", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		prompt := fmt.Sprintf(`Сформируй профессиональное предложение от социального предпринимателя для следующей задачи бизнеса:

Название задачи: %s
Описание: %s
Категория: %s
Бюджет: %s

Предложение должно включать:
1. Краткое описание решения
2. План реализации
3. Ожидаемые социальные и бизнес-результаты
4. Почему именно наша организация подходит для этого

Пиши на русском языке, профессионально и убедительно.`, req.Title, req.Description, req.Category, req.Budget)

		log.Printf("📝 Generating proposal for: %s", req.Title)

		content, err := callGigaChat(prompt)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			sendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"proposal": content})
	})

	// Эндпоинт для классификации задач
	http.HandleFunc("/classify-task", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Description string `json:"description"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			sendJSONError(w, "Failed to read request", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		prompt := fmt.Sprintf(`Классифицируй задачу в одну из категорий: Экология, Образование, Социальное жилье, Обучение ИТ, Помощь пожилым, Инклюзивность.

Задача: %s

Ответь только одним словом - названием категории.`, req.Description)

		log.Printf("📝 Classifying task: %s", req.Description)

		category, err := callGigaChat(prompt)
		if err != nil {
			log.Printf("❌ Error: %v", err)
			sendJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"category": category})
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "ai"})
	})

	port := os.Getenv("AI_PORT")
	if port == "" {
		port = "3003"
	}

	log.Printf("🚀 AI Service running on :%s", port)
	log.Printf("✅ Endpoints:")
	log.Printf("   POST /ai/ask - задать вопрос AI")
	log.Printf("   POST /ai/generate-proposal - сгенерировать предложение")
	log.Printf("   POST /classify-task - классифицировать задачу")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func sendJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// getAccessToken - получение токена GigaChat
func getAccessToken() (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	// Если токен еще валиден, возвращаем его
	if cachedToken != "" && time.Now().Before(tokenExpiry) {
		return cachedToken, nil
	}

	apiKey := strings.TrimSpace(os.Getenv("GIGACHAT_API_KEY"))
	if apiKey == "" {
		return "", fmt.Errorf("GIGACHAT_API_KEY not set")
	}

	log.Println("🔑 Requesting new GigaChat access token...")

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
	scope := os.Getenv("GIGACHAT_SCOPE")
	if scope == "" {
		scope = "GIGACHAT_API_PERS"
	}
	data := "scope=" + scope

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", newUUID())
	req.Header.Set("Authorization", "Basic "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token error (%d): %s", resp.StatusCode, string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("failed to parse token response: %v", err)
	}

	cachedToken = res.AccessToken
	// Запас 10 минут до истечения (было 5)
	tokenExpiry = time.Unix(res.ExpiresAt/1000, 0).Add(-10 * time.Minute)

	log.Println("✅ New access token received")
	return cachedToken, nil
}

// callGigaChat - отправка запроса к GigaChat API
func callGigaChat(prompt string) (string, error) {
	return callGigaChatWithRetry(prompt, true)
}

func callGigaChatWithRetry(prompt string, allowRetry bool) (string, error) {
	token, err := getAccessToken()
	if err != nil {
		return "", err
	}

	log.Println("🔄 Calling GigaChat API...")

	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	url := "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": "GigaChat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
		"max_tokens":  2000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	log.Printf("📡 GigaChat response status: %d", resp.StatusCode)

	if resp.StatusCode == 401 && allowRetry {
		log.Printf("⚠️ GigaChat: Token expired or invalid (Status 401). Body: %s", string(body))
		log.Println("🔄 Force refreshing token and retrying once...")
		tokenMutex.Lock()
		cachedToken = "" // Инвалидируем токен
		tokenExpiry = time.Time{} // Сбрасываем время
		tokenMutex.Unlock()
		return callGigaChatWithRetry(prompt, false)
	}

	if resp.StatusCode != 200 {
		log.Printf("❌ GigaChat API Error (%d): %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("GigaChat API error (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(result.Choices) > 0 {
		log.Println("✅ GigaChat response received successfully")
		return result.Choices[0].Message.Content, nil
	}

	return "Извините, не удалось получить ответ от AI помощника.", nil
}
