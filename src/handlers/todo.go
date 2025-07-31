package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"monolith-chi-htmx-sqlite/src/models"

	"github.com/go-chi/chi/v5"
)

type TodoHandler struct {
	todoRepo     *models.TodoRepository
	categoryRepo *models.CategoryRepository
	templates    *template.Template
}

func NewTodoHandler(db *sql.DB, templates *template.Template) *TodoHandler {
	return &TodoHandler{
		todoRepo:     models.NewTodoRepository(db),
		categoryRepo: models.NewCategoryRepository(db),
		templates:    templates,
	}
}

func (h *TodoHandler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	categoryFilter := r.URL.Query().Get("category")
	
	var todos []models.Todo
	var err error
	
	if categoryFilter != "" {
		categoryID, err := strconv.Atoi(categoryFilter)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		todos, err = h.todoRepo.GetTodosByCategory(categoryID)
	} else {
		todos, err = h.todoRepo.GetAllTodos()
	}
	
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryRepo.GetAllCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Todos      []models.Todo
		Categories []models.Category
		Filter     string
	}{
		Todos:      todos,
		Categories: categories,
		Filter:     categoryFilter,
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
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	categoryIDsStr := r.Form["categories"]
	var categoryIDs []int
	for _, idStr := range categoryIDsStr {
		if id, err := strconv.Atoi(idStr); err == nil {
			categoryIDs = append(categoryIDs, id)
		}
	}

	if err := h.todoRepo.CreateTodo(title, categoryIDs); err != nil {
		log.Printf("Error creating todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTodoList(w, r)
}

func (h *TodoHandler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	status := r.FormValue("status")
	if status != "complete" && status != "incomplete" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	if err := h.todoRepo.UpdateTodoStatus(id, status); err != nil {
		log.Printf("Error updating todo status: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTodoList(w, r)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := h.todoRepo.DeleteTodo(id); err != nil {
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.renderTodoList(w, r)
}

func (h *TodoHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if err := h.categoryRepo.CreateCategory(title); err != nil {
		log.Printf("Error creating category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryRepo.GetAllCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}

	if err := h.templates.ExecuteTemplate(w, "categories-list", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *TodoHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	if err := h.categoryRepo.DeleteCategory(id); err != nil {
		log.Printf("Error deleting category: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	categories, err := h.categoryRepo.GetAllCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}

	if err := h.templates.ExecuteTemplate(w, "categories-list", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *TodoHandler) renderTodoList(w http.ResponseWriter, r *http.Request) {
	categoryFilter := r.URL.Query().Get("category")
	
	var todos []models.Todo
	var err error
	
	if categoryFilter != "" {
		categoryID, err := strconv.Atoi(categoryFilter)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
		todos, err = h.todoRepo.GetTodosByCategory(categoryID)
	} else {
		todos, err = h.todoRepo.GetAllTodos()
	}
	
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Todos []models.Todo
	}{
		Todos: todos,
	}

	if err := h.templates.ExecuteTemplate(w, "todo-list", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}