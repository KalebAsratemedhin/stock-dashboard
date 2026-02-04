package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/redis"
	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/services"
)

type Worker struct {
	consumer        *Consumer
	stockService    services.StockService
	saleService     services.SaleService
	userEventService services.UserEventService
	financialService services.FinancialMetricService
	redisClient     *redis.Client
	batchSize       int
	batchBuffer     map[string][]interface{}
	batchMutex      sync.Mutex
}

func NewWorker(
	consumer *Consumer,
	stockService services.StockService,
	saleService services.SaleService,
	userEventService services.UserEventService,
	financialService services.FinancialMetricService,
	redisClient *redis.Client,
) *Worker {
	return &Worker{
		consumer:        consumer,
		stockService:    stockService,
		saleService:     saleService,
		userEventService: userEventService,
		financialService: financialService,
		redisClient:     redisClient,
		batchSize:       100,
		batchBuffer:     make(map[string][]interface{}),
	}
}

func (w *Worker) StartStockWorker() error {
	return w.consumer.ConsumeJSON(StockQueue, func(data interface{}) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		var quote models.StockQuote
		if err := json.Unmarshal(jsonData, &quote); err != nil {
			return fmt.Errorf("failed to unmarshal stock quote: %w", err)
		}

		if err := w.stockService.CreateQuote(&quote); err != nil {
			return err
		}

		if w.redisClient != nil {
			w.redisClient.Publish(redis.StockChannel, &quote)
		}

		return nil
	})
}

func (w *Worker) StartSaleWorker() error {
	return w.consumer.ConsumeJSON(SalesQueue, func(data interface{}) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		var sale models.Sale
		if err := json.Unmarshal(jsonData, &sale); err != nil {
			return fmt.Errorf("failed to unmarshal sale: %w", err)
		}

		if err := w.saleService.CreateSale(&sale); err != nil {
			return err
		}

		if w.redisClient != nil {
			w.redisClient.Publish(redis.SalesChannel, &sale)
		}

		return nil
	})
}

func (w *Worker) StartUserEventWorker() error {
	return w.consumer.ConsumeJSON(UserEventsQueue, func(data interface{}) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		var event models.UserEvent
		if err := json.Unmarshal(jsonData, &event); err != nil {
			return fmt.Errorf("failed to unmarshal user event: %w", err)
		}

		if err := w.userEventService.CreateEvent(&event); err != nil {
			return err
		}

		if w.redisClient != nil {
			w.redisClient.Publish(redis.UserEventsChannel, &event)
		}

		return nil
	})
}

func (w *Worker) StartFinancialWorker() error {
	return w.consumer.ConsumeJSON(FinancialQueue, func(data interface{}) error {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}

		var metric models.FinancialMetric
		if err := json.Unmarshal(jsonData, &metric); err != nil {
			return fmt.Errorf("failed to unmarshal financial metric: %w", err)
		}

		if err := w.financialService.CreateMetric(&metric); err != nil {
			return err
		}

		if w.redisClient != nil {
			w.redisClient.Publish(redis.FinancialChannel, &metric)
		}

		return nil
	})
}

func (w *Worker) StartBatchStockWorker() error {
	return w.consumer.Consume(StockQueue, func(body []byte) error {
		var quote models.StockQuote
		if err := json.Unmarshal(body, &quote); err != nil {
			return fmt.Errorf("failed to unmarshal stock quote: %w", err)
		}

		w.batchMutex.Lock()
		w.batchBuffer[StockQueue] = append(w.batchBuffer[StockQueue], &quote)
		batch := w.batchBuffer[StockQueue]
		w.batchMutex.Unlock()

		if len(batch) >= w.batchSize {
			quotes := make([]*models.StockQuote, len(batch))
			for i, v := range batch {
				quotes[i] = v.(*models.StockQuote)
			}
			if err := w.stockService.BatchCreateQuotes(quotes); err != nil {
				return err
			}
			w.batchMutex.Lock()
			w.batchBuffer[StockQueue] = w.batchBuffer[StockQueue][:0]
			w.batchMutex.Unlock()
		}

		return nil
	})
}

func (w *Worker) StartAll() error {
	if err := w.StartStockWorker(); err != nil {
		return fmt.Errorf("failed to start stock worker: %w", err)
	}
	if err := w.StartSaleWorker(); err != nil {
		return fmt.Errorf("failed to start sale worker: %w", err)
	}
	if err := w.StartUserEventWorker(); err != nil {
		return fmt.Errorf("failed to start user event worker: %w", err)
	}
	if err := w.StartFinancialWorker(); err != nil {
		return fmt.Errorf("failed to start financial worker: %w", err)
	}

	log.Println("All workers started")
	return nil
}