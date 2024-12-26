package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/config"
	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/service"
)

type Server struct {
	config  *config.Config
	service *service.InventoryService
	router  *mux.Router
}

func NewServer(cfg *config.Config, svc *service.InventoryService) *Server {
	s := &Server{
		config:  cfg,
		service: svc,
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

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

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
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
		duration := time.Since(start)
		_ = duration
	})
}

func (s *Server) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		_ = duration
	})
}

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

func (s *Server) handleGetInventory(w http.ResponseWriter, r *http.Request) {
	mockData := []map[string]interface{}{
		{
			"id":       1,
			"name":     "Apples",
			"category": "fruits",
			"quantity": 150,
			"price":    1.99,
		},
		{
			"id":       2,
			"name":     "Bananas",
			"category": "fruits",
			"quantity": 200,
			"price":    0.99,
		},
		{
			"id":       3,
			"name":     "Carrots",
			"category": "vegetables",
			"quantity": 80,
			"price":    1.49,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockData)
}

func (s *Server) handleCreateInventory(w http.ResponseWriter, r *http.Request) {
	var item map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (s *Server) handleGetInventoryItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mockItem := map[string]interface{}{
		"id":       id,
		"name":     "Test Item",
		"category": "test",
		"quantity": 100,
		"price":    9.99,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockItem)
}

func (s *Server) handleUpdateInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item["id"] = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (s *Server) handleDeleteInventory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
