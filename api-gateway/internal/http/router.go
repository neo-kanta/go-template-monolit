package http

import (
	"go-transfer-agent/api-gateway/internal/config"
	"go-transfer-agent/api-gateway/internal/http/handlers"
	"go-transfer-agent/api-gateway/internal/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates and configures the Gin engine with all routes.
func NewRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.Default()

	// ---------- Swagger UI ----------
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ---------- API v1 ----------
	v1 := r.Group("/api/v1")

	// Public routes
	v1.GET("/health", handlers.Health)

	auth := handlers.NewAuthHandler(cfg)
	v1.POST("/auth/login", auth.Login)

	// Protected routes
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		protected.GET("/auth/me", auth.Me)
		protected.POST("/fnd/transactions", handlers.CreateTransaction)
	}

	return r
}
