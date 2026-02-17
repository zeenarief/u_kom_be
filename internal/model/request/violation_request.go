package request

import "time"

type CreateViolationCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type UpdateViolationCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateViolationTypeRequest struct {
	CategoryID    string `json:"category_id" validate:"required,uuid"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	DefaultPoints int    `json:"default_points" validate:"required,min=0"`
}

type UpdateViolationTypeRequest struct {
	CategoryID    string `json:"category_id" validate:"omitempty,uuid"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DefaultPoints *int   `json:"default_points" validate:"omitempty,min=0"`
}

type CreateStudentViolationRequest struct {
	StudentID       string    `json:"student_id" validate:"required,uuid"`
	ViolationTypeID string    `json:"violation_type_id" validate:"required,uuid"`
	ViolationDate   time.Time `json:"violation_date" validate:"required"`
	ActionTaken     string    `json:"action_taken"`
	Notes           string    `json:"notes"`
}

type UpdateStudentViolationRequest struct {
	ViolationTypeID string     `json:"violation_type_id" validate:"omitempty,uuid"`
	ViolationDate   *time.Time `json:"violation_date"`
	Points          *int       `json:"points"`
	ActionTaken     string     `json:"action_taken"`
	Notes           string     `json:"notes"`
}
