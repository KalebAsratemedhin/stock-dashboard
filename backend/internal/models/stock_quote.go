package models

import (
	"time"
)

type StockQuote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Symbol    string    `gorm:"type:varchar(10);not null;index" json:"symbol"`
	Timestamp time.Time `gorm:"type:timestamptz;not null;index" json:"timestamp"`
	Open      float64   `gorm:"type:decimal(10,2);not null" json:"open"`
	High      float64   `gorm:"type:decimal(10,2);not null" json:"high"`
	Low       float64   `gorm:"type:decimal(10,2);not null" json:"low"`
	Close     float64   `gorm:"type:decimal(10,2);not null" json:"close"`
	Volume    int64     `gorm:"not null" json:"volume"`
	Bid       float64   `gorm:"type:decimal(10,2)" json:"bid,omitempty"`
	Ask       float64   `gorm:"type:decimal(10,2)" json:"ask,omitempty"`
	Change    float64   `gorm:"type:decimal(10,2)" json:"change,omitempty"`
	ChangePct float64   `gorm:"type:decimal(5,2)" json:"change_pct,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (StockQuote) TableName() string {
	return "stock_quotes"
}