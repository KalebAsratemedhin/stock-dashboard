package repository

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"gorm.io/gorm"
)

type SaleRepository interface {
	Create(sale *models.Sale) error
	BatchCreate(sales []*models.Sale) error
	GetByID(id uint) (*models.Sale, error)
	GetByTimeRange(start, end time.Time) ([]*models.Sale, error)
	GetByCategory(category string, limit int) ([]*models.Sale, error)
	GetByRegion(region string, limit int) ([]*models.Sale, error)
	GetRevenueByTimeRange(start, end time.Time) (float64, error)
	GetRevenueByCategory(start, end time.Time) ([]map[string]interface{}, error)
}

type saleRepository struct {
	db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) SaleRepository {
	return &saleRepository{db: db}
}

func (r *saleRepository) Create(sale *models.Sale) error {
	return r.db.Create(sale).Error
}

func (r *saleRepository) BatchCreate(sales []*models.Sale) error {
	if len(sales) == 0 {
		return nil
	}
	return r.db.CreateInBatches(sales, 1000).Error
}

func (r *saleRepository) GetByID(id uint) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.Where("id = ?", id).First(&sale).Error
	if err != nil {
		return nil, err
	}
	return &sale, nil
}

func (r *saleRepository) GetByTimeRange(start, end time.Time) ([]*models.Sale, error) {
	var sales []*models.Sale
	err := r.db.Where("timestamp >= ? AND timestamp <= ?", start, end).
		Order("timestamp ASC").
		Find(&sales).Error
	return sales, err
}

func (r *saleRepository) GetByCategory(category string, limit int) ([]*models.Sale, error) {
	var sales []*models.Sale
	query := r.db.Where("category = ?", category).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&sales).Error
	return sales, err
}

func (r *saleRepository) GetByRegion(region string, limit int) ([]*models.Sale, error) {
	var sales []*models.Sale
	query := r.db.Where("region = ?", region).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&sales).Error
	return sales, err
}

func (r *saleRepository) GetRevenueByTimeRange(start, end time.Time) (float64, error) {
	var total float64
	err := r.db.Model(&models.Sale{}).
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Select("COALESCE(SUM(revenue), 0)").
		Scan(&total).Error
	return total, err
}

func (r *saleRepository) GetRevenueByCategory(start, end time.Time) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := r.db.Model(&models.Sale{}).
		Select("category, SUM(revenue) as total_revenue, COUNT(*) as count").
		Where("timestamp >= ? AND timestamp <= ?", start, end).
		Group("category").
		Order("total_revenue DESC").
		Find(&results).Error
	return results, err
}