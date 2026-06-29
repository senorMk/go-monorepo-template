// Package httpx contains shared HTTP helpers: a typed error envelope and
// request decoding/validation, mirroring the NestJS template's behaviour.
package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// statusToError maps HTTP status codes to the stable error codes returned in
// the JSON envelope (matches the TypeScript HttpExceptionFilter).
var statusToError = map[int]string{
	http.StatusBadRequest:          "BAD_REQUEST",
	http.StatusUnauthorized:        "UNAUTHORIZED",
	http.StatusForbidden:           "FORBIDDEN",
	http.StatusNotFound:            "NOT_FOUND",
	http.StatusConflict:            "CONFLICT",
	http.StatusUnprocessableEntity: "UNPROCESSABLE",
	http.StatusInternalServerError: "INTERNAL_ERROR",
}

func errorCode(status int) string {
	if code, ok := statusToError[status]; ok {
		return code
	}
	return "INTERNAL_ERROR"
}

// AppError is an HTTP-aware error carrying a status code and client message.
type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string { return e.Message }

// Error constructors mirror the NestJS HttpException helpers.
func NewError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}
func BadRequest(message string) *AppError   { return NewError(http.StatusBadRequest, message) }
func Unauthorized(message string) *AppError { return NewError(http.StatusUnauthorized, message) }
func Forbidden(message string) *AppError    { return NewError(http.StatusForbidden, message) }
func NotFound(message string) *AppError     { return NewError(http.StatusNotFound, message) }
func Conflict(message string) *AppError     { return NewError(http.StatusConflict, message) }
func Unprocessable(message string) *AppError {
	return NewError(http.StatusUnprocessableEntity, message)
}

type errorBody struct {
	StatusCode int    `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`
}

// WriteError renders err as the standard JSON error envelope. Unknown errors
// are masked as 500 and logged.
func WriteError(w http.ResponseWriter, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		slog.Error("unhandled error", "err", err)
		appErr = NewError(http.StatusInternalServerError, "Internal server error")
	}
	WriteJSON(w, appErr.Status, errorBody{
		StatusCode: appErr.Status,
		Error:      errorCode(appErr.Status),
		Message:    appErr.Message,
	})
}

// WriteJSON writes v as a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "err", err)
	}
}

// NoContent writes a 204 response.
func NoContent(w http.ResponseWriter) { w.WriteHeader(http.StatusNoContent) }
