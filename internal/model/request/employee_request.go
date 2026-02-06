package request

import "time"

// DTO untuk Create Employee
type EmployeeCreateRequest struct {
	FullName         string     `json:"full_name" binding:"required"`
	NIP              *string    `json:"nip" binding:"omitempty"`
	JobTitle         string     `json:"job_title" binding:"required"`
	NIK              string     `json:"nik" binding:"required"`
	Gender           string     `json:"gender,omitempty"`
	PhoneNumber      *string    `json:"phone_number" binding:"omitempty"`
	Address          string     `json:"address,omitempty"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	JoinDate         *time.Time `json:"join_date,omitempty"`
	EmploymentStatus string     `json:"employment_status,omitempty"`
}

// DTO untuk Update Employee
type EmployeeUpdateRequest struct {
	FullName         string     `json:"full_name,omitempty"`
	NIP              *string    `json:"nip,omitempty"`
	JobTitle         string     `json:"job_title,omitempty"`
	NIK              string     `json:"nik,omitempty"`
	Gender           string     `json:"gender,omitempty"`
	PhoneNumber      *string    `json:"phone_number,omitempty"`
	Address          string     `json:"address,omitempty"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	JoinDate         *time.Time `json:"join_date,omitempty"`
	EmploymentStatus string     `json:"employment_status,omitempty"`
}
