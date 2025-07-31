package models

import (
	"database/sql"
	"log"
	"time"
)

// Pagination represents pagination information
// @Description Pagination information for list responses
type Pagination struct {
	Page       int `json:"page" example:"1"`        // Current page number
	PageSize   int `json:"page_size" example:"10"`  // Number of items per page
	Total      int `json:"total" example:"25"`      // Total number of items
	TotalPages int `json:"total_pages" example:"3"` // Total number of pages
}

// PaginatedTodos represents paginated todo results
// @Description Paginated list of todos with pagination metadata
type PaginatedTodos struct {
	Todos      []Todo     `json:"todos"`      // List of todo items
	Pagination Pagination `json:"pagination"` // Pagination information
}

// Todo represents a todo item
// @Description A todo item with its properties and associated categories
type Todo struct {
	ID         int        `json:"id" example:"1"`                            // Unique identifier
	Title      string     `json:"title" example:"Buy groceries"`             // Todo title
	Status     string     `json:"status" example:"incomplete"`               // Todo status (incomplete/complete)
	CreatedAt  time.Time  `json:"created_at" example:"2024-01-01T12:00:00Z"` // Creation timestamp
	UpdatedAt  time.Time  `json:"updated_at" example:"2024-01-01T12:00:00Z"` // Last update timestamp
	Categories []Category `json:"categories"`                                // Associated categories
}

// Category represents a category for organizing todos
// @Description A category used to organize and group todo items
type Category struct {
	ID        int       `json:"id" example:"1"`                            // Unique identifier
	Title     string    `json:"title" example:"Work"`                      // Category title
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"` // Creation timestamp
	TodoCount int       `json:"todo_count" example:"5"`                    // Number of todos in this category
}

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) GetAllTodos() ([]Todo, error) {
	query := `
		SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at
		FROM todos t
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories, err := r.GetCategoriesForTodo(todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Categories = categories

		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) GetTodosByCategory(categoryID int) ([]Todo, error) {
	query := `
		SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at
		FROM todos t
		JOIN todo_category tc ON t.id = tc.todo_id
		WHERE tc.category_id = ?
		ORDER BY t.created_at DESC
	`
	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories, err := r.GetCategoriesForTodo(todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Categories = categories

		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepository) CreateTodo(title string, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO todos (title) VALUES (?)", title)
	if err != nil {
		return err
	}

	todoID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	for _, categoryID := range categoryIDs {
		_, err = tx.Exec("INSERT INTO todo_category (todo_id, category_id) VALUES (?, ?)", todoID, categoryID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *TodoRepository) UpdateTodoStatus(id int, status string) error {
	_, err := r.db.Exec("UPDATE todos SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", status, id)
	return err
}

func (r *TodoRepository) DeleteTodo(id int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func (r *TodoRepository) GetCategoriesForTodo(todoID int) ([]Category, error) {
	query := `
		SELECT c.id, c.title, c.created_at
		FROM categories c
		JOIN todo_category tc ON c.id = tc.category_id
		WHERE tc.todo_id = ?
	`
	rows, err := r.db.Query(query, todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Title, &category.CreatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	// Debug: Log categories for todo
	if len(categories) > 0 {
		log.Printf("Todo %d has %d categories: %v", todoID, len(categories), categories)
	}

	return categories, nil
}

func (r *TodoRepository) GetAllTodosPaginated(page, pageSize int) (*PaginatedTodos, error) {
	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(DISTINCT t.id) FROM todos t").Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (page - 1) * pageSize
	totalPages := (total + pageSize - 1) / pageSize

	// Get paginated todos
	query := `
		SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at
		FROM todos t
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories, err := r.GetCategoriesForTodo(todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Categories = categories

		todos = append(todos, todo)
	}

	return &PaginatedTodos{
		Todos: todos,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (r *TodoRepository) GetTodosByCategoryPaginated(categoryID, page, pageSize int) (*PaginatedTodos, error) {
	// Get total count for category
	var total int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT t.id)
		FROM todos t
		JOIN todo_category tc ON t.id = tc.todo_id
		WHERE tc.category_id = ?
	`, categoryID).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (page - 1) * pageSize
	totalPages := (total + pageSize - 1) / pageSize

	// Get paginated todos by category
	query := `
		SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at
		FROM todos t
		JOIN todo_category tc ON t.id = tc.todo_id
		WHERE tc.category_id = ?
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, categoryID, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		categories, err := r.GetCategoriesForTodo(todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Categories = categories

		todos = append(todos, todo)
	}

	return &PaginatedTodos{
		Todos: todos,
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories() ([]Category, error) {
	query := `
		SELECT c.id, c.title, c.created_at, COUNT(tc.todo_id) as todo_count
		FROM categories c
		LEFT JOIN todo_category tc ON c.id = tc.category_id
		GROUP BY c.id, c.title, c.created_at
		ORDER BY c.title
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.ID, &category.Title, &category.CreatedAt, &category.TodoCount)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryRepository) CreateCategory(title string) error {
	_, err := r.db.Exec("INSERT INTO categories (title) VALUES (?)", title)
	return err
}

func (r *CategoryRepository) DeleteCategory(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = ?", id)
	return err
}

func (r *TodoRepository) GetTodoByID(id int) (*Todo, error) {
	query := `
		SELECT t.id, t.title, t.status, t.created_at, t.updated_at
		FROM todos t
		WHERE t.id = ?
	`
	var todo Todo
	err := r.db.QueryRow(query, id).Scan(&todo.ID, &todo.Title, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	categories, err := r.GetCategoriesForTodo(todo.ID)
	if err != nil {
		return nil, err
	}
	todo.Categories = categories

	return &todo, nil
}

func (r *TodoRepository) UpdateTodo(id int, title string, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update todo title
	_, err = tx.Exec("UPDATE todos SET title = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", title, id)
	if err != nil {
		return err
	}

	// Delete existing category associations
	_, err = tx.Exec("DELETE FROM todo_category WHERE todo_id = ?", id)
	if err != nil {
		return err
	}

	// Add new category associations
	for _, categoryID := range categoryIDs {
		_, err = tx.Exec("INSERT INTO todo_category (todo_id, category_id) VALUES (?, ?)", id, categoryID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
