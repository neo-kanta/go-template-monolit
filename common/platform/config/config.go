package config

import "os"

// Config holds shared platform configuration.
type Config struct {
	GRPCAddr    string
	HTTPAddr    string
	DBDsn       string
	LogLevel    string
	JWTSecret   string
	JWTIssuer   string // Matches C# Identity.Issuer (ValidateIssuer)
	JWTAudience string // Matches C# Identity.Audience (ValidateAudience)
}

// Load reads platform config from environment variables with defaults.
func Load() *Config {
	return &Config{
		GRPCAddr:    getEnv("GRPC_ADDR", "0.0.0.0:50051"),
		HTTPAddr:    getEnv("HTTP_ADDR", "0.0.0.0:8081"),
		DBDsn:       getEnv("DB_DSN", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-poc-key-change-me"),
		JWTIssuer:   getEnv("JWT_ISSUER", ""),   // empty = skip issuer validation
		JWTAudience: getEnv("JWT_AUDIENCE", ""), // empty = skip audience validation
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
