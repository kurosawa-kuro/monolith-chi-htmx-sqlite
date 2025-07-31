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

func TestNewCategoryRepository(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)
	assert.NotNil(t, repo)
}

func TestCategoryRepository_GetAllCategories(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)
	mockTime := helpers.MockTime()

	rows := sqlmock.NewRows([]string{"id", "title", "created_at", "todo_count"}).
		AddRow(1, "Test Category", mockTime, 5)

	testDB.Mock.ExpectQuery("SELECT c.id, c.title, c.created_at, COUNT\\(tc.todo_id\\) as todo_count").
		WillReturnRows(rows)

	categories, err := repo.GetAllCategories()

	require.NoError(t, err)
	assert.Len(t, categories, 1)
	assert.Equal(t, 1, categories[0].ID)
	assert.Equal(t, "Test Category", categories[0].Title)
	assert.Equal(t, 5, categories[0].TodoCount)

	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestCategoryRepository_GetAllCategories_Error(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)

	testDB.Mock.ExpectQuery("SELECT c.id, c.title, c.created_at, COUNT\\(tc.todo_id\\) as todo_count").
		WillReturnError(sql.ErrConnDone)

	categories, err := repo.GetAllCategories()

	assert.Error(t, err)
	assert.Nil(t, categories)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestCategoryRepository_CreateCategory(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)

	testDB.Mock.ExpectExec("INSERT INTO categories \\(title\\) VALUES \\(\\?\\)").
		WithArgs("Test Category").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.CreateCategory("Test Category")

	require.NoError(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestCategoryRepository_CreateCategory_Error(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)

	testDB.Mock.ExpectExec("INSERT INTO categories \\(title\\) VALUES \\(\\?\\)").
		WillReturnError(sql.ErrConnDone)

	err := repo.CreateCategory("Test Category")

	assert.Error(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestCategoryRepository_DeleteCategory(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)

	testDB.Mock.ExpectExec("DELETE FROM categories WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteCategory(1)

	require.NoError(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}

func TestCategoryRepository_DeleteCategory_Error(t *testing.T) {
	testDB := helpers.NewTestDB(t)
	defer testDB.Close()

	repo := models.NewCategoryRepository(testDB.DB)

	testDB.Mock.ExpectExec("DELETE FROM categories WHERE id = \\?").
		WillReturnError(sql.ErrConnDone)

	err := repo.DeleteCategory(1)

	assert.Error(t, err)
	require.NoError(t, testDB.Mock.ExpectationsWereMet())
}
