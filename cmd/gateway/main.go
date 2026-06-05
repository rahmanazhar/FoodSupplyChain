package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/auth"
)

type Config struct {
	Port             int
	InventoryService string
	ShipmentService  string
	JWTSecret        string
	TokenTTL         time.Duration
	ReadTimeout      time.Duration
	WriteTimeout     time.Duration
	IdleTimeout      time.Duration
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := &Config{
		Port:             3000,
		InventoryService: getEnv("INVENTORY_SERVICE_URL", "http://localhost:8080"),
		ShipmentService:  getEnv("SHIPMENT_SERVICE_URL", "http://localhost:8081"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-here"),
		TokenTTL:         parseDuration(getEnv("TOKEN_TTL", "1h"), time.Hour),
		ReadTimeout:      5 * time.Second,
		WriteTimeout:     10 * time.Second,
		IdleTimeout:      120 * time.Second,
	}

	authManager := auth.NewManager(cfg.JWTSecret, cfg.TokenTTL)

	router := mux.NewRouter()
	router.Use(corsMiddleware)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet, http.MethodOptions)

	// Issues a JWT for the selected role. This is a development-grade sign-in
	// (no password store); the gateway and shipment service share JWT_SECRET.
	router.HandleFunc("/auth/login", loginHandler(authManager, cfg.TokenTTL)).Methods(http.MethodPost, http.MethodOptions)

	setupProxyRoutes(router, cfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Printf("Starting gateway server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func setupProxyRoutes(router *mux.Router, cfg *Config) {
	inventoryURL, err := url.Parse(cfg.InventoryService)
	if err != nil {
		log.Fatalf("Invalid inventory service URL: %v", err)
	}
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}

	// The inventory service owns inventory, products and locations.
	inventoryProxy := createReverseProxy(inventoryURL, "")
	router.PathPrefix("/inventory").Handler(inventoryProxy).Methods(methods...)
	router.PathPrefix("/products").Handler(inventoryProxy).Methods(methods...)
	router.PathPrefix("/locations").Handler(inventoryProxy).Methods(methods...)

	shipmentURL, err := url.Parse(cfg.ShipmentService)
	if err != nil {
		log.Fatalf("Invalid shipment service URL: %v", err)
	}
	// The shipment service serves under /api/v1, so /shipments/* is rewritten
	// to /api/v1/shipments/* before being forwarded.
	router.PathPrefix("/shipments").Handler(createReverseProxy(shipmentURL, "/api/v1")).Methods(methods...)
}

// createReverseProxy builds a reverse proxy to target. If pathPrefix is set, it
// is prepended to the request path before forwarding (used to map the gateway's
// public paths onto a service's internal route prefix).
func createReverseProxy(target *url.URL, pathPrefix string) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		if pathPrefix != "" {
			req.URL.Path = pathPrefix + req.URL.Path
		}
		originalDirector(req)
		req.Header.Set("X-Proxy-Gateway", "api-gateway")
	}

	// The gateway is the single source of CORS headers; strip any set by the
	// backend so responses don't carry duplicate Access-Control-* headers
	// (which browsers reject).
	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Allow-Headers")
		resp.Header.Del("Access-Control-Allow-Credentials")
		return nil
	}

	return proxy
}

// loginHandler issues a JWT for a requested role. Roles are validated against
// the known set; the subject defaults to the role name when no username is given.
func loginHandler(manager *auth.Manager, ttl time.Duration) http.HandlerFunc {
	allowed := map[string]bool{
		auth.RoleAdmin:    true,
		auth.RoleManager:  true,
		auth.RoleOperator: true,
		auth.RoleViewer:   true,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Username string `json:"username"`
			Role     string `json:"role"`
			Tenant   string `json:"tenant"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}
		if body.Role == "" {
			body.Role = auth.RoleViewer
		}
		if !allowed[body.Role] {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "unknown role"})
			return
		}
		subject := body.Username
		if subject == "" {
			subject = body.Role
		}

		token, err := manager.GenerateToken(subject, body.Role, body.Tenant)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"token":      token,
			"role":       body.Role,
			"username":   subject,
			"expires_in": int(ttl.Seconds()),
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func parseDuration(value string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(value); err == nil {
		return d
	}
	return fallback
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
