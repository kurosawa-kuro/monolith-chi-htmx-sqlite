package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"monolith-chi-htmx-sqlite/src/middleware"
	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/services"
	"monolith-chi-htmx-sqlite/src/validation"

	"github.com/go-chi/chi/v5"
)

type TodoHandler struct {
	todoService *services.TodoService
	templates   *template.Template
	config      *Config
}

type Config struct {
	DefaultPageSize int
	MaxPageSize     int
}

func NewTodoHandler(db *sql.DB, templates *template.Template, config *Config) *TodoHandler {
	todoRepo := models.NewTodoRepository(db)
	categoryRepo := models.NewCategoryRepository(db)
	todoService := services.NewTodoService(todoRepo, categoryRepo)

	return &TodoHandler{
		todoService: todoService,
		templates:   templates,
		config:      config,
	}
}

func (h *TodoHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	categoryFilter := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	// Validate pagination parameters
	page, pageSize, validationResult := validation.ValidatePagination(pageStr, pageSizeStr, h.config.MaxPageSize)
	if !validationResult.IsValid {
		middleware.HandleValidationError(w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Get todos using service
	paginatedTodos, err := h.todoService.GetTodos(categoryFilter, page, pageSize)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	// Get categories
	categories, err := h.todoService.GetCategories()
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	data := struct {
		Todos      []models.Todo
		Categories []models.Category
		Filter     string
		Pagination models.Pagination
	}{
		Todos:      paginatedTodos.Todos,
		Categories: categories,
		Filter:     categoryFilter,
		Pagination: paginatedTodos.Pagination,
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	categoryIDsStr := r.Form["categories"]

	// Convert category IDs to integers
	var categoryIDs []int
	for _, idStr := range categoryIDsStr {
		if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
			categoryIDs = append(categoryIDs, id)
		}
	}

	// Create todo using service
	err := h.todoService.CreateTodo(title, categoryIDs)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	// Redirect to index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *TodoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, validationResult := validation.ValidateID(todoIDStr)
	if !validationResult.IsValid {
		middleware.HandleValidationError(w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	status := r.FormValue("status")

	// Update status using service
	err := h.todoService.UpdateTodoStatus(todoID, status)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	// Return success response for HTMX
	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, validationResult := validation.ValidateID(todoIDStr)
	if !validationResult.IsValid {
		middleware.HandleValidationError(w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Delete todo using service
	err := h.todoService.DeleteTodo(todoID)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	// Return success response for HTMX
	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))

	// Create category using service
	err := h.todoService.CreateCategory(title)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	// Redirect to index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *TodoHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Get category ID from URL
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, validationResult := validation.ValidateID(categoryIDStr)
	if !validationResult.IsValid {
		middleware.HandleValidationError(w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Delete category using service
	err := h.todoService.DeleteCategory(categoryID)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	// Return success response for HTMX
	w.WriteHeader(http.StatusOK)
}

func (h *TodoHandler) renderTodoList(w http.ResponseWriter, r *http.Request) {
	categoryFilter := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	// Validate pagination parameters
	page, pageSize, validationResult := validation.ValidatePagination(pageStr, pageSizeStr, h.config.MaxPageSize)
	if !validationResult.IsValid {
		middleware.HandleValidationError(w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Get todos using service
	paginatedTodos, err := h.todoService.GetTodos(categoryFilter, page, pageSize)
	if err != nil {
		middleware.HandleAppError(w, err)
		return
	}

	data := struct {
		Todos      []models.Todo
		Pagination models.Pagination
	}{
		Todos:      paginatedTodos.Todos,
		Pagination: paginatedTodos.Pagination,
	}

	if err := h.templates.ExecuteTemplate(w, "todo_list.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
