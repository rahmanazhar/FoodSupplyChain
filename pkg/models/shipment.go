package models

import (
	"time"
)

// Shipment represents a shipment in the supply chain
type Shipment struct {
	ID               string     `json:"id" gorm:"primaryKey"`
	OrderID          string     `json:"order_id" gorm:"index;not null"`
	Status           string     `json:"status" gorm:"not null"` // pending, in_transit, delivered, cancelled
	Origin           string     `json:"origin" gorm:"not null"`
	Destination      string     `json:"destination" gorm:"not null"`
	EstimatedArrival time.Time  `json:"estimated_arrival"`
	ActualArrival    *time.Time `json:"actual_arrival,omitempty"`
	CarrierID        string     `json:"carrier_id" gorm:"index"`
	TrackingNumber   string     `json:"tracking_number"`
	Notes            string     `json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// ShipmentEvent represents events in a shipment's lifecycle
type ShipmentEvent struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ShipmentID  string    `json:"shipment_id" gorm:"index;not null"`
	Shipment    Shipment  `json:"shipment" gorm:"foreignKey:ShipmentID"`
	Type        string    `json:"type" gorm:"not null"` // status_changed, location_updated, etc.
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Carrier represents a shipping carrier
type Carrier struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null"`
	ContactInfo string    `json:"contact_info"`
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShipmentAlert represents notifications for shipment-related events
type ShipmentAlert struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	ShipmentID string    `json:"shipment_id" gorm:"index;not null"`
	Shipment   Shipment  `json:"shipment" gorm:"foreignKey:ShipmentID"`
	Type       string    `json:"type" gorm:"not null"` // delay, damage, delivery_attempt, etc.
	Message    string    `json:"message" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null"` // new, acknowledged, resolved
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
