package shared

import "context"

// ctxKey is the typed context key for JWT claims to avoid string collisions.
type ctxKey string

// CtxKeyClaims is the context key used to store authenticated JWT claims.
const CtxKeyClaims ctxKey = "claims"

// AuthClaims holds the authenticated user identity extracted from JWT.
type AuthClaims struct {
	Username string
	Role     string
}

// GetUsernameFromCtx extracts the authenticated username from the gRPC context.
// Returns empty string if no claims are present (e.g., in tests without auth).
func GetUsernameFromCtx(ctx context.Context) string {
	claims, ok := ctx.Value(CtxKeyClaims).(*AuthClaims)
	if !ok || claims == nil {
		return ""
	}
	return claims.Username
}

// GetRoleFromCtx extracts the authenticated user role from the gRPC context.
func GetRoleFromCtx(ctx context.Context) string {
	claims, ok := ctx.Value(CtxKeyClaims).(*AuthClaims)
	if !ok || claims == nil {
		return ""
	}
	return claims.Role
}

// NewAuthContext creates a context with the given AuthClaims.
// Useful for tests and when injecting claims from the auth interceptor.
func NewAuthContext(ctx context.Context, username, role string) context.Context {
	return context.WithValue(ctx, CtxKeyClaims, &AuthClaims{
		Username: username,
		Role:     role,
	})
}
