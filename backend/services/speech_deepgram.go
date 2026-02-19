// backend/services/speech_deepgram.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DeepgramService transcribes audio using the Deepgram Nova cloud API.
type DeepgramService struct {
	apiKey string
	client *http.Client
}

// NewDeepgramService creates a new DeepgramService.
func NewDeepgramService(apiKey string) *DeepgramService {
	return &DeepgramService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Transcribe sends audio bytes to Deepgram and returns the transcript text.
func (s *DeepgramService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("Deepgram API key not configured")
	}

	url := fmt.Sprintf("https://api.deepgram.com/v1/listen?model=nova&language=%s", options.Language)

	req, err := http.NewRequest("POST", url, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+s.apiKey)
	req.Header.Set("Content-Type", "audio/wav")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Deepgram API error: %s", string(body))
	}

	var result struct {
		Results struct {
			Channels []struct {
				Alternatives []struct {
					Transcript string `json:"transcript"`
				} `json:"alternatives"`
			} `json:"channels"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Results.Channels) > 0 && len(result.Results.Channels[0].Alternatives) > 0 {
		return result.Results.Channels[0].Alternatives[0].Transcript, nil
	}

	return "", fmt.Errorf("no transcript found in Deepgram response")
}
