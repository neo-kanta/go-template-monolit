// Package main is the entry point for the API Gateway.
//
//	@title						Transfer Agent API Gateway
//	@version					0.1.0
//	@description				POC API Gateway with JWT authentication. Calls FND service via gRPC.
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
	"go-transfer-agent/api-gateway/internal/grpcclient"
	apphttp "go-transfer-agent/api-gateway/internal/http"

	_ "go-transfer-agent/api-gateway/docs"
)

func main() {
	cfg := config.Load()

	// Connect to FND gRPC service
	fndClient := grpcclient.NewFNDClient(cfg.FNDServiceAddr)
	defer fndClient.Close()

	router := apphttp.NewRouter(cfg, fndClient)

	log.Printf("🚀 API Gateway starting on :%s", cfg.Port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Printf("🔗 FND gRPC: %s", cfg.FNDServiceAddr)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
