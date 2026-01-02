package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
	"time"

	"github.com/google/uuid"
	"github.com/ollama/ollama/api"
)

// HardwareInfo represents detected hardware capabilities
type HardwareInfo struct {
	RAMGB     int  `json:"ram_gb"`
	HasGPU    bool `json:"has_gpu"`
	CPUCores  int  `json:"cpu_cores"`
	CPUSpeed  int  `json:"cpu_speed_mhz"` // CPU speed in MHz
	DiskSpace int  `json:"disk_space_gb"` // Available disk space in GB
}

// RecommendedModels represents model recommendations for different hardware tiers
type RecommendedModels struct {
	HardwareTier string            `json:"hardware_tier"`
	OCR          string            `json:"ocr_model"`
	Writing      string            `json:"writing_model"`
	Analysis     string            `json:"analysis_model"`
	Reasoning    string            `json:"reasoning"`
	AllModels    map[string]string `json:"all_models"`
}

// AIService handles AI model operations
type AIService struct {
	config *config.AIConfig
}

// NewAIService creates a new AI service instance
func NewAIService(cfg *config.AIConfig) *AIService {
	return &AIService{
		config: cfg,
	}
}

// AIModelInfo represents information about an AI model
type AIModelInfo struct {
	Name         string   `json:"name"`
	Size         string   `json:"size"`
	ModifiedAt   string   `json:"modified_at"`
	Digest       string   `json:"digest"`
	Capabilities []string `json:"capabilities"`
}

// OllamaClient manages Ollama API interactions
type OllamaClient struct {
	baseURL string
	client  *http.Client
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(host, port string) *OllamaClient {
	return &OllamaClient{
		baseURL: fmt.Sprintf("http://%s:%s", host, port),
		client: &http.Client{
			Timeout: 5 * time.Minute, // Long timeout for AI processing
		},
	}
}

// ListModels gets available models from Ollama
func (c *OllamaClient) ListModels() ([]AIModelInfo, error) {
	resp, err := c.client.Get(c.baseURL + "/api/tags")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ollama API error: %s", resp.Status)
	}

	var result struct {
		Models []struct {
			Name       string `json:"name"`
			Size       int64  `json:"size"`
			ModifiedAt string `json:"modified_at"`
			Digest     string `json:"digest"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Ollama response: %w", err)
	}

	var models []AIModelInfo
	for _, model := range result.Models {
		models = append(models, AIModelInfo{
			Name:         model.Name,
			Size:         formatBytes(model.Size),
			ModifiedAt:   model.ModifiedAt,
			Digest:       model.Digest,
			Capabilities: inferCapabilities(model.Name),
		})
	}

	return models, nil
}

// PullModel downloads a model from Ollama
func (c *OllamaClient) PullModel(modelName string) error {
	req := map[string]string{"name": modelName}
	jsonData, _ := json.Marshal(req)

	resp, err := c.client.Post(c.baseURL+"/api/pull", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to pull model %s: %s", modelName, string(body))
	}

	return nil
}

// GenerateResponse generates a response using Ollama
func (c *OllamaClient) GenerateResponse(modelName, prompt string) (string, error) {
	stream := false
	req := api.GenerateRequest{
		Model:  modelName,
		Prompt: prompt,
		Stream: &stream,
	}

	jsonData, _ := json.Marshal(req)

	resp, err := c.client.Post(c.baseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("generation failed: %s", string(body))
	}

	var result api.GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Response, nil
}

// OpenRouterClient manages OpenRouter API interactions
type OpenRouterClient struct {
	apiKey  string
	client  *http.Client
	baseURL string
}

// NewOpenRouterClient creates a new OpenRouter client
func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 5 * time.Minute},
		baseURL: "https://openrouter.ai/api/v1",
	}
}

// GenerateResponse generates a response using OpenRouter
func (c *OpenRouterClient) GenerateResponse(modelName, prompt string) (string, error) {
	req := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	jsonData, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
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

		client := NewOllamaClient(s.config.OllamaHost, s.config.OllamaPort)
		response, err = client.GenerateResponse(model.ModelID, request.Input)

	case "openrouter":
		if !s.config.EnableOnlineAI {
			request.Status = "failed"
			request.Error = "Online AI is disabled"
			return
		}

		client := NewOpenRouterClient(s.config.OpenRouterKey)
		response, err = client.GenerateResponse(model.ModelID, request.Input)

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
	client := NewOllamaClient(s.config.OllamaHost, s.config.OllamaPort)
	availableModels, err := client.ListModels()
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

// extractModelSize extracts model size from model name (simplified)
func extractModelSize(modelName string) float64 {
	// Extract size from model name like "qwen2.5:7b" or "llama3.2:3b"
	// This is a simplified implementation
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

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
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

func inferCapabilities(modelName string) []string {
	name := strings.ToLower(modelName)
	capabilities := []string{}

	// Infer capabilities based on model name
	if strings.Contains(name, "vision") || strings.Contains(name, "llava") {
		capabilities = append(capabilities, "vision", "ocr")
	}

	if strings.Contains(name, "coder") || strings.Contains(name, "code") {
		capabilities = append(capabilities, "coding", "analysis")
	}

	if strings.Contains(name, "llama") || strings.Contains(name, "mistral") {
		capabilities = append(capabilities, "writing", "analysis", "general")
	}

	if strings.Contains(name, "phi") || strings.Contains(name, "gemma") {
		capabilities = append(capabilities, "writing", "tasks", "forms")
	}

	// Default capabilities if none inferred
	if len(capabilities) == 0 {
		capabilities = []string{"general"}
	}

	return capabilities
}
