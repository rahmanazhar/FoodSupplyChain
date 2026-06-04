// Package auth provides JWT issuance/validation and role-based access control
// middleware for the supply chain services. Tokens are signed with HMAC-SHA256
// using only the Go standard library, so the package adds no dependencies.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// Roles recognised by the access-control layer, from most to least privileged.
const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

// Sentinel errors returned by ValidateToken.
var (
	ErrInvalidToken = errors.New("auth: invalid token")
	ErrExpiredToken = errors.New("auth: token expired")
)

// Claims is the JWT payload describing an authenticated principal.
type Claims struct {
	Subject   string `json:"sub"`
	Role      string `json:"role"`
	TenantID  string `json:"tenant,omitempty"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// jwtHeader is the fixed JOSE header for HS256 tokens.
type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Manager issues and validates JWTs using a shared HMAC secret.
type Manager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

// NewManager creates a token manager. A non-positive ttl defaults to one hour.
func NewManager(secret string, ttl time.Duration) *Manager {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &Manager{secret: []byte(secret), ttl: ttl, now: time.Now}
}

func encode(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

func (m *Manager) sign(signingInput string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(signingInput))
	return encode(mac.Sum(nil))
}

// GenerateToken issues a signed JWT for the given subject, role and tenant.
func (m *Manager) GenerateToken(subject, role, tenantID string) (string, error) {
	if subject == "" {
		return "", errors.New("auth: subject is required")
	}
	now := m.now()
	claims := Claims{
		Subject:   subject,
		Role:      role,
		TenantID:  tenantID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.ttl).Unix(),
	}
	headerBytes, err := json.Marshal(jwtHeader{Alg: "HS256", Typ: "JWT"})
	if err != nil {
		return "", err
	}
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signingInput := encode(headerBytes) + "." + encode(claimsBytes)
	return signingInput + "." + m.sign(signingInput), nil
}

// ValidateToken verifies the signature and expiry of a JWT and returns its claims.
func (m *Manager) ValidateToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signingInput := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(m.sign(signingInput)), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt > 0 && m.now().Unix() >= claims.ExpiresAt {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}
