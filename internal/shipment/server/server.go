package server

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/config"
	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/service"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/auth"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/httpx"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/metrics"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// ShipmentService is the behaviour the HTTP layer requires from the service.
// Defining it here (the consumer) keeps the handlers testable with a fake.
type ShipmentService interface {
	ListShipments(ctx context.Context, limit, offset int, search, status string) ([]models.Shipment, int, error)
	GetShipment(ctx context.Context, id string) (*models.Shipment, error)
	CreateShipment(ctx context.Context, shipment *models.Shipment) error
	UpdateShipment(ctx context.Context, id string, update *models.Shipment) (*models.Shipment, error)
	DeleteShipment(ctx context.Context, id string) error
	UpdateShipmentStatus(ctx context.Context, id, status, location string) error
	ListShipmentEvents(ctx context.Context, id string) ([]models.ShipmentEvent, error)
}

// Server exposes the shipment service over HTTP.
type Server struct {
	config  *config.Config
	service ShipmentService
	auth    *auth.Manager
	router  *mux.Router
	logger  *slog.Logger
	metrics *metrics.Collector
}

// NewServer wires the routes and returns a ready-to-serve Server. When the
// configured JWT secret is non-empty the /api/v1 routes require authentication.
// A nil logger falls back to the slog default so tests need no setup.
func NewServer(cfg *config.Config, svc ShipmentService, logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.Default()
	}
	s := &Server{
		config:  cfg,
		service: svc,
		router:  mux.NewRouter(),
		logger:  logger,
		metrics: metrics.NewCollector(),
	}
	if cfg != nil && cfg.Auth.JWTSecret != "" {
		s.auth = auth.NewManager(cfg.Auth.JWTSecret, cfg.Auth.TokenExpiry)
	}
	s.setupRoutes()
	return s
}

// Router returns the configured router.
func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) setupRoutes() {
	// Shared, structured middleware replaces the previous ad-hoc logging and
	// atomic request counter.
	s.router.Use(httpx.RequestID)
	s.router.Use(httpx.Logger(s.logger))
	s.router.Use(httpx.Recoverer(s.logger))
	s.router.Use(httpx.SecurityHeaders)
	s.router.Use(s.metrics.Instrument)
	s.router.Use(s.corsMiddleware)

	s.router.HandleFunc("/health", s.healthCheckHandler).Methods(http.MethodGet, http.MethodOptions)
	s.router.Handle("/metrics", s.metrics.Handler()).Methods(http.MethodGet, http.MethodOptions)

	api := s.router.PathPrefix("/api/v1").Subrouter()
	if s.auth != nil {
		api.Use(s.auth.Middleware)
	}

	api.HandleFunc("/shipments", s.handleGetShipments).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/shipments", s.handleCreateShipment).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/shipments/{id}", s.handleGetShipment).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/shipments/{id}", s.handleUpdateShipment).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/shipments/{id}/status", s.handleUpdateShipmentStatus).Methods(http.MethodPut, http.MethodOptions)
	api.HandleFunc("/shipments/{id}/track", s.handleTrackShipment).Methods(http.MethodGet, http.MethodOptions)

	// Deletion is restricted to elevated roles when auth is enabled.
	if s.auth != nil {
		api.Handle("/shipments/{id}",
			auth.RequireRole(auth.RoleAdmin, auth.RoleManager)(http.HandlerFunc(s.handleDeleteShipment)),
		).Methods(http.MethodDelete, http.MethodOptions)
	} else {
		api.HandleFunc("/shipments/{id}", s.handleDeleteShipment).Methods(http.MethodDelete, http.MethodOptions)
	}
}

// Middleware

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Handlers

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleGetShipments(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, offset := httpx.ParsePagination(q)
	search := q.Get("search")
	status := q.Get("status")

	shipments, total, err := s.service.ListShipments(r.Context(), limit, offset, search, status)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, httpx.Page{Data: shipments, Total: total, Limit: limit, Offset: offset})
}

func (s *Server) handleCreateShipment(w http.ResponseWriter, r *http.Request) {
	var shipment models.Shipment
	if err := json.NewDecoder(r.Body).Decode(&shipment); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := s.service.CreateShipment(r.Context(), &shipment); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.writeJSON(w, http.StatusCreated, shipment)
}

func (s *Server) handleGetShipment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	shipment, err := s.service.GetShipment(r.Context(), id)
	if err != nil {
		s.writeServiceError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, shipment)
}

func (s *Server) handleUpdateShipment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var update models.Shipment
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	shipment, err := s.service.UpdateShipment(r.Context(), id, &update)
	if err != nil {
		s.writeServiceError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, shipment)
}

func (s *Server) handleDeleteShipment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := s.service.DeleteShipment(r.Context(), id); err != nil {
		s.writeServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpdateShipmentStatus(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		Status   string `json:"status"`
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if body.Status == "" {
		s.writeError(w, http.StatusBadRequest, "status is required")
		return
	}
	if err := s.service.UpdateShipmentStatus(r.Context(), id, body.Status, body.Location); err != nil {
		s.writeServiceError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"id":       id,
		"status":   body.Status,
		"location": body.Location,
	})
}

func (s *Server) handleTrackShipment(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	events, err := s.service.ListShipmentEvents(r.Context(), id)
	if err != nil {
		s.writeServiceError(w, err)
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"shipment_id": id,
		"events":      events,
	})
}

// Helpers

func (s *Server) writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrNotFound) {
		s.writeError(w, http.StatusNotFound, "shipment not found")
		return
	}
	s.writeError(w, http.StatusInternalServerError, err.Error())
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}
