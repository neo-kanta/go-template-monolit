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

	// Connect to Sample gRPC service
	sampleClient := grpcclient.NewSampleClient(cfg.SampleServiceAddr)
	defer sampleClient.Close()

	router := apphttp.NewRouter(cfg, sampleClient)

	log.Printf("🚀 API Gateway starting on :%s", cfg.Port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)
	log.Printf("🔗 Sample gRPC: %s", cfg.SampleServiceAddr)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
