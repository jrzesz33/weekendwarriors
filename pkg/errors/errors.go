package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Validation errors
	ErrValidation        ErrorCode = "validation_error"
	ErrInvalidHandicap   ErrorCode = "invalid_handicap"
	ErrInvalidHoleNumber ErrorCode = "invalid_hole_number"
	ErrInvalidScoreValues ErrorCode = "invalid_score_values"
	ErrMissingRequiredField ErrorCode = "missing_required_field"

	// Business logic errors
	ErrGameNotStarted        ErrorCode = "game_not_started"
	ErrGameAlreadyCompleted  ErrorCode = "game_already_completed"
	ErrPlayerLimitExceeded   ErrorCode = "player_limit_exceeded"
	ErrDuplicatePlayerName   ErrorCode = "duplicate_player_name"
	ErrScoreAlreadyExists    ErrorCode = "score_already_exists"
	ErrFutureHole           ErrorCode = "future_hole"
	ErrInvalidGameState     ErrorCode = "invalid_game_state"

	// Resource errors
	ErrResourceNotFound ErrorCode = "resource_not_found"
	ErrGameNotFound     ErrorCode = "game_not_found"
	ErrPlayerNotFound   ErrorCode = "player_not_found"
	ErrScoreNotFound    ErrorCode = "score_not_found"

	// Side bet errors
	ErrSideBetNotEnabled  ErrorCode = "side_bet_not_enabled"
	ErrInsufficientHoles  ErrorCode = "insufficient_holes"
	ErrCardsAlreadyDealt  ErrorCode = "cards_already_dealt"
	ErrInvalidPuttCount   ErrorCode = "invalid_putt_count"
	ErrGameNotCompleted   ErrorCode = "game_not_completed"

	// Authentication errors
	ErrInvalidToken    ErrorCode = "invalid_token"
	ErrExpiredToken    ErrorCode = "expired_token"
	ErrInsufficientPermissions ErrorCode = "insufficient_permissions"

	// Rate limiting errors
	ErrRateLimitExceeded ErrorCode = "rate_limit_exceeded"

	// Server errors
	ErrInternalServer ErrorCode = "internal_server_error"
	ErrServiceUnavailable ErrorCode = "service_unavailable"
)

// APIError represents a structured API error
type APIError struct {
	Code      ErrorCode   `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ErrorResponse represents the full error response structure
type ErrorResponse struct {
	Error *APIError `json:"error"`
}

// New creates a new APIError
func New(code ErrorCode, message string) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NewWithDetails creates a new APIError with additional details
func NewWithDetails(code ErrorCode, message string, details interface{}) *APIError {
	return &APIError{
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// ValidationError creates a validation error with field details
func ValidationError(field, value, constraint string) *APIError {
	return NewWithDetails(ErrValidation, fmt.Sprintf("Invalid value for field '%s'", field), map[string]interface{}{
		"field":      field,
		"value":      value,
		"constraint": constraint,
	})
}

// ValidationErrorWithAllowedValues creates a validation error with allowed values
func ValidationErrorWithAllowedValues(field, value string, allowedValues []interface{}) *APIError {
	return NewWithDetails(ErrValidation, fmt.Sprintf("Invalid value for field '%s'", field), map[string]interface{}{
		"field":          field,
		"value":          value,
		"allowed_values": allowedValues,
	})
}

// ResourceNotFoundError creates a resource not found error
func ResourceNotFoundError(resourceType, resourceID string) *APIError {
	return NewWithDetails(ErrResourceNotFound, fmt.Sprintf("%s not found", resourceType), map[string]interface{}{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	})
}

// BusinessLogicError creates a business logic error with state information
func BusinessLogicError(code ErrorCode, message string, currentState, requiredState string) *APIError {
	return NewWithDetails(code, message, map[string]interface{}{
		"current_state":  currentState,
		"required_state": requiredState,
	})
}

// GetHTTPStatus returns the appropriate HTTP status code for an error
func (e *APIError) GetHTTPStatus() int {
	switch e.Code {
	case ErrValidation, ErrInvalidHandicap, ErrInvalidHoleNumber, ErrInvalidScoreValues, ErrMissingRequiredField:
		return http.StatusBadRequest
	case ErrGameNotStarted, ErrGameAlreadyCompleted, ErrFutureHole, ErrInvalidGameState:
		return http.StatusBadRequest
	case ErrSideBetNotEnabled, ErrInsufficientHoles, ErrCardsAlreadyDealt, ErrInvalidPuttCount, ErrGameNotCompleted:
		return http.StatusBadRequest
	case ErrPlayerLimitExceeded, ErrDuplicatePlayerName, ErrScoreAlreadyExists:
		return http.StatusConflict
	case ErrResourceNotFound, ErrGameNotFound, ErrPlayerNotFound, ErrScoreNotFound:
		return http.StatusNotFound
	case ErrInvalidToken, ErrExpiredToken:
		return http.StatusUnauthorized
	case ErrInsufficientPermissions:
		return http.StatusForbidden
	case ErrRateLimitExceeded:
		return http.StatusTooManyRequests
	case ErrServiceUnavailable:
		return http.StatusServiceUnavailable
	case ErrInternalServer:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// WriteHTTPError writes an APIError as an HTTP response
func WriteHTTPError(w http.ResponseWriter, apiErr *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.GetHTTPStatus())

	response := ErrorResponse{Error: apiErr}
	json.NewEncoder(w).Encode(response)
}

// FromError converts a regular error to an APIError
func FromError(err error) *APIError {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}
	return New(ErrInternalServer, "An unexpected error occurred")
}