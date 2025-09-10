package response

import "time"

// BaseResponse is the standardized response structure
type BaseResponse struct {
	Status    string      `json:"status"`              // "success" or "error"
	Message   string      `json:"message"`             // Short description
	Data      interface{} `json:"data,omitempty"`      // Response data for success
	Error     interface{} `json:"error,omitempty"`     // Error details for error responses
	Timestamp interface{} `json:"timestamp,omitempty"` // Timestamp for error responses
}

// Success creates a success response
func Success(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

// Error creates an error response
func Error(message string, errorData interface{}) BaseResponse {
	return BaseResponse{
		Status:    "error",
		Message:   message,
		Error:     errorData,
		Timestamp: time.Now(),
	}
}

// ErrorDetail for detailed error information
type ErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ValidationError for form validation errors
type ValidationError struct {
	Errors []ErrorDetail `json:"errors"`
}

// SimpleError for simple error messages
type SimpleError struct {
	Message string `json:"message"`
}
