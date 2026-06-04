package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/config"
	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/service"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// fakeInventoryService implements InventoryService backed by an in-memory map,
// so the HTTP handlers can be tested without a database.
type fakeInventoryService struct {
	items map[string]*models.Inventory
}

func newFake() *fakeInventoryService {
	return &fakeInventoryService{items: map[string]*models.Inventory{}}
}

func (f *fakeInventoryService) ListInventory(ctx context.Context) ([]models.Inventory, error) {
	out := make([]models.Inventory, 0, len(f.items))
	for _, v := range f.items {
		out = append(out, *v)
	}
	return out, nil
}

func (f *fakeInventoryService) GetInventory(ctx context.Context, id string) (*models.Inventory, error) {
	if v, ok := f.items[id]; ok {
		return v, nil
	}
	return nil, service.ErrNotFound
}

func (f *fakeInventoryService) CreateInventory(ctx context.Context, inv *models.Inventory) error {
	if inv.ID == "" {
		inv.ID = "generated-id"
	}
	f.items[inv.ID] = inv
	return nil
}

func (f *fakeInventoryService) UpdateInventory(ctx context.Context, id string, quantity int) error {
	v, ok := f.items[id]
	if !ok {
		return service.ErrNotFound
	}
	v.Quantity = quantity
	return nil
}

func (f *fakeInventoryService) DeleteInventory(ctx context.Context, id string) error {
	if _, ok := f.items[id]; !ok {
		return service.ErrNotFound
	}
	delete(f.items, id)
	return nil
}

func newTestServer(svc InventoryService) *Server {
	return NewServer(&config.Config{}, svc)
}

func TestListInventory(t *testing.T) {
	fake := newFake()
	fake.items["a"] = &models.Inventory{ID: "a", Quantity: 5}
	srv := newTestServer(fake)

	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/inventory", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var got []models.Inventory
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestGetInventoryItemNotFound(t *testing.T) {
	srv := newTestServer(newFake())

	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/inventory/missing", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestCreateInventory(t *testing.T) {
	fake := newFake()
	srv := newTestServer(fake)

	body := strings.NewReader(`{"id":"x","product_id":"p1","quantity":10}`)
	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/inventory", body))

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if _, ok := fake.items["x"]; !ok {
		t.Fatal("item was not stored by the service")
	}
}

func TestDeleteInventory(t *testing.T) {
	fake := newFake()
	fake.items["x"] = &models.Inventory{ID: "x"}
	srv := newTestServer(fake)

	rec := httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/inventory/x", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}

	rec = httptest.NewRecorder()
	srv.Router().ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/inventory/x", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("second delete status = %d, want 404", rec.Code)
	}
}
