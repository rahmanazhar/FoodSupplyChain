package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"

	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/config"
	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/service"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// InventoryService is the behaviour the HTTP layer requires from the service.
// Defining it here (the consumer) keeps the handlers testable with a fake.
type InventoryService interface {
	ListInventory(ctx context.Context) ([]models.Inventory, error)
	GetInventory(ctx context.Context, id string) (*models.Inventory, error)
	CreateInventory(ctx context.Context, inv *models.Inventory) error
	UpdateInventory(ctx context.Context, id string, quantity int) error
	DeleteInventory(ctx context.Context, id string) error
}

// Server exposes the inventory service over HTTP.
type Server struct {
	config       *config.Config
	service      InventoryService
	router       *mux.Router
	requestCount int64
}

// NewServer wires the routes and returns a ready-to-serve Server.
func NewServer(cfg *config.Config, svc InventoryService) *Server {
	s := &Server{
		config:  cfg,
		service: svc,
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

// Router returns the configured router.
func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) setupRoutes() {
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.metricsMiddleware)
	s.router.Use(s.corsMiddleware)

	s.router.HandleFunc("/health", s.healthCheckHandler).Methods(http.MethodGet, http.MethodOptions)

	s.router.HandleFunc("/inventory", s.handleGetInventory).Methods(http.MethodGet, http.MethodOptions)
	s.router.HandleFunc("/inventory", s.handleCreateInventory).Methods(http.MethodPost, http.MethodOptions)
	s.router.HandleFunc("/inventory/{id}", s.handleGetInventoryItem).Methods(http.MethodGet, http.MethodOptions)
	s.router.HandleFunc("/inventory/{id}", s.handleUpdateInventory).Methods(http.MethodPut, http.MethodOptions)
	s.router.HandleFunc("/inventory/{id}", s.handleDeleteInventory).Methods(http.MethodDelete, http.MethodOptions)
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

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("inventory %s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func (s *Server) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&s.requestCount, 1)
		next.ServeHTTP(w, r)
	})
}

// Handlers

func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":          "ok",
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
		"requests_served": atomic.LoadInt64(&s.requestCount),
	})
}

func (s *Server) handleGetInventory(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.ListInventory(r.Context())
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleCreateInventory(w http.ResponseWriter, r *http.Request) {
	var inv models.Inventory
	if err := json.NewDecoder(r.Body).Decode(&inv); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := s.service.CreateInventory(r.Context(), &inv); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.writeJSON(w, http.StatusCreated, inv)
}

func (s *Server) handleGetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	item, err := s.service.GetInventory(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			s.writeError(w, http.StatusNotFound, "inventory item not found")
			return
		}
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, item)
}

func (s *Server) handleUpdateInventory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := s.service.UpdateInventory(r.Context(), id, body.Quantity); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			s.writeError(w, http.StatusNotFound, "inventory item not found")
			return
		}
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{"id": id, "quantity": body.Quantity})
}

func (s *Server) handleDeleteInventory(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := s.service.DeleteInventory(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			s.writeError(w, http.StatusNotFound, "inventory item not found")
			return
		}
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Helpers

func (s *Server) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string) {
	s.writeJSON(w, status, map[string]string{"error": message})
}
