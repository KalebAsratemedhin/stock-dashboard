package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/queue"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/pkg/generator"
)

func main() {
	rmq, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	publisher := queue.NewPublisher(rmq)

	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NVDA", "NFLX", "AMD", "INTC"}

	stockGen := generator.NewStockGenerator(symbols)
	salesGen := generator.NewSalesGenerator()
	userEventGen := generator.NewUserEventGenerator()
	financialGen := generator.NewFinancialGenerator()

	prevQuotes := make(map[string]*models.StockQuote)

	stopChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down simulator...")
		close(stopChan)
	}()

	log.Println("Starting real-time data simulator...")

	go simulateStockData(stockGen, publisher, symbols, stopChan, prevQuotes)
	go simulateSalesData(salesGen, publisher, stopChan)
	go simulateUserEvents(userEventGen, publisher, stopChan)
	go simulateFinancialMetrics(financialGen, publisher, stopChan)

	<-stopChan
	log.Println("Simulator stopped")
}

func simulateStockData(gen *generator.StockGenerator, publisher *queue.Publisher, symbols []string, stopChan chan struct{}, prevQuotes map[string]*models.StockQuote) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			respectMarketHours := os.Getenv("SIMULATE_MARKET_HOURS") != "false"
			if respectMarketHours && !isMarketOpen(time.Now()) {
				continue
			}

			for _, symbol := range symbols {
				prevQuote := prevQuotes[symbol]
				quote := gen.GenerateQuote(symbol, prevQuote, time.Now())
				prevQuotes[symbol] = quote
				
				if err := publisher.PublishStockQuote(queue.StockQueue, quote); err != nil {
					log.Printf("Error publishing stock quote: %v", err)
				}
			}
		}
	}
}

func simulateSalesData(gen *generator.SalesGenerator, publisher *queue.Publisher, stopChan chan struct{}) {
	interval := 60 + time.Duration(time.Now().Unix()%60)*time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			sale := gen.GenerateSale(time.Now())
			if err := publisher.PublishSale(queue.SalesQueue, sale); err != nil {
				log.Printf("Error publishing sale: %v", err)
			}
			
			interval = 60 + time.Duration(time.Now().Unix()%60)*time.Second
			ticker.Reset(interval)
		}
	}
}

func simulateUserEvents(gen *generator.UserEventGenerator, publisher *queue.Publisher, stopChan chan struct{}) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			eventsCount := 1 + int(time.Now().Unix()%3)
			for i := 0; i < eventsCount; i++ {
				event := gen.GenerateEvent(time.Now())
				if err := publisher.PublishUserEvent(queue.UserEventsQueue, event); err != nil {
					log.Printf("Error publishing user event: %v", err)
				}
			}
		}
	}
}

func simulateFinancialMetrics(gen *generator.FinancialGenerator, publisher *queue.Publisher, stopChan chan struct{}) {
	interval := 5 + time.Duration(time.Now().Unix()%5)*time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			metric := gen.GenerateMetric(time.Now())
			if err := publisher.PublishFinancialMetric(queue.FinancialQueue, metric); err != nil {
				log.Printf("Error publishing financial metric: %v", err)
			}
			
			interval = 5 + time.Duration(time.Now().Unix()%5)*time.Minute
			ticker.Reset(interval)
		}
	}
}

func isMarketOpen(t time.Time) bool {
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}
	hour := t.Hour()
	return hour >= 9 && hour < 16
}

