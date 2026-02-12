package config

import "os"

// Config holds shared platform configuration.
type Config struct {
	GRPCAddr string
	DBDsn    string
	LogLevel string
}

// Load reads platform config from environment variables with defaults.
func Load() *Config {
	return &Config{
		GRPCAddr: getEnv("GRPC_ADDR", "0.0.0.0:50051"),
		DBDsn:    getEnv("DB_DSN", ""),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
