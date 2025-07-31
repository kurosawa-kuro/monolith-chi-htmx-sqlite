# =============================================================================
# Go Monolith Application Makefile
# =============================================================================

# Application Configuration
APP_NAME := todo-app
APP_BINARY := bin/$(APP_NAME)
MAIN_FILE := src/main.go
SOURCE_DIR := src

# Database Configuration
DB_PATH := src/db/todo.db
DB_SCHEMA := src/db/schema.sql

# Development Configuration
DEFAULT_PORT := 8080
DEFAULT_HOST := localhost
DEFAULT_ENV := development

# Build Configuration
BUILD_FLAGS := -v
TEST_FLAGS := -v -race -cover

# Scripts
KILL_PORT_SCRIPT := ./script/kill-port.sh
SEED_DATA_SCRIPT := script/seed_data.go

# Air Configuration
AIR_CONFIG := .air.toml

# =============================================================================
# Development Commands
# =============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Go Monolith Application - Available Commands:"
	@echo ""
	@echo "Development Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | grep -v "test-e2e" | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "E2E Testing Commands:"
	@grep -E '^test-e2e[a-zA-Z_-]*:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[32m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Environment Variables:"
	@echo "  PORT    - Server port (default: $(DEFAULT_PORT))"
	@echo "  HOST    - Server host (default: $(DEFAULT_HOST))"
	@echo "  ENV     - Environment (default: $(DEFAULT_ENV))"
	@echo "  DB_PATH - Database path (default: $(DB_PATH))"

.PHONY: build
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	@mkdir -p bin
	@go build $(BUILD_FLAGS) -o $(APP_BINARY) $(MAIN_FILE)
	@echo "Build completed: $(APP_BINARY)"

.PHONY: run
run: ## Run the application directly
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_FILE)

.PHONY: run-dev
run-dev: ## Run with development configuration
	@echo "Running $(APP_NAME) in development mode..."
	@ENV=$(DEFAULT_ENV) PORT=$(DEFAULT_PORT) HOST=$(DEFAULT_HOST) go run $(MAIN_FILE)

.PHONY: dev
dev: build ## Build and run the application
	@echo "Starting $(APP_NAME)..."
	@$(KILL_PORT_SCRIPT) || true
	@$(APP_BINARY)

.PHONY: hot
hot: install-air ## Run with hot reload using Air
	@echo "Starting hot reload with Air..."
	@$(KILL_PORT_SCRIPT) || true
	@air -c $(AIR_CONFIG)

# =============================================================================
# Build and Clean Commands
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -f $(APP_BINARY)
	@rm -rf tmp/
	@rm -f build-errors.log
	@rm -rf test-results/
	@rm -f coverage.out
	@rm -f coverage.html
	@echo "Clean completed"

.PHONY: clean-all
clean-all: clean ## Clean all artifacts including database
	@echo "Cleaning all artifacts..."
	@rm -f $(DB_PATH)
	@rm -rf bin/
	@echo "All artifacts cleaned"

# =============================================================================
# Dependency Management
# =============================================================================

.PHONY: deps
deps: ## Install Go dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

.PHONY: install-air
install-air: ## Install Air for hot reload
	@echo "Checking Air installation..."
	@if ! command -v air > /dev/null; then \
		echo "Installing Air..."; \
		go install github.com/cosmtrek/air@latest; \
		echo "Air installed successfully"; \
	else \
		echo "Air is already installed"; \
	fi

# =============================================================================
# Code Quality Commands
# =============================================================================

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./$(SOURCE_DIR)/...
	@echo "Code formatting completed"

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./$(SOURCE_DIR)/...; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./$(SOURCE_DIR)/...
	@echo "Go vet completed"

.PHONY: check
check: fmt lint vet ## Run all code quality checks

# =============================================================================
# Testing Commands
# =============================================================================

.PHONY: test
test: ## Run Go tests
	@echo "Running Go tests..."
	@go test $(TEST_FLAGS) ./$(SOURCE_DIR)/...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test $(TEST_FLAGS) ./src/tests/models/... ./src/tests/services/... ./src/tests/handlers/...

.PHONY: test-integration
test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test $(TEST_FLAGS) ./src/tests/integration/...

.PHONY: test-all-go
test-all-go: test-unit test-integration ## Run all Go tests (unit + integration)

.PHONY: test-short
test-short: ## Run Go tests with short flag
	@echo "Running short Go tests..."
	@go test -short ./$(SOURCE_DIR)/...

.PHONY: test-coverage
test-coverage: ## Run Go tests with coverage report
	@echo "Running Go tests with coverage..."
	@go test $(TEST_FLAGS) -coverprofile=coverage.out ./$(SOURCE_DIR)/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-bench
test-bench: ## Run Go benchmark tests
	@echo "Running Go benchmark tests..."
	@go test -bench=. ./$(SOURCE_DIR)/...

# =============================================================================
# E2E Testing Commands (Playwright)
# =============================================================================

.PHONY: test-e2e
test-e2e: ## Run E2E tests with Playwright
	@echo "Running E2E tests..."
	@npm run test:e2e

.PHONY: test-e2e-simple
test-e2e-simple: ## Run simple E2E tests
	@echo "Running simple E2E tests..."
	@npm run test:e2e:simple

.PHONY: test-e2e-ui
test-e2e-ui: ## Run E2E tests with UI mode
	@echo "Running E2E tests with UI..."
	@npm run test:e2e:ui

.PHONY: test-e2e-headed
test-e2e-headed: ## Run E2E tests in headed mode
	@echo "Running E2E tests in headed mode..."
	@npm run test:e2e:headed

.PHONY: test-e2e-debug
test-e2e-debug: ## Run E2E tests in debug mode
	@echo "Running E2E tests in debug mode..."
	@npm run test:e2e:debug

.PHONY: test-e2e-report
test-e2e-report: ## Show E2E test report
	@echo "Opening E2E test report..."
	@npm run test:e2e:report

.PHONY: test-e2e-install
test-e2e-install: ## Install Playwright browsers
	@echo "Installing Playwright browsers..."
	@npm run test:e2e:install

.PHONY: test-e2e-codegen
test-e2e-codegen: ## Generate E2E test code
	@echo "Starting Playwright codegen..."
	@npm run test:e2e:codegen

.PHONY: test-e2e-update
test-e2e-update: ## Update E2E test snapshots
	@echo "Updating E2E test snapshots..."
	@npm run test:e2e:update

.PHONY: test-e2e-parallel
test-e2e-parallel: ## Run E2E tests in parallel
	@echo "Running E2E tests in parallel..."
	@npm run test:e2e:parallel

.PHONY: test-e2e-retry
test-e2e-retry: ## Run E2E tests with retry
	@echo "Running E2E tests with retry..."
	@npm run test:e2e:retry

.PHONY: test-all
test-all: test test-e2e ## Run all tests (Go + E2E)

# =============================================================================
# Database Commands
# =============================================================================

.PHONY: db-init
db-init: ## Initialize database
	@echo "Initializing database..."
	@mkdir -p $(dir $(DB_PATH))
	@if [ -f $(DB_SCHEMA) ]; then \
		sqlite3 $(DB_PATH) < $(DB_SCHEMA); \
		echo "Database initialized: $(DB_PATH)"; \
	else \
		echo "Schema file not found: $(DB_SCHEMA)"; \
		exit 1; \
	fi

.PHONY: db-reset
db-reset: ## Reset database (delete and reinitialize)
	@echo "Resetting database..."
	@rm -f $(DB_PATH)
	@$(MAKE) db-init

.PHONY: seed-data
seed-data: db-init ## Seed sample data for testing
	@echo "Seeding sample data..."
	@if [ -f $(SEED_DATA_SCRIPT) ]; then \
		go run $(SEED_DATA_SCRIPT); \
		echo "Sample data seeded"; \
	else \
		echo "Seed data script not found: $(SEED_DATA_SCRIPT)"; \
		exit 1; \
	fi

# =============================================================================
# Utility Commands
# =============================================================================

.PHONY: kill-port
kill-port: ## Kill process running on default port
	@echo "Killing process on port $(DEFAULT_PORT)..."
	@$(KILL_PORT_SCRIPT) || true

.PHONY: status
status: ## Show application status
	@echo "Application Status:"
	@echo "  Binary: $(APP_BINARY)"
	@echo "  Database: $(DB_PATH)"
	@echo "  Port: $(DEFAULT_PORT)"
	@echo "  Environment: $(DEFAULT_ENV)"
	@if [ -f $(APP_BINARY) ]; then \
		echo "  Binary exists: Yes"; \
	else \
		echo "  Binary exists: No"; \
	fi
	@if [ -f $(DB_PATH) ]; then \
		echo "  Database exists: Yes"; \
	else \
		echo "  Database exists: No"; \
	fi

.PHONY: logs
logs: ## Show recent logs (if available)
	@echo "Recent logs:"
	@if [ -f build-errors.log ]; then \
		echo "Build errors:"; \
		tail -10 build-errors.log; \
	else \
		echo "No build error log found"; \
	fi

# =============================================================================
# Development Workflow
# =============================================================================

.PHONY: setup
setup: deps install-air db-init seed-data ## Complete development setup
	@echo "Development environment setup completed!"

.PHONY: setup-e2e
setup-e2e: ## Setup E2E testing environment
	@echo "Setting up E2E testing environment..."
	@npm install
	@npm run test:e2e:install
	@echo "E2E testing environment setup completed!"

.PHONY: setup-full
setup-full: setup setup-e2e ## Complete setup including E2E testing
	@echo "Full development environment setup completed!"

.PHONY: start
start: build db-init ## Quick start (build and run with database)
	@echo "Starting application..."
	@$(KILL_PORT_SCRIPT) || true
	@$(APP_BINARY)

.PHONY: restart
restart: kill-port start ## Restart the application

# =============================================================================
# Default Target
# =============================================================================

.DEFAULT_GOAL := help
