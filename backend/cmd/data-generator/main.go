package main

import (
	"log"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/database"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/services"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/pkg/generator"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX", "AMD", "INTC"}

	stockRepo := repository.NewStockRepository(database.DB)
	saleRepo := repository.NewSaleRepository(database.DB)
	userEventRepo := repository.NewUserEventRepository(database.DB)
	financialRepo := repository.NewFinancialMetricRepository(database.DB)

	stockService := services.NewStockService(stockRepo)
	saleService := services.NewSaleService(saleRepo)
	userEventService := services.NewUserEventService(userEventRepo)
	financialService := services.NewFinancialMetricService(financialRepo)

	stockGen := generator.NewStockGenerator(symbols)
	salesGen := generator.NewSalesGenerator()
	userEventGen := generator.NewUserEventGenerator()
	financialGen := generator.NewFinancialGenerator()

	months := 6
	startDate := time.Now().AddDate(0, -months, 0)

	log.Printf("Generating %d months of historical data starting from %s", months, startDate.Format("2006-01-02"))

	generateStockData(stockGen, stockService, symbols, startDate)
	generateSalesData(salesGen, saleService, startDate)
	generateUserEventData(userEventGen, userEventService, startDate)
	generateFinancialData(financialGen, financialService, startDate)

	log.Println("Historical data generation completed!")
}

func generateStockData(gen *generator.StockGenerator, service services.StockService, symbols []string, startDate time.Time) {
	log.Println("Generating stock data...")
	
	current := startDate
	end := time.Now()
	prevQuotes := make(map[string]*models.StockQuote)
	batch := make([]*models.StockQuote, 0, 1000)

	for current.Before(end) {
		if isMarketOpen(current) {
			for _, symbol := range symbols {
				prevQuote := prevQuotes[symbol]
				quote := gen.GenerateQuote(symbol, prevQuote, current)
				prevQuotes[symbol] = quote
				batch = append(batch, quote)

				if len(batch) >= 1000 {
					if err := service.BatchCreateQuotes(batch); err != nil {
						log.Printf("Error batch creating quotes: %v", err)
					}
					batch = batch[:0]
				}
			}
		}

		current = current.Add(1 * time.Minute)
	}

	if len(batch) > 0 {
		if err := service.BatchCreateQuotes(batch); err != nil {
			log.Printf("Error batch creating quotes: %v", err)
		}
	}

	log.Println("Stock data generation completed")
}

func generateSalesData(gen *generator.SalesGenerator, service services.SaleService, startDate time.Time) {
	log.Println("Generating sales data...")
	
	current := startDate
	end := time.Now()
	batch := make([]*models.Sale, 0, 1000)

	for current.Before(end) {
		salesPerDay := 50 + (current.Hour() * 2)
		for i := 0; i < salesPerDay; i++ {
			timestamp := current.Add(time.Duration(i*86400/salesPerDay) * time.Second)
			sale := gen.GenerateSale(timestamp)
			batch = append(batch, sale)

			if len(batch) >= 1000 {
				if err := service.BatchCreateSales(batch); err != nil {
					log.Printf("Error batch creating sales: %v", err)
				}
				batch = batch[:0]
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	if len(batch) > 0 {
		if err := service.BatchCreateSales(batch); err != nil {
			log.Printf("Error batch creating sales: %v", err)
		}
	}

	log.Println("Sales data generation completed")
}

func generateUserEventData(gen *generator.UserEventGenerator, service services.UserEventService, startDate time.Time) {
	log.Println("Generating user event data...")
	
	current := startDate
	end := time.Now()
	batch := make([]*models.UserEvent, 0, 1000)

	for current.Before(end) {
		eventsPerDay := 1000 + (current.Hour() * 100)
		for i := 0; i < eventsPerDay; i++ {
			timestamp := current.Add(time.Duration(i*86400/eventsPerDay) * time.Second)
			event := gen.GenerateEvent(timestamp)
			batch = append(batch, event)

			if len(batch) >= 1000 {
				if err := service.BatchCreateEvents(batch); err != nil {
					log.Printf("Error batch creating events: %v", err)
				}
				batch = batch[:0]
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	if len(batch) > 0 {
		if err := service.BatchCreateEvents(batch); err != nil {
			log.Printf("Error batch creating events: %v", err)
		}
	}

	log.Println("User event data generation completed")
}

func generateFinancialData(gen *generator.FinancialGenerator, service services.FinancialMetricService, startDate time.Time) {
	log.Println("Generating financial metrics data...")
	
	current := startDate
	end := time.Now()
	batch := make([]*models.FinancialMetric, 0, 1000)

	for current.Before(end) {
		metricsPerDay := 20
		for i := 0; i < metricsPerDay; i++ {
			timestamp := current.Add(time.Duration(i*86400/metricsPerDay) * time.Second)
			metric := gen.GenerateMetric(timestamp)
			batch = append(batch, metric)

			if len(batch) >= 1000 {
				if err := service.BatchCreateMetrics(batch); err != nil {
					log.Printf("Error batch creating metrics: %v", err)
				}
				batch = batch[:0]
			}
		}

		current = current.AddDate(0, 0, 1)
	}

	if len(batch) > 0 {
		if err := service.BatchCreateMetrics(batch); err != nil {
			log.Printf("Error batch creating metrics: %v", err)
		}
	}

	log.Println("Financial metrics data generation completed")
}

func isMarketOpen(t time.Time) bool {
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}
	hour := t.Hour()
	return hour >= 9 && hour < 16
}

