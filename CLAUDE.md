# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a fully-implemented monolithic ToDo management application with category functionality, built using server-side rendering.

**Tech Stack:**
- Backend: Go 1.21+ with [chi](https://github.com/go-chi/chi) router
- Frontend: HTML templates + [htmx](https://htmx.org/) for dynamic interactions
- Database: SQLite with comprehensive schema
- Testing: Playwright for E2E testing
- Build automation: Makefile with comprehensive commands

## Development Commands

### Primary Development Workflow
```bash
# Complete development setup
make setup-full          # Sets up Go + E2E testing environment

# Development with hot reload
make hot                 # Requires Air (auto-installed)

# Standard development
make run-dev             # Run with development configuration
make dev                 # Build and run
```

### Build and Run Commands
```bash
make build               # Build application to bin/todo-app
make run                 # Run application directly
make start               # Quick start with database initialization
make restart             # Kill port and restart
```

### Testing Commands
```bash
# Go Tests
make test                # Run all Go tests with race detection
make test-coverage       # Generate coverage report (coverage.html)
make test-bench          # Run benchmark tests

# E2E Tests (Playwright)
make test-e2e            # Run all E2E tests
make test-e2e-ui         # Run with interactive UI
make test-e2e-headed     # Run in headed browser mode
make test-e2e-debug      # Debug mode
make test-e2e-codegen    # Generate test code

# Combined
make test-all            # Run both Go and E2E tests
```

### Code Quality Commands
```bash
make fmt                 # Format Go code
make lint                # Run golangci-lint (if installed)
make vet                 # Run go vet
make check               # Run all quality checks
```

### Database Management
```bash
make db-init             # Initialize SQLite database
make db-reset            # Delete and reinitialize database
make seed-data           # Add sample data for testing
```

## Architecture

### Database Schema (SQLite)
The application uses a normalized three-table structure:
- `todos`: Core todo entities with status, timestamps
- `categories`: Category entities with unique titles
- `todo_category`: Many-to-many junction table with foreign key constraints

### Application Structure
```
src/
├── main.go              # Application entry point with router setup
├── config/              # Configuration management
├── handlers/            # HTTP handlers (TodoHandler)
├── services/            # Business logic layer (TodoService)
├── models/              # Data models and repository patterns
├── middleware/          # Error handling, logging middleware
├── templates/           # Go HTML templates with template loader
├── validation/          # Input validation utilities
├── static/             # Static CSS assets
└── db/                 # Database schema and SQLite file
```

### Key Architectural Patterns
- **Repository Pattern**: `TodoRepository`, `CategoryRepository` for data access
- **Service Layer**: `TodoService` for business logic
- **Handler Layer**: HTTP request handling with proper error management
- **Middleware Chain**: Error handling, logging, compression, recovery
- **Template Loader**: Centralized template management with hot reloading support

## Core Features

### Todo Management
- Paginated todo listing with configurable page sizes
- CRUD operations with proper transaction handling
- Status toggling (incomplete ↔ complete)
- Category associations (many-to-many)

### Category Management
- Category creation and deletion
- Todo filtering by category
- Category usage count tracking

### Pagination & Filtering
- Configurable pagination (DEFAULT_PAGE_SIZE, MAX_PAGE_SIZE)
- Category-based filtering with pagination
- URL parameter validation

## Configuration

The application uses environment variables with sensible defaults:
- `HOST`: Server host (default: localhost)
- `PORT`: Server port (default: 8080)
- `ENV`: Environment mode (default: development)
- `DB_PATH`: Database file path (default: src/db/todo.db)

## Testing Strategy

### E2E Testing with Playwright
Comprehensive E2E test suite covering:
- Home page loading and navigation
- Todo CRUD operations
- Category management
- Filtering and pagination
- Dynamic HTMX interactions

### Test Data Attributes
Use these `data-testid` attributes for reliable test targeting:
- `todo-title-input`, `create-todo-button`
- `todo-item`, `status-toggle`, `delete-todo-button`
- `category-filter`, `category-name-input`, `create-category-button`
- `pagination`

## Development Best Practices

### Error Handling
- Custom error types in `src/errors/`
- Centralized error handling middleware
- Proper HTTP status codes and user-friendly messages

### Database Operations
- Transaction handling for multi-table operations
- Proper connection management and cleanup
- Foreign key constraints and cascade deletes

### Template Management
- Template loader with error handling
- Separation of layout and content templates
- Template hot reloading in development

## Utility Commands
```bash
make status              # Show application status and file existence
make logs                # Show recent build error logs
make kill-port           # Kill process on default port (8080)
make clean               # Clean build artifacts
make clean-all           # Clean everything including database
```