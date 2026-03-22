// utils/errors.go
package utils

import "net/http"

type AppError struct {
	Code    int    // HTTP status code
	Message string // error message
}

func (e *AppError) Error() string {
	return e.Message
}

// helper functions to create common errors
func NotFound(msg string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: msg}
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func InternalServerError(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func UnprocessableEntity(msg string) *AppError {
	return &AppError{Code: http.StatusUnprocessableEntity, Message: msg}
}
