package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv     string `env:"APP_ENV" envDefault:"dev"`
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	JWTSecret             string `env:"JWT_SECRET" envDefault:"secret"`
	JWTAccessExpireMin    int    `env:"JWT_ACCESS_EXPIRE_MIN" envDefault:"15"`
	JWTRefreshExpireHours int    `env:"JWT_REFRESH_EXPIRE_HOURS" envDefault:"168"` // 7 days

	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"password"`
	DBName     string `env:"DB_NAME" envDefault:"mws_ai"`
	DBSSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`

	UploadDir string `env:"UPLOAD_DIR" envDefault:"uploads"`

	// --- EXTERNAL SERVICES ---
	HeuristicURL string `env:"HEURISTIC_URL" envDefault:"http://localhost:8081"`
	MLURL        string `env:"ML_URL" envDefault:"http://localhost:8082"`
	LLMURL       string `env:"LLM_URL" envDefault:"http://localhost:8083"`
}

// Load config from env
func Load() (*Config, []string, error) {
	warnings := []string{}

	cfg := &Config{
		// --- APP ---
		AppEnv:     getEnvWithWarn("APP_ENV", "dev", &warnings),
		ServerPort: getEnvWithWarn("SERVER_PORT", "8080", &warnings),
		LogLevel:   getEnvWithWarn("LOG_LEVEL", "info", &warnings),

		// --- JWT ---
		JWTSecret:             os.Getenv("JWT_SECRET"), // обязательное
		JWTAccessExpireMin:    getEnvIntWithWarn("JWT_ACCESS_EXPIRE_MIN", 15, &warnings),
		JWTRefreshExpireHours: getEnvIntWithWarn("JWT_REFRESH_EXPIRE_HOURS", 168, &warnings), // 7 days

		// --- DATABASE ---
		DBHost:     getEnvWithWarn("DB_HOST", "localhost", &warnings),
		DBPort:     getEnvWithWarn("DB_PORT", "5432", &warnings),
		DBUser:     getEnvWithWarn("DB_USER", "postgres", &warnings),
		DBPassword: getEnvWithWarn("DB_PASSWORD", "password", &warnings),
		DBName:     getEnvWithWarn("DB_NAME", "mws_ai", &warnings),
		DBSSLMode:  getEnvWithWarn("DB_SSL_MODE", "disable", &warnings),

		// --- FILE STORAGE ---
		UploadDir: getEnvWithWarn("UPLOAD_DIR", "uploads", &warnings),

		// --- EXTERNAL SERVICES ---
		HeuristicURL: getEnvWithWarn("HEURISTIC_URL", "http://localhost:8081", &warnings),
		MLURL:        getEnvWithWarn("ML_URL", "http://localhost:8082", &warnings),
		LLMURL:       getEnvWithWarn("LLM_URL", "http://localhost:8083", &warnings),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, warnings, err
	}

	return cfg, warnings, nil
}

func getEnvIntWithWarn(key string, defaultVal int, warnings *[]string) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		*warnings = append(*warnings, fmt.Sprintf("ENV %s is missing, using default: %d", key, defaultVal))
		return defaultVal
	}

	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		*warnings = append(*warnings, fmt.Sprintf("ENV %s invalid (%s), using default: %d", key, valStr, defaultVal))
		return defaultVal
	}

	return valInt
}

func getEnvWithWarn(key, defaultValue string, warnings *[]string) string {
	v := os.Getenv(key)
	if v == "" {
		*warnings = append(*warnings, fmt.Sprintf("%s not set — using default '%s'", key, defaultValue))
		return defaultValue
	}
	return v
}

// Validate
func (c *Config) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required but missing")
	}

	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required but missing")
	}

	if c.DBUser == "" {
		return fmt.Errorf("DB_USER is required but missing")
	}

	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required but missing")
	}

	if c.HeuristicURL == "" {
		return fmt.Errorf("HEURISTIC_URL is required but missing")
	}

	if c.MLURL == "" {
		return fmt.Errorf("ML_URL is required but missing")
	}

	if c.LLMURL == "" {
		return fmt.Errorf("LLM_URL is required but missing")
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
