package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"monolith-chi-htmx-sqlite/src/config"
	"monolith-chi-htmx-sqlite/src/handlers"
	"monolith-chi-htmx-sqlite/src/middleware"
	"monolith-chi-htmx-sqlite/src/templates"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := initDB(cfg.Database.Path)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Load templates
	templateLoader := templates.NewLoader()
	if err := templateLoader.Load(); err != nil {
		log.Fatal("Failed to load templates:", err)
	}
	templates := templateLoader.GetTemplates()

	// Initialize handlers
	handlerConfig := &handlers.Config{
		DefaultPageSize: cfg.App.DefaultPageSize,
		MaxPageSize:     cfg.App.MaxPageSize,
	}
	todoHandler := handlers.NewTodoHandler(db, templates, handlerConfig)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.ErrorHandler)
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

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
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
