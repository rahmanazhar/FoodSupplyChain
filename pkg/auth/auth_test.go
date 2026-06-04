package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGenerateAndValidate(t *testing.T) {
	m := NewManager("test-secret", time.Hour)

	token, err := m.GenerateToken("user-1", RoleManager, "tenant-1")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := m.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Errorf("subject = %q, want user-1", claims.Subject)
	}
	if claims.Role != RoleManager {
		t.Errorf("role = %q, want %q", claims.Role, RoleManager)
	}
	if claims.TenantID != "tenant-1" {
		t.Errorf("tenant = %q, want tenant-1", claims.TenantID)
	}
}

func TestGenerateRequiresSubject(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	if _, err := m.GenerateToken("", RoleViewer, ""); err == nil {
		t.Fatal("expected error for empty subject")
	}
}

func TestExpiredToken(t *testing.T) {
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	m := NewManager("test-secret", time.Hour)
	m.now = func() time.Time { return base }

	token, err := m.GenerateToken("user-1", RoleViewer, "")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	m.now = func() time.Time { return base.Add(2 * time.Hour) }
	if _, err := m.ValidateToken(token); err != ErrExpiredToken {
		t.Fatalf("err = %v, want ErrExpiredToken", err)
	}
}

func TestTamperedToken(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	token, err := m.GenerateToken("user-1", RoleAdmin, "")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	// Validating with a different secret must fail signature verification.
	other := NewManager("other-secret", time.Hour)
	if _, err := other.ValidateToken(token); err != ErrInvalidToken {
		t.Fatalf("err = %v, want ErrInvalidToken", err)
	}

	if _, err := m.ValidateToken("not.a.jwt.token"); err != ErrInvalidToken {
		t.Fatalf("err = %v, want ErrInvalidToken for malformed token", err)
	}
}

func TestMiddlewareRejectsMissingToken(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestMiddlewareAcceptsValidToken(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	token, _ := m.GenerateToken("user-1", RoleViewer, "")

	var gotRole string
	handler := m.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if claims, ok := ClaimsFromContext(r.Context()); ok {
			gotRole = claims.Role
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if gotRole != RoleViewer {
		t.Fatalf("role from context = %q, want %q", gotRole, RoleViewer)
	}
}

func TestRequireRole(t *testing.T) {
	m := NewManager("test-secret", time.Hour)
	protected := m.Middleware(RequireRole(RoleAdmin)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	cases := []struct {
		role string
		want int
	}{
		{RoleAdmin, http.StatusOK},
		{RoleViewer, http.StatusForbidden},
	}
	for _, tc := range cases {
		token, _ := m.GenerateToken("user-1", tc.role, "")
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()
		protected.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Errorf("role %q: status = %d, want %d", tc.role, rec.Code, tc.want)
		}
	}
}
