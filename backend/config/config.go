package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadAndValidate loads configuration and validates required fields
func LoadAndValidate() *Config {
	cfg := Load()
	
	// Validate required configuration
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required. Please set a secure secret key (min 32 characters recommended)")
	}
	
	if len(cfg.JWT.Secret) < 32 {
		log.Println("WARNING: JWT_SECRET should be at least 32 characters for security")
	}
	
	return cfg
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Email    EmailConfig
	AI       AIConfig
	Speech   SpeechConfig
	P2P      P2PConfig
}

type ServerConfig struct {
	Host string
	Port string
	Mode string
}

type DatabaseConfig struct {
	Type     string // "sqlite" or "postgres"
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Path     string // SQLite database file path
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	ExpiryHour int
}

type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}

type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	Username string
	Password string
}

type AIConfig struct {
	OllamaHost      string
	OllamaPort      string
	OpenRouterKey   string
	DefaultModels   []string
	EnableLocalAI   bool
	EnableOnlineAI  bool
}

type SpeechConfig struct {
	ElevenLabsKey   string
	OpenAIKey       string
	ReplicateKey    string
	PiperKey        string
	AssemblyAIKey   string
	DeepgramKey     string
	EnableLocalTTS  bool
	EnableLocalSTT  bool
	EnableCloudTTS  bool
	EnableCloudSTT  bool
	DefaultVoice    string
	DefaultLanguage string
}

type P2PConfig struct {
	Enabled             bool   // Enable P2P mode
	ServiceName         string // mDNS service name
	ServiceType         string // mDNS service type
	Port                int    // P2P listening port
	DiscoveryTimeout    int    // Peer discovery timeout (seconds)
	SyncInterval        int    // Sync interval (seconds)
	MaxPeers            int    // Maximum number of peers to connect to
	EnableEncryption    bool   // Encrypt P2P communications
	EnableCompression   bool   // Compress data transfers
	ConflictStrategy    string // "last_write_wins" or "manual_merge"

	// Remote access configuration
	AllowRemoteAccess   bool   // Allow remote connections via cloud relay
	RemoteAccessMode    string // "automatic", "relay_only", "vpn_required"
	CloudRelayEnabled   bool   // Use cloud relay for remote connections
	RelayServerURL      string // Cloud relay server URL
	AutoDetectTopology  bool   // Automatically choose P2P vs relay
	PreferredTopology   string // "mesh", "star", "hybrid"
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "sqlite"), // Default to SQLite
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "tpt_user"),
			Password: getEnv("DB_PASSWORD", "tpt_password"),
			DBName:   getEnv("DB_NAME", "tpt_titan"),
			Path:     getEnv("DB_PATH", "./data/tpt-titan.db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			ExpiryHour: getEnvAsInt("JWT_EXPIRY_HOUR", 24),
		},
		Redis: RedisConfig{
			Enabled:  getEnvAsBool("REDIS_ENABLED", false),
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Email: EmailConfig{
			SMTPHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort: getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
		},
		AI: AIConfig{
			OllamaHost:     getEnv("OLLAMA_HOST", "localhost"),
			OllamaPort:     getEnv("OLLAMA_PORT", "11434"),
			OpenRouterKey:  getEnv("OPENROUTER_API_KEY", ""),
			DefaultModels:  []string{"llama3.2:3b", "phi3:3.8b", "deepseek-coder:6.7b"},
			EnableLocalAI:  getEnvAsBool("ENABLE_LOCAL_AI", true),
			EnableOnlineAI: getEnvAsBool("ENABLE_ONLINE_AI", false),
		},
		Speech: SpeechConfig{
			ElevenLabsKey:   getEnv("ELEVENLABS_API_KEY", ""),
			OpenAIKey:       getEnv("OPENAI_API_KEY", ""),
			AssemblyAIKey:   getEnv("ASSEMBLYAI_API_KEY", ""),
			DeepgramKey:     getEnv("DEEPGRAM_API_KEY", ""),
			EnableLocalTTS:  getEnvAsBool("ENABLE_LOCAL_TTS", true),
			EnableLocalSTT:  getEnvAsBool("ENABLE_LOCAL_STT", true),
			EnableCloudTTS:  getEnvAsBool("ENABLE_CLOUD_TTS", false),
			EnableCloudSTT:  getEnvAsBool("ENABLE_CLOUD_STT", false),
			DefaultVoice:    getEnv("DEFAULT_VOICE", "alloy"),
			DefaultLanguage: getEnv("DEFAULT_LANGUAGE", "en"),
		},
		P2P: P2PConfig{
			Enabled:           getEnvAsBool("P2P_ENABLED", false),
			ServiceName:       getEnv("P2P_SERVICE_NAME", "TPT Titan"),
			ServiceType:       getEnv("P2P_SERVICE_TYPE", "_tpt-titan._tcp"),
			Port:              getEnvAsInt("P2P_PORT", 8081),
			DiscoveryTimeout:  getEnvAsInt("P2P_DISCOVERY_TIMEOUT", 30),
			SyncInterval:      getEnvAsInt("P2P_SYNC_INTERVAL", 60),
			MaxPeers:          getEnvAsInt("P2P_MAX_PEERS", 10),
			EnableEncryption:  getEnvAsBool("P2P_ENCRYPTION", true),
			EnableCompression: getEnvAsBool("P2P_COMPRESSION", true),
			ConflictStrategy:  getEnv("P2P_CONFLICT_STRATEGY", "last_write_wins"),

			// Remote access - enabled by default for ease of use
			AllowRemoteAccess: getEnvAsBool("P2P_ALLOW_REMOTE", true),
			RemoteAccessMode:  getEnv("P2P_REMOTE_MODE", "automatic"),
			CloudRelayEnabled: getEnvAsBool("P2P_CLOUD_RELAY", true),
			RelayServerURL:    getEnv("P2P_RELAY_URL", "https://relay.tpt-titan.com"),
			AutoDetectTopology: getEnvAsBool("P2P_AUTO_TOPOLOGY", true),
			PreferredTopology: getEnv("P2P_TOPOLOGY", "hybrid"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
