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

type StudentViolationListResponse struct {
	ID                string    `json:"id"`
	StudentID         string    `json:"student_id"`
	StudentName       string    `json:"student_name"`
	ViolationDate     time.Time `json:"violation_date"`
	Points            int       `json:"points"`
	ViolationName     string    `json:"violation_name"`
	ViolationCategory string    `json:"violation_category"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type StudentViolationDetailResponse struct {
	ID                string    `json:"id"`
	StudentID         string    `json:"student_id"`
	StudentName       string    `json:"student_name"`
	StudentNIM        *string   `json:"student_nim"`
	StudentNISN       *string   `json:"student_nisn"`
	StudentClass      string    `json:"student_class,omitempty"`
	ViolationTypeID   string    `json:"violation_type_id"`
	ViolationName     string    `json:"violation_name"`
	ViolationCategory string    `json:"violation_category"`
	ViolationDate     time.Time `json:"violation_date"`
	Points            int       `json:"points"`
	ActionTaken       string    `json:"action_taken"`
	Notes             string    `json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
