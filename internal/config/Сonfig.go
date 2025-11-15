package config

import (
	"os"
)

type DbConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type AppConfig struct {
	ServerPort string
}

func LoadDbConfig() *DbConfig {
	return &DbConfig{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5433"),
		DBUser:     getEnv("POSTGRES_USER", "pr-service-admin"),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", "postgres"),
	}
}

func LoadAppConfig() *AppConfig {
	return &AppConfig{
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
