package api

import (
	"net/http"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/websocket"
)

func NewRouter(h *Handler, wsHub *websocket.Hub) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/api/stocks", h.GetStocks)
	mux.HandleFunc("/api/stocks/range", h.GetStocksByTimeRange)
	mux.HandleFunc("/api/sales", h.GetSales)
	mux.HandleFunc("/api/sales/revenue", h.GetSalesRevenue)
	mux.HandleFunc("/api/events", h.GetUserEvents)
	mux.HandleFunc("/api/metrics", h.GetFinancialMetrics)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleWebSocket(wsHub, w, r)
	})

	handler := ErrorHandler(Logging(Gzip(CORS(mux))))

	return handler
}

