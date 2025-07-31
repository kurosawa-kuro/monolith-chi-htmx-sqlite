package integration_test

import (
	"testing"

	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/services"
	"monolith-chi-htmx-sqlite/src/tests/helpers"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoServiceIntegration(t *testing.T) {
	// Create temporary database
	db, dbPath := helpers.CreateTempDB(t)
	defer db.Close()
	defer helpers.CleanupTempDB(dbPath)

	// Initialize repositories
	todoRepo := models.NewTodoRepository(db)
	categoryRepo := models.NewCategoryRepository(db)

	// Initialize service
	todoService := services.NewTodoServiceWithRepositories(todoRepo, categoryRepo)

	// Test creating a category
	err := todoService.CreateCategory("Test Category")
	require.NoError(t, err)

	// Test getting categories
	categories, err := todoService.GetCategories()
	require.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, "Test Category", categories[0].Title)

	// Test creating a todo
	err = todoService.CreateTodo("Test Todo", []int{categories[0].ID})
	require.NoError(t, err)

	// Test getting todos
	todos, err := todoService.GetTodos("", 1, 10)
	require.NoError(t, err)
	assert.Len(t, todos.Todos, 1)
	assert.Equal(t, "Test Todo", todos.Todos[0].Title)
	assert.Equal(t, "incomplete", todos.Todos[0].Status)

	// Test updating todo status
	err = todoService.UpdateTodoStatus(todos.Todos[0].ID, "complete")
	require.NoError(t, err)

	// Verify status was updated
	updatedTodos, err := todoService.GetTodos("", 1, 10)
	require.NoError(t, err)
	assert.Equal(t, "complete", updatedTodos.Todos[0].Status)

	// Test filtering by category
	filteredTodos, err := todoService.GetTodos("1", 1, 10)
	require.NoError(t, err)
	assert.Len(t, filteredTodos.Todos, 1)

	// Test deleting todo
	err = todoService.DeleteTodo(todos.Todos[0].ID)
	require.NoError(t, err)

	// Verify todo was deleted
	deletedTodos, err := todoService.GetTodos("", 1, 10)
	require.NoError(t, err)
	assert.Len(t, deletedTodos.Todos, 0)

	// Test deleting category
	err = todoService.DeleteCategory(categories[0].ID)
	require.NoError(t, err)

	// Verify category was deleted
	deletedCategories, err := todoService.GetCategories()
	require.NoError(t, err)
	assert.Len(t, deletedCategories, 0)
}

func TestTodoRepositoryIntegration(t *testing.T) {
	// Create temporary database
	db, dbPath := helpers.CreateTempDB(t)
	defer db.Close()
	defer helpers.CleanupTempDB(dbPath)

	// Initialize repository
	todoRepo := models.NewTodoRepository(db)

	// Test creating a todo
	err := todoRepo.CreateTodo("Test Todo", []int{})
	require.NoError(t, err)

	// Test getting todos
	todos, err := todoRepo.GetAllTodos()
	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, "Test Todo", todos[0].Title)

	// Test updating todo status
	err = todoRepo.UpdateTodoStatus(todos[0].ID, "complete")
	require.NoError(t, err)

	// Verify status was updated
	updatedTodos, err := todoRepo.GetAllTodos()
	require.NoError(t, err)
	assert.Equal(t, "complete", updatedTodos[0].Status)

	// Test deleting todo
	err = todoRepo.DeleteTodo(todos[0].ID)
	require.NoError(t, err)

	// Verify todo was deleted
	deletedTodos, err := todoRepo.GetAllTodos()
	require.NoError(t, err)
	assert.Len(t, deletedTodos, 0)
}

func TestCategoryRepositoryIntegration(t *testing.T) {
	// Create temporary database
	db, dbPath := helpers.CreateTempDB(t)
	defer db.Close()
	defer helpers.CleanupTempDB(dbPath)

	// Initialize repository
	categoryRepo := models.NewCategoryRepository(db)

	// Test creating a category
	err := categoryRepo.CreateCategory("Test Category")
	require.NoError(t, err)

	// Test getting categories
	categories, err := categoryRepo.GetAllCategories()
	require.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, "Test Category", categories[0].Title)

	// Test deleting category
	err = categoryRepo.DeleteCategory(categories[0].ID)
	require.NoError(t, err)

	// Verify category was deleted
	deletedCategories, err := categoryRepo.GetAllCategories()
	require.NoError(t, err)
	assert.Len(t, deletedCategories, 0)
}
