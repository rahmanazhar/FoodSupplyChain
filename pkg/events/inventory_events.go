package events

import "time"

// InventoryEventType defines the types of inventory events
type InventoryEventType string

const (
	// Event types for inventory
	InventoryCreated     InventoryEventType = "inventory.created"
	InventoryUpdated     InventoryEventType = "inventory.updated"
	InventoryDeleted     InventoryEventType = "inventory.deleted"
	StockLevelChanged    InventoryEventType = "inventory.stock.changed"
	LowStockAlert        InventoryEventType = "inventory.alert.low_stock"
	OverstockAlert       InventoryEventType = "inventory.alert.overstock"
	StockReorderRequired InventoryEventType = "inventory.alert.reorder"
)

// BaseEvent represents the common fields for all events
type BaseEvent struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	Source      string    `json:"source"`
	TraceID     string    `json:"trace_id,omitempty"`
	TenantID    string    `json:"tenant_id,omitempty"`
}

// InventoryEvent represents an event related to inventory changes
type InventoryEvent struct {
	BaseEvent
	Data struct {
		InventoryID string `json:"inventory_id"`
		ProductID   string `json:"product_id"`
		LocationID  string `json:"location_id"`
		Quantity    int    `json:"quantity"`
		PrevQuantity int   `json:"prev_quantity,omitempty"`
		Reason      string `json:"reason,omitempty"`
	} `json:"data"`
}

// InventoryAlertEvent represents an event for inventory alerts
type InventoryAlertEvent struct {
	BaseEvent
	Data struct {
		AlertID     string `json:"alert_id"`
		InventoryID string `json:"inventory_id"`
		AlertType   string `json:"alert_type"`
		Message     string `json:"message"`
		Threshold   int    `json:"threshold,omitempty"`
		CurrentLevel int   `json:"current_level"`
	} `json:"data"`
}

// InventoryTransactionEvent represents an event for inventory transactions
type InventoryTransactionEvent struct {
	BaseEvent
	Data struct {
		TransactionID string `json:"transaction_id"`
		InventoryID   string `json:"inventory_id"`
		Type          string `json:"type"`
		Quantity      int    `json:"quantity"`
		Reference     string `json:"reference,omitempty"`
		Notes         string `json:"notes,omitempty"`
	} `json:"data"`
}

// EventPublisher defines the interface for publishing events
type EventPublisher interface {
	Publish(subject string, event interface{}) error
}

// EventSubscriber defines the interface for subscribing to events
type EventSubscriber interface {
	Subscribe(subject string, handler func(event interface{}) error) error
}

// EventProcessor defines the interface for processing events
type EventProcessor interface {
	Process(event interface{}) error
}
