# Library System - Makefile
#
# Usage:
#   make help          Show this help message
#   make setup         Install dependencies and tools
#   make run           Run the API server
#   make test          Run all tests
#   make migrate-up    Apply all migrations
#   make migrate-down  Rollback last migration

# Database configuration
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= library
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Colors for output
GREEN  := \033[0;32m
YELLOW := \033[0;33m
NC     := \033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "Library System - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# Setup
# =============================================================================

.PHONY: setup
setup: ## Install dependencies and tools
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(YELLOW)Installing golang-migrate...$(NC)"
	@which migrate > /dev/null || brew install golang-migrate
	@echo "$(GREEN)Setup complete!$(NC)"

# =============================================================================
# Development
# =============================================================================

.PHONY: run
run: ## Run the API server
	go run cmd/api/main.go

.PHONY: build
build: ## Build the application
	go build -o bin/api cmd/api/main.go

.PHONY: test
test: ## Run all tests
	go test ./... -v

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

.PHONY: lint
lint: ## Run linter
	@command -v golangci-lint > /dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(shell go env GOPATH)/bin/golangci-lint run

# =============================================================================
# Database
# =============================================================================

.PHONY: db-start
db-start: ## Start PostgreSQL (primary + 2 replicas)
	docker-compose up -d postgres-primary postgres-replica-1 postgres-replica-2
	@echo "$(YELLOW)Waiting for PostgreSQL cluster to be ready...$(NC)"
	@sleep 10
	@echo "$(GREEN)PostgreSQL cluster is ready!$(NC)"
	@echo "  Primary:   localhost:5432"
	@echo "  Replica 1: localhost:5433"
	@echo "  Replica 2: localhost:5434"

.PHONY: db-stop
db-stop: ## Stop PostgreSQL cluster
	docker-compose down

.PHONY: db-reset
db-reset: ## Reset database (drop and recreate)
	docker-compose down -v
	docker-compose up -d postgres-primary postgres-replica-1 postgres-replica-2
	@sleep 10
	@make migrate-up

# =============================================================================
# Migrations
# =============================================================================

.PHONY: migrate-up
migrate-up: ## Apply all migrations
	@echo "$(YELLOW)Applying migrations...$(NC)"
	migrate -path migrations -database "$(DB_URL)" up
	@echo "$(GREEN)Migrations applied!$(NC)"

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	migrate -path migrations -database "$(DB_URL)" down 1
	@echo "$(GREEN)Rollback complete!$(NC)"

.PHONY: migrate-down-all
migrate-down-all: ## Rollback all migrations
	@echo "$(YELLOW)Rolling back all migrations...$(NC)"
	migrate -path migrations -database "$(DB_URL)" down -all
	@echo "$(GREEN)All migrations rolled back!$(NC)"

.PHONY: migrate-version
migrate-version: ## Show current migration version
	migrate -path migrations -database "$(DB_URL)" version

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	migrate create -ext sql -dir migrations -seq $(name)
	@echo "$(GREEN)Created migration: $(name)$(NC)"

.PHONY: migrate-force
migrate-force: ## Force set migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=N"; exit 1; fi
	migrate -path migrations -database "$(DB_URL)" force $(version)

# =============================================================================
# Load Testing
# =============================================================================

.PHONY: load-smoke
load-smoke: ## Run smoke test (1 user, 10s)
	K6_WEB_DASHBOARD=true k6 run tests/load/smoke.js

.PHONY: load-test
load-test: ## Run load test (up to 1000 users)
	K6_WEB_DASHBOARD=true k6 run tests/load/load.js

.PHONY: load-stress
load-stress: ## Run stress test (up to 10000 users)
	K6_WEB_DASHBOARD=true k6 run tests/load/stress.js

# With Grafana dashboard (requires: make grafana-start)
.PHONY: load-smoke-grafana
load-smoke-grafana: ## Run smoke test with Grafana output
	k6 run --out influxdb=http://localhost:8086/k6 tests/load/smoke.js

.PHONY: load-test-grafana
load-test-grafana: ## Run load test with Grafana output
	k6 run --out influxdb=http://localhost:8086/k6 tests/load/load.js

.PHONY: load-stress-grafana
load-stress-grafana: ## Run stress test with Grafana output
	k6 run --out influxdb=http://localhost:8086/k6 tests/load/stress.js

# =============================================================================
# Grafana & InfluxDB
# =============================================================================

.PHONY: grafana-start
grafana-start: ## Start InfluxDB and Grafana
	docker-compose up -d influxdb grafana
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	@sleep 5
	@echo "$(GREEN)Grafana available at http://localhost:3000 (admin/admin)$(NC)"
	@echo "$(GREEN)InfluxDB available at http://localhost:8086$(NC)"

.PHONY: grafana-stop
grafana-stop: ## Stop InfluxDB and Grafana
	docker-compose stop influxdb grafana

# =============================================================================
# Docker
# =============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t library-system:latest .

.PHONY: docker-run
docker-run: ## Run with Docker Compose
	docker-compose up -d
	@make migrate-up
	@echo "$(GREEN)Application running at http://localhost:8080$(NC)"

.PHONY: docker-down
docker-down: ## Stop all containers
	docker-compose down

.PHONY: docker-logs
docker-logs: ## Show container logs
	docker-compose logs -f
