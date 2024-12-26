package models

import (
	"time"
)

// Product represents a product in the supply chain
type Product struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	SKU         string    `json:"sku" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UnitPrice   float64   `json:"unit_price" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Inventory represents the current stock level of a product at a location
type Inventory struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	ProductID string    `json:"product_id" gorm:"index;not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	LocationID string   `json:"location_id" gorm:"index;not null"`
	Location   Location `json:"location" gorm:"foreignKey:LocationID"`
	Quantity   int      `json:"quantity" gorm:"not null"`
	MinQuantity int     `json:"min_quantity" gorm:"not null"`
	MaxQuantity int     `json:"max_quantity" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Location represents a physical location in the supply chain
type Location struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"` // warehouse, store, distribution center
	Address   string    `json:"address"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// InventoryTransaction represents a change in inventory levels
type InventoryTransaction struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	InventoryID string    `json:"inventory_id" gorm:"index;not null"`
	Inventory   Inventory `json:"inventory" gorm:"foreignKey:InventoryID"`
	Type        string    `json:"type" gorm:"not null"` // received, shipped, adjusted
	Quantity    int       `json:"quantity" gorm:"not null"`
	Reference   string    `json:"reference"` // PO number, shipment ID, etc.
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// InventoryAlert represents notifications for inventory-related events
type InventoryAlert struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	InventoryID string    `json:"inventory_id" gorm:"index;not null"`
	Inventory   Inventory `json:"inventory" gorm:"foreignKey:InventoryID"`
	Type        string    `json:"type" gorm:"not null"` // low_stock, overstock, reorder
	Message     string    `json:"message" gorm:"not null"`
	Status      string    `json:"status" gorm:"not null"` // new, acknowledged, resolved
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
