package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"platforma-sp/internal/db"
	"platforma-sp/internal/shared"
)

func AskAIHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Prompt string `json:"prompt"` }
	json.NewDecoder(r.Body).Decode(&req)
	
	userID := r.Header.Get("User-ID")
	var uid uint
	fmt.Sscanf(userID, "%d", &uid)

	db.DB.Create(&shared.AIChatMessage{UserID: uid, Role: "user", Content: req.Prompt})

	aiResponse, err := CallAI(req.Prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	db.DB.Create(&shared.AIChatMessage{UserID: uid, Role: "assistant", Content: aiResponse})
	json.NewEncoder(w).Encode(map[string]string{"response": aiResponse})
}

func GenerateProposalHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title, Description, Category, Budget string
	}
	json.NewDecoder(r.Body).Decode(&req)
	prompt := fmt.Sprintf(`Сформируй профессиональное предложение от социального предпринимателя для задачи: %s. Описание: %s. Категория: %s. Бюджет: %s.`, req.Title, req.Description, req.Category, req.Budget)
	
	content, err := CallAI(prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"proposal": content})
}

func ClassifyTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Description string `json:"description"` }
	json.NewDecoder(r.Body).Decode(&req)

	prompt := fmt.Sprintf(`Классифицируй задачу в одну из категорий: Экология, Образование, Социальное жилье, Обучение ИТ, Помощь пожилым, Инклюзивность. Один результат.
Задача: %s`, req.Description)

	category, err := CallAI(prompt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"category": strings.TrimSpace(category)})
}

func CallAI(prompt string) (string, error) {
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
	json.NewDecoder(resp.Body).Decode(&resData)
	if len(resData.Choices) > 0 {
		return resData.Choices[0].Message.Content, nil
	}
	return "No AI Response", nil
}
