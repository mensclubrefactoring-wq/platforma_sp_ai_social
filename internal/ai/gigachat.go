package ai

import (
	"bytes"
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
	"crypto/rand"
)

var (
	tokenMutex  sync.Mutex
	cachedToken string
	tokenExpiry time.Time
)

func newUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func getAccessToken() (string, error) {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	if cachedToken != "" && time.Now().Before(tokenExpiry) {
		return cachedToken, nil
	}

	apiKey := strings.TrimSpace(os.Getenv("GIGACHAT_API_KEY"))
	if apiKey == "" {
		return "", fmt.Errorf("GIGACHAT_API_KEY not set")
	}

	log.Println("🔑 Requesting GigaChat token...")

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	scope := os.Getenv("GIGACHAT_SCOPE")
	if scope == "" {
		scope = "GIGACHAT_API_PERS"
	}

	req, _ := http.NewRequest("POST", "https://ngw.devices.sberbank.ru:9443/api/v2/oauth", bytes.NewBufferString("scope="+scope))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", newUUID())
	req.Header.Set("Authorization", "Basic "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token error: %s", string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}
	json.Unmarshal(body, &res)

	cachedToken = res.AccessToken
	tokenExpiry = time.Unix(res.ExpiresAt/1000, 0).Add(-10 * time.Minute)
	return cachedToken, nil
}

func CallGigaChat(prompt string) (string, error) {
	return callWithRetry(prompt, true)
}

func callWithRetry(prompt string, allowRetry bool) (string, error) {
	token, err := getAccessToken()
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	requestBody := map[string]interface{}{
		"model": "GigaChat",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.7,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "https://gigachat.devices.sberbank.ru/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 401 && allowRetry {
		tokenMutex.Lock()
		cachedToken = ""
		tokenMutex.Unlock()
		return callWithRetry(prompt, false)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GigaChat error: %s", string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	json.Unmarshal(body, &result)

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "No response", nil
}
