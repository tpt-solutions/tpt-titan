package models

import (
	"time"

	"github.com/google/uuid"
)

// AIModel represents a configured AI model
type AIModel struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Type         string    `gorm:"size:50;not null" json:"type"`      // "local", "cloud"
	Provider     string    `gorm:"size:100;not null" json:"provider"` // "ollama", "openrouter"
	ModelID      string    `gorm:"size:255;not null" json:"model_id"`
	Capabilities []string  `gorm:"serializer:json" json:"capabilities"`
	IsSystem     bool      `gorm:"default:false" json:"is_system"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	APIKey       []byte    `json:"-"` // Encrypted
	Endpoint     string    `gorm:"size:500" json:"endpoint,omitempty"`
	Config       string    `gorm:"type:jsonb" json:"config,omitempty"`
	Priority     int       `gorm:"default:0" json:"priority"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AITask represents a configured AI task/use-case for a user
type AITask struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid" json:"user_id,omitempty"`
	ModelID     uuid.UUID `gorm:"type:uuid" json:"model_id,omitempty"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:100;not null" json:"category"`
	Priority    int       `gorm:"default:1" json:"priority"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Config      string    `gorm:"type:jsonb" json:"config,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AIRequest represents a single AI inference request
type AIRequest struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ModelID     uuid.UUID  `gorm:"type:uuid;not null" json:"model_id"`
	TaskID      uuid.UUID  `gorm:"type:uuid" json:"task_id,omitempty"`
	InputType   string     `gorm:"size:50;default:'text'" json:"input_type"` // "text", "image", "audio"
	Input       string     `gorm:"type:text" json:"input"`
	Output      string     `gorm:"type:text" json:"output"`
	Status      string     `gorm:"size:20;default:'pending'" json:"status"` // "pending","processing","completed","failed"
	Error       string     `gorm:"type:text" json:"error,omitempty"`
	Tokens      int        `gorm:"default:0" json:"tokens"`
	Cost        float64    `gorm:"default:0" json:"cost"`
	Duration    int        `gorm:"default:0" json:"duration_ms"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// DocumentAnalysis represents AI-processed document data
type DocumentAnalysis struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	DocumentID  uuid.UUID `gorm:"type:uuid;not null" json:"document_id"` // Reference to documents table
	FileName    string    `gorm:"size:255;not null" json:"file_name"`
	FileType    string    `gorm:"size:50;not null" json:"file_type"` // "pdf", "image", etc.
	TextContent string    `gorm:"type:text" json:"text_content"`     // Full extracted text
	Confidence  float64   `gorm:"default:0" json:"confidence"`       // Overall confidence score (0-1)

	// Structured data extraction
	Fields   string `gorm:"type:jsonb" json:"fields"`   // JSON object with extracted fields
	Tables   string `gorm:"type:jsonb" json:"tables"`   // JSON array of extracted tables
	Entities string `gorm:"type:jsonb" json:"entities"` // JSON object with named entities

	// Processing metadata
	Pages          int    `gorm:"default:1" json:"pages"` // Number of pages processed
	Language       string `gorm:"size:10;default:'en'" json:"language"`
	ProcessingTime int    `gorm:"default:0" json:"processing_time_ms"` // Processing time in milliseconds

	Status string `gorm:"size:50;default:'pending'" json:"status"` // "pending", "processing", "completed", "failed"
	Error  string `gorm:"type:text" json:"error"`                  // Error message if failed

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExtractedField represents a single extracted field from a document
type ExtractedField struct {
	Name        string       `json:"name"`
	Value       string       `json:"value"`
	Confidence  float64      `json:"confidence"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

// BoundingBox represents coordinates of text in a document
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ExtractedTable represents a table extracted from a document
type ExtractedTable struct {
	PageIndex   int          `json:"page_index"`
	Headers     []string     `json:"headers"`
	Rows        [][]string   `json:"rows"`
	Confidence  float64      `json:"confidence"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

// DocumentAnalysisResult represents the complete analysis result
type DocumentAnalysisResult struct {
	TextContent string           `json:"text_content"`
	Fields      []ExtractedField `json:"fields"`
	Tables      []ExtractedTable `json:"tables"`
	Confidence  float64          `json:"confidence"`
	Pages       int              `json:"pages"`
	Language    string           `json:"language"`
}

// AIUsage tracks usage statistics for billing/analytics
type AIUsage struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	ModelID   uuid.UUID `gorm:"type:uuid;not null" json:"model_id"`
	Date      time.Time `gorm:"type:date;not null" json:"date"` // Usage date
	Tokens    int       `gorm:"default:0" json:"tokens"`        // Total tokens used
	Requests  int       `gorm:"default:0" json:"requests"`      // Total requests
	Cost      float64   `gorm:"default:0" json:"cost"`          // Total cost in USD
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AIUpgradeCheck represents a manual upgrade evaluation
type AIUpgradeCheck struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CheckedAt      time.Time `json:"checked_at"`
	HardwareInfo   string    `gorm:"type:text" json:"hardware_info"`   // JSON hardware info
	UpgradeOptions string    `gorm:"type:text" json:"upgrade_options"` // JSON upgrade options
	Status         string    `gorm:"size:50;default:'completed'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// UpgradeOption represents a potential model upgrade
type UpgradeOption struct {
	ID              string   `json:"id"`
	CurrentModel    string   `json:"current_model"`
	NewModel        string   `json:"new_model"`
	Provider        string   `json:"provider"` // "ollama" or "openrouter"
	SizeGB          float64  `json:"size_gb"`
	Capabilities    []string `json:"capabilities"`
	PerformanceGain string   `json:"performance_gain"` // "better", "much_better", "similar"
	Compatibility   bool     `json:"compatibility"`
	Requirements    string   `json:"requirements"` // JSON hardware requirements
	Reasoning       string   `json:"reasoning"`
	RiskLevel       string   `json:"risk_level"` // "low", "medium", "high"
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

// Speech Service Models

// SpeechProvider represents different TTS/STT service providers
type SpeechProvider string

const (
	SpeechProviderLocal      SpeechProvider = "local"
	SpeechProviderElevenLabs SpeechProvider = "elevenlabs"
	SpeechProviderOpenAI     SpeechProvider = "openai"
	SpeechProviderReplicate  SpeechProvider = "replicate"
	SpeechProviderPiper      SpeechProvider = "piper"
	SpeechProviderAssemblyAI SpeechProvider = "assemblyai"
	SpeechProviderDeepgram   SpeechProvider = "deepgram"
)

// SpeechModel represents a TTS/STT model configuration
type SpeechModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Provider  SpeechProvider `gorm:"size:50;not null" json:"provider"`
	ModelID   string         `gorm:"size:255;not null" json:"model_id"`
	Type      string         `gorm:"size:20;not null" json:"type"` // "tts" or "stt"
	Language  string         `gorm:"size:10;default:'en'" json:"language"`
	Voice     string         `gorm:"size:100" json:"voice,omitempty"`           // For TTS
	Quality   string         `gorm:"size:20;default:'standard'" json:"quality"` // "low", "standard", "high"
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	IsSystem  bool           `gorm:"default:false" json:"is_system"`
	UserID    uuid.UUID      `gorm:"type:uuid" json:"user_id,omitempty"`
	APIKey    []byte         `json:"-"` // Encrypted API key
	Endpoint  string         `gorm:"size:500" json:"endpoint,omitempty"`
	Config    string         `gorm:"type:jsonb" json:"config,omitempty"` // Provider-specific config
	Priority  int            `gorm:"default:0" json:"priority"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// SpeechRequest represents a TTS/STT processing request
type SpeechRequest struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ModelID        uuid.UUID  `gorm:"type:uuid;not null" json:"model_id"`
	RequestType    string     `gorm:"size:20;not null" json:"request_type"`       // "tts" or "stt"
	InputText      string     `gorm:"type:text" json:"input_text,omitempty"`      // For TTS
	InputAudio     []byte     `json:"-"`                                          // For STT (stored encrypted)
	OutputText     string     `gorm:"type:text" json:"output_text,omitempty"`     // For STT
	OutputAudio    []byte     `json:"-"`                                          // For TTS (stored encrypted)
	Status         string     `gorm:"size:20;default:'processing'" json:"status"` // "processing", "completed", "failed"
	Error          string     `gorm:"type:text" json:"error,omitempty"`
	ProcessingTime int        `gorm:"default:0" json:"processing_time_ms"`
	AudioFormat    string     `gorm:"size:20;default:'mp3'" json:"audio_format"` // "mp3", "wav", "ogg"
	Language       string     `gorm:"size:10;default:'en'" json:"language"`
	Voice          string     `gorm:"size:100" json:"voice,omitempty"`
	Speed          float64    `gorm:"default:1.0" json:"speed"` // TTS speed multiplier
	Pitch          float64    `gorm:"default:1.0" json:"pitch"` // TTS pitch multiplier
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

// SpeechSettings represents user speech preferences
type SpeechSettings struct {
	UserID            uuid.UUID `gorm:"type:uuid;primary_key" json:"user_id"`
	EnableTTS         bool      `gorm:"default:true" json:"enable_tts"`
	EnableSTT         bool      `gorm:"default:true" json:"enable_stt"`
	DefaultTTSModel   uuid.UUID `gorm:"type:uuid" json:"default_tts_model,omitempty"`
	DefaultSTTModel   uuid.UUID `gorm:"type:uuid" json:"default_stt_model,omitempty"`
	DefaultVoice      string    `gorm:"size:100;default:'alloy'" json:"default_voice"`
	DefaultLanguage   string    `gorm:"size:10;default:'en'" json:"default_language"`
	TTSSpeed          float64   `gorm:"default:1.0" json:"tts_speed"`
	TTSVolume         float64   `gorm:"default:1.0" json:"tts_volume"`
	STTLanguage       string    `gorm:"size:10;default:'en'" json:"stt_language"`
	AutoPlayTTS       bool      `gorm:"default:false" json:"auto_play_tts"`
	ShowSTTTranscript bool      `gorm:"default:true" json:"show_stt_transcript"`
	KeyboardShortcut  string    `gorm:"size:50;default:'ctrl+shift+s'" json:"keyboard_shortcut"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Default Speech Models
var DefaultSpeechModels = []SpeechModel{
	// Local TTS (system voices)
	{
		Name:     "System TTS (English)",
		Provider: SpeechProviderLocal,
		ModelID:  "system-en",
		Type:     "tts",
		Language: "en",
		IsSystem: true,
		IsActive: true,
		Priority: 10,
	},
	{
		Name:     "System TTS (Spanish)",
		Provider: SpeechProviderLocal,
		ModelID:  "system-es",
		Type:     "tts",
		Language: "es",
		IsSystem: true,
		IsActive: true,
		Priority: 9,
	},

	// Local STT (system recognition)
	{
		Name:     "System STT (English)",
		Provider: SpeechProviderLocal,
		ModelID:  "system-stt-en",
		Type:     "stt",
		Language: "en",
		IsSystem: true,
		IsActive: true,
		Priority: 10,
	},

	// ElevenLabs TTS
	{
		Name:     "ElevenLabs - Rachel",
		Provider: SpeechProviderElevenLabs,
		ModelID:  "21m00Tcm4TlvDq8ikWAM",
		Type:     "tts",
		Language: "en",
		Voice:    "Rachel",
		Quality:  "high",
		IsSystem: true,
		IsActive: true,
		Priority: 8,
	},
	{
		Name:     "ElevenLabs - Drew",
		Provider: SpeechProviderElevenLabs,
		ModelID:  "29vD33N1CtxCmqQRPOHJ",
		Type:     "tts",
		Language: "en",
		Voice:    "Drew",
		Quality:  "high",
		IsSystem: true,
		IsActive: true,
		Priority: 8,
	},

	// OpenAI TTS
	{
		Name:     "OpenAI - Alloy",
		Provider: SpeechProviderOpenAI,
		ModelID:  "tts-1",
		Type:     "tts",
		Language: "en",
		Voice:    "alloy",
		Quality:  "standard",
		IsSystem: true,
		IsActive: true,
		Priority: 7,
	},
	{
		Name:     "OpenAI - Echo",
		Provider: SpeechProviderOpenAI,
		ModelID:  "tts-1",
		Type:     "tts",
		Language: "en",
		Voice:    "echo",
		Quality:  "standard",
		IsSystem: true,
		IsActive: true,
		Priority: 7,
	},

	// OpenAI Whisper STT
	{
		Name:     "OpenAI Whisper",
		Provider: SpeechProviderOpenAI,
		ModelID:  "whisper-1",
		Type:     "stt",
		Language: "en",
		IsSystem: true,
		IsActive: true,
		Priority: 9,
	},

	// AssemblyAI STT
	{
		Name:     "AssemblyAI Premium",
		Provider: SpeechProviderAssemblyAI,
		ModelID:  "premium",
		Type:     "stt",
		Language: "en",
		IsSystem: true,
		IsActive: true,
		Priority: 8,
	},

	// Deepgram STT
	{
		Name:     "Deepgram Nova",
		Provider: SpeechProviderDeepgram,
		ModelID:  "nova",
		Type:     "stt",
		Language: "en",
		IsSystem: true,
		IsActive: true,
		Priority: 8,
	},
}

// Workflow Automation Models

// Workflow represents an automated workflow
type Workflow struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsTemplate  bool      `gorm:"default:false" json:"is_template"`
	Category    string    `gorm:"size:100" json:"category"`             // "invoice_processing", "lead_management", etc.
	TriggerType string    `gorm:"size:50;not null" json:"trigger_type"` // "manual", "scheduled", "event"
	Schedule    string    `gorm:"size:500" json:"schedule,omitempty"`   // Cron expression for scheduled workflows

	// Visual layout
	CanvasData string `gorm:"type:jsonb" json:"canvas_data"` // Node positions, connections

	// Metadata
	Version      int        `gorm:"default:1" json:"version"`
	LastRunAt    *time.Time `json:"last_run_at"`
	RunCount     int        `gorm:"default:0" json:"run_count"`
	SuccessCount int        `gorm:"default:0" json:"success_count"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// WorkflowNode represents a node in the workflow
type WorkflowNode struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;not null" json:"workflow_id"`
	NodeID     string    `gorm:"size:100;not null" json:"node_id"`  // Unique within workflow
	NodeType   string    `gorm:"size:50;not null" json:"node_type"` // "trigger", "action", "condition", "delay"

	// Node configuration
	Name        string `gorm:"size:255;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Config      string `gorm:"type:jsonb" json:"config"` // Node-specific configuration

	// Visual properties
	PositionX int `gorm:"default:0" json:"position_x"`
	PositionY int `gorm:"default:0" json:"position_y"`
	Width     int `gorm:"default:200" json:"width"`
	Height    int `gorm:"default:100" json:"height"`

	// Execution state
	IsEnabled bool       `gorm:"default:true" json:"is_enabled"`
	LastRunAt *time.Time `json:"last_run_at"`
	RunCount  int        `gorm:"default:0" json:"run_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkflowConnection represents a connection between workflow nodes
type WorkflowConnection struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;not null" json:"workflow_id"`
	FromNodeID string    `gorm:"size:100;not null" json:"from_node_id"`
	ToNodeID   string    `gorm:"size:100;not null" json:"to_node_id"`
	FromPort   string    `gorm:"size:50;default:'output'" json:"from_port"`
	ToPort     string    `gorm:"size:50;default:'input'" json:"to_port"`

	// Connection properties
	Label     string `gorm:"size:255" json:"label"`
	IsEnabled bool   `gorm:"default:true" json:"is_enabled"`
	Condition string `gorm:"type:text" json:"condition,omitempty"` // Conditional logic

	CreatedAt time.Time `json:"created_at"`
}

// WorkflowExecution represents a workflow execution instance
type WorkflowExecution struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	WorkflowID uuid.UUID `gorm:"type:uuid;not null" json:"workflow_id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`

	// Execution details
	Status      string `gorm:"size:20;default:'running'" json:"status"` // "running", "completed", "failed", "paused"
	TriggerType string `gorm:"size:50" json:"trigger_type"`
	TriggerData string `gorm:"type:jsonb" json:"trigger_data"`  // Data that triggered execution
	IsDryRun    bool   `gorm:"default:false" json:"is_dry_run"` // If true, action nodes only preview their effect

	// Execution state
	CurrentNodeID string `gorm:"size:100" json:"current_node_id"`
	NodeStates    string `gorm:"type:jsonb" json:"node_states"` // State of each node

	// Results
	OutputData   string     `gorm:"type:jsonb" json:"output_data"`
	ErrorMessage string     `gorm:"type:text" json:"error_message"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	Duration     int        `gorm:"default:0" json:"duration_ms"` // Execution time in milliseconds

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WorkflowTemplate represents a reusable workflow template
type WorkflowTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:100;not null" json:"category"`
	Icon        string    `gorm:"size:100" json:"icon"`
	Color       string    `gorm:"size:7;default:'#007bff'" json:"color"`

	// Template data
	TemplateData string `gorm:"type:jsonb;not null" json:"template_data"` // Complete workflow definition
	IsSystem     bool   `gorm:"default:false" json:"is_system"`
	IsPublic     bool   `gorm:"default:true" json:"is_public"`

	// Usage statistics
	UseCount int `gorm:"default:0" json:"use_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IntegrationConnector represents available integration connectors
type IntegrationConnector struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	AppName       string    `gorm:"size:100;not null" json:"app_name"`      // "forms", "email", "calendar", "tasks", etc.
	ConnectorType string    `gorm:"size:50;not null" json:"connector_type"` // "trigger", "action"

	// Configuration
	ConfigSchema string `gorm:"type:jsonb" json:"config_schema"` // JSON schema for configuration
	Icon         string `gorm:"size:100" json:"icon"`
	Color        string `gorm:"size:7;default:'#007bff'" json:"color"`

	IsActive bool `gorm:"default:true" json:"is_active"`
	IsSystem bool `gorm:"default:true" json:"is_system"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Default Workflow Templates.
//
// Each TemplateData is a real, executable canvas (nodes/connections) built
// entirely from connectors that already exist in workflow_service.go — none
// of these rely on integrations that aren't implemented yet. Every node's
// "config" must be filled in with real IDs (form_id, spreadsheet_id, etc.)
// by the user after instantiating from the template; CreateWorkflowFromTemplate
// ships the resulting workflow inactive by default so nothing runs until the
// user explicitly reviews and enables it (ideally via a dry run first).
var DefaultWorkflowTemplates = []WorkflowTemplate{
	{
		Name:        "Form Response Triage",
		Description: "When a form is submitted, route it to an urgent or normal-priority task depending on a field value.",
		Category:    "form_triage",
		Icon:        "🚦",
		Color:       "#dc3545",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 120}, "config": {"connector": "forms.submission", "form_id": ""}},
				{"id": "condition-1", "type": "condition", "position": {"x": 320, "y": 120}, "config": {"field": "priority", "operator": "equals", "value": "high"}},
				{"id": "action-urgent", "type": "action", "position": {"x": 600, "y": 40}, "config": {"connector": "tasks.create", "title": "Urgent: review new form response", "description": "Auto-created by the Form Response Triage preset workflow.", "priority": "urgent"}},
				{"id": "action-medium", "type": "action", "position": {"x": 600, "y": 200}, "config": {"connector": "tasks.create", "title": "Review new form response", "description": "Auto-created by the Form Response Triage preset workflow.", "priority": "medium"}}
			],
			"connections": [
				{"from": "trigger-1", "to": "condition-1"},
				{"from": "condition-1", "to": "action-urgent", "fromPort": "true"},
				{"from": "condition-1", "to": "action-medium", "fromPort": "false"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "Form Response to Spreadsheet Log",
		Description: "When a form is submitted, log the response into a spreadsheet.",
		Category:    "form_logging",
		Icon:        "📋",
		Color:       "#007bff",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 120}, "config": {"connector": "forms.submission", "form_id": ""}},
				{"id": "action-log", "type": "action", "position": {"x": 340, "y": 120}, "config": {"connector": "spreadsheet.update", "spreadsheet_id": "", "range": "A:A", "values": []}}
			],
			"connections": [
				{"from": "trigger-1", "to": "action-log"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "New Client Onboarding Checklist",
		Description: "Manually run this after signing a new client: creates the setup tasks and schedules a kickoff call.",
		Category:    "client_onboarding",
		Icon:        "🚀",
		Color:       "#6f42c1",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 140}, "config": {"connector": "manual"}},
				{"id": "task-folder", "type": "action", "position": {"x": 340, "y": 40}, "config": {"connector": "tasks.create", "title": "Set up client folder and draft contract", "priority": "high"}},
				{"id": "task-agenda", "type": "action", "position": {"x": 340, "y": 140}, "config": {"connector": "tasks.create", "title": "Prepare kickoff call agenda", "priority": "medium"}},
				{"id": "calendar-kickoff", "type": "action", "position": {"x": 340, "y": 240}, "config": {"connector": "calendar.create_event", "title": "Client Kickoff Call", "duration": 60}}
			],
			"connections": [
				{"from": "trigger-1", "to": "task-folder"},
				{"from": "trigger-1", "to": "task-agenda"},
				{"from": "trigger-1", "to": "calendar-kickoff"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "Overdue Task Digest",
		Description: "Manually run to email yourself a digest of open high-priority work for the week.",
		Category:    "digests",
		Icon:        "📨",
		Color:       "#17a2b8",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 120}, "config": {"connector": "manual"}},
				{"id": "task-digest", "type": "action", "position": {"x": 340, "y": 120}, "config": {"connector": "email.send", "to": "", "subject": "Your open high-priority tasks", "body": "Reminder to review this week's open high-priority tasks."}},
				{"id": "task-followup", "type": "action", "position": {"x": 620, "y": 120}, "config": {"connector": "tasks.create", "title": "Triage this week's open tasks", "description": "Created by the Overdue Task Digest preset.", "priority": "medium"}}
			],
			"connections": [
				{"from": "trigger-1", "to": "task-digest"},
				{"from": "trigger-1", "to": "task-followup"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "Meeting Follow-up Action Items",
		Description: "After a meeting, capture the agreed action items as tasks and schedule a review.",
		Category:    "follow_up",
		Icon:        "📝",
		Color:       "#20c997",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 140}, "config": {"connector": "manual"}},
				{"id": "task-actions", "type": "action", "position": {"x": 340, "y": 40}, "config": {"connector": "tasks.create", "title": "Complete meeting action items", "description": "Owner-assigned follow-ups from the meeting.", "priority": "high"}},
				{"id": "task-share", "type": "action", "position": {"x": 340, "y": 160}, "config": {"connector": "tasks.create", "title": "Share meeting notes with attendees", "priority": "low"}},
				{"id": "calendar-review", "type": "action", "position": {"x": 620, "y": 160}, "config": {"connector": "calendar.create_event", "title": "Follow-up review meeting", "duration": 30}}
			],
			"connections": [
				{"from": "trigger-1", "to": "task-actions"},
				{"from": "trigger-1", "to": "task-share"},
				{"from": "task-share", "to": "calendar-review"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "Spreadsheet Threshold Alert",
		Description: "When a form is submitted, alert if a numeric field exceeds a threshold (e.g. budget overage).",
		Category:    "alerting",
		Icon:        "🔔",
		Color:       "#fd7e14",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 120}, "config": {"connector": "forms.submission", "form_id": ""}},
				{"id": "condition-1", "type": "condition", "position": {"x": 320, "y": 120}, "config": {"field": "amount", "operator": "greater_than", "value": "1000"}},
				{"id": "notify-over", "type": "action", "position": {"x": 600, "y": 40}, "config": {"connector": "notifications.send", "title": "Threshold exceeded", "message": "A submitted value exceeded the configured threshold.", "type": "warning"}},
				{"id": "task-over", "type": "action", "position": {"x": 600, "y": 200}, "config": {"connector": "tasks.create", "title": "Review threshold-exceeding submission", "priority": "high"}}
			],
			"connections": [
				{"from": "trigger-1", "to": "condition-1"},
				{"from": "condition-1", "to": "notify-over", "fromPort": "true"},
				{"from": "condition-1", "to": "task-over", "fromPort": "true"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "Delayed Escalation Reminder",
		Description: "Send an initial reminder, wait, then escalate to a task if still unaddressed.",
		Category:    "escalation",
		Icon:        "⏰",
		Color:       "#e83e8c",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 120}, "config": {"connector": "manual"}},
				{"id": "notify-first", "type": "action", "position": {"x": 320, "y": 120}, "config": {"connector": "notifications.send", "title": "Reminder", "message": "Please address the pending item.", "type": "info"}},
				{"id": "delay-1", "type": "action", "position": {"x": 600, "y": 120}, "config": {"connector": "logic.delay", "delay_seconds": 86400}},
				{"id": "task-escalate", "type": "action", "position": {"x": 880, "y": 120}, "config": {"connector": "tasks.create", "title": "Escalate: item still unaddressed", "priority": "high"}}
			],
			"connections": [
				{"from": "trigger-1", "to": "notify-first"},
				{"from": "notify-first", "to": "delay-1"},
				{"from": "delay-1", "to": "task-escalate"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
	{
		Name:        "ERP Bridge (Form → signed ERP POST)",
		Description: "When a form is submitted, POST its data to an external ERP REST API (e.g. tpt-free-erp) with an HMAC-SHA256 signature so the ERP can verify the call came from Titan. Fill in the ERP url and shared secret before activating.",
		Category:    "integration",
		Icon:        "🔗",
		Color:       "#6f42c1",
		IsSystem:    true,
		IsPublic:    true,
		TemplateData: `{
			"nodes": [
				{"id": "trigger-1", "type": "trigger", "position": {"x": 60, "y": 120}, "config": {"connector": "forms.submission", "form_id": ""}},
				{"id": "erp-post", "type": "action", "position": {"x": 360, "y": 120}, "config": {"connector": "http.request", "method": "POST", "url": "https://erp.example.com/api/v1/records", "headers": {"Content-Type": "application/json"}, "body": "{\"source\": \"tpt-titan\", \"form_id\": \"{{form_id}}\", \"data\": {{submission}}", "signing_secret": "REPLACE_WITH_ERP_SHARED_SECRET", "retry_attempts": 2}},
				{"id": "erp-task", "type": "action", "position": {"x": 680, "y": 120}, "config": {"connector": "tasks.create", "title": "Investigate ERP bridge failure", "priority": "high"}}
			],
			"connections": [
				{"from": "trigger-1", "to": "erp-post"},
				{"from": "erp-post", "to": "erp-task", "fromPort": "false"}
			],
			"viewport": {"x": 0, "y": 0, "zoom": 1}
		}`,
	},
}

// Default Integration Connectors
var DefaultIntegrationConnectors = []IntegrationConnector{
	// Triggers
	{
		Name:          "Form Submission",
		Description:   "Trigger when a form is submitted",
		AppName:       "forms",
		ConnectorType: "trigger",
		Icon:          "📝",
		Color:         "#007bff",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Email Received",
		Description:   "Trigger when a new email is received",
		AppName:       "email",
		ConnectorType: "trigger",
		Icon:          "📧",
		Color:         "#dc3545",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Calendar Event",
		Description:   "Trigger on calendar events",
		AppName:       "calendar",
		ConnectorType: "trigger",
		Icon:          "📅",
		Color:         "#28a745",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Task Created/Updated",
		Description:   "Trigger when tasks are created or updated",
		AppName:       "tasks",
		ConnectorType: "trigger",
		Icon:          "✅",
		Color:         "#ffc107",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Document Uploaded",
		Description:   "Trigger when documents are uploaded",
		AppName:       "documents",
		ConnectorType: "trigger",
		Icon:          "📎",
		Color:         "#6c757d",
		IsSystem:      true,
		IsActive:      true,
	},

	// Actions
	{
		Name:          "Send Email",
		Description:   "Send an automated email",
		AppName:       "email",
		ConnectorType: "action",
		Icon:          "📤",
		Color:         "#dc3545",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Create Calendar Event",
		Description:   "Schedule a calendar event",
		AppName:       "calendar",
		ConnectorType: "action",
		Icon:          "📅",
		Color:         "#28a745",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Create/Update Task",
		Description:   "Create or update a task",
		AppName:       "tasks",
		ConnectorType: "action",
		Icon:          "📋",
		Color:         "#ffc107",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Update Spreadsheet",
		Description:   "Update spreadsheet data",
		AppName:       "spreadsheet",
		ConnectorType: "action",
		Icon:          "📊",
		Color:         "#17a2b8",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Send Notification",
		Description:   "Send in-app notification",
		AppName:       "notifications",
		ConnectorType: "action",
		Icon:          "🔔",
		Color:         "#6f42c1",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Conditional Logic",
		Description:   "Add conditional branching logic",
		AppName:       "logic",
		ConnectorType: "action",
		Icon:          "🔀",
		Color:         "#fd7e14",
		IsSystem:      true,
		IsActive:      true,
	},
	{
		Name:          "Delay/Timer",
		Description:   "Add time delay to workflow",
		AppName:       "logic",
		ConnectorType: "action",
		Icon:          "⏱️",
		Color:         "#20c997",
		IsSystem:      true,
		IsActive:      true,
	},
}

// AI Settings Model
type AISettings struct {
	UserID                uuid.UUID `gorm:"type:uuid;primary_key" json:"user_id"`
	EnableAIFeatures      bool      `gorm:"default:true" json:"enable_ai_features"`
	EnableOCR             bool      `gorm:"default:true" json:"enable_ocr"`
	EnableSpeech          bool      `gorm:"default:true" json:"enable_speech"`
	EnableWorkflows       bool      `gorm:"default:true" json:"enable_workflows"`
	EnableLocalAI         bool      `gorm:"default:true" json:"enable_local_ai"`
	EnableCloudAI         bool      `gorm:"default:false" json:"enable_cloud_ai"`
	DefaultProvider       string    `gorm:"size:50;default:'ollama'" json:"default_provider"`
	MaxConcurrentRequests int       `gorm:"default:3" json:"max_concurrent_requests"`
	RequestTimeout        int       `gorm:"default:30" json:"request_timeout"`
	HardwareAcceleration  bool      `gorm:"default:true" json:"hardware_acceleration"`
	LowPowerMode          bool      `gorm:"default:false" json:"low_power_mode"`

	// API Keys (encrypted in database)
	OpenAIKey     []byte `json:"-"` // Encrypted
	ElevenLabsKey []byte `json:"-"` // Encrypted
	ReplicateKey  []byte `json:"-"` // Encrypted
	AssemblyAIKey []byte `json:"-"` // Encrypted
	DeepgramKey   []byte `json:"-"` // Encrypted

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AI Usage Statistics
type AIUsageStats struct {
	UserID    uuid.UUID `gorm:"type:uuid;primary_key" json:"user_id"`
	Date      time.Time `gorm:"type:date;primary_key" json:"date"`
	Provider  string    `gorm:"size:50;not null" json:"provider"`
	Service   string    `gorm:"size:50;not null" json:"service"` // "tts", "stt", "ocr", "chat", etc.
	Tokens    int       `gorm:"default:0" json:"tokens"`
	Requests  int       `gorm:"default:0" json:"requests"`
	Cost      float64   `gorm:"default:0" json:"cost"`
	CreatedAt time.Time `json:"created_at"`
}

// Graceful Degradation Settings
type GracefulDegradation struct {
	UserID             uuid.UUID `gorm:"type:uuid;primary_key" json:"user_id"`
	OfflineMode        bool      `gorm:"default:false" json:"offline_mode"`
	FallbackToLocal    bool      `gorm:"default:true" json:"fallback_to_local"`
	ReduceQuality      bool      `gorm:"default:false" json:"reduce_quality"`
	DisableAdvanced    bool      `gorm:"default:false" json:"disable_advanced"`
	ShowDegradedNotice bool      `gorm:"default:true" json:"show_degraded_notice"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TableName overrides for better naming
func (AIModel) TableName() string          { return "ai_models" }
func (AITask) TableName() string           { return "ai_tasks" }
func (AIRequest) TableName() string        { return "ai_requests" }
func (AIUsage) TableName() string          { return "ai_usage" }
func (AIUpgradeCheck) TableName() string   { return "ai_upgrade_checks" }
func (DocumentAnalysis) TableName() string { return "document_analyses" }
func (SpeechModel) TableName() string      { return "speech_models" }
func (SpeechRequest) TableName() string    { return "speech_requests" }
func (SpeechSettings) TableName() string   { return "speech_settings" }
func (AISettings) TableName() string       { return "ai_settings" }
func (AIUsageStats) TableName() string     { return "ai_usage_stats" }
