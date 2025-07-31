package helpers

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// TestDB provides a test database with mock
type TestDB struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
}

// NewTestDB creates a new test database with mock
func NewTestDB(t *testing.T) *TestDB {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return &TestDB{
		DB:   db,
		Mock: mock,
	}
}

// Close closes the test database
func (tdb *TestDB) Close() {
	tdb.DB.Close()
}

// CreateTempDB creates a temporary SQLite database for integration tests
func CreateTempDB(t *testing.T) (*sql.DB, string) {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "test-*.db")
	require.NoError(t, err)
	defer tmpfile.Close()

	dbPath := tmpfile.Name()

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	require.NoError(t, err)

	// Read and execute schema
	schema, err := os.ReadFile("src/db/schema.sql")
	if err != nil {
		// Try relative path from current working directory
		schema, err = os.ReadFile("../../db/schema.sql")
	}
	if err != nil {
		// Try absolute path from project root
		schema, err = os.ReadFile("/home/wsl/dev_my_study/monolith-chi-htmx-sqlite/src/db/schema.sql")
	}
	require.NoError(t, err)

	_, err = db.Exec(string(schema))
	require.NoError(t, err)

	return db, dbPath
}

// CleanupTempDB removes the temporary database file
func CleanupTempDB(dbPath string) {
	os.Remove(dbPath)
}

// MockTime provides a fixed time for testing
func MockTime() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

// SampleTodoData provides sample todo data for testing
func SampleTodoData() map[string]interface{} {
	return map[string]interface{}{
		"id":         1,
		"title":      "Test Todo",
		"status":     "pending",
		"created_at": MockTime(),
		"updated_at": MockTime(),
	}
}

// SampleCategoryData provides sample category data for testing
func SampleCategoryData() map[string]interface{} {
	return map[string]interface{}{
		"id":         1,
		"title":      "Test Category",
		"created_at": MockTime(),
		"todo_count": 0,
	}
}
