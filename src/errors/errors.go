package errors

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors
type AppError struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new application error
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors
func NewBadRequestError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

func NewNotFoundError(message string, err error) *AppError {
	return NewAppError(http.StatusNotFound, message, err)
}

func NewInternalServerError(message string, err error) *AppError {
	return NewAppError(http.StatusInternalServerError, message, err)
}

func NewValidationError(message string, err error) *AppError {
	return NewAppError(http.StatusBadRequest, message, err)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetAppError extracts AppError from error
func GetAppError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return nil
}
