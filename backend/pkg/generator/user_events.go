package generator

import (
	"math/rand"
	"time"

	"github.com/KalebAsratemedhi/stock-dashboard/backend/internal/models"
)

type UserEventGenerator struct {
	pages      []string
	devices    []string
	browsers   []string
	countries  []string
	eventTypes []string
	rng        *rand.Rand
}

func NewUserEventGenerator() *UserEventGenerator {
	return &UserEventGenerator{
		pages: []string{
			"/", "/products", "/about", "/contact", "/dashboard",
			"/login", "/signup", "/cart", "/checkout", "/profile",
		},
		devices: []string{"desktop", "mobile", "tablet"},
		browsers: []string{"Chrome", "Firefox", "Safari", "Edge"},
		countries: []string{"US", "UK", "CA", "DE", "FR", "JP", "AU", "BR"},
		eventTypes: []string{"page_view", "click", "purchase", "signup", "login"},
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ueg *UserEventGenerator) GenerateEvent(timestamp time.Time) *models.UserEvent {
	eventType := ueg.eventTypes[ueg.rng.Intn(len(ueg.eventTypes))]
	page := ueg.pages[ueg.rng.Intn(len(ueg.pages))]
	device := ueg.devices[ueg.rng.Intn(len(ueg.devices))]
	browser := ueg.browsers[ueg.rng.Intn(len(ueg.browsers))]
	country := ueg.countries[ueg.rng.Intn(len(ueg.countries))]

	userID := ""
	sessionID := ""
	if ueg.rng.Float64() < 0.7 {
		userID = "USER" + generateRandomID(8)
		sessionID = "SESS" + generateRandomID(12)
	}

	city := generateCity(country)

	return &models.UserEvent{
		Timestamp: timestamp,
		EventType: eventType,
		UserID:    userID,
		SessionID: sessionID,
		Page:      page,
		Device:    device,
		Browser:   browser,
		Country:   country,
		City:      city,
		Referrer:  generateReferrer(),
		Metadata:  "{}",
	}
}

