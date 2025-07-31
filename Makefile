# Build the application
build:
	go build -o bin/todo-app src/main.go

# Run the application
run:
	go run src/main.go

# Build and run
dev: build
	./script/kill-port.sh
	./bin/todo-app

# Hot reload with Air
hot: install-air
	./script/kill-port.sh
	air

# Clean build artifacts
clean:
	rm -f bin/todo-app

# Install dependencies
deps:
	go mod download

# Install Air for hot reload
install-air:
	@if ! command -v air > /dev/null; then \
		echo "Installing Air..."; \
		go install github.com/cosmtrek/air@latest; \
	else \
		echo "Air is already installed"; \
	fi

# Format code
fmt:
	go fmt ./src/...

# Run tests
test:
	go test ./src/...

# Database operations
db-init:
	sqlite3 src/db/todo.db < src/db/schema.sql

# Seed sample data
seed-data:
	go run script/seed_data.go

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  run         - Run the application directly"
	@echo "  dev         - Build and run the application"
	@echo "  hot         - Run with hot reload using Air"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  install-air - Install Air for hot reload"
	@echo "  fmt         - Format code"
	@echo "  test        - Run tests"
	@echo "  db-init     - Initialize database"
	@echo "  seed-data   - Add sample data for testing"
	@echo "  help        - Show this help"

.PHONY: build run dev hot clean deps install-air fmt test db-init seed-data help
