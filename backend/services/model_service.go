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

// RecommendedModels represents model recommendations for different hardware tiers
type RecommendedModels struct {
	HardwareTier string            `json:"hardware_tier"`
	OCR          string            `json:"ocr_model"`
	Writing      string            `json:"writing_model"`
	Analysis     string            `json:"analysis_model"`
	Reasoning    string            `json:"reasoning"`
	AllModels    map[string]string `json:"all_models"`
}

// ModelService handles AI model management and recommendations
type ModelService struct {
	hardwareService *HardwareService
	ollamaService   *OllamaService
	openRouterService *OpenRouterService
}

// NewModelService creates a new model service
func NewModelService(hardwareService *HardwareService, ollamaService *OllamaService, openRouterService *OpenRouterService) *ModelService {
	return &ModelService{
		hardwareService: hardwareService,
		ollamaService:   ollamaService,
		openRouterService: openRouterService,
	}
}

// GetAvailableModels returns all available models for a user
func (s *ModelService) GetAvailableModels(userID uuid.UUID) ([]models.AIModel, error) {
	var models []models.AIModel
	err := config.DB.Where("(user_id = ? OR is_system = ?) AND is_active = ?", userID, true, true).Find(&models).Error
	return models, err
}

// GetRecommendedModels returns model recommendations based on hardware
func (s *ModelService) GetRecommendedModels(hardware *HardwareInfo) *RecommendedModels {
	recommendations := &RecommendedModels{
		AllModels: make(map[string]string),
	}

	// Determine hardware tier
	tier := s.hardwareService.GetHardwareTier(hardware)
	recommendations.HardwareTier = tier

	// Set model recommendations based on hardware tier
	switch tier {
	case "high-end":
		recommendations.OCR = "qwen3:8b"
		recommendations.Writing = "qwen3:14b-instruct"
		recommendations.Analysis = "qwen3:30b-coder"
		recommendations.Reasoning = "qwen3:30b-coder"
	case "standard":
		recommendations.OCR = "qwen2.5-vl:7b"
		recommendations.Writing = "qwen2.5:7b-instruct"
		recommendations.Analysis = "qwen2.5-coder:7b-instruct"
		recommendations.Reasoning = "qwen2.5:7b-instruct"
	case "low-resource":
		recommendations.OCR = "qwen2.5-vl:2b"
		recommendations.Writing = "qwen2.5:3b-instruct"
		recommendations.Analysis = "qwen2.5-coder:3b-instruct"
		recommendations.Reasoning = "qwen2.5:3b-instruct"
	case "minimal":
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
func (s *ModelService) SetupRecommendedModels(userID uuid.UUID, recommendations *RecommendedModels) error {
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

// CheckForUpgrades performs a manual upgrade check for AI models
func (s *ModelService) CheckForUpgrades(userID uuid.UUID) (*models.AIUpgradeCheck, error) {
	// Detect hardware
	hardware, err := s.hardwareService.DetectHardware()
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

	// Check for upgrades from available models
	upgradeOptions, err = s.checkAvailableUpgrades(currentModels, hardware)
	if err != nil {
		log.Printf("Warning: Failed to check for upgrades: %v", err)
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

// checkAvailableUpgrades checks for available model upgrades
func (s *ModelService) checkAvailableUpgrades(currentModels []models.AIModel, hardware *HardwareInfo) ([]models.UpgradeOption, error) {
	var upgrades []models.UpgradeOption

	// Get available Ollama models if service is available
	if s.ollamaService != nil {
		ollamaModels, err := s.ollamaService.ListModels()
		if err == nil {
			for _, current := range currentModels {
				if current.Provider == "ollama" {
					for _, available := range ollamaModels {
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
								Compatibility:  s.hardwareService.CheckCompatibility(sizeGB, hardware),
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
			}
		}
	}

	return upgrades, nil
}

// Helper functions for upgrade evaluation

func (s *ModelService) isUpgradeCandidate(currentModelID, newModelName string, hardware *HardwareInfo) bool {
	// Simple logic: check if new model is different and compatible
	if currentModelID == newModelName {
		return false
	}

	// Check basic compatibility (placeholder - would need size info)
	return true
}

func (s *ModelService) assessPerformanceGain(currentModel, newModel string) string {
	// Simplified performance assessment
	currentSize := extractModelSize(currentModel)
	newSize := extractModelSize(newModel)

	if newSize > currentSize*1.5 {
		return "much_better"
	} else if newSize > currentSize {
		return "better"
	}
	return "similar"
}

func (s *ModelService) generateUpgradeReasoning(currentModel, newModel string) string {
	currentSize := extractModelSize(currentModel)
	newSize := extractModelSize(newModel)

	if newSize > currentSize {
		return fmt.Sprintf("Upgrading from %.1fB to %.1fB model for better performance and capabilities", currentSize, newSize)
	}
	return fmt.Sprintf("Model update to %s with improved training data", newModel)
}

// extractModelSize extracts model size from model name (simplified)
func extractModelSize(modelName string) float64 {
	// Extract size from model name like "qwen2.5:7b" or "llama3.2:3b"
	sizes := map[string]float64{
		"1b": 1, "2b": 2, "3b": 3, "7b": 7, "8b": 8, "14b": 14, "30b": 30,
	}

	for suffix, size := range sizes {
		if strings.Contains(strings.ToLower(modelName), suffix) {
			return size
		}
	}

	return 1.0 // Default fallback
}

// parseSizeToGB parses a size string (e.g., "4.1 GB") to GB as float64
func parseSizeToGB(sizeStr string) float64 {
	// Parse strings like "4.1 GB", "500 MB", etc.
	var value float64
	var unit string
	fmt.Sscanf(sizeStr, "%f %s", &value, &unit)

	// Convert to GB
	switch unit {
	case "KB":
		return value / (1024 * 1024)
	case "MB":
		return value / 1024
	case "GB":
		return value
	case "TB":
		return value * 1024
	default:
		return 0
	}
}
