package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/ai/ask", func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Prompt string `json:"prompt"` }
		json.NewDecoder(r.Body).Decode(&req)

		apiKey := os.Getenv("GEMINI_API_KEY")
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

	http.HandleFunc("/ai/generate-proposal", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Budget      string `json:"budget"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		prompt := fmt.Sprintf(`Сформируй профессиональное предложение от социального предпринимателя для следующей задачи бизнеса:
Название задачи: %s
Описание: %s
Категория: %s
Бюджет: %s

Предложение должно включать:
1. Краткое описание решения.
2. План реализации.
3. Ожидаемые социальные и бизнес-результаты.
4. Почему именно наша организация подходит для этого.
Пеши на русском языке.`, req.Title, req.Description, req.Category, req.Budget)

		// Reusing base AI call logic would be better but for MVP this is fine
		apiKey := os.Getenv("GEMINI_API_KEY")
		baseURL := os.Getenv("AI_BASE_URL")
		if baseURL == "" { baseURL = "https://generativelanguage.googleapis.com/v1beta/openai" }
		model := os.Getenv("AI_MODEL")
		if model == "" { model = "gpt-4o" }

		url := fmt.Sprintf("%s/v1/chat/completions", baseURL)
		body, _ := json.Marshal(map[string]interface{}{
			"model":    model,
			"messages": []map[string]string{{"role": "user", "content": prompt}},
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

		var resData struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}
		json.NewDecoder(resp.Body).Decode(&resData)
		
		content := ""
		if len(resData.Choices) > 0 {
			content = resData.Choices[0].Message.Content
		}

		json.NewEncoder(w).Encode(map[string]string{"proposal": content})
	})

	fmt.Println("AI Service running on :3003")
	log.Fatal(http.ListenAndServe(":3003", nil))
}
