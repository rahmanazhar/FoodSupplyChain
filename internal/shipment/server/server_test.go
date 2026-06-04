package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/config"
	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/service"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/auth"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

const testSecret = "test-secret"

// fakeShipmentService implements ShipmentService backed by an in-memory map.
type fakeShipmentService struct {
	items map[string]*models.Shipment
}

func newFake() *fakeShipmentService {
	return &fakeShipmentService{items: map[string]*models.Shipment{}}
}

func (f *fakeShipmentService) ListShipments(ctx context.Context) ([]models.Shipment, error) {
	out := make([]models.Shipment, 0, len(f.items))
	for _, v := range f.items {
		out = append(out, *v)
	}
	return out, nil
}

func (f *fakeShipmentService) GetShipment(ctx context.Context, id string) (*models.Shipment, error) {
	if v, ok := f.items[id]; ok {
		return v, nil
	}
	return nil, errNotFound()
}

func (f *fakeShipmentService) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	if shipment.ID == "" {
		shipment.ID = "generated-id"
	}
	f.items[shipment.ID] = shipment
	return nil
}

func (f *fakeShipmentService) UpdateShipment(ctx context.Context, id string, update *models.Shipment) (*models.Shipment, error) {
	v, ok := f.items[id]
	if !ok {
		return nil, errNotFound()
	}
	if update.Status != "" {
		v.Status = update.Status
	}
	return v, nil
}

func (f *fakeShipmentService) DeleteShipment(ctx context.Context, id string) error {
	if _, ok := f.items[id]; !ok {
		return errNotFound()
	}
	delete(f.items, id)
	return nil
}

func (f *fakeShipmentService) UpdateShipmentStatus(ctx context.Context, id, status, location string) error {
	v, ok := f.items[id]
	if !ok {
		return errNotFound()
	}
	v.Status = status
	return nil
}

func (f *fakeShipmentService) ListShipmentEvents(ctx context.Context, id string) ([]models.ShipmentEvent, error) {
	return []models.ShipmentEvent{}, nil
}

func errNotFound() error { return service.ErrNotFound }

func newTestServer(svc ShipmentService) *Server {
	cfg := &config.Config{}
	cfg.Auth.JWTSecret = testSecret
	cfg.Auth.TokenExpiry = time.Hour
	return NewServer(cfg, svc)
}

func tokenFor(t *testing.T, role string) string {
	t.Helper()
	m := auth.NewManager(testSecret, time.Hour)
	tok, err := m.GenerateToken("user-1", role, "tenant-1")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	return tok
}

func TestHealthIsPublic(t *testing.T) {
	srv := newTestServer(newFake())
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestListShipmentsRequiresAuth(t *testing.T) {
	srv := newTestServer(newFake())
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/shipments", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestListShipmentsWithToken(t *testing.T) {
	srv := newTestServer(newFake())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/shipments", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, auth.RoleViewer))
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}

func TestCreateShipment(t *testing.T) {
	fake := newFake()
	srv := newTestServer(fake)
	body := strings.NewReader(`{"id":"s1","order_id":"o1","status":"pending","origin":"A","destination":"B"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shipments", body)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, auth.RoleOperator))
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if _, ok := fake.items["s1"]; !ok {
		t.Fatal("shipment was not stored by the service")
	}
}

func TestGetShipmentNotFound(t *testing.T) {
	srv := newTestServer(newFake())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/shipments/missing", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, auth.RoleViewer))
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestDeleteRequiresElevatedRole(t *testing.T) {
	fake := newFake()
	fake.items["s1"] = &models.Shipment{ID: "s1"}
	srv := newTestServer(fake)

	// A viewer is authenticated but lacks permission to delete.
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/shipments/s1", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, auth.RoleViewer))
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("viewer delete status = %d, want 403", rec.Code)
	}

	// An admin can delete.
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/shipments/s1", nil)
	req.Header.Set("Authorization", "Bearer "+tokenFor(t, auth.RoleAdmin))
	rec = httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("admin delete status = %d, want 204", rec.Code)
	}
}
