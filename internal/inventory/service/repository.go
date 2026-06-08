package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/events"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// ErrNotFound is returned when a requested record does not exist. Handlers map
// it to an HTTP 404 response.
var ErrNotFound = errors.New("record not found")

// ListInventory returns a page of inventory records (with their product and
// location preloaded) and the total number of matching records. When search is
// non-empty it filters by product name or SKU, case-insensitively.
func (s *InventoryService) ListInventory(ctx context.Context, limit, offset int, search string) ([]models.Inventory, int, error) {
	// Base query joins Product so the search predicate and ordering can use its
	// columns; Product/Location are still preloaded into the returned items.
	base := s.db.WithContext(ctx).Model(&models.Inventory{}).
		Joins("LEFT JOIN products ON products.id = inventories.product_id")
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		base = base.Where("LOWER(products.name) LIKE ? OR LOWER(products.sku) LIKE ?", like, like)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count inventory: %w", err)
	}

	var items []models.Inventory
	if err := base.
		Preload("Product").
		Preload("Location").
		Order("products.name asc").
		Limit(limit).
		Offset(offset).
		Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list inventory: %w", err)
	}
	return items, int(total), nil
}

// GetInventory returns a single inventory record by ID, or ErrNotFound.
func (s *InventoryService) GetInventory(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory
	if err := s.db.WithContext(ctx).Preload("Product").Preload("Location").First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return &item, nil
}

// ListProducts returns all products.
func (s *InventoryService) ListProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	if err := s.db.WithContext(ctx).Order("name asc").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

// DeleteProduct removes a product by ID, or returns ErrNotFound.
func (s *InventoryService) DeleteProduct(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ListLocations returns all locations.
func (s *InventoryService) ListLocations(ctx context.Context) ([]models.Location, error) {
	var locations []models.Location
	if err := s.db.WithContext(ctx).Order("name asc").Find(&locations).Error; err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}
	return locations, nil
}

// CreateLocation persists a new location.
func (s *InventoryService) CreateLocation(ctx context.Context, location *models.Location) error {
	if location.ID == "" {
		location.ID = uuid.New().String()
	}
	now := time.Now()
	location.CreatedAt = now
	location.UpdatedAt = now
	if err := s.db.WithContext(ctx).Create(location).Error; err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}
	return nil
}

// DeleteLocation removes a location by ID, or returns ErrNotFound.
func (s *InventoryService) DeleteLocation(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Delete(&models.Location{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete location: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateInventory persists a new inventory record and emits a created event.
// Associations (Product, Location) are referenced by ID and not upserted here.
func (s *InventoryService) CreateInventory(ctx context.Context, inv *models.Inventory) error {
	if inv.ID == "" {
		inv.ID = uuid.New().String()
	}
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now

	if err := s.db.WithContext(ctx).Omit("Product", "Location").Create(inv).Error; err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}

	event := &events.InventoryEvent{
		BaseEvent: events.BaseEvent{
			ID:        uuid.New().String(),
			Type:      string(events.InventoryCreated),
			Timestamp: now,
			Version:   "1.0",
			Source:    s.config.App.Name,
		},
	}
	event.Data.InventoryID = inv.ID
	event.Data.ProductID = inv.ProductID
	event.Data.LocationID = inv.LocationID
	event.Data.Quantity = inv.Quantity

	if err := s.publishEvent(fmt.Sprintf("%s.inventory.created", s.config.NATS.SubjectPrefix), event); err != nil {
		return fmt.Errorf("failed to publish inventory created event: %w", err)
	}
	return nil
}

// DeleteInventory removes an inventory record and its dependent transactions
// and alerts in a single transaction, or returns ErrNotFound.
func (s *InventoryService) DeleteInventory(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("inventory_id = ?", id).Delete(&models.InventoryTransaction{}).Error; err != nil {
			return fmt.Errorf("failed to delete inventory transactions: %w", err)
		}
		if err := tx.Where("inventory_id = ?", id).Delete(&models.InventoryAlert{}).Error; err != nil {
			return fmt.Errorf("failed to delete inventory alerts: %w", err)
		}
		result := tx.Delete(&models.Inventory{}, "id = ?", id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete inventory: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}
