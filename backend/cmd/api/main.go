package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/api"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/database"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/redis"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/repository"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/services"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/websocket"
)

func main() {
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	stockRepo := repository.NewStockRepository(database.DB)
	saleRepo := repository.NewSaleRepository(database.DB)
	userEventRepo := repository.NewUserEventRepository(database.DB)
	financialRepo := repository.NewFinancialMetricRepository(database.DB)

	stockService := services.NewStockService(stockRepo)
	saleService := services.NewSaleService(saleRepo)
	userEventService := services.NewUserEventService(userEventRepo)
	financialService := services.NewFinancialMetricService(financialRepo)

	handler := api.NewHandler(stockService, saleService, userEventService, financialService)
	
	wsHub := websocket.NewHub()
	go wsHub.Run()

	redisClient, err := redis.NewClient()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
	} else {
		defer redisClient.Close()
		redisBridge := websocket.NewRedisBridge(wsHub, redisClient)
		go redisBridge.Start()
	}

	router := api.NewRouter(handler, wsHub)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}

