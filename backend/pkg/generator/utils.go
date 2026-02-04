package generator

import (
	"math/rand"
)

var referrers = []string{
	"https://google.com",
	"https://bing.com",
	"https://direct",
	"https://facebook.com",
	"https://twitter.com",
	"",
}

var cityMap = map[string][]string{
	"US": {"New York", "Los Angeles", "Chicago", "Houston", "Phoenix"},
	"UK": {"London", "Manchester", "Birmingham", "Liverpool", "Leeds"},
	"CA": {"Toronto", "Vancouver", "Montreal", "Calgary", "Ottawa"},
	"DE": {"Berlin", "Munich", "Hamburg", "Frankfurt", "Cologne"},
	"FR": {"Paris", "Lyon", "Marseille", "Toulouse", "Nice"},
	"JP": {"Tokyo", "Osaka", "Yokohama", "Nagoya", "Sapporo"},
	"AU": {"Sydney", "Melbourne", "Brisbane", "Perth", "Adelaide"},
	"BR": {"São Paulo", "Rio de Janeiro", "Brasília", "Salvador", "Fortaleza"},
}

func generateRandomID(length int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func generateReferrer() string {
	return referrers[rand.Intn(len(referrers))]
}

func generateCity(country string) string {
	cities, ok := cityMap[country]
	if !ok || len(cities) == 0 {
		return "Unknown"
	}
	return cities[rand.Intn(len(cities))]
}

