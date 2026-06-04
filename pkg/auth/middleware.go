package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey string

const claimsContextKey contextKey = "auth.claims"

// ClaimsFromContext returns the authenticated claims previously stored by
// Middleware, if the request was authenticated.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}

// bearerToken extracts a token from an "Authorization: Bearer <token>" header.
func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

// Middleware authenticates requests using a Bearer token. On success the
// validated claims are stored in the request context; otherwise it responds 401.
func (m *Manager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		claims, err := m.ValidateToken(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns middleware that permits a request only if its
// authenticated claims carry one of the supplied roles. It must be installed
// after Middleware so that claims are present in the context.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				writeError(w, http.StatusUnauthorized, "authentication required")
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				writeError(w, http.StatusForbidden, "insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
