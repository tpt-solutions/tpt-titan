package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"

	"github.com/google/uuid"
)

// AIService handles AI model operations and orchestration
type AIService struct {
	config          *config.AIConfig
	hardwareService *HardwareService
	modelService    *ModelService
	ollamaService   *OllamaService
	openRouterService *OpenRouterService
}

// NewAIService creates a new AI service instance
func NewAIService(cfg *config.AIConfig) *AIService {
	hardwareService := NewHardwareService()
	ollamaService := NewOllamaService(cfg.OllamaHost, cfg.OllamaPort)
	openRouterService := NewOpenRouterService(cfg.OpenRouterKey)
	modelService := NewModelService(hardwareService, ollamaService, openRouterService)

	return &AIService{
		config:            cfg,
		hardwareService:   hardwareService,
		modelService:      modelService,
		ollamaService:     ollamaService,
		openRouterService: openRouterService,
	}
}

// InitializeSystemModels sets up the default AI models
func (s *AIService) InitializeSystemModels() error {
	// This would typically be called during system setup
	// For now, we'll handle this in the database migration
	log.Println("AI system initialized with default models")
	return nil
}

// ProcessRequest processes an AI request
func (s *AIService) ProcessRequest(userID uuid.UUID, taskID uuid.UUID, modelID uuid.UUID, input string, inputType string) (*models.AIRequest, error) {
	// Create request record
	now := time.Now()
	request := &models.AIRequest{
		UserID:    userID,
		TaskID:    taskID,
		ModelID:   modelID,
		Input:     input,
		InputType: inputType,
		Status:    "processing",
		StartedAt: &now,
	}

	// Save to database
	if err := config.DB.Create(request).Error; err != nil {
		return nil, fmt.Errorf("failed to create request record: %w", err)
	}

	// Process the request asynchronously
	go s.processAIRequest(request)

	return request, nil
}

// processAIRequest handles the actual AI processing
func (s *AIService) processAIRequest(request *models.AIRequest) {
	defer func() {
		request.UpdatedAt = time.Now()
		config.DB.Save(request)
	}()

	// Get model information
	var model models.AIModel
	if err := config.DB.First(&model, "id = ?", request.ModelID).Error; err != nil {
		request.Status = "failed"
		request.Error = "Model not found"
		return
	}

	var response string
	var err error

	// Route to appropriate AI service
	switch model.Provider {
	case "ollama":
		if !s.config.EnableLocalAI {
			request.Status = "failed"
			request.Error = "Local AI is disabled"
			return
		}

		response, err = s.ollamaService.GenerateResponse(model.ModelID, request.Input)

	case "openrouter":
		if !s.config.EnableOnlineAI {
			request.Status = "failed"
			request.Error = "Online AI is disabled"
			return
		}

		response, err = s.openRouterService.GenerateResponse(model.ModelID, request.Input)

	default:
		request.Status = "failed"
		request.Error = "Unsupported AI provider"
		return
	}

	if err != nil {
		request.Status = "failed"
		request.Error = err.Error()
		return
	}

	// Success
	completedAt := time.Now()
	request.Status = "completed"
	request.Output = response
	request.CompletedAt = &completedAt

	// Estimate token usage (rough approximation)
	request.Tokens = len(strings.Fields(request.Input)) + len(strings.Fields(response))

	// Update usage statistics
	s.updateUsageStatistics(request.UserID, request.ModelID, request.Tokens, 0.0)
}

// updateUsageStatistics updates the usage tracking
func (s *AIService) updateUsageStatistics(userID, modelID uuid.UUID, tokens int, cost float64) {
	today := time.Now().Truncate(24 * time.Hour)

	usage := &models.AIUsage{
		UserID:   userID,
		ModelID:  modelID,
		Date:     today,
		Tokens:   tokens,
		Requests: 1,
		Cost:     cost,
	}

	// Upsert usage record
	config.DB.Where(models.AIUsage{UserID: userID, ModelID: modelID, Date: today}).
		Assign(models.AIUsage{
			Tokens:   usage.Tokens,
			Requests: usage.Requests,
			Cost:     usage.Cost,
		}).
		FirstOrCreate(usage)
}

// GetAvailableModels returns all available models for a user
func (s *AIService) GetAvailableModels(userID uuid.UUID) ([]models.AIModel, error) {
	var models []models.AIModel
	err := config.DB.Where("(user_id = ? OR is_system = ?) AND is_active = ?", userID, true, true).Find(&models).Error
	return models, err
}

// CheckForUpgrades performs a manual upgrade check for AI models
func (s *AIService) CheckForUpgrades(userID uuid.UUID) (*models.AIUpgradeCheck, error) {
	// Detect hardware
	hardware, err := s.DetectHardware()
	if err != nil {
		return nil, fmt.Errorf("failed to detect hardware: %w", err)
	}

	// Get current user models
	currentModels, err := s.GetAvailableModels(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current models: %w", err)
	}

	// Collect upgrade options
	var upgradeOptions []models.UpgradeOption

	// Check Ollama models
	if s.config.EnableLocalAI {
		ollamaOptions, err := s.checkOllamaUpgrades(currentModels, hardware)
		if err != nil {
			log.Printf("Warning: Failed to check Ollama upgrades: %v", err)
		} else {
			upgradeOptions = append(upgradeOptions, ollamaOptions...)
		}
	}

	// Check OpenRouter models if API key is configured
	if s.config.EnableOnlineAI && s.config.OpenRouterKey != "" {
		openRouterOptions, err := s.checkOpenRouterUpgrades(currentModels, hardware)
		if err != nil {
			log.Printf("Warning: Failed to check OpenRouter upgrades: %v", err)
		} else {
			upgradeOptions = append(upgradeOptions, openRouterOptions...)
		}
	}

	// Create upgrade check record
	check := &models.AIUpgradeCheck{
		UserID:    userID,
		CheckedAt: time.Now(),
		Status:    "completed",
	}

	// Serialize hardware info and upgrade options
	hardwareJSON, _ := json.Marshal(hardware)
	check.HardwareInfo = string(hardwareJSON)

	optionsJSON, _ := json.Marshal(upgradeOptions)
	check.UpgradeOptions = string(optionsJSON)

	// Save to database
	if err := config.DB.Create(check).Error; err != nil {
		return nil, fmt.Errorf("failed to save upgrade check: %w", err)
	}

	return check, nil
}

// checkOllamaUpgrades checks for Ollama model upgrades
func (s *AIService) checkOllamaUpgrades(currentModels []models.AIModel, hardware *HardwareInfo) ([]models.UpgradeOption, error) {
	availableModels, err := s.ollamaService.ListModels()
	if err != nil {
		return nil, fmt.Errorf("failed to list Ollama models: %w", err)
	}

	var upgrades []models.UpgradeOption

	// Compare current models with available models
	for _, current := range currentModels {
		if current.Provider != "ollama" {
			continue
		}

		// Find potential upgrades
		for _, available := range availableModels {
			if s.isUpgradeCandidate(current.ModelID, available.Name, hardware) {
				sizeGB := parseSizeToGB(available.Size)
				upgrade := models.UpgradeOption{
					ID:             fmt.Sprintf("ollama-%s-%s", current.ModelID, available.Name),
					CurrentModel:   current.ModelID,
					NewModel:       available.Name,
					Provider:       "ollama",
					SizeGB:         sizeGB,
					Capabilities:   available.Capabilities,
					PerformanceGain: s.assessPerformanceGain(current.ModelID, available.Name),
					Compatibility:  s.checkCompatibility(sizeGB, hardware),
					Reasoning:      s.generateUpgradeReasoning(current.ModelID, available.Name),
					RiskLevel:      "low", // Ollama upgrades are generally safe
				}

				// Serialize requirements
				reqJSON, _ := json.Marshal(hardware)
				upgrade.Requirements = string(reqJSON)

				upgrades = append(upgrades, upgrade)
			}
		}
	}

	return upgrades, nil
}

// checkOpenRouterUpgrades checks for OpenRouter model upgrades
func (s *AIService) checkOpenRouterUpgrades(currentModels []models.AIModel, hardware *HardwareInfo) ([]models.UpgradeOption, error) {
	// This would query OpenRouter API for available models
	// For now, return a placeholder - in production, this would make API calls
	var upgrades []models.UpgradeOption

	// Example: Check if user has access to GPT-4, Claude 3, etc.
	// This is simplified - real implementation would query OpenRouter API

	for _, current := range currentModels {
		if current.Provider == "openrouter" {
			// Add potential OpenRouter upgrades here
			// This would compare current model against available OpenRouter models
		}
	}

	return upgrades, nil
}

// Helper functions for upgrade evaluation

func (s *AIService) isUpgradeCandidate(currentModelID, newModelName string, hardware *HardwareInfo) bool {
	// Simple logic: check if new model is different and compatible
	if currentModelID == newModelName {
		return false
	}

	// Check basic compatibility
	return s.checkCompatibility(1.0, hardware) // Placeholder size check
}

func (s *AIService) assessPerformanceGain(currentModel, newModel string) string {
	// Simplified performance assessment
	// In production, this could use benchmarks, user feedback, etc.

	// Check model sizes as proxy for performance
	currentSize := extractModelSize(currentModel)
	newSize := extractModelSize(newModel)

	if newSize > currentSize*1.5 {
		return "much_better"
	} else if newSize > currentSize {
		return "better"
	}
	return "similar"
}

func (s *AIService) checkCompatibility(modelSizeGB float64, hardware *HardwareInfo) bool {
	// Estimate RAM requirements (rough approximation)
	estimatedRAM := modelSizeGB * 1.5 // Model size * 1.5 for inference overhead

	// Add buffer for system
	requiredRAM := estimatedRAM + 2.0 // 2GB system buffer

	return float64(hardware.RAMGB) >= requiredRAM
}

func (s *AIService) generateUpgradeReasoning(currentModel, newModel string) string {
	currentSize := extractModelSize(currentModel)
	newSize := extractModelSize(newModel)

	if newSize > currentSize {
		return fmt.Sprintf("Upgrading from %.1fB to %.1fB model for better performance and capabilities", currentSize, newSize)
	}
	return fmt.Sprintf("Model update to %s with improved training data", newModel)
}



// DetectHardware attempts to detect system hardware capabilities
func (s *AIService) DetectHardware() (*HardwareInfo, error) {
	// This is a basic implementation - in production, you'd use system calls
	// For now, we'll return conservative defaults
	// TODO: Implement actual hardware detection using runtime.NumCPU(), syscall, etc.

	info := &HardwareInfo{
		RAMGB:     8,  // Conservative default - most users have at least 8GB
		HasGPU:    false, // Conservative - assume no GPU unless detected
		CPUCores:  4,  // Most systems have at least 4 cores
		CPUSpeed:  2000, // 2GHz base speed
		DiskSpace: 50, // 50GB free space minimum
	}

	// Try to detect actual RAM (simplified)
	// In production, use: runtime.MemStats, syscall.Sysinfo, etc.
	// For now, return defaults that work for most systems

	return info, nil
}

// GetRecommendedModels returns model recommendations based on hardware
func (s *AIService) GetRecommendedModels(hardware *HardwareInfo) *RecommendedModels {
	recommendations := &RecommendedModels{
		AllModels: make(map[string]string),
	}

	// Determine hardware tier
	if hardware.RAMGB >= 32 && hardware.HasGPU {
		recommendations.HardwareTier = "high-end"
		recommendations.OCR = "qwen3:8b"
		recommendations.Writing = "qwen3:14b-instruct"
		recommendations.Analysis = "qwen3:30b-coder"
		recommendations.Reasoning = "qwen3:30b-coder"
	} else if hardware.RAMGB >= 16 {
		recommendations.HardwareTier = "standard"
		recommendations.OCR = "qwen2.5-vl:7b"
		recommendations.Writing = "qwen2.5:7b-instruct"
		recommendations.Analysis = "qwen2.5-coder:7b-instruct"
		recommendations.Reasoning = "qwen2.5:7b-instruct"
	} else if hardware.RAMGB >= 8 {
		recommendations.HardwareTier = "low-resource"
		recommendations.OCR = "qwen2.5-vl:2b"
		recommendations.Writing = "qwen2.5:3b-instruct"
		recommendations.Analysis = "qwen2.5-coder:3b-instruct"
		recommendations.Reasoning = "qwen2.5:3b-instruct"
	} else {
		recommendations.HardwareTier = "minimal"
		recommendations.OCR = "qwen2.5:1.5b"
		recommendations.Writing = "qwen2.5:1.5b"
		recommendations.Analysis = "qwen2.5-coder:1.5b"
		recommendations.Reasoning = "qwen2.5:1.5b"
	}

	// Populate all models map
	recommendations.AllModels["ocr"] = recommendations.OCR
	recommendations.AllModels["writing"] = recommendations.Writing
	recommendations.AllModels["analysis"] = recommendations.Analysis
	recommendations.AllModels["forms"] = recommendations.Writing
	recommendations.AllModels["tasks"] = recommendations.Reasoning

	return recommendations
}

// SetupRecommendedModels creates default tasks with recommended models for a user
func (s *AIService) SetupRecommendedModels(userID uuid.UUID, recommendations *RecommendedModels) error {
	// Get or create tasks for each category
	categories := []string{"ocr", "writing", "analysis", "forms", "tasks"}

	for _, category := range categories {
		modelIDStr := recommendations.AllModels[category]
		if modelIDStr == "" {
			continue
		}

		// Find the system model by model_id
		var systemModel models.AIModel
		if err := config.DB.Where("model_id = ? AND is_system = ?", modelIDStr, true).First(&systemModel).Error; err != nil {
			log.Printf("System model %s not found, skipping task setup", modelIDStr)
			continue
		}

		// Check if task already exists for this user and category
		var existingTask models.AITask
		err := config.DB.Where("user_id = ? AND category = ?", userID, category).First(&existingTask).Error
		if err == nil {
			// Update existing task
			existingTask.ModelID = systemModel.ID
			config.DB.Save(&existingTask)
		} else {
			// Create new task
			task := models.AITask{
				UserID:   userID,
				Name:     strings.Title(category) + " Assistant",
				Category: category,
				ModelID:  systemModel.ID,
				Priority: 1,
				IsActive: true,
			}

			// Set descriptions based on category
			switch category {
			case "ocr":
				task.Description = "Convert images and PDFs to editable text"
				task.Priority = 2
			case "writing":
				task.Description = "Grammar checking and content suggestions"
			case "analysis":
				task.Description = "Data analysis and formula suggestions"
				task.Priority = 2
			case "forms":
				task.Description = "Smart form field detection and validation"
			case "tasks":
				task.Description = "Task prioritization and scheduling suggestions"
			}

			if err := config.DB.Create(&task).Error; err != nil {
				log.Printf("Failed to create task for category %s: %v", category, err)
			}
		}
	}

	return nil
}

// Helper functions
