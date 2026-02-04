package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Get database URL from environment or construct it
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Construct from individual environment variables
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")
		dbName := getEnv("DB_NAME", "stock_dashboard")
		sslMode := getEnv("DB_SSLMODE", "disable")
		
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbName, sslMode)
	}

	// Get migrations directory (relative to where the binary is run from)
	migrationsPath := getEnv("MIGRATIONS_PATH", "migrations")
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		log.Fatalf("Migrations directory not found: %s", migrationsPath)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No migrations to apply - database is up to date")
			return
		}
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}