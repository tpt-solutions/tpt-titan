// backend/services/speech_elevenlabs.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ElevenLabsService synthesises speech via the ElevenLabs cloud API.
type ElevenLabsService struct {
	apiKey string
	client *http.Client
}

// NewElevenLabsService creates a new ElevenLabsService.
func NewElevenLabsService(apiKey string) *ElevenLabsService {
	return &ElevenLabsService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Synthesize converts text to MP3 audio bytes using ElevenLabs.
// voiceID is the ElevenLabs voice identifier (e.g. "21m00Tcm4TlvDq8ikWAM").
func (s *ElevenLabsService) Synthesize(text, voiceID string, options SpeechOptions) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("ElevenLabs API key not configured")
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)

	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_monolingual_v1",
		"voice_settings": map[string]interface{}{
			"stability":        0.5,
			"similarity_boost": 0.5,
		},
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ElevenLabs API error: %s", string(body))
	}

	return io.ReadAll(resp.Body)
}
