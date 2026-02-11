package response

import (
	"time"
	"u_kom_be/internal/utils"
)

// EmployeeListResponse adalah DTO untuk tampilan list (ringkas)
// TIDAK mengandung NIK atau alamat lengkap
type EmployeeListResponse struct {
	ID               string  `json:"id"`
	FullName         string  `json:"full_name"`
	NIP              *string `json:"nip,omitempty"`
	JobTitle         *string `json:"job_title,omitempty"` // Changed to pointer
	PhoneNumber      *string `json:"phone_number,omitempty"`
	EmploymentStatus *string `json:"employment_status,omitempty"` // Changed to pointer
	UserID           *string `json:"user_id,omitempty"`           // ID akun user yang terhubung
}

// EmployeeDetailResponse adalah DTO untuk tampilan detail (lengkap)
type EmployeeDetailResponse struct {
	ID               string              `json:"id"`
	User             *UserLinkedResponse `json:"user"` // ID akun user yang terhubung
	FullName         string              `json:"full_name"`
	NIP              *string             `json:"nip,omitempty"`
	JobTitle         *string             `json:"job_title,omitempty"` // Changed to pointer
	NIK              string              `json:"nik,omitempty"`       // Akan berisi plaintext
	Gender           *string             `json:"gender,omitempty"`    // Changed to pointer
	PhoneNumber      *string             `json:"phone_number,omitempty"`
	Address          *string             `json:"address,omitempty"` // Changed to pointer
	DateOfBirth      *utils.Date         `json:"date_of_birth,omitempty"`
	JoinDate         *utils.Date         `json:"join_date,omitempty"`
	EmploymentStatus *string             `json:"employment_status,omitempty"` // Changed to pointer
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	// Kita bisa tambahkan UserInfo (dari user_id) di sini nanti jika perlu
}
