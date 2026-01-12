package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenRouterService handles OpenRouter API interactions
type OpenRouterService struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

// NewOpenRouterService creates a new OpenRouter service
func NewOpenRouterService(apiKey string) *OpenRouterService {
	return &OpenRouterService{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 5 * time.Minute},
		baseURL: "https://openrouter.ai/api/v1",
	}
}

// GenerateResponse generates a response using OpenRouter
func (s *OpenRouterService) GenerateResponse(modelName, prompt string) (string, error) {
	req := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonData, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", s.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenRouter API error: %s", string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	return result.Choices[0].Message.Content, nil
}
