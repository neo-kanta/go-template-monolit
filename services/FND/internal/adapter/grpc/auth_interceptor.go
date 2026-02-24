package grpc

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserClaims holds the JWT claims. Matches the format issued by the API Gateway.
type UserClaims struct {
	Sub  string `json:"sub"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// AuthInterceptor creates a gRPC UnaryServerInterceptor that validates JWTs.
func AuthInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Optional: allow certain health check methods to bypass JWT auth
		if info.FullMethod == "/grpc.health.v1.Health/Check" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md["authorization"]
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		authHeader := values[0]
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization format, expected: Bearer <token>")
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
			jwt.WithValidMethods([]string{"HS256"}),
			// Note: not enforcing strictly WithExpirationRequired() here just in case API gateway does it or not
		)

		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token: %v", err)
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok || !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token claims")
		}

		// Inject the authenticated claims into the context for down-stream handlers
		newCtx := context.WithValue(ctx, "claims", claims)

		return handler(newCtx, req)
	}
}
