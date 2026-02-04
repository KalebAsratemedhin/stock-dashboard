package repository

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"gorm.io/gorm"
)

type UserEventRepository interface {
	Create(event *models.UserEvent) error
	BatchCreate(events []*models.UserEvent) error
	GetByID(id uint) (*models.UserEvent, error)
	GetByTimeRange(start, end time.Time) ([]*models.UserEvent, error)
	GetByEventType(eventType string, limit int) ([]*models.UserEvent, error)
	GetByUserID(userID string, limit int) ([]*models.UserEvent, error)
	GetEventCountsByType(start, end time.Time) ([]map[string]interface{}, error)
	GetPageViewsByTimeRange(start, end time.Time, interval string) ([]map[string]interface{}, error)
}

type userEventRepository struct {
	db *gorm.DB
}

func NewUserEventRepository(db *gorm.DB) UserEventRepository {
	return &userEventRepository{db: db}
}

func (r *userEventRepository) Create(event *models.UserEvent) error {
	return r.db.Create(event).Error
}

func (r *userEventRepository) BatchCreate(events []*models.UserEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.db.CreateInBatches(events, 1000).Error
}

func (r *userEventRepository) GetByID(id uint) (*models.UserEvent, error) {
	var event models.UserEvent
	err := r.db.Where("id = ?", id).First(&event).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *userEventRepository) GetByTimeRange(start, end time.Time) ([]*models.UserEvent, error) {
	var events []*models.UserEvent
	err := r.db.Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp ASC").
		Find(&events).Error
	return events, err
}

func (r *userEventRepository) GetByEventType(eventType string, limit int) ([]*models.UserEvent, error) {
	var events []*models.UserEvent
	query := r.db.Where("event_type = ?", eventType).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&events).Error
	return events, err
}

func (r *userEventRepository) GetByUserID(userID string, limit int) ([]*models.UserEvent, error) {
	var events []*models.UserEvent
	query := r.db.Where("user_id = ?", userID).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&events).Error
	return events, err
}

func (r *userEventRepository) GetEventCountsByType(start, end time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&models.UserEvent{}).
		Select("event_type, COUNT(*) as count").
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Group("event_type").
		Order("count DESC").
		Find(&results).Error
	return results, err
}

func (r *userEventRepository) GetPageViewsByTimeRange(start, end time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT 
			time_bucket(?, timestamp) AS bucket,
			COUNT(*) AS page_views
		FROM user_events
		WHERE event_type = 'page_view' AND timestamp >= ? AND timestamp <= ?
		GROUP BY bucket
		ORDER BY bucket ASC
	`
	err := r.db.Raw(query, interval, start, end).Scan(&results).Error
	return results, err
}