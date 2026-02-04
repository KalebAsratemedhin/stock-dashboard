package models

import (
	"time"
)

// Sale represents a sales transaction
type Sale struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Timestamp   time.Time `gorm:"type:timestamptz;not null;index" json:"timestamp"`
	ProductID   string    `gorm:"type:varchar(50);not null;index" json:"product_id"`
	ProductName string    `gorm:"type:varchar(255);not null" json:"product_name"`
	Category    string    `gorm:"type:varchar(100);index" json:"category"`
	CustomerID  string    `gorm:"type:varchar(50);index" json:"customer_id"`
	Region      string    `gorm:"type:varchar(100);index" json:"region"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	UnitPrice   float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	Discount    float64   `gorm:"type:decimal(10,2);default:0" json:"discount"`
	Revenue     float64   `gorm:"type:decimal(10,2);not null" json:"revenue"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name
func (Sale) TableName() string {
	return "sales"
}