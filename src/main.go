package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"monolith-chi-htmx-sqlite/src/handlers"
	"monolith-chi-htmx-sqlite/src/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Load templates
	templates, err := loadTemplates()
	if err != nil {
		log.Fatal("Failed to load templates:", err)
	}

	// Initialize handlers
	todoHandler := handlers.NewTodoHandler(db, templates)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.FilteredLogger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Compress(5))

	// Static files
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "src/static"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	// Routes
	r.Get("/", todoHandler.IndexHandler)
	r.Post("/todos", todoHandler.CreateTodo)
	r.Post("/todos/{id}/status", todoHandler.UpdateTodoStatus)
	r.Delete("/todos/{id}", todoHandler.DeleteTodo)
	r.Post("/categories", todoHandler.CreateCategory)
	r.Delete("/categories/{id}", todoHandler.DeleteCategory)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./src/db/todo.db")
	if err != nil {
		return nil, err
	}

	// Read and execute schema
	schema, err := ioutil.ReadFile("src/db/schema.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadTemplates() (*template.Template, error) {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"min": func(a, b int) int {
			return int(math.Min(float64(a), float64(b)))
		},
		"max": func(a, b int) int {
			return int(math.Max(float64(a), float64(b)))
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
		"le": func(a, b int) bool {
			return a <= b
		},
		"gt": func(a, b int) bool {
			return a > b
		},
		"lt": func(a, b int) bool {
			return a < b
		},
		"sequence": func(n int) []int {
			result := make([]int, n)
			for i := range result {
				result[i] = i + 1
			}
			return result
		},
	}

	tmpl := template.New("").Funcs(funcMap)

	// Load all template files
	err := filepath.Walk("src/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = tmpl.New(filepath.Base(path)).Parse(string(content))
			return err
		}
		return nil
	})

	return tmpl, err
}
