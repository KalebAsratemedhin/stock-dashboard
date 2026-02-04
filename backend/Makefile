.PHONY: run-api run-generator run-simulator run-fetcher build test migrate migrate-up migrate-down migrate-version migrate-create

# Database connection string - can be overridden via environment
DB_URL ?= postgres://postgres:postgres@localhost:5432/stock_dashboard?sslmode=disable
MIGRATIONS_PATH ?= migrations

# Run services
run-api:
	go run ./cmd/api

run-generator:
	go run ./cmd/data-generator

run-simulator:
	go run ./cmd/data-simulator

run-fetcher:
	go run ./cmd/stock-fetcher

# Build
build:
	go build -o bin/api ./cmd/api
	go build -o bin/generator ./cmd/data-generator
	go build -o bin/simulator ./cmd/data-simulator
	go build -o bin/fetcher ./cmd/stock-fetcher

# Test
test:
	go test ./...

# Migrate - uses Go program (reads from .env)
migrate:
	go run ./cmd/migrate

# Direct migrate commands (uses migrate CLI tool)
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down 1

migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URL)" version

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-clean:
	docker exec -i stock-dashboard-db psql -U postgres -d stock_dashboard -c "DROP TABLE IF EXISTS schema_migrations"

run-worker:
	go run ./cmd/worker

# Dependencies
deps:
	go mod download
	go mod tidy