// Package config provides application configuration management.
package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Environment represents the application runtime environment.
type Environment string

const (
	// Dev is the development environment.
	Dev Environment = "dev"
	// Prod is the production environment.
	Prod Environment = "prod"
)

// Config holds all application configuration values.
type Config struct {
	Environment   Environment
	Host          string
	Port          string
	DBPath        string
	LogLevel      slog.Level
	SessionSecret string
}

var (
	// Global holds the singleton configuration instance.
	Global *Config
	once   sync.Once
)

func init() {
	once.Do(func() {
		Global = Load()
	})
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func loadBase() *Config {
	_ = godotenv.Load()

	return &Config{
		Host:   getEnv("HOST", "0.0.0.0"),
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", "./data/todos.db"),
		LogLevel: func() slog.Level {
			switch os.Getenv("LOG_LEVEL") {
			case "DEBUG":
				return slog.LevelDebug
			case "INFO":
				return slog.LevelInfo
			case "WARN":
				return slog.LevelWarn
			case "ERROR":
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		}(),
		SessionSecret: getEnv("SESSION_SECRET", "session-secret"),
	}
}
