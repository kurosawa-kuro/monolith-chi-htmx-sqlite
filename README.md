# Go Monolith Application with Chi, HTMX, and SQLite

A modern monolithic web application built with Go, Chi router, HTMX for dynamic interactions, and SQLite for data persistence.

## Features

- **Go Backend**: Fast and efficient server-side logic
- **Chi Router**: Lightweight HTTP router
- **HTMX**: Dynamic UI interactions without JavaScript
- **SQLite**: Lightweight database
- **E2E Testing**: Comprehensive end-to-end testing with Playwright
- **Hot Reload**: Development with Air for automatic reloading

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher (for E2E testing)
- SQLite3

## Quick Start

### 1. Clone the repository
```bash
git clone <repository-url>
cd monolith-chi-htmx-sqlite
```

### 2. Setup development environment
```bash
# Complete setup (Go + E2E testing)
make setup-full

# Or setup step by step
make setup      # Go development environment
make setup-e2e  # E2E testing environment
```

### 3. Run the application
```bash
# Development mode with hot reload
make hot

# Or run directly
make run-dev
```

### 4. Access the application
Open your browser and navigate to `http://localhost:8080`

## Development

### Available Commands

#### Application Management
```bash
make build          # Build the application
make run            # Run the application
make run-dev        # Run with development configuration
make dev            # Build and run
make hot            # Run with hot reload
make start          # Quick start with database
make restart        # Restart the application
```

#### Testing
```bash
# Go Tests
make test           # Run Go tests
make test-short     # Run short Go tests
make test-coverage  # Run tests with coverage
make test-bench     # Run benchmark tests

# E2E Tests (Playwright)
make test-e2e       # Run E2E tests
make test-e2e-ui    # Run E2E tests with UI
make test-e2e-headed # Run E2E tests in headed mode
make test-e2e-debug # Run E2E tests in debug mode
make test-e2e-report # Show E2E test report
make test-e2e-codegen # Generate E2E test code

# All Tests
make test-all       # Run all tests (Go + E2E)
```

#### Database Management
```bash
make db-init        # Initialize database
make db-reset       # Reset database
make seed-data      # Seed sample data
```

#### Code Quality
```bash
make fmt            # Format Go code
make lint           # Run linter
make vet            # Run go vet
make check          # Run all code quality checks
```

#### Utilities
```bash
make status         # Show application status
make logs           # Show recent logs
make kill-port      # Kill process on default port
make clean          # Clean build artifacts
make clean-all      # Clean all artifacts including database
```

### E2E Testing with Playwright

This project includes comprehensive E2E testing using Playwright. The tests cover:

- **Home Page**: Basic page loading and navigation
- **Todo CRUD**: Create, read, update, delete operations
- **Category Management**: Category creation and deletion
- **Filtering**: Todo filtering by category
- **Pagination**: Page navigation

#### Running E2E Tests

```bash
# Install Playwright browsers (first time only)
make test-e2e-install

# Run all E2E tests
make test-e2e

# Run with UI mode (interactive)
make test-e2e-ui

# Run in headed mode (see browser)
make test-e2e-headed

# Debug mode
make test-e2e-debug

# Generate test code
make test-e2e-codegen
```

#### E2E Test Structure

```
tests/e2e/
├── global-setup.ts      # Global test setup
├── global-teardown.ts   # Global test cleanup
├── home.spec.ts         # Home page tests
├── todo-crud.spec.ts    # Todo CRUD tests
├── category-management.spec.ts  # Category tests
└── utils/
    └── test-helpers.ts  # Test helper functions
```

#### Test Data Attributes

The application uses `data-testid` attributes for reliable test selection:

- `todo-title-input`: Todo title input field
- `create-todo-button`: Create todo button
- `todo-item`: Individual todo item
- `status-toggle`: Todo status toggle
- `delete-todo-button`: Delete todo button
- `category-filter`: Category filter dropdown
- `category-name-input`: Category name input
- `create-category-button`: Create category button
- `pagination`: Pagination container

## Project Structure

```
.
├── src/                    # Go source code
│   ├── config/            # Configuration management
│   ├── errors/            # Error handling
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models
│   ├── services/          # Business logic
│   ├── templates/         # Template management
│   ├── validation/        # Input validation
│   ├── db/               # Database files
│   ├── static/           # Static assets
│   └── main.go           # Application entry point
├── tests/                 # Test files
│   └── e2e/              # E2E tests (Playwright)
├── script/               # Utility scripts
├── bin/                  # Build artifacts
├── docs/                 # Documentation
├── Makefile              # Build automation
├── package.json          # Node.js dependencies
├── playwright.config.ts  # Playwright configuration
└── README.md            # This file
```

## Configuration

The application can be configured using environment variables:

```bash
# Server Configuration
HOST=localhost
PORT=8080

# Environment
ENV=development

# Database Configuration
DB_PATH=./src/db/todo.db

# Application Configuration
DEFAULT_PAGE_SIZE=10
MAX_PAGE_SIZE=100
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test-all`
5. Submit a pull request

## License

MIT License - see LICENSE file for details.