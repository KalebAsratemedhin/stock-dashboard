package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/services"
)

type Handler struct {
	stockService    services.StockService
	saleService     services.SaleService
	userEventService services.UserEventService
	financialService services.FinancialMetricService
}

func NewHandler(
	stockService services.StockService,
	saleService services.SaleService,
	userEventService services.UserEventService,
	financialService services.FinancialMetricService,
) *Handler {
	return &Handler{
		stockService:    stockService,
		saleService:     saleService,
		userEventService: userEventService,
		financialService: financialService,
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) GetStocks(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	limitStr := r.URL.Query().Get("limit")
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	var quotes interface{}
	var err error

	if symbol != "" {
		quotes, err = h.stockService.GetQuotesBySymbol(symbol, limit)
	} else {
		quotes, err = h.stockService.GetLatestQuotes([]string{})
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, quotes)
}

func (h *Handler) GetStocksByTimeRange(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if symbol == "" || startStr == "" || endStr == "" {
		jsonError(w, http.StatusBadRequest, "Missing required parameters: symbol, start, end")
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid start time format")
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid end time format")
		return
	}

	quotes, err := h.stockService.GetQuotesByTimeRange(symbol, start, end)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, quotes)
}

func (h *Handler) GetSales(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	category := r.URL.Query().Get("category")
	region := r.URL.Query().Get("region")

	var sales interface{}
	var err error

	if startStr != "" && endStr != "" {
		start, err1 := time.Parse(time.RFC3339, startStr)
		end, err2 := time.Parse(time.RFC3339, endStr)
		if err1 != nil || err2 != nil {
			jsonError(w, http.StatusBadRequest, "Invalid time format")
			return
		}
		sales, err = h.saleService.GetSalesByTimeRangeWithFilters(start, end, category, region)
	} else {
		limitStr := r.URL.Query().Get("limit")
		limit := 100
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}
		sales, err = h.saleService.GetRecentSales(limit)
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, sales)
}

func (h *Handler) GetSalesRevenue(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		jsonError(w, http.StatusBadRequest, "Missing required parameters: start, end")
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid start time format")
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		jsonError(w, http.StatusBadRequest, "Invalid end time format")
		return
	}

	revenue, err := h.saleService.GetTotalRevenue(start, end)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"revenue": revenue})
}

func (h *Handler) GetUserEvents(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	eventType := r.URL.Query().Get("type")

	var events interface{}
	var err error

	if startStr != "" && endStr != "" {
		start, err1 := time.Parse(time.RFC3339, startStr)
		end, err2 := time.Parse(time.RFC3339, endStr)
		if err1 != nil || err2 != nil {
			jsonError(w, http.StatusBadRequest, "Invalid time format")
			return
		}
		events, err = h.userEventService.GetEventsByTimeRange(start, end)
	} else {
		limitStr := r.URL.Query().Get("limit")
		limit := 100
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}
		if eventType != "" {
			events, err = h.userEventService.GetEventsByType(eventType, limit)
		} else {
			events, err = h.userEventService.GetRecentEvents(limit)
		}
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, events)
}

func (h *Handler) GetFinancialMetrics(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	metricType := r.URL.Query().Get("type")
	department := r.URL.Query().Get("department")

	var metrics interface{}
	var err error

	if startStr != "" && endStr != "" {
		start, err1 := time.Parse(time.RFC3339, startStr)
		end, err2 := time.Parse(time.RFC3339, endStr)
		if err1 != nil || err2 != nil {
			jsonError(w, http.StatusBadRequest, "Invalid time format")
			return
		}
		metrics, err = h.financialService.GetMetricsByTimeRange(start, end)
	} else {
		limitStr := r.URL.Query().Get("limit")
		limit := 100
		if limitStr != "" {
			limit, _ = strconv.Atoi(limitStr)
		}
		if metricType != "" {
			metrics, err = h.financialService.GetMetricsByType(metricType, limit)
		} else if department != "" {
			metrics, err = h.financialService.GetMetricsByDepartment(department, limit)
		} else {
			metrics, err = h.financialService.GetRecentMetrics(limit)
		}
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse(w, http.StatusOK, metrics)
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

