package services

import (
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
)

type StockService interface {
	CreateQuote(quote *models.StockQuote) error
	BatchCreateQuotes(quotes []*models.StockQuote) error
	GetQuoteByID(id uint) (*models.StockQuote, error)
	GetQuotesBySymbol(symbol string, limit int) ([]*models.StockQuote, error)
	GetQuotesByTimeRange(symbol string, start, end time.Time) ([]*models.StockQuote, error)
	GetLatestQuote(symbol string) (*models.StockQuote, error)
	GetLatestQuotes(symbols []string) ([]*models.StockQuote, error)
	GetAggregatedData(symbol string, start, end time.Time, interval string) ([]map[string]interface{}, error)
}

type stockService struct {
	repo repository.StockRepository
}

func NewStockService(repo repository.StockRepository) StockService {
	return &stockService{repo: repo}
}

func (s *stockService) CreateQuote(quote *models.StockQuote) error {
	// Validate
	if quote.Symbol == "" {
		return ErrInvalidInput
	}
	if quote.Timestamp.IsZero() {
		quote.Timestamp = time.Now()
	}
	
	// Calculate change if not set
	if quote.Change == 0 && quote.Close > 0 {
		// This would need previous close price - simplified for now
	}
	
	return s.repo.Create(quote)
}

func (s *stockService) BatchCreateQuotes(quotes []*models.StockQuote) error {
	// Validate all quotes
	for _, quote := range quotes {
		if quote.Symbol == "" {
			return ErrInvalidInput
		}
		if quote.Timestamp.IsZero() {
			quote.Timestamp = time.Now()
		}
	}
	return s.repo.BatchCreate(quotes)
}

func (s *stockService) GetQuoteByID(id uint) (*models.StockQuote, error) {
	return s.repo.GetByID(id)
}

func (s *stockService) GetQuotesBySymbol(symbol string, limit int) ([]*models.StockQuote, error) {
	return s.repo.GetBySymbol(symbol, limit)
}

func (s *stockService) GetQuotesByTimeRange(symbol string, start, end time.Time) ([]*models.StockQuote, error) {
	return s.repo.GetByTimeRange(symbol, start, end)
}

func (s *stockService) GetLatestQuote(symbol string) (*models.StockQuote, error) {
	return s.repo.GetLatest(symbol)
}

func (s *stockService) GetLatestQuotes(symbols []string) ([]*models.StockQuote, error) {
	return s.repo.GetLatestAll(100)
}

func (s *stockService) GetAggregatedData(symbol string, start, end time.Time, interval string) ([]map[string]interface{}, error) {
	return s.repo.AggregateByTimeRange(symbol, start, end, interval)
}