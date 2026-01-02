package models

import (
	"time"

	"github.com/google/uuid"
)

// AIModel represents an AI model configuration
type AIModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` // null for system models
	Name        string    `gorm:"size:255;not null" json:"name"`
	Type        string    `gorm:"size:50;not null" json:"type"` // "local", "openrouter", "custom"
	Provider    string    `gorm:"size:100;not null" json:"provider"` // "ollama", "openrouter", etc.
	ModelID     string    `gorm:"size:255;not null" json:"model_id"` // Model identifier
	Capabilities []string `gorm:"type:text[]" json:"capabilities"` // ["ocr", "writing", "analysis"]
	APIKey      []byte    `gorm:"type:bytea" json:"api_key"` // Encrypted API key
	Endpoint    string    `gorm:"size:500" json:"endpoint"` // Custom endpoint URL
	Config      string    `gorm:"type:text" json:"config"` // JSON config for model parameters
	IsSystem    bool      `gorm:"default:false" json:"is_system"` // System-provided model
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AITask represents a task that can be performed by AI models
type AITask struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:100;not null" json:"category"` // "ocr", "writing", "analysis", "forms", "tasks"
	ModelID     uuid.UUID `gorm:"type:uuid;not null" json:"model_id"` // Assigned model
	Priority    int       `gorm:"default:1" json:"priority"` // 1=low, 2=medium, 3=high
	Config      string    `gorm:"type:text" json:"config"` // JSON config for task parameters
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AIRequest represents an AI processing request
type AIRequest struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null" json:"task_id"`
	ModelID   uuid.UUID `gorm:"type:uuid;not null" json:"model_id"`
	Input     string    `gorm:"type:text;not null" json:"input"` // Input data (text, base64 image, etc.)
	InputType string    `gorm:"size:50;not null" json:"input_type"` // "text", "image", "document"
	Output    string    `gorm:"type:text" json:"output"` // AI response
	Status    string    `gorm:"size:50;default:'pending'" json:"status"` // "pending", "processing", "completed", "failed"
	Error     string    `gorm:"type:text" json:"error"` // Error message if failed
	Tokens    int       `gorm:"default:0" json:"tokens"` // Token usage
	Cost      float64   `gorm:"default:0" json:"cost"` // API cost in USD
	StartedAt *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIUsage tracks usage statistics for billing/analytics
type AIUsage struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ModelID    uuid.UUID `gorm:"type:uuid;not null" json:"model_id"`
	Date       time.Time `gorm:"type:date;not null" json:"date"` // Usage date
	Tokens     int       `gorm:"default:0" json:"tokens"` // Total tokens used
	Requests   int       `gorm:"default:0" json:"requests"` // Total requests
	Cost       float64   `gorm:"default:0" json:"cost"` // Total cost in USD
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AIUpgradeCheck represents a manual upgrade evaluation
type AIUpgradeCheck struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CheckedAt      time.Time `json:"checked_at"`
	HardwareInfo   string    `gorm:"type:text" json:"hardware_info"` // JSON hardware info
	UpgradeOptions string    `gorm:"type:text" json:"upgrade_options"` // JSON upgrade options
	Status         string    `gorm:"size:50;default:'completed'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// UpgradeOption represents a potential model upgrade
type UpgradeOption struct {
	ID             string `json:"id"`
	CurrentModel   string `json:"current_model"`
	NewModel       string `json:"new_model"`
	Provider       string `json:"provider"` // "ollama" or "openrouter"
	SizeGB         float64 `json:"size_gb"`
	Capabilities   []string `json:"capabilities"`
	PerformanceGain string `json:"performance_gain"` // "better", "much_better", "similar"
	Compatibility  bool   `json:"compatibility"`
	Requirements   string `json:"requirements"` // JSON hardware requirements
	Reasoning      string `json:"reasoning"`
	RiskLevel      string `json:"risk_level"` // "low", "medium", "high"
}

// Default system models - tiered for different hardware capabilities
var DefaultAIModels = []AIModel{
	// High-end models (16GB+ RAM, GPU recommended)
	{
		Name:         "Qwen3 8B Vision",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen3:8b",
		Capabilities: []string{"ocr", "vision", "writing", "analysis", "forms"},
		IsSystem:     true,
		IsActive:     true,
	},
	{
		Name:         "Qwen3 14B Instruct",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen3:14b-instruct",
		Capabilities: []string{"writing", "analysis", "tasks", "forms", "reasoning"},
		IsSystem:     true,
		IsActive:     true,
	},
	{
		Name:         "Qwen3 30B Coder",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen3:30b-coder",
		Capabilities: []string{"analysis", "coding", "data", "automation"},
		IsSystem:     true,
		IsActive:     true,
	},
	// Standard models (8-16GB RAM)
	{
		Name:         "Qwen2.5 Vision 7B",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen2.5-vl:7b",
		Capabilities: []string{"ocr", "vision", "writing", "analysis", "forms"},
		IsSystem:     true,
		IsActive:     true,
	},
	{
		Name:         "Qwen2.5 7B Instruct",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen2.5:7b-instruct",
		Capabilities: []string{"writing", "analysis", "tasks", "forms", "reasoning"},
		IsSystem:     true,
		IsActive:     true,
	},
	{
		Name:         "Qwen2.5 Coder 7B",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen2.5-coder:7b-instruct",
		Capabilities: []string{"analysis", "coding", "data", "automation"},
		IsSystem:     true,
		IsActive:     true,
	},
	// Low-resource models (8GB RAM, CPU-only)
	{
		Name:         "Qwen2.5 Vision 2B",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen2.5-vl:2b",
		Capabilities: []string{"ocr", "vision", "writing", "forms"},
		IsSystem:     true,
		IsActive:     true,
	},
	{
		Name:         "Qwen2.5 3B Instruct",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen2.5:3b-instruct",
		Capabilities: []string{"writing", "analysis", "tasks", "forms"},
		IsSystem:     true,
		IsActive:     true,
	},
	{
		Name:         "Qwen2.5 Coder 3B",
		Type:         "local",
		Provider:     "ollama",
		ModelID:      "qwen2.5-coder:3b-instruct",
		Capabilities: []string{"analysis", "coding", "data"},
		IsSystem:     true,
		IsActive:     true,
	},
}

// Default task configurations
var DefaultAITasks = []AITask{
	{
		Name:        "Document OCR",
		Description: "Convert images and PDFs to editable text",
		Category:    "ocr",
		Priority:    2,
		IsActive:    true,
	},
	{
		Name:        "Writing Assistant",
		Description: "Grammar checking and content suggestions",
		Category:    "writing",
		Priority:    1,
		IsActive:    true,
	},
	{
		Name:        "Data Analysis",
		Description: "Spreadsheet insights and formula suggestions",
		Category:    "analysis",
		Priority:    2,
		IsActive:    true,
	},
	{
		Name:        "Form Intelligence",
		Description: "Smart form field detection and validation",
		Category:    "forms",
		Priority:    1,
		IsActive:    true,
	},
	{
		Name:        "Task Management",
		Description: "Task prioritization and scheduling suggestions",
		Category:    "tasks",
		Priority:    1,
		IsActive:    true,
	},
}

// TableName overrides for better naming
func (AIModel) TableName() string { return "ai_models" }
func (AITask) TableName() string { return "ai_tasks" }
func (AIRequest) TableName() string { return "ai_requests" }
func (AIUsage) TableName() string { return "ai_usage" }
func (AIUpgradeCheck) TableName() string { return "ai_upgrade_checks" }
