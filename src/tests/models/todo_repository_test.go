package models_test

import (
	"database/sql"
	"testing"

	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/tests/helpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTodoRepository(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)
	assert.NotNil(t, repo)
	// Note: db field is private, so we can't directly test it
}

func TestTodoRepository_GetAllTodos(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)
	mockTime := helpers.MockTime()

	// Mock query expectations
	rows := sqlmock.NewRows([]string{"id", "title", "status", "created_at", "updated_at"}).
		AddRow(1, "Test Todo", "pending", mockTime, mockTime)

	testDB.Mock.ExpectQuery("SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at").
		WillReturnRows(rows)

	// Mock categories query
	categoryRows := sqlmock.NewRows([]string{"id", "title", "created_at"})
	testDB.Mock.ExpectQuery("SELECT c.id, c.title, c.created_at").
		WithArgs(1).
		WillReturnRows(categoryRows)

	todos, err := repo.GetAllTodos()

	require.NoError(t, err)
	assert.Len(t, todos, 1)
	assert.Equal(t, 1, todos[0].ID)
	assert.Equal(t, "Test Todo", todos[0].Title)
	assert.Equal(t, "pending", todos[0].Status)

	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestTodoRepository_GetAllTodos_Error(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)

	testDB.Mock.ExpectQuery("SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at").
		WillReturnError(sql.ErrConnDone)

	todos, err := repo.GetAllTodos()

	assert.Error(t, err)
	assert.Nil(t, todos)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestTodoRepository_CreateTodo(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)

	// Mock transaction expectations
	testDB.Mock.ExpectBegin()
	testDB.Mock.ExpectExec("INSERT INTO todos \\(title\\) VALUES \\(\\?\\)").
		WithArgs("Test Todo").
		WillReturnResult(sqlmock.NewResult(1, 1))

	testDB.Mock.ExpectExec("INSERT INTO todo_category \\(todo_id, category_id\\) VALUES \\(\\?, \\?\\)").
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	testDB.Mock.ExpectCommit()

	err := repo.CreateTodo("Test Todo", []int{1})

	require.NoError(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestTodoRepository_CreateTodo_TransactionError(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)

	testDB.Mock.ExpectBegin()
	testDB.Mock.ExpectExec("INSERT INTO todos").
		WillReturnError(sql.ErrConnDone)
	testDB.Mock.ExpectRollback()

	err := repo.CreateTodo("Test Todo", []int{1})

	assert.Error(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestTodoRepository_UpdateTodoStatus(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)

	testDB.Mock.ExpectExec("UPDATE todos SET status = \\?, updated_at = CURRENT_TIMESTAMP WHERE id = \\?").
		WithArgs("completed", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateTodoStatus(1, "completed")

	require.NoError(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestTodoRepository_DeleteTodo(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)

	testDB.Mock.ExpectExec("DELETE FROM todos WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteTodo(1)

	require.NoError(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestTodoRepository_GetAllTodosPaginated(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewTodoRepository(testDB.DB)
	mockTime := helpers.MockTime()

	// Mock count query
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	testDB.Mock.ExpectQuery("SELECT COUNT\\(DISTINCT t.id\\)").
		WillReturnRows(countRows)

	// Mock data query
	dataRows := sqlmock.NewRows([]string{"id", "title", "status", "created_at", "updated_at"}).
		AddRow(1, "Test Todo", "pending", mockTime, mockTime)

	testDB.Mock.ExpectQuery("SELECT DISTINCT t.id, t.title, t.status, t.created_at, t.updated_at").
		WillReturnRows(dataRows)

	// Mock categories query
	categoryRows := sqlmock.NewRows([]string{"id", "title", "created_at"})
	testDB.Mock.ExpectQuery("SELECT c.id, c.title, c.created_at").
		WithArgs(1).
		WillReturnRows(categoryRows)

	result, err := repo.GetAllTodosPaginated(1, 10)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Todos, 1)
	assert.Equal(t, 1, result.Pagination.Total)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 10, result.Pagination.PageSize)

	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}
