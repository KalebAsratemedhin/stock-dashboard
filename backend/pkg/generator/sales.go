package generator

import (
	"math/rand"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
)

type SalesGenerator struct {
	products   []Product
	regions    []string
	categories []string
	rng        *rand.Rand
}

type Product struct {
	ID       string
	Name     string
	Category string
	BasePrice float64
}

func NewSalesGenerator() *SalesGenerator {
	products := []Product{
		{ID: "PROD001", Name: "Laptop Pro", Category: "Electronics", BasePrice: 1299.99},
		{ID: "PROD002", Name: "Wireless Mouse", Category: "Electronics", BasePrice: 29.99},
		{ID: "PROD003", Name: "Office Chair", Category: "Furniture", BasePrice: 299.99},
		{ID: "PROD004", Name: "Desk Lamp", Category: "Furniture", BasePrice: 49.99},
		{ID: "PROD005", Name: "Notebook Set", Category: "Stationery", BasePrice: 19.99},
		{ID: "PROD006", Name: "Pen Pack", Category: "Stationery", BasePrice: 9.99},
	}

	regions := []string{"North America", "Europe", "Asia", "South America", "Africa", "Oceania"}
	categories := []string{"Electronics", "Furniture", "Stationery", "Clothing", "Food"}

	return &SalesGenerator{
		products:   products,
		regions:    regions,
		categories: categories,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (sg *SalesGenerator) GenerateSale(timestamp time.Time) *models.Sale {
	product := sg.products[sg.rng.Intn(len(sg.products))]
	region := sg.regions[sg.rng.Intn(len(sg.regions))]
	
	quantity := 1 + sg.rng.Intn(5)
	priceVariation := 0.8 + sg.rng.Float64()*0.4
	unitPrice := product.BasePrice * priceVariation
	
	discount := 0.0
	if sg.rng.Float64() < 0.3 {
		discount = 5.0 + sg.rng.Float64()*20.0
	}
	
	revenue := float64(quantity) * unitPrice * (1 - discount/100)
	customerID := "CUST" + generateRandomID(6)

	return &models.Sale{
		Timestamp:   timestamp,
		ProductID:   product.ID,
		ProductName: product.Name,
		Category:    product.Category,
		CustomerID:  customerID,
		Region:      region,
		Quantity:    quantity,
		UnitPrice:   unitPrice,
		Discount:    discount,
		Revenue:     revenue,
	}
}

