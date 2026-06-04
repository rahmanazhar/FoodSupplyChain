package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/events"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// ErrNotFound is returned when a requested record does not exist. Handlers map
// it to an HTTP 404 response.
var ErrNotFound = errors.New("record not found")

// ListInventory returns all inventory records.
func (s *InventoryService) ListInventory(ctx context.Context) ([]models.Inventory, error) {
	var items []models.Inventory
	if err := s.db.WithContext(ctx).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list inventory: %w", err)
	}
	return items, nil
}

// GetInventory returns a single inventory record by ID, or ErrNotFound.
func (s *InventoryService) GetInventory(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory
	if err := s.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	return &item, nil
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

// DeleteInventory removes an inventory record by ID, or returns ErrNotFound.
func (s *InventoryService) DeleteInventory(ctx context.Context, id string) error {
	result := s.db.WithContext(ctx).Delete(&models.Inventory{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete inventory: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
