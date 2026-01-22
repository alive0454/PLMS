package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name    string
	Version string
	Port    string
	Env     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "PLMS"),
			Version: getEnv("APP_VERSION", "1.0.0"),
			Port:    getEnv("APP_PORT", "8080"),
			Env:     getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "47.93.2.51"),
			Port:     getEnv("DB_PORT", "23306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "HvOI*Vzb2uTC6V45"),
			Name:     getEnv("DB_NAME", "plms"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
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
