package generator

import (
	"math"
	"math/rand"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
)

type StockGenerator struct {
	symbols      []string
	basePrices   map[string]float64
	volatilities map[string]float64
	rng          *rand.Rand
}

func NewStockGenerator(symbols []string) *StockGenerator {
	sg := &StockGenerator{
		symbols:      symbols,
		basePrices:   make(map[string]float64),
		volatilities: make(map[string]float64),
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, symbol := range symbols {
		sg.basePrices[symbol] = 50 + sg.rng.Float64()*200
		// Per-step volatility (0.04–0.12) gives ~0.5–2% moves per tick for more realistic charts
		sg.volatilities[symbol] = 0.04 + sg.rng.Float64()*0.08
	}

	return sg
}

func (sg *StockGenerator) GenerateQuote(symbol string, prevQuote *models.StockQuote, timestamp time.Time) *models.StockQuote {
	var price float64
	var volume int64

	if prevQuote == nil {
		price = sg.basePrices[symbol]
		volume = int64(1000000 + sg.rng.Intn(5000000))
	} else {
		price = sg.nextPrice(symbol, prevQuote.Close)
		volume = sg.nextVolume(prevQuote.Volume, price, prevQuote.Close)
	}

	volatility := sg.volatilities[symbol]
	// Wider high/low range for more realistic bar charts
	barRange := volatility * 1.5
	var open, high, low float64
	if prevQuote == nil {
		open = price
		high = price * (1 + barRange)
		low = price * (1 - barRange)
	} else {
		open = prevQuote.Close
		high = math.Max(open, price) * (1 + barRange*0.5)
		low = math.Min(open, price) * (1 - barRange*0.5)
	}
	if high < price {
		high = price
	}
	if low > price {
		low = price
	}

	changeAmount := price - open
	changePct := (changeAmount / open) * 100

	bid := price * 0.9999
	ask := price * 1.0001

	return &models.StockQuote{
		Symbol:    symbol,
		Timestamp: timestamp,
		Open:      open,
		High:      high,
		Low:       low,
		Close:     price,
		Volume:    volume,
		Bid:       bid,
		Ask:       ask,
		Change:    changeAmount,
		ChangePct: changePct,
	}
}

func (sg *StockGenerator) nextPrice(symbol string, prevPrice float64) float64 {
	volatility := sg.volatilities[symbol]
	// Light mean reversion so prices don't drift to infinity but still show clear moves
	meanReversion := 0.02
	basePrice := sg.basePrices[symbol]

	// Scale volatility so each 30s step has visible movement (~0.5–2% typical)
	stepVol := volatility * 1.2
	random := sg.rng.NormFloat64() * stepVol
	drift := meanReversion * (basePrice - prevPrice) / basePrice

	newPrice := prevPrice * math.Exp(drift + random)

	if newPrice < 1.0 {
		newPrice = 1.0
	}

	return newPrice
}

func (sg *StockGenerator) nextVolume(prevVolume int64, price, prevPrice float64) int64 {
	priceChange := math.Abs((price - prevPrice) / prevPrice)
	volumeMultiplier := 1.0 + priceChange*2.0

	baseVolume := float64(prevVolume) * volumeMultiplier
	noise := 0.8 + sg.rng.Float64()*0.4

	newVolume := int64(baseVolume * noise)
	if newVolume < 100000 {
		newVolume = 100000
	}
	if newVolume > 10000000 {
		newVolume = 10000000
	}

	return newVolume
}

