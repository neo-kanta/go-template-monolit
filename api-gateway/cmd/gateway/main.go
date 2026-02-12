// Package main is the entry point for the API Gateway.
//
//	@title						Transfer Agent API Gateway
//	@version					0.1.0
//	@description				POC API Gateway with JWT authentication for the Transfer Agent platform.
//	@host						localhost:8080
//	@BasePath					/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token: **Bearer &lt;token&gt;**
package main

import (
	"log"

	"go-transfer-agent/api-gateway/internal/config"
	apphttp "go-transfer-agent/api-gateway/internal/http"

	// blank import so the generated docs are linked into the binary
	_ "go-transfer-agent/api-gateway/docs"
)

func main() {
	cfg := config.Load()

	router := apphttp.NewRouter(cfg)

	log.Printf("🚀 API Gateway starting on :%s", cfg.Port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
