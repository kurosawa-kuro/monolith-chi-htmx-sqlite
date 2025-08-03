package handlers_test

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"monolith-chi-htmx-sqlite/src/gen"
	"monolith-chi-htmx-sqlite/src/handlers"
	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/tests/helpers"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// MockTodoService for testing
type MockTodoService struct {
	todos      []models.Todo
	categories []models.Category
	err        error
}

func (m *MockTodoService) GetTodos(categoryFilter string, page, pageSize int) (*models.PaginatedTodos, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &models.PaginatedTodos{
		Todos: m.todos,
		Pagination: models.Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      len(m.todos),
			TotalPages: 1,
		},
	}, nil
}

func (m *MockTodoService) CreateTodo(title string, categoryIDs []int) error {
	return m.err
}

func (m *MockTodoService) UpdateTodoStatus(id int, status string) error {
	return m.err
}

func (m *MockTodoService) DeleteTodo(id int) error {
	return m.err
}

func (m *MockTodoService) GetCategories() ([]models.Category, error) {
	return m.categories, m.err
}

func (m *MockTodoService) CreateCategory(title string) error {
	return m.err
}

func (m *MockTodoService) DeleteCategory(id int) error {
	return m.err
}

func (m *MockTodoService) GetTodoByID(id int) (*models.Todo, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Find todo by ID
	for _, todo := range m.todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, nil
}

func (m *MockTodoService) UpdateTodo(id int, title string, categoryIDs []int) error {
	return m.err
}

func setupTestOAPIHandler() (gen.ServerInterface, *MockTodoService) {
	mockService := &MockTodoService{}

	// Create test templates
	tmpl := template.Must(template.New("index.html").Parse(`
		{{range .Todos}}
			<div class="todo">{{.Title}}</div>
		{{end}}
		{{range .Categories}}
			<div class="category">{{.Title}}</div>
		{{end}}
	`))

	// Add todo_detail.html template
	tmpl.New("todo_detail.html").Parse(`
		<div class="todo-detail">
			<h1>{{.Todo.Title}}</h1>
			<p>Status: {{.Todo.Status}}</p>
			{{range .Categories}}
				<div class="category">{{.Title}}</div>
			{{end}}
		</div>
	`)

	config := &handlers.Config{
		DefaultPageSize: 10,
		MaxPageSize:     100,
	}

	handler := handlers.NewOAPIHandler(tmpl, config, mockService, logrus.New())
	return handler, mockService
}

func TestOAPIHandler_GetTodoListPage_Success(t *testing.T) {
	handler, mockService := setupTestOAPIHandler()
	mockTime := helpers.MockTime()

	mockService.todos = []models.Todo{
		{
			ID:        1,
			Title:     "Test Todo",
			Status:    "incomplete",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}

	mockService.categories = []models.Category{
		{
			ID:        1,
			Title:     "Test Category",
			CreatedAt: mockTime,
			TodoCount: 1,
		},
	}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	params := gen.GetTodoListPageParams{}
	handler.GetTodoListPage(w, req, params)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Test Todo")
	assert.Contains(t, body, "Test Category")
}

func TestOAPIHandler_GetTodoListPage_WithCategoryFilter(t *testing.T) {
	handler, mockService := setupTestOAPIHandler()
	mockTime := helpers.MockTime()

	mockService.todos = []models.Todo{
		{
			ID:        1,
			Title:     "Filtered Todo",
			Status:    "incomplete",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}

	req := httptest.NewRequest("GET", "/?category=1", nil)
	w := httptest.NewRecorder()

	categoryFilter := "1"
	params := gen.GetTodoListPageParams{
		Category: &categoryFilter,
	}
	handler.GetTodoListPage(w, req, params)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Filtered Todo")
}

func TestOAPIHandler_GetTodoListPage_WithPagination(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	req := httptest.NewRequest("GET", "/?page=2&page_size=5", nil)
	w := httptest.NewRecorder()

	page := 2
	pageSize := 5
	params := gen.GetTodoListPageParams{
		Page:     &page,
		PageSize: &pageSize,
	}
	handler.GetTodoListPage(w, req, params)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_CreateTodo_Success(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	data := url.Values{}
	data.Set("title", "New Todo")
	data.Set("categories", "1")
	data.Set("categories", "2")

	req := httptest.NewRequest("POST", "/todos", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
}

func TestOAPIHandler_CreateTodo_EmptyTitle(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	data := url.Values{}
	data.Set("title", "")
	data.Set("categories", "1")

	req := httptest.NewRequest("POST", "/todos", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.CreateTodo(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
}

func TestOAPIHandler_UpdateTodoStatus_Success(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	data := url.Values{}
	data.Set("status", "complete")

	req := httptest.NewRequest("POST", "/todos/1/status", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateTodoStatus(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_DeleteTodo_Success(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	req := httptest.NewRequest("DELETE", "/todos/1", nil)
	w := httptest.NewRecorder()

	handler.DeleteTodo(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_CreateCategory_Success(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	data := url.Values{}
	data.Set("title", "New Category")

	req := httptest.NewRequest("POST", "/categories", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
}

func TestOAPIHandler_CreateCategory_EmptyTitle(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	data := url.Values{}
	data.Set("title", "")

	req := httptest.NewRequest("POST", "/categories", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
}

func TestOAPIHandler_DeleteCategory_Success(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	req := httptest.NewRequest("DELETE", "/categories/1", nil)
	w := httptest.NewRecorder()

	handler.DeleteCategory(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_GetTodoDetail_Success(t *testing.T) {
	handler, mockService := setupTestOAPIHandler()
	mockTime := helpers.MockTime()

	// Setup mock data
	mockService.todos = []models.Todo{
		{
			ID:        1,
			Title:     "Test Todo Detail",
			Status:    "incomplete",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
			Categories: []models.Category{
				{
					ID:        1,
					Title:     "Work",
					CreatedAt: mockTime,
				},
				{
					ID:        2,
					Title:     "Personal",
					CreatedAt: mockTime,
				},
			},
		},
	}

	mockService.categories = []models.Category{
		{
			ID:        1,
			Title:     "Work",
			CreatedAt: mockTime,
		},
		{
			ID:        2,
			Title:     "Personal",
			CreatedAt: mockTime,
		},
		{
			ID:        3,
			Title:     "Shopping",
			CreatedAt: mockTime,
		},
	}

	req := httptest.NewRequest("GET", "/todos/1", nil)
	w := httptest.NewRecorder()

	handler.GetTodoDetail(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Test Todo Detail")
	assert.Contains(t, body, "incomplete")
	assert.Contains(t, body, "Work")
	assert.Contains(t, body, "Personal")
	assert.Contains(t, body, "Shopping")
}

func TestOAPIHandler_GetTodoDetail_NotFound(t *testing.T) {
	handler, mockService := setupTestOAPIHandler()

	// Setup empty todos to simulate not found
	mockService.todos = []models.Todo{}

	req := httptest.NewRequest("GET", "/todos/999", nil)
	w := httptest.NewRecorder()

	handler.GetTodoDetail(w, req, 999)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestOAPIHandler_UpdateTodo_Success(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	// Setup form data
	formData := url.Values{}
	formData.Set("title", "Updated Todo Title")
	formData.Set("categories", "1,2")

	req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_UpdateTodo_EmptyTitle(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	// Setup form data with empty title
	formData := url.Values{}
	formData.Set("title", "")
	formData.Set("categories", "1,2")

	req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_UpdateTodo_NoCategories(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	// Setup form data without categories
	formData := url.Values{}
	formData.Set("title", "Updated Todo Title")
	formData.Set("categories", "")

	req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOAPIHandler_UpdateTodo_InvalidCategoryIDs(t *testing.T) {
	handler, _ := setupTestOAPIHandler()

	// Setup form data with invalid category IDs
	formData := url.Values{}
	formData.Set("title", "Updated Todo Title")
	formData.Set("categories", "invalid,abc,123")

	req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.UpdateTodo(w, req, 1)

	assert.Equal(t, http.StatusOK, w.Code)
}
