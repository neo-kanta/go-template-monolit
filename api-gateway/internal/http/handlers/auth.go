package handlers

import (
	"net/http"
	"time"

	"go-transfer-agent/api-gateway/internal/config"
	"go-transfer-agent/api-gateway/internal/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// --- request / response models ---

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
}

// LoginResponse is returned on successful authentication.
type LoginResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	TokenType   string `json:"token_type" example:"bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

// MeResponse returns the authenticated user's claims.
type MeResponse struct {
	Subject   string `json:"sub" example:"admin"`
	Role      string `json:"role" example:"admin"`
	ExpiresAt string `json:"expires_at" example:"2026-02-12T12:00:00Z"`
}

// ErrorResponse is a generic error envelope.
type ErrorResponse struct {
	Error string `json:"error" example:"invalid credentials"`
}

// --- hardcoded POC users ---

var pocUsers = map[string]struct {
	Password string
	Role     string
}{
	"admin": {Password: "admin123", Role: "admin"},
	"user":  {Password: "user123", Role: "user"},
}

// --- handlers ---

// AuthHandler holds dependencies for auth-related endpoints.
type AuthHandler struct {
	Cfg *config.Config
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{Cfg: cfg}
}

// Login godoc
//
//	@Summary		Authenticate user
//	@Description	Validates credentials and returns a JWT access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Router			/api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "username and password are required"})
		return
	}

	u, ok := pocUsers[req.Username]
	if !ok || u.Password != req.Password {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	expiry := time.Duration(h.Cfg.JWTExpiryMinutes) * time.Minute
	now := time.Now()

	claims := middleware.UserClaims{
		Sub:  req.Username,
		Role: u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   req.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.Cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken: signed,
		TokenType:   "bearer",
		ExpiresIn:   int(expiry.Seconds()),
	})
}

// Me godoc
//
//	@Summary		Current user info
//	@Description	Returns the claims of the authenticated user
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	MeResponse
//	@Failure		401	{object}	ErrorResponse
//	@Security		BearerAuth
//	@Router			/api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	raw, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "no claims in context"})
		return
	}

	claims, ok := raw.(*middleware.UserClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid claims type"})
		return
	}

	expiresAt := ""
	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, MeResponse{
		Subject:   claims.Sub,
		Role:      claims.Role,
		ExpiresAt: expiresAt,
	})
}
