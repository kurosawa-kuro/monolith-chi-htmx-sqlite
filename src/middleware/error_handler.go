package middleware

import (
	"net/http"
	"strings"

	"monolith-chi-htmx-sqlite/src/errors"

	"github.com/sirupsen/logrus"
)

// ErrorHandler middleware handles application errors
func ErrorHandler(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.WithFields(logrus.Fields{
						"error":  err,
						"path":   r.URL.Path,
						"method": r.Method,
					}).Error("Panic recovered")
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// HandleAppError handles application-specific errors
func HandleAppError(logger *logrus.Logger, w http.ResponseWriter, err error) {
	if appErr := errors.GetAppError(err); appErr != nil {
		logger.WithFields(logrus.Fields{
			"error_code":    appErr.Code,
			"error_message": appErr.Message,
		}).Error("Application error")
		http.Error(w, appErr.Message, appErr.Code)
		return
	}

	// Handle generic errors
	logger.WithFields(logrus.Fields{
		"error": err,
	}).Error("Unexpected error")
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

// HandleValidationError handles validation errors
func HandleValidationError(logger *logrus.Logger, w http.ResponseWriter, validationResult *ValidationResult) {
	if !validationResult.IsValid {
		errorMessage := "Validation failed: " + strings.Join(validationResult.Errors, "; ")
		logger.WithFields(logrus.Fields{
			"validation_errors": validationResult.Errors,
		}).Warn("Validation failed")
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}
}

// ValidationResult represents validation result
type ValidationResult struct {
	IsValid bool
	Errors  []string
}
