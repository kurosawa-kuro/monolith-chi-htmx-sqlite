package middleware

import (
	"log"
	"net/http"
	"strings"

	"monolith-chi-htmx-sqlite/src/errors"
)

// ErrorHandler middleware handles application errors
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// HandleAppError handles application-specific errors
func HandleAppError(w http.ResponseWriter, err error) {
	if appErr := errors.GetAppError(err); appErr != nil {
		log.Printf("Application error: %v", appErr)
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	// Handle generic errors
	log.Printf("Unexpected error: %v", err)
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// HandleValidationError handles validation errors
func HandleValidationError(w http.ResponseWriter, validationResult *ValidationResult) {
	if !validationResult.IsValid {
		errorMessage := "Validation failed: " + strings.Join(validationResult.Errors, "; ")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
}

// ValidationResult represents validation result
type ValidationResult struct {
	IsValid bool
	Errors  []string
}
