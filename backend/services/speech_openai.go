// backend/services/speech_openai.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAISpeechService handles both TTS (via /audio/speech) and STT (via
// /audio/transcriptions with Whisper) using the OpenAI API.
type OpenAISpeechService struct {
	apiKey string
	client *http.Client
}

// NewOpenAISpeechService creates a new OpenAISpeechService.
func NewOpenAISpeechService(apiKey string) *OpenAISpeechService {
	return &OpenAISpeechService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Synthesize converts text to audio using the OpenAI TTS API.
// voice should be one of: alloy, echo, fable, onyx, nova, shimmer.
func (s *OpenAISpeechService) Synthesize(text, voice string, options SpeechOptions) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	url := "https://api.openai.com/v1/audio/speech"

	payload := map[string]interface{}{
		"model":           "tts-1",
		"input":           text,
		"voice":           voice,
		"response_format": options.AudioFormat,
		"speed":           options.Speed,
	}

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	return io.ReadAll(resp.Body)
}

// Transcribe converts audio bytes to text using the OpenAI Whisper model.
func (s *OpenAISpeechService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	url := "https://api.openai.com/v1/audio/transcriptions"

	// Build multipart form data manually.
	var b bytes.Buffer
	w := io.Writer(&b)

	w.Write([]byte("------FormBoundary\r\n"))
	w.Write([]byte("Content-Disposition: form-data; name=\"file\"; filename=\"audio.wav\"\r\n"))
	w.Write([]byte("Content-Type: audio/wav\r\n\r\n"))
	w.Write(audioData)
	w.Write([]byte("\r\n"))

	w.Write([]byte("------FormBoundary\r\n"))
	w.Write([]byte("Content-Disposition: form-data; name=\"model\"\r\n\r\n"))
	w.Write([]byte("whisper-1\r\n"))

	if options.Language != "" {
		w.Write([]byte("------FormBoundary\r\n"))
		w.Write([]byte("Content-Disposition: form-data; name=\"language\"\r\n\r\n"))
		w.Write([]byte(options.Language + "\r\n"))
	}

	w.Write([]byte("------FormBoundary--\r\n"))

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----FormBoundary")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil
}
