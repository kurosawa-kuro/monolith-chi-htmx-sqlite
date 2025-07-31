package models

import (
	"database/sql"
	"time"
)

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedTodos represents paginated todo results
type PaginatedTodos struct {
	Todos      []Todo     `json:"todos"`
	Pagination Pagination `json:"pagination"`
}

type Todo struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Categories []Category `json:"categories"`
}

type Category struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	TodoCount int       `json:"todo_count"`
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
