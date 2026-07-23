// backend/services/speech_replicate.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ReplicateService runs TTS and STT models hosted on Replicate.com.
// Both operations follow the same pattern: create a prediction → poll until done.
type ReplicateService struct {
	apiKey string
	client *http.Client
}

// NewReplicateService creates a new ReplicateService.
func NewReplicateService(apiKey string) *ReplicateService {
	return &ReplicateService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// Synthesize converts text to audio using a Replicate TTS model.
// modelName should be the full "owner/name:version" identifier.
func (s *ReplicateService) Synthesize(text, modelName string, options SpeechOptions) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("Replicate API key not configured")
	}

	model := "cjwbw/tts:dc1e2c35c85fb6b33c18e8d4b209a452b1a10e1c8e8c0b5b5c5c5c5c5c5c5c5c5c5c5c"
	if modelName != "" {
		model = modelName
	}

	payload := map[string]interface{}{
		"version": model,
		"input": map[string]interface{}{
			"text":          text,
			"speaker":       options.Voice,
			"language":      options.Language,
			"cleanup_voice": true,
		},
	}

	predictionURL, err := s.createPrediction(payload)
	if err != nil {
		return nil, err
	}

	return s.pollPredictionAudio(predictionURL)
}

// Transcribe converts audio bytes to text using a Replicate Whisper model.
// modelName should be the full "owner/name:version" identifier.
func (s *ReplicateService) Transcribe(audioData []byte, modelName string, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("Replicate API key not configured")
	}

	// Step 1: Upload audio to Replicate file storage.
	uploadedURL, err := s.uploadAudio(audioData)
	if err != nil {
		return "", err
	}

	model := "openai/whisper:91ee9c0c3df30478510ff8c8a3a545add29ad025ca642e1c606e318e6eab18b25"
	if modelName != "" {
		model = modelName
	}

	payload := map[string]interface{}{
		"version": model,
		"input": map[string]interface{}{
			"audio":                             uploadedURL,
			"model":                             "large-v3",
			"language":                          options.Language,
			"translate":                         false,
			"temperature":                       0,
			"transcription":                     "plain text",
			"suppress_tokens":                   "-1",
			"log_probability_threshold":         -1.0,
			"no_speech_threshold":               0.6,
			"condition_on_previous_text":        true,
			"compression_ratio_threshold":       2.4,
			"temperature_increment_on_fallback": 0.2,
		},
	}

	predictionURL, err := s.createPrediction(payload)
	if err != nil {
		return "", err
	}

	return s.pollPredictionText(predictionURL)
}

// createPrediction posts a prediction job and returns the polling URL.
func (s *ReplicateService) createPrediction(payload map[string]interface{}) (string, error) {
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Replicate prediction error: %s", string(body))
	}

	var prediction struct {
		URLs struct {
			Get string `json:"get"`
		} `json:"urls"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return "", err
	}
	return prediction.URLs.Get, nil
}

// uploadAudio uploads audio bytes to Replicate file storage and returns the URL.
func (s *ReplicateService) uploadAudio(audioData []byte) (string, error) {
	uploadPayload := map[string]interface{}{
		"content_type": "audio/wav",
		"filename":     "audio.wav",
	}
	jsonData, _ := json.Marshal(uploadPayload)

	req, err := http.NewRequest("POST", "https://api.replicate.com/v1/uploads", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Token "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}

	var uploadResult struct {
		UploadURL string `json:"upload_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&uploadResult); err != nil {
		resp.Body.Close()
		return "", err
	}
	resp.Body.Close()

	// PUT the audio bytes to the signed URL.
	putReq, err := http.NewRequest("PUT", uploadResult.UploadURL, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}
	putReq.Header.Set("Content-Type", "audio/wav")

	putResp, err := s.client.Do(putReq)
	if err != nil {
		return "", err
	}
	putResp.Body.Close()

	if putResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Replicate upload failed with status: %d", putResp.StatusCode)
	}

	return uploadResult.UploadURL, nil
}

// pollPredictionAudio polls until a prediction succeeds, then downloads the audio.
func (s *ReplicateService) pollPredictionAudio(predictionURL string) ([]byte, error) {
	for {
		req, _ := http.NewRequest("GET", predictionURL, nil)
		req.Header.Set("Authorization", "Token "+s.apiKey)

		resp, err := s.client.Do(req)
		if err != nil {
			return nil, err
		}

		var result struct {
			Status string `json:"status"`
			Output string `json:"output"`
			Error  string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		switch result.Status {
		case "succeeded":
			audioResp, err := http.Get(result.Output)
			if err != nil {
				return nil, err
			}
			defer audioResp.Body.Close()
			return io.ReadAll(audioResp.Body)
		case "failed":
			return nil, fmt.Errorf("Replicate prediction failed: %s", result.Error)
		}

		time.Sleep(2 * time.Second)
	}
}

// pollPredictionText polls until a prediction succeeds, then returns the text output.
func (s *ReplicateService) pollPredictionText(predictionURL string) (string, error) {
	for {
		req, _ := http.NewRequest("GET", predictionURL, nil)
		req.Header.Set("Authorization", "Token "+s.apiKey)

		resp, err := s.client.Do(req)
		if err != nil {
			return "", err
		}

		var result struct {
			Status string `json:"status"`
			Output struct {
				Text string `json:"text"`
			} `json:"output"`
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return "", err
		}
		resp.Body.Close()

		switch result.Status {
		case "succeeded":
			return result.Output.Text, nil
		case "failed":
			return "", fmt.Errorf("Replicate transcription failed: %s", result.Error)
		}

		time.Sleep(2 * time.Second)
	}
}
