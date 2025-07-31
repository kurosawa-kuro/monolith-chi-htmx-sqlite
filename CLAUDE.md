# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a monolithic ToDo management application with category functionality, built using server-side rendering.

**Tech Stack:**
- Backend: Go with [chi](https://github.com/go-chi/chi) router
- Frontend: HTML + [htmx](https://htmx.org/) + [Tailwind CSS](https://tailwindcss.com/)
- Database: SQLite
- Templates: Go standard `html/template`

## Development Commands

### Frontend Development (Tailwind CSS)
```bash
# Watch mode for CSS development
cd template-admin/html
npm run dev

# Build CSS for production
cd template-admin/html
npm run build
```

### Go Development
The project appears to be in early development with no Go files yet created. Common Go commands will be:
```bash
# Run the application
go run main.go

# Build for production
go build -o app main.go

# Run tests
go test ./...

# Install dependencies
go mod tidy
```

## Architecture

### Database Schema
The application uses a three-table structure:
- `todos`: Task entities (id, title, status)
- `categories`: Category entities (id, title)
- `todo_category`: Many-to-many relationship table (todo_id, category_id)

### Planned Directory Structure
```
/src
├── main.go                 # Application entry point
├── templates/              # HTML templates
│   ├── layout.html
│   ├── index.html
│   └── edit.html
├── static/                 # Static assets
│   └── tailwind.css
├── models/                 # Data models
│   └── models.go
├── handlers/               # HTTP handlers
│   └── todo.go
└── db/                     # Database schemas
    └── schema.sql
```

## Design System

**Important**: Must use the Tailwind template designs located in `/template-admin/html/` for UI consistency. Custom HTML/CSS outside of this template system is prohibited per project specifications.

Available template components in `template-admin/html/`:
- Forms, buttons, modals, alerts
- Dashboard layouts and components
- Authentication pages
- Various UI components (accordion, badges, breadcrumbs, etc.)

## Core Features

### Todo Management
- List todos with category associations
- CRUD operations (Create, Read, Update, Delete)
- Status updates (incomplete ↔ complete)

### Category Management
- Add/edit/delete categories
- Filter todos by category

### UX Enhancements
- htmx for partial page updates without full reloads
- Tailwind CSS for responsive styling

## Key Considerations

- Server-side rendering with Go templates
- SQLite for simple, file-based data persistence
- htmx enables dynamic interactions without JavaScript frameworks
- Monolithic architecture keeps everything in a single deployable unit