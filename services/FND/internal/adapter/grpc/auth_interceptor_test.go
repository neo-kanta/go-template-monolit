package grpc

import (
	"context"
	"testing"
	"time"

	"go-transfer-agent/common/platform/config"
	"go-transfer-agent/services/fnd/shared"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const testSecret = "test-jwt-secret"

// testCfg returns a Config with only JWTSecret set (no Issuer/Audience validation).
func testCfg() *config.Config {
	return &config.Config{JWTSecret: testSecret}
}

// testCfgWithIssuerAudience returns a Config that enforces Issuer and Audience.
func testCfgWithIssuerAudience() *config.Config {
	return &config.Config{
		JWTSecret:   testSecret,
		JWTIssuer:   "https://flowapi.systemweb.co.th",
		JWTAudience: "ta-platform",
	}
}

// signToken creates a signed JWT with the given claims and secret.
func signToken(t *testing.T, claims *UserClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

// validClaims returns UserClaims that expire 1 hour from now.
func validClaims() *UserClaims {
	return &UserClaims{
		Sub:  "admin",
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "admin",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
}

// validClaimsWithIssuerAudience returns claims that include Issuer and Audience.
func validClaimsWithIssuerAudience() *UserClaims {
	return &UserClaims{
		Sub:  "admin",
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "admin",
			Issuer:    "https://flowapi.systemweb.co.th",
			Audience:  jwt.ClaimStrings{"ta-platform"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
}

// expiredClaims returns UserClaims that expired 1 hour ago.
func expiredClaims() *UserClaims {
	return &UserClaims{
		Sub:  "admin",
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "admin",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
}

// noopHandler is a trivial gRPC handler used in tests.
var noopHandler grpc.UnaryHandler = func(ctx context.Context, req interface{}) (interface{}, error) {
	return "ok", nil
}

// dummyInfo returns a UnaryServerInfo for a non-health-check method.
func dummyInfo() *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{FullMethod: "/fnd.v1.FNDService/QueryFundInfo"}
}

// healthInfo returns a UnaryServerInfo for the gRPC health-check method.
func healthInfo() *grpc.UnaryServerInfo {
	return &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}
}

func TestAuthInterceptor_BearerToken(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())
	token := signToken(t, validClaims(), testSecret)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got: %v", resp)
	}
}

func TestAuthInterceptor_RawToken(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())
	token := signToken(t, validClaims(), testSecret)

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err != nil {
		t.Fatalf("expected no error for raw token, got: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got: %v", resp)
	}
}

func TestAuthInterceptor_AuthorizeHeader(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())
	token := signToken(t, validClaims(), testSecret)

	// Use "authorize" header (lowercase) — matches C# fallback
	md := metadata.Pairs("authorize", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err != nil {
		t.Fatalf("expected no error for 'authorize' header, got: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got: %v", resp)
	}
}

func TestAuthInterceptor_MissingMetadata(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())

	_, err := interceptor(context.Background(), nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for missing metadata")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_MissingAuthorizationHeader(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())

	md := metadata.Pairs("other-header", "value")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for missing authorization header")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_InvalidToken(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())

	md := metadata.Pairs("authorization", "Bearer not-a-valid-jwt")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_ExpiredToken(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())
	token := signToken(t, expiredClaims(), testSecret)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_WrongSecret(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())
	token := signToken(t, validClaims(), "wrong-secret")

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_HealthCheckBypass(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())

	resp, err := interceptor(context.Background(), nil, healthInfo(), noopHandler)
	if err != nil {
		t.Fatalf("expected health check to bypass auth, got: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got: %v", resp)
	}
}

func TestAuthInterceptor_ClaimsInContext(t *testing.T) {
	interceptor := AuthInterceptor(testCfg())
	token := signToken(t, validClaims(), testSecret)

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		username := shared.GetUsernameFromCtx(ctx)
		if username == "" {
			t.Fatal("expected username in context")
		}
		if username != "admin" {
			t.Fatalf("expected username='admin', got: %s", username)
		}
		role := shared.GetRoleFromCtx(ctx)
		if role != "admin" {
			t.Fatalf("expected role='admin', got: %s", role)
		}
		return "ok", nil
	}

	_, err := interceptor(ctx, nil, dummyInfo(), handler)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

// --- Issuer/Audience validation tests (C# compatibility) ---

func TestAuthInterceptor_IssuerAudience_Valid(t *testing.T) {
	interceptor := AuthInterceptor(testCfgWithIssuerAudience())
	token := signToken(t, validClaimsWithIssuerAudience(), testSecret)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	resp, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err != nil {
		t.Fatalf("expected no error with valid issuer/audience, got: %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got: %v", resp)
	}
}

func TestAuthInterceptor_IssuerAudience_WrongIssuer(t *testing.T) {
	cfg := testCfgWithIssuerAudience()
	interceptor := AuthInterceptor(cfg)

	claims := validClaimsWithIssuerAudience()
	claims.Issuer = "https://wrong-issuer.com"
	token := signToken(t, claims, testSecret)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for wrong issuer")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_IssuerAudience_WrongAudience(t *testing.T) {
	cfg := testCfgWithIssuerAudience()
	interceptor := AuthInterceptor(cfg)

	claims := validClaimsWithIssuerAudience()
	claims.Audience = jwt.ClaimStrings{"wrong-audience"}
	token := signToken(t, claims, testSecret)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error for wrong audience")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}

func TestAuthInterceptor_IssuerAudience_MissingInToken(t *testing.T) {
	cfg := testCfgWithIssuerAudience()
	interceptor := AuthInterceptor(cfg)

	// Token has no issuer/audience but config requires them
	token := signToken(t, validClaims(), testSecret)

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	_, err := interceptor(ctx, nil, dummyInfo(), noopHandler)
	if err == nil {
		t.Fatal("expected error when token missing required issuer/audience")
	}
	if st, ok := status.FromError(err); !ok || st.Code() != codes.Unauthenticated {
		t.Fatalf("expected Unauthenticated, got: %v", err)
	}
}
