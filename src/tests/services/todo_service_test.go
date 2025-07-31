package services_test

import (
	"testing"

	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/services"
	"monolith-chi-htmx-sqlite/src/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock repositories for testing
type MockTodoRepository struct {
	todos []models.Todo
	err   error
}

func (m *MockTodoRepository) GetAllTodos() ([]models.Todo, error) {
	return m.todos, m.err
}

func (m *MockTodoRepository) GetTodosByCategory(categoryID int) ([]models.Todo, error) {
	return m.todos, m.err
}

func (m *MockTodoRepository) CreateTodo(title string, categoryIDs []int) error {
	return m.err
}

func (m *MockTodoRepository) UpdateTodoStatus(id int, status string) error {
	return m.err
}

func (m *MockTodoRepository) DeleteTodo(id int) error {
	return m.err
}

func (m *MockTodoRepository) GetCategoriesForTodo(todoID int) ([]models.Category, error) {
	return []models.Category{}, m.err
}

func (m *MockTodoRepository) GetAllTodosPaginated(page, pageSize int) (*models.PaginatedTodos, error) {
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

func (m *MockTodoRepository) GetTodosByCategoryPaginated(categoryID, page, pageSize int) (*models.PaginatedTodos, error) {
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

type MockCategoryRepository struct {
	categories []models.Category
	err        error
}

func (m *MockCategoryRepository) GetAllCategories() ([]models.Category, error) {
	return m.categories, m.err
}

func (m *MockCategoryRepository) CreateCategory(title string) error {
	return m.err
}

func (m *MockCategoryRepository) DeleteCategory(id int) error {
	return m.err
}

func TestNewTodoService(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)
	assert.NotNil(t, service)
}

func TestTodoService_GetTodos_Success(t *testing.T) {
	mockTime := helpers.MockTime()
	todos := []models.Todo{
		{
			ID:        1,
			Title:     "Test Todo",
			Status:    "pending",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}

	todoRepo := &MockTodoRepository{todos: todos}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	result, err := service.GetTodos("", 1, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Todos, 1)
	assert.Equal(t, "Test Todo", result.Todos[0].Title)
}

func TestTodoService_GetTodos_WithCategoryFilter(t *testing.T) {
	mockTime := helpers.MockTime()
	todos := []models.Todo{
		{
			ID:        1,
			Title:     "Test Todo",
			Status:    "pending",
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		},
	}

	todoRepo := &MockTodoRepository{todos: todos}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	result, err := service.GetTodos("1", 1, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Todos, 1)
}

func TestTodoService_GetTodos_InvalidCategoryFilter(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	result, err := service.GetTodos("invalid", 1, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestTodoService_CreateTodo_Success(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.CreateTodo("Test Todo", []int{1})

	require.NoError(t, err)
}

func TestTodoService_CreateTodo_EmptyTitle(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.CreateTodo("", []int{1})

	assert.Error(t, err)
}

func TestTodoService_UpdateTodoStatus_Success(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.UpdateTodoStatus(1, "completed")

	require.NoError(t, err)
}

func TestTodoService_UpdateTodoStatus_InvalidStatus(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.UpdateTodoStatus(1, "invalid_status")

	assert.Error(t, err)
}

func TestTodoService_DeleteTodo_Success(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.DeleteTodo(1)

	require.NoError(t, err)
}

func TestTodoService_GetCategories_Success(t *testing.T) {
	mockTime := helpers.MockTime()
	categories := []models.Category{
		{
			ID:        1,
			Title:     "Test Category",
			CreatedAt: mockTime,
			TodoCount: 5,
		},
	}

	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{categories: categories}

	service := services.NewTodoService(todoRepo, categoryRepo)

	result, err := service.GetCategories()

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Test Category", result[0].Title)
}

func TestTodoService_CreateCategory_Success(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.CreateCategory("Test Category")

	require.NoError(t, err)
}

func TestTodoService_CreateCategory_EmptyTitle(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.CreateCategory("")

	assert.Error(t, err)
}

func TestTodoService_DeleteCategory_Success(t *testing.T) {
	todoRepo := &MockTodoRepository{}
	categoryRepo := &MockCategoryRepository{}

	service := services.NewTodoService(todoRepo, categoryRepo)

	err := service.DeleteCategory(1)

	require.NoError(t, err)
} 