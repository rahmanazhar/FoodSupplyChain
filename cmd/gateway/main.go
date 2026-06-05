package main

import (
	"context"
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
)

type Config struct {
	Port             int
	InventoryService string
	ShipmentService  string
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
		ReadTimeout:      5 * time.Second,
		WriteTimeout:     10 * time.Second,
		IdleTimeout:      120 * time.Second,
	}

	router := mux.NewRouter()
	router.Use(corsMiddleware)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet, http.MethodOptions)

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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
