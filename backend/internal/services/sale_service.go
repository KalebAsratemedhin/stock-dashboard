package services

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
)

type SaleService interface {
	CreateSale(sale *models.Sale) error
	BatchCreateSales(sales []*models.Sale) error
	GetSaleByID(id uint) (*models.Sale, error)
	GetRecentSales(limit int) ([]*models.Sale, error)
	GetSalesByTimeRange(start, end time.Time) ([]*models.Sale, error)
	GetSalesByTimeRangeWithFilters(start, end time.Time, category, region string) ([]*models.Sale, error)
	GetSalesByCategory(category string, limit int) ([]*models.Sale, error)
	GetSalesByRegion(region string, limit int) ([]*models.Sale, error)
	GetTotalRevenue(start, end time.Time) (float64, error)
	GetRevenueByCategory(start, end time.Time) ([]map[string]interface{}, error)
}

type saleService struct {
	repo repository.SaleRepository
}

func NewSaleService(repo repository.SaleRepository) SaleService {
	return &saleService{repo: repo}
}

func (s *saleService) CreateSale(sale *models.Sale) error {
	if sale.ProductID == "" {
		return ErrInvalidInput
	}
	if sale.ProductName == "" {
		return ErrInvalidInput
	}
	if sale.Quantity <= 0 {
		return ErrInvalidInput
	}
	if sale.UnitPrice <= 0 {
		return ErrInvalidInput
	}
	if sale.Timestamp.IsZero() {
		sale.Timestamp = time.Now()
	}
	if sale.Revenue == 0 {
		sale.Revenue = float64(sale.Quantity) * sale.UnitPrice * (1 - sale.Discount/100)
	}
	return s.repo.Create(sale)
}

func (s *saleService) BatchCreateSales(sales []*models.Sale) error {
	for _, sale := range sales {
		if sale.ProductID == "" {
			return ErrInvalidInput
		}
		if sale.ProductName == "" {
			return ErrInvalidInput
		}
		if sale.Quantity <= 0 {
			return ErrInvalidInput
		}
		if sale.UnitPrice <= 0 {
			return ErrInvalidInput
		}
		if sale.Timestamp.IsZero() {
			sale.Timestamp = time.Now()
		}
		if sale.Revenue == 0 {
			sale.Revenue = float64(sale.Quantity) * sale.UnitPrice * (1 - sale.Discount/100)
		}
	}
	return s.repo.BatchCreate(sales)
}

func (s *saleService) GetSaleByID(id uint) (*models.Sale, error) {
	return s.repo.GetByID(id)
}

func (s *saleService) GetRecentSales(limit int) ([]*models.Sale, error) {
	if limit <= 0 {
		limit = 100
	}
	sales, err := s.repo.GetByTimeRange(time.Now().AddDate(0, 0, -7), time.Now())
	if err != nil {
		return nil, err
	}
	if len(sales) > limit {
		sales = sales[:limit]
	}
	return sales, nil
}

func (s *saleService) GetSalesByTimeRange(start, end time.Time) ([]*models.Sale, error) {
	return s.repo.GetByTimeRange(start, end)
}

func (s *saleService) GetSalesByTimeRangeWithFilters(start, end time.Time, category, region string) ([]*models.Sale, error) {
	sales, err := s.repo.GetByTimeRange(start, end)
	if err != nil {
		return nil, err
	}
	
	var filtered []*models.Sale
	for _, sale := range sales {
		if category != "" && sale.Category != category {
			continue
		}
		if region != "" && sale.Region != region {
			continue
		}
		filtered = append(filtered, sale)
	}
	
	return filtered, nil
}

func (s *saleService) GetSalesByCategory(category string, limit int) ([]*models.Sale, error) {
	return s.repo.GetByCategory(category, limit)
}

func (s *saleService) GetSalesByRegion(region string, limit int) ([]*models.Sale, error) {
	return s.repo.GetByRegion(region, limit)
}

func (s *saleService) GetTotalRevenue(start, end time.Time) (float64, error) {
	return s.repo.GetRevenueByTimeRange(start, end)
}

func (s *saleService) GetRevenueByCategory(start, end time.Time) ([]map[string]interface{}, error) {
	return s.repo.GetRevenueByCategory(start, end)
}