package apperrors

import (
	"fmt"
	"net/http"
)

type ErrorType string

const (
	NotFound     ErrorType = "NOT_FOUND"
	Conflict     ErrorType = "CONFLICT"
	Internal     ErrorType = "INTERNAL"
	Unauthorized ErrorType = "UNAUTHORIZED"
	BadRequest   ErrorType = "BAD_REQUEST"
	Forbidden    ErrorType = "FORBIDDEN"
)

type AppError struct {
	Type    ErrorType
	Message string
	Code    int // Optional: custom error code
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Factory functions

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFound,
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    Conflict,
		Message: message,
		Code:    http.StatusConflict,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Type:    Internal,
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Type:    BadRequest,
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    Unauthorized,
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    Forbidden,
		Message: message,
		Code:    http.StatusForbidden,
	}
}

// WrapError allows wrapping an existing error with an AppError type
func WrapError(err error, errType ErrorType, message string) *AppError {
	code := http.StatusInternalServerError
	switch errType {
	case NotFound:
		code = http.StatusNotFound
	case Conflict:
		code = http.StatusConflict
	case BadRequest:
		code = http.StatusBadRequest
	case Unauthorized:
		code = http.StatusUnauthorized
	case Forbidden:
		code = http.StatusForbidden
	}

	return &AppError{
		Type:    errType,
		Message: message,
		Code:    code,
		Err:     err,
	}
}
