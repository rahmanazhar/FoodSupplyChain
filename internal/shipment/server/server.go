package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/config"
	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/service"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	service *service.ShipmentService
	router  *mux.Router
}

// NewServer creates a new HTTP server instance
func NewServer(cfg *config.Config, svc *service.ShipmentService) *Server {
	s := &Server{
		config:  cfg,
		service: svc,
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

// Router returns the configured router
func (s *Server) Router() *mux.Router {
	return s.router
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	// Apply common middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.metricsMiddleware)

	// Health check endpoint
	s.router.HandleFunc("/health", s.healthCheckHandler).Methods(http.MethodGet)

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.Use(s.authMiddleware)

	// Shipments
	api.HandleFunc("/shipments", s.handleGetShipments).Methods(http.MethodGet)
	api.HandleFunc("/shipments", s.handleCreateShipment).Methods(http.MethodPost)
	api.HandleFunc("/shipments/{id}", s.handleGetShipment).Methods(http.MethodGet)
	api.HandleFunc("/shipments/{id}", s.handleUpdateShipment).Methods(http.MethodPut)
	api.HandleFunc("/shipments/{id}", s.handleDeleteShipment).Methods(http.MethodDelete)
	api.HandleFunc("/shipments/{id}/status", s.handleUpdateShipmentStatus).Methods(http.MethodPut)
	api.HandleFunc("/shipments/{id}/track", s.handleTrackShipment).Methods(http.MethodGet)
}

// Middleware functions

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		// Log request details
		duration := time.Since(start)
		// TODO: Use proper logging framework
		_ = duration
	})
}

func (s *Server) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		// TODO: Record metrics
		_ = duration
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement proper authentication
		next.ServeHTTP(w, r)
	})
}

// Health check handler
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	}{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

// Handler stubs - to be implemented
func (s *Server) handleGetShipments(w http.ResponseWriter, r *http.Request)         {}
func (s *Server) handleCreateShipment(w http.ResponseWriter, r *http.Request)       {}
func (s *Server) handleGetShipment(w http.ResponseWriter, r *http.Request)          {}
func (s *Server) handleUpdateShipment(w http.ResponseWriter, r *http.Request)       {}
func (s *Server) handleDeleteShipment(w http.ResponseWriter, r *http.Request)       {}
func (s *Server) handleUpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {}
func (s *Server) handleTrackShipment(w http.ResponseWriter, r *http.Request)        {}
