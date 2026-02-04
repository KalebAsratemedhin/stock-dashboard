package repository

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"gorm.io/gorm"
)

type FinancialMetricRepository interface {
	Create(metric *models.FinancialMetric) error
	BatchCreate(metrics []*models.FinancialMetric) error
	GetByID(id uint) (*models.FinancialMetric, error)
	GetByTimeRange(start, end time.Time) ([]*models.FinancialMetric, error)
	GetByMetricType(metricType string, limit int) ([]*models.FinancialMetric, error)
	GetByDepartment(department string, limit int) ([]*models.FinancialMetric, error)
	GetTotalByType(metricType string, start, end time.Time) (float64, error)
	GetMetricsByDepartment(start, end time.Time) ([]map[string]interface{}, error)
}

type financialMetricRepository struct {
	db *gorm.DB
}

func NewFinancialMetricRepository(db *gorm.DB) FinancialMetricRepository {
	return &financialMetricRepository{db: db}
}

func (r *financialMetricRepository) Create(metric *models.FinancialMetric) error {
	return r.db.Create(metric).Error
}

func (r *financialMetricRepository) BatchCreate(metrics []*models.FinancialMetric) error {
	if len(metrics) == 0 {
		return nil
	}
	return r.db.CreateInBatches(metrics, 1000).Error
}

func (r *financialMetricRepository) GetByID(id uint) (*models.FinancialMetric, error) {
	var metric models.FinancialMetric
	err := r.db.Where("id = ?", id).First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func (r *financialMetricRepository) GetByTimeRange(start, end time.Time) ([]*models.FinancialMetric, error) {
	var metrics []*models.FinancialMetric
	err := r.db.Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp ASC").
		Find(&metrics).Error
	return metrics, err
}

func (r *financialMetricRepository) GetByMetricType(metricType string, limit int) ([]*models.FinancialMetric, error) {
	var metrics []*models.FinancialMetric
	query := r.db.Where("metric_type = ?", metricType).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&metrics).Error
	return metrics, err
}

func (r *financialMetricRepository) GetByDepartment(department string, limit int) ([]*models.FinancialMetric, error) {
	var metrics []*models.FinancialMetric
	query := r.db.Where("department = ?", department).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&metrics).Error
	return metrics, err
}

func (r *financialMetricRepository) GetTotalByType(metricType string, start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.FinancialMetric{}).
		Where("metric_type = ? AND timestamp >= ? AND timestamp <= ?", metricType, start, end).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *financialMetricRepository) GetMetricsByDepartment(start, end time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&models.FinancialMetric{}).
		Select("department, metric_type, SUM(amount) as total_amount").
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Group("department, metric_type").
		Order("department, metric_type").
		Find(&results).Error
	return results, err
}