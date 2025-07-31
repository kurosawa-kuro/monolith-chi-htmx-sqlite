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

# Clean build artifacts
clean:
	rm -f bin/todo-app

# Install dependencies
deps:
	go mod download

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
	@echo "  build     - Build the application"
	@echo "  run       - Run the application directly"
	@echo "  dev       - Build and run the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  deps      - Install dependencies"
	@echo "  fmt       - Format code"
	@echo "  test      - Run tests"
	@echo "  db-init   - Initialize database"
	@echo "  seed-data - Add sample data for testing"
	@echo "  help      - Show this help"

.PHONY: build run dev clean deps fmt test db-init seed-data help
