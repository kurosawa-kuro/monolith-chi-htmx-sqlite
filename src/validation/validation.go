package validation

import (
	"strconv"
	"strings"
)

// ValidationResult represents validation result
type ValidationResult struct {
	IsValid bool
	Errors  []string
}

// NewValidationResult creates a new validation result
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		IsValid: true,
		Errors:  make([]string, 0),
	}
}

// AddError adds an error to validation result
func (v *ValidationResult) AddError(message string) {
	v.IsValid = false
	v.Errors = append(v.Errors, message)
}

// ValidateTodo validates todo data
func ValidateTodo(title string, categoryIDs []int) *ValidationResult {
	result := NewValidationResult()

	// Validate title
	if strings.TrimSpace(title) == "" {
		result.AddError("Title is required")
	}

	if len(title) > 255 {
		result.AddError("Title must be less than 255 characters")
	}

	// Validate category IDs
	for _, id := range categoryIDs {
		if id <= 0 {
			result.AddError("Invalid category ID")
			break
		}
	}

	return result
}

// ValidateCategory validates category data
func ValidateCategory(title string) *ValidationResult {
	result := NewValidationResult()

	// Validate title
	if strings.TrimSpace(title) == "" {
		result.AddError("Category title is required")
	}

	if len(title) > 100 {
		result.AddError("Category title must be less than 100 characters")
	}

	return result
}

// ValidatePagination validates pagination parameters
func ValidatePagination(pageStr, pageSizeStr string, maxPageSize int) (int, int, *ValidationResult) {
	result := NewValidationResult()

	// Parse page
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			result.AddError("Invalid page number")
		}
	}

	// Parse page size
	pageSize := 10
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= maxPageSize {
			pageSize = ps
		} else {
			result.AddError("Invalid page size")
		}
	}

	return page, pageSize, result
}

// ValidateID validates ID parameter
func ValidateID(idStr string) (int, *ValidationResult) {
	result := NewValidationResult()

	if idStr == "" {
		result.AddError("ID is required")
		return 0, result
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		result.AddError("Invalid ID")
		return 0, result
	}

	return id, result
}

// ValidateStatus validates todo status
func ValidateStatus(status string) *ValidationResult {
	result := NewValidationResult()

	validStatuses := []string{"pending", "in_progress", "completed"}
	isValid := false

	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		result.AddError("Invalid status. Must be one of: pending, in_progress, completed")
	}

	return result
}
