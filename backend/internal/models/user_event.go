package models

import (
	"time"
)

// UserEvent represents user behavior tracking data
type UserEvent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Timestamp   time.Time `gorm:"type:timestamptz;not null;index" json:"timestamp"`
	EventType   string    `gorm:"type:varchar(50);not null;index" json:"event_type"`
	UserID      string    `gorm:"type:varchar(50);index" json:"user_id"`
	SessionID   string    `gorm:"type:varchar(100);index" json:"session_id"`
	Page        string    `gorm:"type:varchar(255)" json:"page"`
	Device      string    `gorm:"type:varchar(50)" json:"device"`
	Browser     string    `gorm:"type:varchar(50)" json:"browser"`
	Country     string    `gorm:"type:varchar(100);index" json:"country"`
	City        string    `gorm:"type:varchar(100)" json:"city"`
	Referrer    string    `gorm:"type:varchar(500)" json:"referrer"`
	Metadata    string    `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName specifies the table name
func (UserEvent) TableName() string {
	return "user_events"
}