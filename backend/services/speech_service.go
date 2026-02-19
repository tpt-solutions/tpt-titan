// backend/services/speech_service.go
package services

import (
	"fmt"
	"log"
	"time"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"github.com/google/uuid"
)

// SpeechOptions represents options for TTS/STT processing
type SpeechOptions struct {
	Language    string  `json:"language"`
	Voice       string  `json:"voice"`
	AudioFormat string  `json:"audio_format"`
	Speed       float64 `json:"speed"`
	Pitch       float64 `json:"pitch"`
	Volume      float64 `json:"volume"`
}

// SpeechService handles TTS/STT operations and orchestration
type SpeechService struct {
	config            *config.SpeechConfig
	localTTSService   *LocalTTSService
	localSTTService   *LocalSTTService
	elevenLabsService *ElevenLabsService
	openAIService     *OpenAISpeechService
	replicateService  *ReplicateService
	piperService      *PiperService
	assemblyAIService *AssemblyAIService
	deepgramService   *DeepgramService
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
	var model models.SpeechModel
	if err := config.DB.Where("id = ? AND type = ?", modelID, "tts").First(&model).Error; err != nil {
		return nil, fmt.Errorf("TTS model not found: %w", err)
	}

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

	go s.processTTSRequest(request, &model, options)

	return request, nil
}

// SpeechToText converts speech to text
func (s *SpeechService) SpeechToText(userID uuid.UUID, audioData []byte, modelID uuid.UUID, options SpeechOptions) (*models.SpeechRequest, error) {
	var model models.SpeechModel
	if err := config.DB.Where("id = ? AND type = ?", modelID, "stt").First(&model).Error; err != nil {
		return nil, fmt.Errorf("STT model not found: %w", err)
	}

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

	go s.processSTTRequest(request, &model, options)

	return request, nil
}

// processTTSRequest handles TTS processing asynchronously
func (s *SpeechService) processTTSRequest(request *models.SpeechRequest, model *models.SpeechModel, options SpeechOptions) {
	defer func() {
		config.DB.Save(request)
	}()

	var audioData []byte
	var err error

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

	completedAt := time.Now()
	request.Status = "completed"
	request.OutputAudio = audioData
	request.CompletedAt = &completedAt
	request.ProcessingTime = int(time.Since(request.CreatedAt).Milliseconds())
}

// processSTTRequest handles STT processing asynchronously
func (s *SpeechService) processSTTRequest(request *models.SpeechRequest, model *models.SpeechModel, options SpeechOptions) {
	defer func() {
		config.DB.Save(request)
	}()

	var text string
	var err error

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

	completedAt := time.Now()
	request.Status = "completed"
	request.OutputText = text
	request.CompletedAt = &completedAt
	request.ProcessingTime = int(time.Since(request.CreatedAt).Milliseconds())
}

// GetAvailableModels returns all available speech models for a user
func (s *SpeechService) GetAvailableModels(userID uuid.UUID, modelType string) ([]models.SpeechModel, error) {
	var speechModels []models.SpeechModel
	err := config.DB.Where("(user_id = ? OR is_system = ?) AND is_active = ? AND type = ?",
		userID, true, true, modelType).Find(&speechModels).Error
	return speechModels, err
}

// GetSpeechSettings returns user's speech preferences
func (s *SpeechService) GetSpeechSettings(userID uuid.UUID) (*models.SpeechSettings, error) {
	var settings models.SpeechSettings
	err := config.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		settings = models.SpeechSettings{
			UserID:            userID,
			EnableTTS:         true,
			EnableSTT:         true,
			DefaultLanguage:   "en",
			DefaultVoice:      "alloy",
			TTSSpeed:          1.0,
			TTSVolume:         1.0,
			STTLanguage:       "en",
			AutoPlayTTS:       false,
			ShowSTTTranscript: true,
			KeyboardShortcut:  "ctrl+shift+s",
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
