package services_test

import "monolith-chi-htmx-sqlite/src/models"

// TodoRepositoryInterface defines the interface for todo repository
type TodoRepositoryInterface interface {
	GetAllTodos() ([]models.Todo, error)
	GetTodosByCategory(categoryID int) ([]models.Todo, error)
	CreateTodo(title string, categoryIDs []int) error
	UpdateTodoStatus(id int, status string) error
	DeleteTodo(id int) error
	GetCategoriesForTodo(todoID int) ([]models.Category, error)
	GetAllTodosPaginated(page, pageSize int) (*models.PaginatedTodos, error)
	GetTodosByCategoryPaginated(categoryID, page, pageSize int) (*models.PaginatedTodos, error)
}

// CategoryRepositoryInterface defines the interface for category repository
type CategoryRepositoryInterface interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(title string) error
	DeleteCategory(id int) error
} 