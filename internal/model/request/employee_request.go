package request

import "smart_school_be/internal/utils"

// DTO untuk Create Employee
type EmployeeCreateRequest struct {
	FullName         string      `json:"full_name" binding:"required"`
	NIP              *string     `json:"nip" binding:"omitempty"`
	JobTitle         *string     `json:"job_title"` // Changed to pointer for nullable
	NIK              string      `json:"nik" binding:"required"`
	Gender           *string     `json:"gender,omitempty"` // Changed to pointer for nullable
	PhoneNumber      *string     `json:"phone_number" binding:"omitempty"`
	Address          *string     `json:"address,omitempty"` // Changed to pointer for nullable
	DateOfBirth      *utils.Date `json:"date_of_birth,omitempty"`
	JoinDate         *utils.Date `json:"join_date,omitempty"`
	EmploymentStatus *string     `json:"employment_status,omitempty"` // Changed to pointer for nullable
}

// DTO untuk Update Employee
type EmployeeUpdateRequest struct {
	FullName         string      `json:"full_name,omitempty"`
	NIP              *string     `json:"nip,omitempty"`
	JobTitle         *string     `json:"job_title,omitempty"` // Changed to pointer for nullable
	NIK              string      `json:"nik,omitempty"`
	Gender           *string     `json:"gender,omitempty"` // Changed to pointer for nullable
	PhoneNumber      *string     `json:"phone_number,omitempty"`
	Address          *string     `json:"address,omitempty"` // Changed to pointer for nullable
	DateOfBirth      *utils.Date `json:"date_of_birth,omitempty"`
	JoinDate         *utils.Date `json:"join_date,omitempty"`
	EmploymentStatus *string     `json:"employment_status,omitempty"` // Changed to pointer for nullable
}
