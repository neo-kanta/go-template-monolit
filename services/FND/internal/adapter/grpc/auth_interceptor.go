package grpc

import (
	"context"
	"strings"

	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/services/fnd/shared"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserClaims holds the JWT claims.
// Compatible with tokens issued by both the Go API Gateway and the C# FlowAPI.
type UserClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// AuthInterceptor creates a gRPC UnaryServerInterceptor that validates JWTs.
//
// It mirrors the C# JwtBearerEvents.OnMessageReceived logic:
//   - Checks both "authorization" and "authorize" headers (case-insensitive via gRPC metadata)
//   - Strips "Bearer " prefix if present, otherwise uses the raw token
//   - Validates Issuer/Audience if configured (matches C# ValidateIssuer/ValidateAudience)
func AuthInterceptor(cfg *config.Config) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Allow health check to bypass JWT auth
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		// Try "authorization" header first, then "authorize" (matches C# fallback behavior)
		values := md["authorization"]
		if len(values) == 0 {
			values = md["authorize"]
		}
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		// Accept both "Bearer {token}" and raw "{token}" formats
		authHeader := values[0]
		var tokenStr string
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			tokenStr = strings.TrimSpace(parts[1])
		} else {
			tokenStr = strings.TrimSpace(authHeader)
		}

		// Build parser options — matches C# TokenValidationParameters
		parserOpts := []jwt.ParserOption{
			jwt.WithValidMethods([]string{"HS256"}),
		}
		if cfg.JWTIssuer != "" {
			parserOpts = append(parserOpts, jwt.WithIssuer(cfg.JWTIssuer))
		}
		if cfg.JWTAudience != "" {
			parserOpts = append(parserOpts, jwt.WithAudience(cfg.JWTAudience))
		}

		token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		}, parserOpts...)

		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token: %v", err)
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token claims")
		}

		// Inject the authenticated claims into the context for downstream handlers
		newCtx := shared.NewAuthContext(ctx, claims.Sub, claims.Role)

		return handler(newCtx, req)
	}
}
