package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"monolith-chi-htmx-sqlite/src/config"
	"monolith-chi-htmx-sqlite/src/handlers"
	"monolith-chi-htmx-sqlite/src/middleware"
	"monolith-chi-htmx-sqlite/src/models"
	"monolith-chi-htmx-sqlite/src/services"
	"monolith-chi-htmx-sqlite/src/templates"

	_ "monolith-chi-htmx-sqlite/src/docs"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"monolith-chi-htmx-sqlite/src/gen"
)

// @title           Todo API
// @version         1.0
// @description     This is a sample server for a todo app.
// @host      localhost:8080
// @BasePath  /

func main() {
	// Load configuration
	cfg := config.Load()
	loggerConfig := config.LoadLoggerConfig()

	// Setup logger
	logger := config.SetupLogger(loggerConfig)
	logger.WithFields(logrus.Fields{
		"environment": cfg.App.Environment,
		"log_level":   loggerConfig.Level,
		"log_format":  loggerConfig.Format,
	}).Info("Application starting")

	// Initialize database
	db, err := initDB(cfg.Database.Path)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":   err,
			"db_path": cfg.Database.Path,
		}).Fatal("Failed to initialize database")
	}
	defer db.Close()

	// Load templates
	templateLoader := templates.NewLoader()
	if err := templateLoader.Load(); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to load templates")
	}
	templates := templateLoader.GetTemplates()

	// Initialize repositories
	todoRepo := models.NewTodoRepository(db)
	categoryRepo := models.NewCategoryRepository(db)

	// Initialize services
	todoService := services.NewTodoServiceWithRepositories(todoRepo, categoryRepo)

	// Initialize handlers
	handlerConfig := &handlers.Config{
		DefaultPageSize: cfg.App.DefaultPageSize,
		MaxPageSize:     cfg.App.MaxPageSize,
	}
	oapiHandler := handlers.NewOAPIHandler(templates, handlerConfig, todoService, logger)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.FilteredLogger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Compress(5))

	// Static files
	workDir, _ := filepath.Abs(".")
	filesDir := http.Dir(filepath.Join(workDir, "src/static"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	// OAPI generated routes (includes GET / for index page)
	r.Mount("/", gen.HandlerFromMux(oapiHandler, r))
	
	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler())

	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	logger.WithFields(logrus.Fields{
		"host": cfg.Server.Host,
		"port": cfg.Server.Port,
	}).Info("Server starting")

	logger.Fatal(http.ListenAndServe(serverAddr, r))
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
