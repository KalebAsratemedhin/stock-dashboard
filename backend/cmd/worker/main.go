package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/database"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/queue"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/redis"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/services"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	rmq, err := queue.NewRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmq.Close()

	stockRepo := repository.NewStockRepository(database.DB)
	saleRepo := repository.NewSaleRepository(database.DB)
	userEventRepo := repository.NewUserEventRepository(database.DB)
	financialRepo := repository.NewFinancialMetricRepository(database.DB)

	stockService := services.NewStockService(stockRepo)
	saleService := services.NewSaleService(saleRepo)
	userEventService := services.NewUserEventService(userEventRepo)
	financialService := services.NewFinancialMetricService(financialRepo)

	redisClient, err := redis.NewClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		defer redisClient.Close()
	}

	consumer := queue.NewConsumer(rmq)
	worker := queue.NewWorker(
		consumer,
		stockService,
		saleService,
		userEventService,
		financialService,
		redisClient,
	)

	if err := worker.StartAll(); err != nil {
		log.Fatalf("Failed to start workers: %v", err)
	}

	log.Println("Worker service started. Press CTRL+C to exit.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down worker service...")
}