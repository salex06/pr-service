// Package config - пакет, в котором хранятся структуры конфигурации приложения
package config

import (
	"os"
)

// DBConfig представляет набор параметров,
// определяющих конфигурацию базы данных
type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// AppConfig представляет набор параметров,
// определяющих конфигурацию приложения
type AppConfig struct {
	ServerPort string
}

// LoadDBConfig формирует конфигурацию БД
// на основе переменных окружения
func LoadDBConfig() *DBConfig {
	return &DBConfig{
		DBHost:     getEnv("POSTGRES_CONTAINER", "localhost"),
		DBPort:     getEnv("POSTGRES_CONTAINER_PORT", "5433"),
		DBUser:     getEnv("POSTGRES_USER", "pr-service-admin"),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", "postgres"),
	}
}

// LoadAppConfig формирует конфигурацию приложения
// на основе переменных окружения
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
