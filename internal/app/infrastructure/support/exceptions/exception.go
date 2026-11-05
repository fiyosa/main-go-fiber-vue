package exceptions

import (
	"errors"
	"net/http"
)

// Exception is the single error shape understood by the global HTTP handler.
// Code is a stable machine-readable identifier, MessageKey is a translation
// key, and Err keeps the original cause for logging only.
type Exception struct {
	Code       string
	Status     int
	MessageKey string
	Fields     map[string][]string
	Err        error
}

func (e *Exception) Error() string {
	return e.Code + ": " + e.MessageKey
}

func (e *Exception) Unwrap() error {
	return e.Err
}

// New builds an exception. A nil cause is allowed for sentinel errors.
func New(code string, status int, messageKey string, cause error) *Exception {
	return &Exception{
		Code:       code,
		Status:     status,
		MessageKey: messageKey,
		Err:        cause,
	}
}

// Validation builds a 422 exception carrying per-field messages.
func Validation(fields map[string][]string) *Exception {
	return &Exception{
		Code:       "common.validation_failed",
		Status:     http.StatusUnprocessableEntity,
		MessageKey: "common.validation_failed",
		Fields:     fields,
	}
}

// As extracts an Exception from an error chain.
func As(err error) (*Exception, bool) {
	var ex *Exception
	if errors.As(err, &ex) {
		return ex, true
	}
	return nil, false
}

// Common sentinel exceptions. They are immutable in practice: always create a
// new one when a specific cause or message key is needed.
var (
	ErrUnauthenticated = New("auth.unauthenticated", http.StatusUnauthorized, "auth.unauthenticated", nil)
	ErrForbidden       = New("auth.forbidden", http.StatusForbidden, "common.forbidden", nil)
	ErrNotFound        = New("common.not_found", http.StatusNotFound, "common.not_found", nil)
)

// Unauthorized returns a 401 exception using the given message key.
func Unauthorized(messageKey string, cause error) *Exception {
	return New("auth.unauthenticated", http.StatusUnauthorized, messageKey, cause)
}

// Forbidden returns a 403 exception.
func Forbidden(cause error) *Exception {
	return New("auth.forbidden", http.StatusForbidden, "common.forbidden", cause)
}

// BadRequest returns a 400 exception.
func BadRequest(cause error) *Exception {
	return New("common.bad_request", http.StatusBadRequest, "common.bad_request", cause)
}

// NotFound returns a 404 exception.
func NotFound(cause error) *Exception {
	return New("common.not_found", http.StatusNotFound, "common.not_found", cause)
}

// MethodNotAllowed returns a 405 exception.
func MethodNotAllowed(cause error) *Exception {
	return New("common.method_not_allowed", http.StatusMethodNotAllowed, "common.method_not_allowed", cause)
}

// ServiceUnavailable returns a 503 exception.
func ServiceUnavailable(cause error) *Exception {
	return New("common.unavailable", http.StatusServiceUnavailable, "common.try_again_later", cause)
}

// Internal returns a 500 exception and keeps the cause for logging.
func Internal(cause error) *Exception {
	return New("common.internal", http.StatusInternalServerError, "common.try_again_later", cause)
}
