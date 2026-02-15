package response

import (
	"time"
)

type ViolationCategoryResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ViolationTypeResponse struct {
	ID            string                     `json:"id"`
	CategoryID    string                     `json:"category_id"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	DefaultPoints int                        `json:"default_points"`
	Category      *ViolationCategoryResponse `json:"category,omitempty"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
}

type StudentViolationResponse struct {
	ID              string    `json:"id"`
	StudentID       string    `json:"student_id"`
	ViolationTypeID string    `json:"violation_type_id"`
	ViolationDate   time.Time `json:"violation_date"`
	Points          int       `json:"points"`
	ActionTaken     string    `json:"action_taken"`
	Notes           string    `json:"notes"`
	StudentName     string    `json:"student_name,omitempty"`   // Helper for UI
	ViolationName   string    `json:"violation_name,omitempty"` // Helper for UI
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
