// Package gateway provides the API gateway's user authentication: a
// database-backed user store with bcrypt password hashing and JWT issuance.
package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/auth"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

var (
	// ErrInvalidCredentials is returned when a login attempt fails.
	ErrInvalidCredentials = errors.New("invalid username or password")
	// ErrUserExists is returned when a username or email is already taken.
	ErrUserExists = errors.New("username or email already taken")
)

// Auth provides user registration/authentication and issues JWTs.
type Auth struct {
	db     *gorm.DB
	tokens *auth.Manager
}

// NewAuth connects to the database, migrates the users table and seeds a demo
// account per role on first run so the app is usable out of the box.
func NewAuth(dsn string, tokens *auth.Manager) (*Auth, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("auth: failed to connect to database: %w", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("auth: failed to migrate users: %w", err)
	}
	a := &Auth{db: db, tokens: tokens}
	if err := a.seedDemoUsers(); err != nil {
		return nil, err
	}
	return a, nil
}

// RegisterRoutes wires the authentication endpoints onto the router.
func (a *Auth) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", a.handleRegister).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth/login", a.handleLogin).Methods(http.MethodPost, http.MethodOptions)
	router.Handle("/auth/me", a.tokens.Middleware(http.HandlerFunc(a.handleMe))).
		Methods(http.MethodGet, http.MethodOptions)
}

// seedDemoUsers creates one known account per role (e.g. admin/admin123) if it
// does not already exist.
func (a *Auth) seedDemoUsers() error {
	demo := []struct{ username, role string }{
		{"admin", auth.RoleAdmin},
		{"manager", auth.RoleManager},
		{"operator", auth.RoleOperator},
		{"viewer", auth.RoleViewer},
	}
	for _, d := range demo {
		var count int64
		a.db.Model(&models.User{}).Where("username = ?", d.username).Count(&count)
		if count > 0 {
			continue
		}
		if _, err := a.createUser(d.username, d.username+"@example.com", d.username+"123", d.role); err != nil {
			return fmt.Errorf("auth: seed %s: %w", d.username, err)
		}
	}
	return nil
}

func (a *Auth) createUser(username, email, password, role string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := a.db.Create(user).Error; err != nil {
		return nil, ErrUserExists
	}
	return user, nil
}

func (a *Auth) authenticate(username, password string) (*models.User, error) {
	var user models.User
	if err := a.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return &user, nil
}

// Handlers

func (a *Auth) handleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody("invalid request body"))
		return
	}
	body.Username = strings.TrimSpace(body.Username)
	if len(body.Username) < 3 {
		writeJSON(w, http.StatusBadRequest, errBody("username must be at least 3 characters"))
		return
	}
	if len(body.Password) < 6 {
		writeJSON(w, http.StatusBadRequest, errBody("password must be at least 6 characters"))
		return
	}
	// Self-registered users get the least-privileged role.
	user, err := a.createUser(body.Username, body.Email, body.Password, auth.RoleViewer)
	if err != nil {
		writeJSON(w, http.StatusConflict, errBody(err.Error()))
		return
	}
	a.issueToken(w, user)
}

func (a *Auth) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errBody("invalid request body"))
		return
	}
	user, err := a.authenticate(strings.TrimSpace(body.Username), body.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, errBody("invalid username or password"))
		return
	}
	a.issueToken(w, user)
}

func (a *Auth) handleMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, errBody("unauthenticated"))
		return
	}
	var user models.User
	if err := a.db.First(&user, "username = ?", claims.Subject).Error; err != nil {
		writeJSON(w, http.StatusOK, map[string]string{"username": claims.Subject, "role": claims.Role})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

// issueToken signs a JWT for the user (subject = username) and returns it with
// the user record.
func (a *Auth) issueToken(w http.ResponseWriter, user *models.User) {
	token, err := a.tokens.GenerateToken(user.Username, user.Role, "")
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errBody(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

func errBody(message string) map[string]string { return map[string]string{"error": message} }

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
