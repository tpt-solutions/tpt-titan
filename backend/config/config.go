package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Email    EmailConfig
	AI       AIConfig
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
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			ExpiryHour: getEnvAsInt("JWT_EXPIRY_HOUR", 24),
		},
		Redis: RedisConfig{
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
