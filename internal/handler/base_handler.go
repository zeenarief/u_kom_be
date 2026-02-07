package handler

import (
	"net/http"
	"u_kom_be/internal/apperrors"
	"u_kom_be/internal/model/response"

	"github.com/gin-gonic/gin"
)

// HandleError maps AppError types to HTTP responses
func HandleError(c *gin.Context, err error) {
	if appErr, ok := err.(*apperrors.AppError); ok {
		switch appErr.Type {
		case apperrors.NotFound:
			NotFoundError(c, appErr.Message)
		case apperrors.Conflict:
			ErrorResponse(c, http.StatusConflict, appErr.Message, response.SimpleError{Message: appErr.Message})
		case apperrors.BadRequest:
			BadRequestError(c, appErr.Message, response.SimpleError{Message: appErr.Message})
		case apperrors.Unauthorized:
			UnauthorizedError(c, appErr.Message)
		case apperrors.Forbidden:
			ForbiddenError(c, appErr.Message)
		default:
			InternalServerError(c, appErr.Message)
		}
		return
	}

	// Default to 500 for unknown errors
	InternalServerError(c, err.Error())
}

// SuccessResponse sends a standardized success response
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, response.Success(message, data))
}

// CreatedResponse sends a 201 Created response
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, response.Success(message, data))
}

// ErrorResponse sends a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, message string, errorData interface{}) {
	c.JSON(statusCode, response.Error(message, errorData))
}

// ValidationErrorResponse sends validation errors
func ValidationErrorResponse(c *gin.Context, errors []response.ErrorDetail) {
	c.JSON(http.StatusBadRequest, response.Error("Validation failed", response.ValidationError{
		Errors: errors,
	}))
}

// InternalServerError sends a 500 error
func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, response.SimpleError{Message: message})
}

// NotFoundError sends a 404 error
func NotFoundError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, response.SimpleError{Message: message})
}

// UnauthorizedError sends a 401 error
func UnauthorizedError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, response.SimpleError{Message: message})
}

// BadRequestError sends a 400 error
func BadRequestError(c *gin.Context, message string, errorData interface{}) {
	ErrorResponse(c, http.StatusBadRequest, message, errorData)
}

// ForbiddenError sends a 403 error
func ForbiddenError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, response.SimpleError{Message: message})
}
