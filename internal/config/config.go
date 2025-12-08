package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppEnv     string
	ServerPort string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	JWTSecret string
}

// Load config from env
func Load() (*Config, []string, error) {
	warnings := []string{}

	cfg := &Config{
		AppEnv:     getEnvWithWarn("APP_ENV", "dev", &warnings),
		ServerPort: getEnvWithWarn("SERVER_PORT", "8080", &warnings),

		DBHost:     getEnvWithWarn("DB_HOST", "localhost", &warnings),
		DBUser:     getEnvWithWarn("DB_USER", "postgres", &warnings),
		DBPassword: getEnvWithWarn("DB_PASSWORD", "password", &warnings),
		DBName:     getEnvWithWarn("DB_NAME", "app", &warnings),
		DBPort:     getEnvWithWarn("DB_PORT", "5432", &warnings),

		JWTSecret: os.Getenv("JWT_SECRET"), // обязательное поле
	}

	if err := cfg.Validate(); err != nil {
		return nil, warnings, err
	}

	return cfg, warnings, nil
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
