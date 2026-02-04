# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/simulator ./cmd/data-simulator
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/generator ./cmd/data-generator

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy binaries from builder
COPY --from=builder /app/bin/* /app/bin/

# Copy migrations
COPY --from=builder /app/migrations /app/migrations

# Default command (can be overridden in docker-compose)
CMD ["/app/bin/api"]

