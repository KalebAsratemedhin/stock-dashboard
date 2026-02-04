package services

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
)

type UserEventService interface {
	CreateEvent(event *models.UserEvent) error
	BatchCreateEvents(events []*models.UserEvent) error
	GetEventByID(id uint) (*models.UserEvent, error)
	GetRecentEvents(limit int) ([]*models.UserEvent, error)
	GetEventsByTimeRange(start, end time.Time) ([]*models.UserEvent, error)
	GetEventsByType(eventType string, limit int) ([]*models.UserEvent, error)
	GetEventsByUserID(userID string, limit int) ([]*models.UserEvent, error)
	GetEventCountsByType(start, end time.Time) ([]map[string]interface{}, error)
	GetPageViewsByTimeRange(start, end time.Time, interval string) ([]map[string]interface{}, error)
	GetUniqueUsers(start, end time.Time) (int64, error)
	GetTopPages(start, end time.Time, limit int) ([]map[string]interface{}, error)
	GetEventsByCountry(start, end time.Time) ([]map[string]interface{}, error)
}

type userEventService struct {
	repo repository.UserEventRepository
}

func NewUserEventService(repo repository.UserEventRepository) UserEventService {
	return &userEventService{repo: repo}
}

func (s *userEventService) CreateEvent(event *models.UserEvent) error {
	if event.EventType == "" {
		return ErrInvalidInput
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	return s.repo.Create(event)
}

func (s *userEventService) BatchCreateEvents(events []*models.UserEvent) error {
	for _, event := range events {
		if event.EventType == "" {
			return ErrInvalidInput
		}
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}
	}
	return s.repo.BatchCreate(events)
}

func (s *userEventService) GetEventByID(id uint) (*models.UserEvent, error) {
	return s.repo.GetByID(id)
}

func (s *userEventService) GetRecentEvents(limit int) ([]*models.UserEvent, error) {
	if limit <= 0 {
		limit = 100
	}
	events, err := s.repo.GetByTimeRange(time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		return nil, err
	}
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}

func (s *userEventService) GetEventsByTimeRange(start, end time.Time) ([]*models.UserEvent, error) {
	return s.repo.GetByTimeRange(start, end)
}

func (s *userEventService) GetEventsByType(eventType string, limit int) ([]*models.UserEvent, error) {
	return s.repo.GetByEventType(eventType, limit)
}

func (s *userEventService) GetEventsByUserID(userID string, limit int) ([]*models.UserEvent, error) {
	return s.repo.GetByUserID(userID, limit)
}

func (s *userEventService) GetEventCountsByType(start, end time.Time) ([]map[string]interface{}, error) {
	return s.repo.GetEventCountsByType(start, end)
}

func (s *userEventService) GetPageViewsByTimeRange(start, end time.Time, interval string) ([]map[string]interface{}, error) {
	if interval == "" {
		interval = "1 hour"
	}
	return s.repo.GetPageViewsByTimeRange(start, end, interval)
}

func (s *userEventService) GetUniqueUsers(start, end time.Time) (int64, error) {
	return 0, ErrInternal
}

func (s *userEventService) GetTopPages(start, end time.Time, limit int) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, ErrInternal
}

func (s *userEventService) GetEventsByCountry(start, end time.Time) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, ErrInternal
}