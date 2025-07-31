// Package handlers provides HTTP handlers for the todo application
package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"monolith-chi-htmx-sqlite/src/middleware"
	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/services"
	"monolith-chi-htmx-sqlite/src/validation"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type TodoHandler struct {
	todoService services.TodoServiceInterface
	templates   *template.Template
	config      *Config
	logger      *logrus.Logger
}

type Config struct {
	DefaultPageSize int
	MaxPageSize     int
}

func NewTodoHandler(db *sql.DB, templates *template.Template, config *Config, logger *logrus.Logger) *TodoHandler {
	todoRepo := models.NewTodoRepository(db)
	categoryRepo := models.NewCategoryRepository(db)
	todoService := services.NewTodoServiceWithRepositories(todoRepo, categoryRepo)

	return &TodoHandler{
		todoService: todoService,
		templates:   templates,
		config:      config,
		logger:      logger,
	}
}

// NewTodoHandlerWithService creates a new todo handler with an existing service
func NewTodoHandlerWithService(templates *template.Template, config *Config, todoService services.TodoServiceInterface, logger *logrus.Logger) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		templates:   templates,
		config:      config,
		logger:      logger,
	}
}

// IndexHandler godoc
// @Summary Get todo list page
// @Description Display the main todo list page with pagination and category filtering
// @Tags todos
// @Accept html
// @Produce html
// @Param category query string false "Category ID to filter todos"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {string} string "HTML page with todo list"
// @Failure 400 {string} string "Bad request - invalid pagination parameters"
// @Failure 500 {string} string "Internal server error"
// @Router / [get]
func (h *TodoHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	categoryFilter := r.URL.Query().Get("category")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	// Validate pagination parameters
	page, pageSize, validationResult := validation.ValidatePagination(pageStr, pageSizeStr, h.config.MaxPageSize)
	if !validationResult.IsValid {
		middleware.HandleValidationError(h.logger, w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Get todos using service
	paginatedTodos, err := h.todoService.GetTodos(categoryFilter, page, pageSize)
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
		return
	}

	// Get categories
	categories, err := h.todoService.GetCategories()
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
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
		h.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error executing template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo item with title and optional category assignments
// @Tags todos
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param title formData string true "Todo title"
// @Param categories formData []int false "Category IDs to assign to the todo"
// @Success 303 {string} string "Redirect to todo list page"
// @Failure 400 {string} string "Bad request - invalid form data"
// @Failure 500 {string} string "Internal server error"
// @Router /todos [post]
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
		middleware.HandleAppError(h.logger, w, err)
		return
	}

	// Redirect to index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// UpdateTodoStatus godoc
// @Summary Update todo status
// @Description Update the status of a specific todo item (incomplete/complete)
// @Tags todos
// @Accept application/x-www-form-urlencoded
// @Produce text/plain
// @Param id path int true "Todo ID"
// @Param status formData string true "Todo status (incomplete or complete)"
// @Success 200 {string} string "Status updated successfully"
// @Failure 400 {string} string "Bad request - invalid todo ID or status"
// @Failure 500 {string} string "Internal server error"
// @Router /todos/{id}/status [post]
func (h *TodoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, validationResult := validation.ValidateID(todoIDStr)
	if !validationResult.IsValid {
		middleware.HandleValidationError(h.logger, w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	status := r.FormValue("status")

	// Update status using service
	err := h.todoService.UpdateTodoStatus(todoID, status)
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
		return
	}

	// Return success response for HTMX
	w.WriteHeader(http.StatusOK)
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a specific todo item by ID
// @Tags todos
// @Accept */*
// @Produce text/plain
// @Param id path int true "Todo ID"
// @Success 200 {string} string "Todo deleted successfully"
// @Failure 400 {string} string "Bad request - invalid todo ID"
// @Failure 500 {string} string "Internal server error"
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get todo ID from URL
	todoIDStr := chi.URLParam(r, "id")
	todoID, validationResult := validation.ValidateID(todoIDStr)
	if !validationResult.IsValid {
		middleware.HandleValidationError(h.logger, w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Delete todo using service
	err := h.todoService.DeleteTodo(todoID)
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
		return
	}

	// Return success response for HTMX
	w.WriteHeader(http.StatusOK)
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category for organizing todos
// @Tags categories
// @Accept application/x-www-form-urlencoded
// @Produce html
// @Param title formData string true "Category title"
// @Success 303 {string} string "Redirect to todo list page"
// @Failure 400 {string} string "Bad request - invalid form data"
// @Failure 500 {string} string "Internal server error"
// @Router /categories [post]
func (h *TodoHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))

	// Create category using service
	err := h.todoService.CreateCategory(title)
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
		return
	}

	// Redirect to index page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a specific category by ID
// @Tags categories
// @Accept */*
// @Produce text/plain
// @Param id path int true "Category ID"
// @Success 200 {string} string "Category deleted successfully"
// @Failure 400 {string} string "Bad request - invalid category ID"
// @Failure 500 {string} string "Internal server error"
// @Router /categories/{id} [delete]
func (h *TodoHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Get category ID from URL
	categoryIDStr := chi.URLParam(r, "id")
	categoryID, validationResult := validation.ValidateID(categoryIDStr)
	if !validationResult.IsValid {
		middleware.HandleValidationError(h.logger, w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Delete category using service
	err := h.todoService.DeleteCategory(categoryID)
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
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
		middleware.HandleValidationError(h.logger, w, &middleware.ValidationResult{
			IsValid: false,
			Errors:  validationResult.Errors,
		})
		return
	}

	// Get todos using service
	paginatedTodos, err := h.todoService.GetTodos(categoryFilter, page, pageSize)
	if err != nil {
		middleware.HandleAppError(h.logger, w, err)
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
		h.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error executing template")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
