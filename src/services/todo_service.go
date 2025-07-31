package services

import (
	"monolith-chi-htmx-sqlite/src/errors"
	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/validation"
)

// TodoService handles business logic for todos
type TodoService struct {
	todoRepo     *models.TodoRepository
	categoryRepo *models.CategoryRepository
}

// NewTodoService creates a new todo service
func NewTodoService(todoRepo *models.TodoRepository, categoryRepo *models.CategoryRepository) *TodoService {
	return &TodoService{
		todoRepo:     todoRepo,
		categoryRepo: categoryRepo,
	}
}

// GetTodos retrieves todos with pagination and filtering
func (s *TodoService) GetTodos(categoryFilter string, page, pageSize int) (*models.PaginatedTodos, error) {
	var paginatedTodos *models.PaginatedTodos
	var err error

	if categoryFilter != "" {
		categoryID, validationResult := validation.ValidateID(categoryFilter)
		if !validationResult.IsValid {
			return nil, errors.NewValidationError("Invalid category filter", nil)
		}
		paginatedTodos, err = s.todoRepo.GetTodosByCategoryPaginated(categoryID, page, pageSize)
	} else {
		paginatedTodos, err = s.todoRepo.GetAllTodosPaginated(page, pageSize)
	}

	if err != nil {
		return nil, errors.NewInternalServerError("Failed to fetch todos", err)
	}

	return paginatedTodos, nil
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(title string, categoryIDs []int) error {
	// Validate input
	validationResult := validation.ValidateTodo(title, categoryIDs)
	if !validationResult.IsValid {
		return errors.NewValidationError("Invalid todo data", nil)
	}

	// Create todo
	err := s.todoRepo.CreateTodo(title, categoryIDs)
	if err != nil {
		return errors.NewInternalServerError("Failed to create todo", err)
	}

	return nil
}

// UpdateTodoStatus updates todo status
func (s *TodoService) UpdateTodoStatus(id int, status string) error {
	// Validate status
	validationResult := validation.ValidateStatus(status)
	if !validationResult.IsValid {
		return errors.NewValidationError("Invalid status", nil)
	}

	// Update status
	err := s.todoRepo.UpdateTodoStatus(id, status)
	if err != nil {
		return errors.NewInternalServerError("Failed to update todo status", err)
	}

	return nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(id int) error {
	err := s.todoRepo.DeleteTodo(id)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete todo", err)
	}

	return nil
}

// GetCategories retrieves all categories
func (s *TodoService) GetCategories() ([]models.Category, error) {
	categories, err := s.categoryRepo.GetAllCategories()
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to fetch categories", err)
	}

	return categories, nil
}

// CreateCategory creates a new category
func (s *TodoService) CreateCategory(title string) error {
	// Validate input
	validationResult := validation.ValidateCategory(title)
	if !validationResult.IsValid {
		return errors.NewValidationError("Invalid category data", nil)
	}

	// Create category
	err := s.categoryRepo.CreateCategory(title)
	if err != nil {
		return errors.NewInternalServerError("Failed to create category", err)
	}

	return nil
}

// DeleteCategory deletes a category
func (s *TodoService) DeleteCategory(id int) error {
	err := s.categoryRepo.DeleteCategory(id)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete category", err)
	}

	return nil
}
