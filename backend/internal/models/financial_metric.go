package models

import (
	"time"
)

// FinancialMetric represents financial KPIs and metrics
type FinancialMetric struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Timestamp     time.Time `gorm:"type:timestamptz;not null;index" json:"timestamp"`
	MetricType    string    `gorm:"type:varchar(50);not null;index" json:"metric_type"`
	Department    string    `gorm:"type:varchar(100);index" json:"department"`
	Category      string    `gorm:"type:varchar(100);index" json:"category"`
	Amount        float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Budget        float64   `gorm:"type:decimal(12,2)" json:"budget,omitempty"`
	Variance      float64   `gorm:"type:decimal(12,2)" json:"variance,omitempty"`
	VariancePct    float64   `gorm:"type:decimal(5,2)" json:"variance_pct,omitempty"`
	Period        string    `gorm:"type:varchar(20)" json:"period"` 
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName specifies the table name
func (FinancialMetric) TableName() string {
	return "financial_metrics"
}