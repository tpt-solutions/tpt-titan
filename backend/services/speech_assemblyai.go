// backend/services/speech_assemblyai.go
package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AssemblyAIService transcribes audio using the AssemblyAI cloud API.
// The workflow is: upload audio → request transcript → poll until done.
type AssemblyAIService struct {
	apiKey string
	client *http.Client
}

// NewAssemblyAIService creates a new AssemblyAIService.
func NewAssemblyAIService(apiKey string) *AssemblyAIService {
	return &AssemblyAIService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// Transcribe uploads audioData, requests a transcript, and polls until the
// transcription is complete or has failed.
func (s *AssemblyAIService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("AssemblyAI API key not configured")
	}

	// Step 1: Upload the audio file.
	audioURL, err := s.uploadAudio(audioData)
	if err != nil {
		return "", err
	}

	// Step 2: Request a transcription job.
	jobID, err := s.requestTranscript(audioURL, options)
	if err != nil {
		return "", err
	}

	// Step 3: Poll until the job is finished.
	return s.pollTranscript(jobID)
}

// uploadAudio uploads raw audio bytes and returns the hosted upload URL.
func (s *AssemblyAIService) uploadAudio(audioData []byte) (string, error) {
	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AssemblyAI upload error: %s", string(body))
	}

	var result struct {
		UploadURL string `json:"upload_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.UploadURL, nil
}

// requestTranscript submits the transcription job and returns the job ID.
func (s *AssemblyAIService) requestTranscript(audioURL string, options SpeechOptions) (string, error) {
	payload := map[string]interface{}{
		"audio_url":     audioURL,
		"language_code": options.Language,
	}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AssemblyAI transcript request error: %s", string(body))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ID, nil
}

// pollTranscript polls the transcript endpoint until the job completes or fails.
func (s *AssemblyAIService) pollTranscript(jobID string) (string, error) {
	pollURL := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", jobID)

	for {
		req, _ := http.NewRequest("GET", pollURL, nil)
		req.Header.Set("Authorization", s.apiKey)

		resp, err := s.client.Do(req)
		if err != nil {
			return "", err
		}

		var result struct {
			Status string `json:"status"`
			Text   string `json:"text"`
			Error  string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return "", err
		}
		resp.Body.Close()

		switch result.Status {
		case "completed":
			return result.Text, nil
		case "error":
			return "", fmt.Errorf("AssemblyAI transcription error: %s", result.Error)
		}

		time.Sleep(2 * time.Second)
	}
}
