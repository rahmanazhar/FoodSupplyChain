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

	"github.com/rahmanazhar/FoodSupplyChain/internal/shipment/config"
	"github.com/rahmanazhar/FoodSupplyChain/pkg/models"
)

// ShipmentService handles the core business logic for shipment management
type ShipmentService struct {
	config *config.Config
	db     *gorm.DB
	nc     *nats.Conn
	js     nats.JetStreamContext
}

// NewShipmentService creates a new shipment service instance
func NewShipmentService(cfg *config.Config) (*ShipmentService, error) {
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
		&models.Shipment{},
		&models.ShipmentEvent{},
		&models.Carrier{},
		&models.ShipmentAlert{},
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

	return &ShipmentService{
		config: cfg,
		db:     db,
		nc:     nc,
		js:     js,
	}, nil
}

// Close closes all connections
func (s *ShipmentService) Close() error {
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

// CreateShipment creates a new shipment
func (s *ShipmentService) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	shipment.ID = uuid.New().String()
	shipment.CreatedAt = time.Now()
	shipment.UpdatedAt = time.Now()

	if err := s.db.Create(shipment).Error; err != nil {
		return fmt.Errorf("failed to create shipment: %v", err)
	}

	// Create initial shipment event
	event := &models.ShipmentEvent{
		ID:          uuid.New().String(),
		ShipmentID:  shipment.ID,
		Type:        "created",
		Description: "Shipment created",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(event).Error; err != nil {
		return fmt.Errorf("failed to create shipment event: %v", err)
	}

	// Publish event
	if err := s.publishEvent(fmt.Sprintf("%s.shipment.created", s.config.NATS.SubjectPrefix), event); err != nil {
		return fmt.Errorf("failed to publish shipment created event: %v", err)
	}

	return nil
}

// UpdateShipmentStatus updates the status of a shipment
func (s *ShipmentService) UpdateShipmentStatus(ctx context.Context, id string, status string, location string) error {
	var shipment models.Shipment
	if err := s.db.First(&shipment, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to find shipment: %v", err)
	}

	shipment.Status = status
	shipment.UpdatedAt = time.Now()

	if err := s.db.Save(&shipment).Error; err != nil {
		return fmt.Errorf("failed to update shipment: %v", err)
	}

	// Create status update event
	event := &models.ShipmentEvent{
		ID:          uuid.New().String(),
		ShipmentID:  shipment.ID,
		Type:        "status_changed",
		Location:    location,
		Description: fmt.Sprintf("Status updated to: %s", status),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(event).Error; err != nil {
		return fmt.Errorf("failed to create shipment event: %v", err)
	}

	// Publish event
	if err := s.publishEvent(fmt.Sprintf("%s.shipment.status_updated", s.config.NATS.SubjectPrefix), event); err != nil {
		return fmt.Errorf("failed to publish status update event: %v", err)
	}

	return nil
}

// Helper function to publish events
func (s *ShipmentService) publishEvent(subject string, event interface{}) error {
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
