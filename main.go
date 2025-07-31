package main

import (
	"database/sql"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"monolith-chi-htmx-sqlite/src/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

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
	db, err := sql.Open("sqlite3", "./todo.db")
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
	tmpl := template.New("")
	
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