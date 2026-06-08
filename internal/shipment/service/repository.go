package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// ErrNotFound is returned when a requested record does not exist. Handlers map
// it to an HTTP 404 response.
var ErrNotFound = errors.New("record not found")

// ListShipments returns a page of shipments and the total number of matching
// records. When search is non-empty it filters by order_id/origin/destination/id
// (case-insensitively); when status is non-empty it filters by exact status.
func (s *ShipmentService) ListShipments(ctx context.Context, limit, offset int, search, status string) ([]models.Shipment, int, error) {
	base := s.db.WithContext(ctx).Model(&models.Shipment{})
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		base = base.Where(
			"LOWER(order_id) LIKE ? OR LOWER(origin) LIKE ? OR LOWER(destination) LIKE ? OR LOWER(id) LIKE ?",
			like, like, like, like,
		)
	}
	if status != "" {
		base = base.Where("status = ?", status)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count shipments: %w", err)
	}

	var shipments []models.Shipment
	if err := base.
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&shipments).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list shipments: %w", err)
	}
	return shipments, int(total), nil
}

// GetShipment returns a single shipment by ID, or ErrNotFound.
func (s *ShipmentService) GetShipment(ctx context.Context, id string) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := s.db.WithContext(ctx).First(&shipment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}
	return &shipment, nil
}

// UpdateShipment applies the non-empty fields of update to an existing shipment.
func (s *ShipmentService) UpdateShipment(ctx context.Context, id string, update *models.Shipment) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := s.db.WithContext(ctx).First(&shipment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	if update.Status != "" {
		shipment.Status = update.Status
	}
	if update.Origin != "" {
		shipment.Origin = update.Origin
	}
	if update.Destination != "" {
		shipment.Destination = update.Destination
	}
	if update.CarrierID != "" {
		shipment.CarrierID = update.CarrierID
	}
	if update.TrackingNumber != "" {
		shipment.TrackingNumber = update.TrackingNumber
	}
	if update.Notes != "" {
		shipment.Notes = update.Notes
	}
	if !update.EstimatedArrival.IsZero() {
		shipment.EstimatedArrival = update.EstimatedArrival
	}
	shipment.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(&shipment).Error; err != nil {
		return nil, fmt.Errorf("failed to update shipment: %w", err)
	}
	return &shipment, nil
}

// DeleteShipment removes a shipment and its dependent events and alerts in a
// single transaction, or returns ErrNotFound.
func (s *ShipmentService) DeleteShipment(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("shipment_id = ?", id).Delete(&models.ShipmentEvent{}).Error; err != nil {
			return fmt.Errorf("failed to delete shipment events: %w", err)
		}
		if err := tx.Where("shipment_id = ?", id).Delete(&models.ShipmentAlert{}).Error; err != nil {
			return fmt.Errorf("failed to delete shipment alerts: %w", err)
		}
		result := tx.Delete(&models.Shipment{}, "id = ?", id)
		if result.Error != nil {
			return fmt.Errorf("failed to delete shipment: %w", result.Error)
		}
		if result.RowsAffected == 0 {
			return ErrNotFound
		}
		return nil
	})
}

// ListShipmentEvents returns the lifecycle events for a shipment, oldest first.
func (s *ShipmentService) ListShipmentEvents(ctx context.Context, shipmentID string) ([]models.ShipmentEvent, error) {
	var shipmentEvents []models.ShipmentEvent
	if err := s.db.WithContext(ctx).
		Where("shipment_id = ?", shipmentID).
		Order("created_at asc").
		Find(&shipmentEvents).Error; err != nil {
		return nil, fmt.Errorf("failed to list shipment events: %w", err)
	}
	return shipmentEvents, nil
}
