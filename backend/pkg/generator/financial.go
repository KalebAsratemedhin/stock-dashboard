package generator

import (
	"math/rand"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
)

type FinancialGenerator struct {
	departments []string
	metricTypes []string
	categories  []string
	rng         *rand.Rand
}

func NewFinancialGenerator() *FinancialGenerator {
	return &FinancialGenerator{
		departments: []string{"Sales", "Marketing", "Engineering", "Operations", "HR", "Finance"},
		metricTypes: []string{"revenue", "expense", "profit", "margin"},
		categories:  []string{"Salary", "Infrastructure", "Marketing", "R&D", "Operations"},
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (fg *FinancialGenerator) GenerateMetric(timestamp time.Time) *models.FinancialMetric {
	metricType := fg.metricTypes[fg.rng.Intn(len(fg.metricTypes))]
	department := fg.departments[fg.rng.Intn(len(fg.departments))]
	category := fg.categories[fg.rng.Intn(len(fg.categories))]

	var amount float64
	var budget float64

	switch metricType {
	case "revenue":
		amount = 50000 + fg.rng.Float64()*200000
		budget = amount * (0.9 + fg.rng.Float64()*0.2)
	case "expense":
		amount = 10000 + fg.rng.Float64()*50000
		budget = amount * (0.8 + fg.rng.Float64()*0.3)
	case "profit":
		amount = 10000 + fg.rng.Float64()*100000
		budget = amount * (0.85 + fg.rng.Float64()*0.25)
	case "margin":
		amount = 10 + fg.rng.Float64()*40
		budget = amount * (0.9 + fg.rng.Float64()*0.2)
	}

	variance := amount - budget
	var variancePct float64
	if budget != 0 {
		variancePct = (variance / budget) * 100
	}

	period := "daily"
	if timestamp.Weekday() == time.Monday {
		period = "weekly"
	}
	if timestamp.Day() == 1 {
		period = "monthly"
	}

	return &models.FinancialMetric{
		Timestamp:   timestamp,
		MetricType:  metricType,
		Department:  department,
		Category:    category,
		Amount:      amount,
		Budget:      budget,
		Variance:    variance,
		VariancePct: variancePct,
		Period:      period,
	}
}

