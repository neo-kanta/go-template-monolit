package config

import (
	"os"
	"strconv"
)

// Config holds API Gateway configuration.
type Config struct {
	Port             string
	JWTSecret        string
	JWTExpiryMinutes int
	GinMode           string
	SampleServiceAddr string // gRPC address of the Sample service
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "super-secret-poc-key-change-me"),
		JWTExpiryMinutes: getEnvInt("JWT_EXPIRY_MINUTES", 60),
		GinMode:           getEnv("GIN_MODE", "debug"),
		SampleServiceAddr: getEnv("SAMPLE_SERVICE_ADDR", "localhost:50051"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
