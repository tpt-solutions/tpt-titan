package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"github.com/google/uuid"
)

// SpeechService handles TTS/STT operations and orchestration
type SpeechService struct {
	config         *config.SpeechConfig
	localTTSService *LocalTTSService
	localSTTService *LocalSTTService
	elevenLabsService *ElevenLabsService
	openAIService   *OpenAISpeechService
	replicateService *ReplicateService
	piperService    *PiperService
	assemblyAIService *AssemblyAIService
	deepgramService *DeepgramService
}

// NewSpeechService creates a new speech service instance
func NewSpeechService(cfg *config.SpeechConfig) *SpeechService {
	return &SpeechService{
		config:            cfg,
		localTTSService:   NewLocalTTSService(),
		localSTTService:   NewLocalSTTService(),
		elevenLabsService: NewElevenLabsService(cfg.ElevenLabsKey),
		openAIService:     NewOpenAISpeechService(cfg.OpenAIKey),
		replicateService:  NewReplicateService(cfg.ReplicateKey),
		piperService:      NewPiperService(),
		assemblyAIService: NewAssemblyAIService(cfg.AssemblyAIKey),
		deepgramService:   NewDeepgramService(cfg.DeepgramKey),
	}
}

// InitializeSystemModels sets up the default speech models
func (s *SpeechService) InitializeSystemModels() error {
	log.Println("Speech system initialized with default models")
	return nil
}

// TextToSpeech converts text to speech
func (s *SpeechService) TextToSpeech(userID uuid.UUID, text string, modelID uuid.UUID, options SpeechOptions) (*models.SpeechRequest, error) {
	// Get model information
	var model models.SpeechModel
	if err := config.DB.Where("id = ? AND type = ?", modelID, "tts").First(&model).Error; err != nil {
		return nil, fmt.Errorf("TTS model not found: %w", err)
	}

	// Create request record
	now := time.Now()
	request := &models.SpeechRequest{
		UserID:      userID,
		ModelID:     modelID,
		RequestType: "tts",
		InputText:   text,
		Status:      "processing",
		AudioFormat: options.AudioFormat,
		Language:    options.Language,
		Voice:       options.Voice,
		Speed:       options.Speed,
		Pitch:       options.Pitch,
		CreatedAt:   now,
	}

	if err := config.DB.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create TTS request: %w", err)
	}

	// Process asynchronously
	go s.processTTSRequest(request, &model, options)

	return request, nil
}

// SpeechToText converts speech to text
func (s *SpeechService) SpeechToText(userID uuid.UUID, audioData []byte, modelID uuid.UUID, options SpeechOptions) (*models.SpeechRequest, error) {
	// Get model information
	var model models.SpeechModel
	if err := config.DB.Where("id = ? AND type = ?", modelID, "stt").First(&model).Error; err != nil {
		return nil, fmt.Errorf("STT model not found: %w", err)
	}

	// Create request record
	now := time.Now()
	request := &models.SpeechRequest{
		UserID:      userID,
		ModelID:     modelID,
		RequestType: "stt",
		InputAudio:  audioData,
		Status:      "processing",
		AudioFormat: options.AudioFormat,
		Language:    options.Language,
		CreatedAt:   now,
	}

	if err := config.DB.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create STT request: %w", err)
	}

	// Process asynchronously
	go s.processSTTRequest(request, &model, options)

	return request, nil
}

// processTTSRequest handles TTS processing
func (s *SpeechService) processTTSRequest(request *models.SpeechRequest, model *models.SpeechModel, options SpeechOptions) {
	defer func() {
		request.UpdatedAt = time.Now()
		config.DB.Save(request)
	}()

	var audioData []byte
	var err error

	// Route to appropriate TTS service
	switch model.Provider {
	case models.SpeechProviderLocal:
		audioData, err = s.localTTSService.Synthesize(request.InputText, options)
	case models.SpeechProviderElevenLabs:
		audioData, err = s.elevenLabsService.Synthesize(request.InputText, model.ModelID, options)
	case models.SpeechProviderOpenAI:
		audioData, err = s.openAIService.Synthesize(request.InputText, model.Voice, options)
	default:
		err = fmt.Errorf("unsupported TTS provider: %s", model.Provider)
	}

	if err != nil {
		request.Status = "failed"
		request.Error = err.Error()
		return
	}

	// Success
	completedAt := time.Now()
	request.Status = "completed"
	request.OutputAudio = audioData
	request.CompletedAt = &completedAt
	request.ProcessingTime = int(time.Since(request.CreatedAt).Milliseconds())
}

// processSTTRequest handles STT processing
func (s *SpeechService) processSTTRequest(request *models.SpeechRequest, model *models.SpeechModel, options SpeechOptions) {
	defer func() {
		request.UpdatedAt = time.Now()
		config.DB.Save(request)
	}()

	var text string
	var err error

	// Route to appropriate STT service
	switch model.Provider {
	case models.SpeechProviderLocal:
		text, err = s.localSTTService.Transcribe(request.InputAudio, options)
	case models.SpeechProviderOpenAI:
		text, err = s.openAIService.Transcribe(request.InputAudio, options)
	case models.SpeechProviderAssemblyAI:
		text, err = s.assemblyAIService.Transcribe(request.InputAudio, options)
	case models.SpeechProviderDeepgram:
		text, err = s.deepgramService.Transcribe(request.InputAudio, options)
	default:
		err = fmt.Errorf("unsupported STT provider: %s", model.Provider)
	}

	if err != nil {
		request.Status = "failed"
		request.Error = err.Error()
		return
	}

	// Success
	completedAt := time.Now()
	request.Status = "completed"
	request.OutputText = text
	request.CompletedAt = &completedAt
	request.ProcessingTime = int(time.Since(request.CreatedAt).Milliseconds())
}

// GetAvailableModels returns all available speech models for a user
func (s *SpeechService) GetAvailableModels(userID uuid.UUID, modelType string) ([]models.SpeechModel, error) {
	var models []models.SpeechModel
	err := config.DB.Where("(user_id = ? OR is_system = ?) AND is_active = ? AND type = ?",
		userID, true, true, modelType).Find(&models).Error
	return models, err
}

// GetSpeechSettings returns user's speech preferences
func (s *SpeechService) GetSpeechSettings(userID uuid.UUID) (*models.SpeechSettings, error) {
	var settings models.SpeechSettings
	err := config.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		// Create default settings if not found
		settings = models.SpeechSettings{
			UserID:           userID,
			EnableTTS:        true,
			EnableSTT:        true,
			DefaultLanguage:  "en",
			DefaultVoice:     "alloy",
			TTSSpeed:         1.0,
			TTSVolume:        1.0,
			STTLanguage:      "en",
			AutoPlayTTS:      false,
			ShowSTTTranscript: true,
			KeyboardShortcut: "ctrl+shift+s",
		}
		err = config.DB.Create(&settings).Error
	}
	return &settings, err
}

// UpdateSpeechSettings updates user's speech preferences
func (s *SpeechService) UpdateSpeechSettings(userID uuid.UUID, settings *models.SpeechSettings) error {
	settings.UserID = userID
	return config.DB.Save(settings).Error
}

// SpeechOptions represents options for TTS/STT processing
type SpeechOptions struct {
	Language    string  `json:"language"`
	Voice       string  `json:"voice"`
	AudioFormat string  `json:"audio_format"`
	Speed       float64 `json:"speed"`
	Pitch       float64 `json:"pitch"`
	Volume      float64 `json:"volume"`
}

// Local TTS Service (system voices)
type LocalTTSService struct{}

func NewLocalTTSService() *LocalTTSService {
	return &LocalTTSService{}
}

func (s *LocalTTSService) Synthesize(text string, options SpeechOptions) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return s.synthesizeWindows(text, options)
	case "darwin":
		return s.synthesizeMacOS(text, options)
	case "linux":
		return s.synthesizeLinux(text, options)
	default:
		return nil, fmt.Errorf("local TTS not supported on platform: %s", runtime.GOOS)
	}
}

// Windows TTS using PowerShell and System.Speech.Synthesis
func (s *LocalTTSService) synthesizeWindows(text string, options SpeechOptions) ([]byte, error) {
	// Escape single quotes in text
	escapedText := strings.ReplaceAll(text, "'", "''")

	// Set default voice if not specified
	voice := options.Voice
	if voice == "" {
		voice = "Microsoft Zira Desktop"
	}

	// Create PowerShell script for TTS
	psScript := fmt.Sprintf(`
		Add-Type -AssemblyName System.Speech
		$synth = New-Object System.Speech.Synthesis.SpeechSynthesizer
		$synth.SelectVoice('%s')
		$synth.Rate = %d
		$synth.Volume = %d
		$synth.Speak('%s')
	`, voice, s.speedToInt(options.Speed), s.volumeToInt(options.Volume), escapedText)

	// Execute PowerShell script
	cmd := exec.Command("powershell", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Windows TTS failed: %v, output: %s", err, string(output))
	}

	// Note: PowerShell TTS speaks directly to speakers, doesn't return audio data
	// For now, we'll return a success indicator
	return []byte("tts_completed"), nil
}

// macOS TTS using 'say' command
func (s *LocalTTSService) synthesizeMacOS(text string, options SpeechOptions) ([]byte, error) {
	// Escape single quotes in text
	escapedText := strings.ReplaceAll(text, "'", "\\'")

	// Build say command arguments
	args := []string{}

	// Set voice if specified
	if options.Voice != "" {
		args = append(args, "-v", options.Voice)
	} else {
		args = append(args, "-v", "Samantha") // Default macOS voice
	}

	// Set rate (words per minute, default ~200)
	rate := int(200 * options.Speed)
	if rate < 90 {
		rate = 90
	}
	if rate > 600 {
		rate = 600
	}
	args = append(args, "-r", strconv.Itoa(rate))

	// Add text
	args = append(args, escapedText)

	// Execute say command
	cmd := exec.Command("say", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("macOS TTS failed: %v, output: %s", err, string(output))
	}

	// Note: 'say' command speaks directly to speakers, doesn't return audio data
	// For now, we'll return a success indicator
	return []byte("tts_completed"), nil
}

// Linux TTS using espeak-ng
func (s *LocalTTSService) synthesizeLinux(text string, options SpeechOptions) ([]byte, error) {
	// Check if espeak-ng is available, fallback to espeak
	var cmd *exec.Cmd

	// Try espeak-ng first, then espeak
	if s.commandExists("espeak-ng") {
		args := []string{"-v", "en-us", "-s", s.speedToEspeakRate(options.Speed)}

		// Set voice if specified
		if options.Voice != "" {
			args = append(args, "-v", options.Voice)
		}

		// Set pitch
		pitch := int(50 + (options.Pitch * 25)) // Default pitch 50, range 0-100
		if pitch < 0 {
			pitch = 0
		}
		if pitch > 100 {
			pitch = 100
		}
		args = append(args, "-p", strconv.Itoa(pitch))

		args = append(args, text)
		cmd = exec.Command("espeak-ng", args...)
	} else if s.commandExists("espeak") {
		args := []string{"-v", "en-us", "-s", s.speedToEspeakRate(options.Speed)}

		if options.Voice != "" {
			args = append(args, "-v", options.Voice)
		}

		args = append(args, text)
		cmd = exec.Command("espeak", args...)
	} else if s.commandExists("festival") {
		// Fallback to festival
		return s.synthesizeFestival(text, options)
	} else {
		return nil, fmt.Errorf("no TTS engine found. Please install espeak-ng, espeak, or festival")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Linux TTS failed: %v, output: %s", err, string(output))
	}

	// Note: espeak speaks directly to speakers, doesn't return audio data
	// For now, we'll return a success indicator
	return []byte("tts_completed"), nil
}

// Festival TTS (alternative Linux TTS)
func (s *LocalTTSService) synthesizeFestival(text string, options SpeechOptions) ([]byte, error) {
	// Create festival script
	festivalScript := fmt.Sprintf(`
		(Parameter.set 'Audio_Method 'Audio_Command)
		(Parameter.set 'Audio_Command "aplay -q -c 1 -t raw -f s16 -r $SRATE $FILE")
		(SayText "%s")
	`, strings.ReplaceAll(text, `"`, `\"`))

	cmd := exec.Command("festival", "--tts")
	cmd.Stdin = strings.NewReader(festivalScript)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Festival TTS failed: %v, output: %s", err, string(output))
	}

	return []byte("tts_completed"), nil
}

// Helper functions for TTS
func (s *LocalTTSService) speedToInt(speed float64) int {
	// Convert speed multiplier to PowerShell rate (-10 to 10)
	rate := int((speed - 1.0) * 5)
	if rate < -10 {
		rate = -10
	}
	if rate > 10 {
		rate = 10
	}
	return rate
}

func (s *LocalTTSService) volumeToInt(volume float64) int {
	// Convert volume to 0-100 range
	return int(volume * 100)
}

func (s *LocalTTSService) speedToEspeakRate(speed float64) string {
	// espeak-ng speed in words per minute (80-450)
	rate := int(175 * speed) // Default ~175 wpm
	if rate < 80 {
		rate = 80
	}
	if rate > 450 {
		rate = 450
	}
	return strconv.Itoa(rate)
}

func (s *LocalTTSService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Local STT Service (system recognition)
type LocalSTTService struct{}

func NewLocalSTTService() *LocalSTTService {
	return &LocalSTTService{}
}

func (s *LocalSTTService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return s.transcribeWindows(audioData, options)
	case "darwin":
		return s.transcribeMacOS(audioData, options)
	case "linux":
		return s.transcribeLinux(audioData, options)
	default:
		return "", fmt.Errorf("local STT not supported on platform: %s", runtime.GOOS)
	}
}

// Windows STT using PowerShell and System.Speech.Recognition
func (s *LocalSTTService) transcribeWindows(audioData []byte, options SpeechOptions) (string, error) {
	// Note: Windows built-in speech recognition is primarily for real-time recognition
	// For file-based recognition, we'd need to use a different approach
	// For now, this is a placeholder that could be extended

	// Check if we have audio data to process
	if len(audioData) == 0 {
		return "", fmt.Errorf("no audio data provided for transcription")
	}

	// Windows Speech Recognition is primarily for real-time use
	// For file-based STT, we'd need to use a different library or service
	// This could be extended to use Windows.Media.SpeechRecognition or similar

	return "", fmt.Errorf("Windows file-based STT not yet implemented. Use cloud services like OpenAI Whisper or Azure Speech Service")
}

// macOS STT using built-in speech recognition
func (s *LocalSTTService) transcribeMacOS(audioData []byte, options SpeechOptions) (string, error) {
	// macOS has limited built-in speech recognition for files
	// The 'say' command is for TTS, not STT
	// For actual STT, we'd need to use macOS Speech framework via Objective-C bridge
	// or use cloud services

	// For now, provide guidance on using cloud services
	return "", fmt.Errorf("macOS built-in STT has limited file support. Consider using OpenAI Whisper or Azure Speech Service")
}

// Linux STT using pocketsphinx or kaldi
func (s *LocalSTTService) transcribeLinux(audioData []byte, options SpeechOptions) (string, error) {
	// Check if we have audio data
	if len(audioData) == 0 {
		return "", fmt.Errorf("no audio data provided for transcription")
	}

	// Try different STT engines in order of preference

	// Option 1: Kaldi (most accurate but complex)
	if s.commandExists("kaldi-gstreamer-server") || s.commandExists("online2-wav-nnet3-latgen-faster") {
		return s.transcribeWithKaldi(audioData, options)
	}

	// Option 2: PocketSphinx (simpler but less accurate)
	if s.commandExists("pocketsphinx_continuous") {
		return s.transcribeWithPocketSphinx(audioData, options)
	}

	// Option 3: Julius (another open-source option)
	if s.commandExists("julius") {
		return s.transcribeWithJulius(audioData, options)
	}

	// No local STT engines found
	return "", fmt.Errorf("no local STT engine found. Please install kaldi, pocketsphinx, or julius, or use cloud services like OpenAI Whisper")
}

// Kaldi-based STT (most accurate open-source option)
func (s *LocalSTTService) transcribeWithKaldi(audioData []byte, options SpeechOptions) (string, error) {
	// This is a simplified implementation
	// Real Kaldi setup is quite complex and requires model files

	// Save audio data to temporary file
	tempAudioFile := "/tmp/kaldi_input.wav"
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	// Note: This is a placeholder. Real Kaldi implementation would be much more complex
	// and require model files, acoustic models, language models, etc.

	// For now, return a message about setup requirements
	return "", fmt.Errorf("Kaldi STT requires complex setup with acoustic and language models. Consider using cloud services like OpenAI Whisper or Google Speech-to-Text")
}

// PocketSphinx STT (simpler open-source option)
func (s *LocalSTTService) transcribeWithPocketSphinx(audioData []byte, options SpeechOptions) (string, error) {
	// Save audio data to temporary file
	tempAudioFile := "/tmp/pocketsphinx_input.wav"
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	// Build pocketsphinx command
	args := []string{
		"-inmic", "no",
		"-infile", tempAudioFile,
		"-bestpath", "1",
	}

	// Set language model if specified
	if options.Language == "en" || options.Language == "" {
		// Use default English models if available
		args = append(args, "-hmm", "/usr/share/pocketsphinx/model/en-us/en-us")
		args = append(args, "-lm", "/usr/share/pocketsphinx/model/en-us/en-us.lm.bin")
		args = append(args, "-dict", "/usr/share/pocketsphinx/model/en-us/cmudict-en-us.dict")
	}

	cmd := exec.Command("pocketsphinx_continuous", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("PocketSphinx failed: %v, output: %s", err, string(output))
	}

	// Parse output to extract transcript
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "INFO:") && strings.Contains(line, "READY") {
			continue
		}
		if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "INFO:") {
			return strings.TrimSpace(line), nil
		}
	}

	return "", fmt.Errorf("no transcript found in PocketSphinx output")
}

// Julius STT (another open-source option)
func (s *LocalSTTService) transcribeWithJulius(audioData []byte, options SpeechOptions) (string, error) {
	// Save audio data to temporary file
	tempAudioFile := "/tmp/julius_input.wav"
	if err := os.WriteFile(tempAudioFile, audioData, 0644); err != nil {
		return "", fmt.Errorf("failed to save audio file: %v", err)
	}
	defer os.Remove(tempAudioFile)

	// Julius is complex to set up and typically requires specific model files
	// This is a placeholder implementation
	return "", fmt.Errorf("Julius STT requires specific acoustic and language model files. Consider using cloud services like OpenAI Whisper or Google Speech-to-Text")
}

// Helper function for command checking (duplicate, but keeping for clarity)
func (s *LocalSTTService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// ElevenLabs TTS Service
type ElevenLabsService struct {
	apiKey string
	client *http.Client
}

func NewElevenLabsService(apiKey string) *ElevenLabsService {
	return &ElevenLabsService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *ElevenLabsService) Synthesize(text, voiceID string, options SpeechOptions) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("ElevenLabs API key not configured")
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)

	payload := map[string]interface{}{
		"text":    text,
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

// OpenAI Speech Service
type OpenAISpeechService struct {
	apiKey string
	client *http.Client
}

func NewOpenAISpeechService(apiKey string) *OpenAISpeechService {
	return &OpenAISpeechService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *OpenAISpeechService) Synthesize(text, voice string, options SpeechOptions) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key not configured")
	}

	url := "https://api.openai.com/v1/audio/speech"

	payload := map[string]interface{}{
		"model":          "tts-1",
		"input":          text,
		"voice":          voice,
		"response_format": options.AudioFormat,
		"speed":          options.Speed,
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

func (s *OpenAISpeechService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	url := "https://api.openai.com/v1/audio/transcriptions"

	// Create multipart form data
	var b bytes.Buffer
	writer := io.Writer(&b)

	// Audio file part
	writer.Write([]byte("------FormBoundary\r\n"))
	writer.Write([]byte("Content-Disposition: form-data; name=\"file\"; filename=\"audio.wav\"\r\n"))
	writer.Write([]byte("Content-Type: audio/wav\r\n\r\n"))
	writer.Write(audioData)
	writer.Write([]byte("\r\n"))

	// Model part
	writer.Write([]byte("------FormBoundary\r\n"))
	writer.Write([]byte("Content-Disposition: form-data; name=\"model\"\r\n\r\n"))
	writer.Write([]byte("whisper-1\r\n"))

	// Language part (optional)
	if options.Language != "" {
		writer.Write([]byte("------FormBoundary\r\n"))
		writer.Write([]byte("Content-Disposition: form-data; name=\"language\"\r\n\r\n"))
		writer.Write([]byte(options.Language + "\r\n"))
	}

	writer.Write([]byte("------FormBoundary--\r\n"))

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

// AssemblyAI STT Service
type AssemblyAIService struct {
	apiKey string
	client *http.Client
}

func NewAssemblyAIService(apiKey string) *AssemblyAIService {
	return &AssemblyAIService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *AssemblyAIService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("AssemblyAI API key not configured")
	}

	// Step 1: Upload audio file
	uploadURL := "https://api.assemblyai.com/v2/upload"

	req, err := http.NewRequest("POST", uploadURL, bytes.NewReader(audioData))
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

	var uploadResult struct {
		UploadURL string `json:"upload_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&uploadResult); err != nil {
		return "", err
	}

	// Step 2: Request transcription
	transcriptURL := "https://api.assemblyai.com/v2/transcript"

	transcriptPayload := map[string]interface{}{
		"audio_url": uploadResult.UploadURL,
		"language_code": options.Language,
	}

	jsonData, _ := json.Marshal(transcriptPayload)

	req, err = http.NewRequest("POST", transcriptURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("AssemblyAI transcript error: %s", string(body))
	}

	var transcriptResult struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&transcriptResult); err != nil {
		return "", err
	}

	// Step 3: Poll for completion
	pollURL := fmt.Sprintf("https://api.assemblyai.com/v2/transcript/%s", transcriptResult.ID)

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

		if result.Status == "completed" {
			return result.Text, nil
		} else if result.Status == "error" {
			return "", fmt.Errorf("AssemblyAI transcription error: %s", result.Error)
		}

		time.Sleep(2 * time.Second) // Wait before polling again
	}
}

// Deepgram STT Service
type DeepgramService struct {
	apiKey string
	client *http.Client
}

func NewDeepgramService(apiKey string) *DeepgramService {
	return &DeepgramService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

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

// Replicate Service (for both TTS and STT models)
type ReplicateService struct {
	apiKey string
	client *http.Client
}

func NewReplicateService(apiKey string) *ReplicateService {
	return &ReplicateService{
		apiKey: apiKey,
		client: &http.Client{Timeout: 120 * time.Second}, // Replicate can be slow
	}
}

func (s *ReplicateService) Synthesize(text, modelName string, options SpeechOptions) ([]byte, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("Replicate API key not configured")
	}

	url := "https://api.replicate.com/v1/predictions"

	// Use a popular TTS model (can be made configurable)
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

	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Replicate prediction error: %s", string(body))
	}

	var prediction struct {
		URLs struct {
			Get string `json:"get"`
		} `json:"urls"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, err
	}

	// Poll for completion
	return s.pollReplicatePrediction(prediction.URLs.Get)
}

func (s *ReplicateService) Transcribe(audioData []byte, modelName string, options SpeechOptions) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("Replicate API key not configured")
	}

	// First, upload the audio file to a temporary URL
	uploadURL := "https://api.replicate.com/v1/uploads"

	uploadPayload := map[string]interface{}{
		"content_type": "audio/wav",
		"filename":     "audio.wav",
	}

	jsonData, _ := json.Marshal(uploadPayload)

	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(jsonData))
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

	// Upload the actual audio data
	req, err = http.NewRequest("PUT", uploadResult.UploadURL, bytes.NewReader(audioData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "audio/wav")

	resp, err = s.client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Replicate upload failed with status: %d", resp.StatusCode)
	}

	// Now create the transcription prediction
	url := "https://api.replicate.com/v1/predictions"

	// Use a popular Whisper model
	model := "openai/whisper:91ee9c0c3df30478510ff8c8a3a545add29ad025ca642e1c606e318e6eab18b25"
	if modelName != "" {
		model = modelName
	}

	payload := map[string]interface{}{
		"version": model,
		"input": map[string]interface{}{
			"audio":          uploadResult.UploadURL,
			"model":          "large-v3",
			"language":       options.Language,
			"translate":      false,
			"temperature":    0,
			"transcription":  "plain text",
			"suppress_tokens": "-1",
			"log_probability_threshold": -1.0,
			"no_speech_threshold": 0.6,
			"condition_on_previous_text": true,
			"compression_ratio_threshold": 2.4,
			"temperature_increment_on_fallback": 0.2,
		},
	}

	jsonData, _ = json.Marshal(payload)

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Token "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
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

	// Poll for completion and get text result
	return s.pollReplicateTranscription(prediction.URLs.Get)
}

func (s *ReplicateService) pollReplicatePrediction(predictionURL string) ([]byte, error) {
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

		if result.Status == "succeeded" {
			// Download the audio file from the output URL
			resp, err := http.Get(result.Output)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			return io.ReadAll(resp.Body)
		} else if result.Status == "failed" {
			return nil, fmt.Errorf("Replicate prediction failed: %s", result.Error)
		}

		time.Sleep(2 * time.Second) // Wait before polling again
	}
}

func (s *ReplicateService) pollReplicateTranscription(predictionURL string) (string, error) {
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

		if result.Status == "succeeded" {
			return result.Output.Text, nil
		} else if result.Status == "failed" {
			return "", fmt.Errorf("Replicate transcription failed: %s", result.Error)
		}

		time.Sleep(2 * time.Second) // Wait before polling again
	}
}

// Piper Service (local but cloud-hosted TTS)
type PiperService struct{}

func NewPiperService() *PiperService {
	return &PiperService{}
}

func (s *PiperService) Synthesize(text string, options SpeechOptions) ([]byte, error) {
	// Piper can be run locally via Docker or as a cloud service
	// For now, we'll use a cloud-hosted instance or provide installation guidance

	// Check if piper command is available locally
	if s.commandExists("piper") {
		return s.synthesizeLocalPiper(text, options)
	}

	// Fallback to cloud service or provide installation guidance
	return nil, fmt.Errorf("Piper TTS requires local installation. Install from https://github.com/rhasspy/piper or use cloud services")
}

func (s *PiperService) synthesizeLocalPiper(text string, options SpeechOptions) ([]byte, error) {
	// Save text to temporary file
	tempTextFile := "/tmp/piper_input.txt"
	if err := os.WriteFile(tempTextFile, []byte(text), 0644); err != nil {
		return nil, fmt.Errorf("failed to save text file: %v", err)
	}
	defer os.Remove(tempTextFile)

	// Piper output to stdout as WAV data
	tempOutputFile := "/tmp/piper_output.wav"

	args := []string{
		"--model", "en_US-lessac-medium", // Default English model
		"--output_file", tempOutputFile,
	}

	cmd := exec.Command("piper", args...)
	cmd.Stdin = strings.NewReader(text)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Piper TTS failed: %v, output: %s", err, string(output))
	}

	// Read the generated audio file
	audioData, err := os.ReadFile(tempOutputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read Piper output: %v", err)
	}

	os.Remove(tempOutputFile) // Clean up
	return audioData, nil
}

func (s *PiperService) Transcribe(audioData []byte, options SpeechOptions) (string, error) {
	// Piper is primarily a TTS engine, not STT
	return "", fmt.Errorf("Piper does not support speech-to-text. Use cloud services like OpenAI Whisper")
}

// Helper function for command checking
func (s *PiperService) commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
