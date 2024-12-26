package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rahmanazhar/FoodSupplyChain/internal/inventory/config"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/events"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// InventoryService handles the core business logic for inventory management
type InventoryService struct {
	config *config.Config
	db     *gorm.DB
	nc     *nats.Conn
	js     nats.JetStreamContext
}

// NewInventoryService creates a new inventory service instance
func NewInventoryService(cfg *config.Config) (*InventoryService, error) {
	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Auto-migrate database schemas
	if err := db.AutoMigrate(
		&models.Product{},
		&models.Location{},
		&models.Inventory{},
		&models.InventoryTransaction{},
		&models.InventoryAlert{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	// Initialize NATS connection
	nc, err := nats.Connect(cfg.NATS.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %v", err)
	}

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to create JetStream context: %v", err)
	}

	return &InventoryService{
		config: cfg,
		db:     db,
		nc:     nc,
		js:     js,
	}, nil
}

// Close closes all connections
func (s *InventoryService) Close() error {
	if s.nc != nil {
		s.nc.Close()
	}

	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return fmt.Errorf("failed to get database instance: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %v", err)
		}
	}

	return nil
}

// CreateProduct creates a new product
func (s *InventoryService) CreateProduct(ctx context.Context, product *models.Product) error {
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	if err := s.db.Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %v", err)
	}

	// Publish event
	event := &events.InventoryEvent{
		BaseEvent: events.BaseEvent{
			ID:        uuid.New().String(),
			Type:      string(events.InventoryCreated),
			Timestamp: time.Now(),
			Version:   "1.0",
			Source:    s.config.App.Name,
		},
	}
	event.Data.ProductID = product.ID

	if err := s.publishEvent(fmt.Sprintf("%s.product.created", s.config.NATS.SubjectPrefix), event); err != nil {
		return fmt.Errorf("failed to publish product created event: %v", err)
	}

	return nil
}

// UpdateInventory updates inventory levels and generates alerts if needed
func (s *InventoryService) UpdateInventory(ctx context.Context, id string, quantity int) error {
	var inventory models.Inventory
	if err := s.db.First(&inventory, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to find inventory: %v", err)
	}

	prevQuantity := inventory.Quantity
	inventory.Quantity = quantity
	inventory.UpdatedAt = time.Now()

	if err := s.db.Save(&inventory).Error; err != nil {
		return fmt.Errorf("failed to update inventory: %v", err)
	}

	// Create transaction record
	transaction := &models.InventoryTransaction{
		ID:          uuid.New().String(),
		InventoryID: inventory.ID,
		Type:        "adjusted",
		Quantity:    quantity - prevQuantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(transaction).Error; err != nil {
		return fmt.Errorf("failed to create transaction: %v", err)
	}

	// Check for alerts
	if quantity <= inventory.MinQuantity {
		alert := &models.InventoryAlert{
			ID:          uuid.New().String(),
			InventoryID: inventory.ID,
			Type:        "low_stock",
			Message:     fmt.Sprintf("Low stock alert for product %s", inventory.ProductID),
			Status:      "new",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.db.Create(alert).Error; err != nil {
			return fmt.Errorf("failed to create alert: %v", err)
		}

		// Publish alert event
		alertEvent := &events.InventoryAlertEvent{
			BaseEvent: events.BaseEvent{
				ID:        uuid.New().String(),
				Type:      string(events.LowStockAlert),
				Timestamp: time.Now(),
				Version:   "1.0",
				Source:    s.config.App.Name,
			},
		}
		alertEvent.Data.AlertID = alert.ID
		alertEvent.Data.InventoryID = inventory.ID
		alertEvent.Data.AlertType = "low_stock"
		alertEvent.Data.CurrentLevel = quantity

		if err := s.publishEvent(fmt.Sprintf("%s.inventory.alert", s.config.NATS.SubjectPrefix), alertEvent); err != nil {
			return fmt.Errorf("failed to publish alert event: %v", err)
		}
	}

	return nil
}

// Helper function to publish events
func (s *InventoryService) publishEvent(subject string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	_, err = s.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish event: %v", err)
	}

	return nil
}
