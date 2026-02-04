package repository

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"gorm.io/gorm"
)

type StockRepository interface {
	Create(quote *models.StockQuote) error
	BatchCreate(quotes []*models.StockQuote) error
	GetByID(id uint) (*models.StockQuote, error)
	GetBySymbol(symbol string, limit int) ([]*models.StockQuote, error)
	GetByTimeRange(symbol string, start, end time.Time) ([]*models.StockQuote, error)
	GetLatest(symbol string) (*models.StockQuote, error)
	GetLatestAll(limit int) ([]*models.StockQuote, error)
	AggregateByTimeRange(symbol string, start, end time.Time, interval string) ([]map[string]interface{}, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(quote *models.StockQuote) error {
	return r.db.Create(quote).Error
}

func (r *stockRepository) BatchCreate(quotes []*models.StockQuote) error {
	if len(quotes) == 0 {
		return nil
	}
	return r.db.CreateInBatches(quotes, 1000).Error
}

func (r *stockRepository) GetByID(id uint) (*models.StockQuote, error) {
	var quote models.StockQuote
	err := r.db.Where("id = ?", id).First(&quote).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *stockRepository) GetBySymbol(symbol string, limit int) ([]*models.StockQuote, error) {
	var quotes []*models.StockQuote
	query := r.db.Where("symbol = ?", symbol).Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&quotes).Error
	return quotes, err
}

func (r *stockRepository) GetByTimeRange(symbol string, start, end time.Time) ([]*models.StockQuote, error) {
	var quotes []*models.StockQuote
	err := r.db.Where("symbol = ? AND timestamp >= ? AND timestamp <= ?", symbol, start, end).
		Order("timestamp ASC").
		Find(&quotes).Error
	return quotes, err
}

func (r *stockRepository) GetLatest(symbol string) (*models.StockQuote, error) {
	var quote models.StockQuote
	err := r.db.Where("symbol = ?", symbol).
		Order("timestamp DESC").
		First(&quote).Error
	if err != nil {
		return nil, err
	}
	return &quote, nil
}

func (r *stockRepository) GetLatestAll(limit int) ([]*models.StockQuote, error) {
	var quotes []*models.StockQuote
	query := r.db.Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&quotes).Error
	return quotes, err
}

func (r *stockRepository) AggregateByTimeRange(symbol string, start, end time.Time, interval string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// Use TimescaleDB time_bucket for aggregation
	query := `
		SELECT 
			time_bucket(?, timestamp) AS bucket,
			symbol,
			FIRST(open, timestamp) AS open,
			MAX(high) AS high,
			MIN(low) AS low,
			LAST(close, timestamp) AS close,
			SUM(volume) AS volume
		FROM stock_quotes
		WHERE symbol = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY bucket, symbol
		ORDER BY bucket ASC
	`
	
	err := r.db.Raw(query, interval, symbol, start, end).Scan(&results).Error
	return results, err
}