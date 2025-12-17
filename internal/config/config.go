package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv     string
	ServerPort string
	LogLevel   string

	JWTSecret             string
	JWTAccessExpireMin    int
	JWTRefreshExpireHours int

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	UploadDir string

	// External services
	HeuristicURL string
	MLURL        string
	LLMURL       string
}

func Load() (*Config, []string, error) {
	warnings := []string{}

	cfg := &Config{
		AppEnv:     getEnvWithWarn("APP_ENV", "dev", &warnings),
		ServerPort: getEnvWithWarn("SERVER_PORT", "8080", &warnings),
		LogLevel:   getEnvWithWarn("LOG_LEVEL", "info", &warnings),

		JWTSecret:             os.Getenv("JWT_SECRET"),
		JWTAccessExpireMin:    getEnvIntWithWarn("JWT_ACCESS_EXPIRE_MIN", 15, &warnings),
		JWTRefreshExpireHours: getEnvIntWithWarn("JWT_REFRESH_EXPIRE_HOURS", 168, &warnings),

		DBHost:     getEnvWithWarn("DB_HOST", "localhost", &warnings),
		DBPort:     getEnvWithWarn("DB_PORT", "5432", &warnings),
		DBUser:     getEnvWithWarn("DB_USER", "postgres", &warnings),
		DBPassword: getEnvWithWarn("DB_PASSWORD", "password", &warnings),
		DBName:     getEnvWithWarn("DB_NAME", "mws_ai", &warnings),

		UploadDir: getEnvWithWarn("UPLOAD_DIR", "uploads", &warnings),

		HeuristicURL: getEnvWithWarn("HEURISTIC_URL", "http://localhost:8081", &warnings),
		MLURL:        getEnvWithWarn("ML_URL", "http://localhost:8082", &warnings),
		LLMURL:       getEnvWithWarn("LLM_URL", "http://localhost:8083", &warnings),
	}

	if err := cfg.Validate(); err != nil {
		return nil, warnings, err
	}

	return cfg, warnings, nil
}

func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.DBHost == "" || c.DBUser == "" || c.DBName == "" {
		return fmt.Errorf("database configuration is incomplete")
	}
	if c.HeuristicURL == "" || c.MLURL == "" || c.LLMURL == "" {
		return fmt.Errorf("external service URLs are required")
	}
	return nil
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}

func (c *Config) JWTSecretBytes() []byte {
	return []byte(c.JWTSecret)
}

func getEnvWithWarn(key, defaultValue string, warnings *[]string) string {
	v := os.Getenv(key)
	if v == "" {
		*warnings = append(*warnings, fmt.Sprintf("%s not set — using default '%s'", key, defaultValue))
		return defaultValue
	}
	return v
}

func getEnvIntWithWarn(key string, defaultVal int, warnings *[]string) int {
	v := os.Getenv(key)
	if v == "" {
		*warnings = append(*warnings, fmt.Sprintf("%s not set — using default %d", key, defaultVal))
		return defaultVal
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		*warnings = append(*warnings, fmt.Sprintf("%s invalid — using default %d", key, defaultVal))
		return defaultVal
	}
	return i
}
