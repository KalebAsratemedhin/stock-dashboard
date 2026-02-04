package services

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
)

type FinancialMetricService interface {
	CreateMetric(metric *models.FinancialMetric) error
	BatchCreateMetrics(metrics []*models.FinancialMetric) error
	GetMetricByID(id uint) (*models.FinancialMetric, error)
	GetRecentMetrics(limit int) ([]*models.FinancialMetric, error)
	GetMetricsByTimeRange(start, end time.Time) ([]*models.FinancialMetric, error)
	GetMetricsByType(metricType string, limit int) ([]*models.FinancialMetric, error)
	GetMetricsByDepartment(department string, limit int) ([]*models.FinancialMetric, error)
	GetTotalByType(metricType string, start, end time.Time) (float64, error)
	GetMetricsGroupedByDepartment(start, end time.Time) ([]map[string]interface{}, error)
	GetBudgetVsActual(start, end time.Time) ([]map[string]interface{}, error)
	GetVarianceByDepartment(start, end time.Time) ([]map[string]interface{}, error)
	GetMetricsByTimeRangeAggregated(start, end time.Time, interval string) ([]map[string]interface{}, error)
}

type financialMetricService struct {
	repo repository.FinancialMetricRepository
}

func NewFinancialMetricService(repo repository.FinancialMetricRepository) FinancialMetricService {
	return &financialMetricService{repo: repo}
}

func (s *financialMetricService) CreateMetric(metric *models.FinancialMetric) error {
	if metric.MetricType == "" {
		return ErrInvalidInput
	}
	if metric.Amount == 0 {
		return ErrInvalidInput
	}
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}
	if metric.Budget > 0 {
		metric.Variance = metric.Amount - metric.Budget
		metric.VariancePct = (metric.Variance / metric.Budget) * 100
	}
	return s.repo.Create(metric)
}

func (s *financialMetricService) BatchCreateMetrics(metrics []*models.FinancialMetric) error {
	for _, metric := range metrics {
		if metric.MetricType == "" {
			return ErrInvalidInput
		}
		if metric.Amount == 0 {
			return ErrInvalidInput
		}
		if metric.Timestamp.IsZero() {
			metric.Timestamp = time.Now()
		}
		if metric.Budget > 0 {
			metric.Variance = metric.Amount - metric.Budget
			metric.VariancePct = (metric.Variance / metric.Budget) * 100
		}
	}
	return s.repo.BatchCreate(metrics)
}

func (s *financialMetricService) GetMetricByID(id uint) (*models.FinancialMetric, error) {
	return s.repo.GetByID(id)
}

func (s *financialMetricService) GetRecentMetrics(limit int) ([]*models.FinancialMetric, error) {
	if limit <= 0 {
		limit = 100
	}
	metrics, err := s.repo.GetByTimeRange(time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		return nil, err
	}
	if len(metrics) > limit {
		metrics = metrics[:limit]
	}
	return metrics, nil
}

func (s *financialMetricService) GetMetricsByTimeRange(start, end time.Time) ([]*models.FinancialMetric, error) {
	return s.repo.GetByTimeRange(start, end)
}

func (s *financialMetricService) GetMetricsByType(metricType string, limit int) ([]*models.FinancialMetric, error) {
	return s.repo.GetByMetricType(metricType, limit)
}

func (s *financialMetricService) GetMetricsByDepartment(department string, limit int) ([]*models.FinancialMetric, error) {
	return s.repo.GetByDepartment(department, limit)
}

func (s *financialMetricService) GetTotalByType(metricType string, start, end time.Time) (float64, error) {
	return s.repo.GetTotalByType(metricType, start, end)
}

func (s *financialMetricService) GetMetricsGroupedByDepartment(start, end time.Time) ([]map[string]interface{}, error) {
	return s.repo.GetMetricsByDepartment(start, end)
}

func (s *financialMetricService) GetBudgetVsActual(start, end time.Time) ([]map[string]interface{}, error) {
	metrics, err := s.repo.GetByTimeRange(start, end)
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string]float64)
	for _, m := range metrics {
		if m.Department == "" {
			continue
		}
		if result[m.Department] == nil {
			result[m.Department] = make(map[string]float64)
		}
		result[m.Department]["actual"] += m.Amount
		if m.Budget > 0 {
			result[m.Department]["budget"] += m.Budget
		}
	}
	var output []map[string]interface{}
	for dept, values := range result {
		variance := values["actual"] - values["budget"]
		var variancePct float64
		if values["budget"] > 0 {
			variancePct = (variance / values["budget"]) * 100
		}
		output = append(output, map[string]interface{}{
			"department":   dept,
			"actual":       values["actual"],
			"budget":       values["budget"],
			"variance":     variance,
			"variance_pct": variancePct,
		})
	}
	return output, nil
}

func (s *financialMetricService) GetVarianceByDepartment(start, end time.Time) ([]map[string]interface{}, error) {
	return s.GetBudgetVsActual(start, end)
}

func (s *financialMetricService) GetMetricsByTimeRangeAggregated(start, end time.Time, interval string) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, ErrInternal
}