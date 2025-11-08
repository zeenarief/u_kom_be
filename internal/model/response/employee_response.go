package response

import "time"

// EmployeeListResponse adalah DTO untuk tampilan list (ringkas)
// TIDAK mengandung NIK atau alamat lengkap
type EmployeeListResponse struct {
	ID               string  `json:"id"`
	FullName         string  `json:"full_name"`
	NIP              *string `json:"nip,omitempty"`
	JobTitle         string  `json:"job_title"`
	PhoneNumber      *string `json:"phone_number,omitempty"`
	EmploymentStatus string  `json:"employment_status,omitempty"`
	UserID           *string `json:"user_id,omitempty"` // ID akun user yang terhubung
}

// EmployeeDetailResponse adalah DTO untuk tampilan detail (lengkap)
type EmployeeDetailResponse struct {
	ID               string     `json:"id"`
	UserID           *string    `json:"user_id,omitempty"` // ID akun user yang terhubung
	FullName         string     `json:"full_name"`
	NIP              *string    `json:"nip,omitempty"`
	JobTitle         string     `json:"job_title"`
	NIK              string     `json:"nik,omitempty"` // Akan berisi plaintext
	Gender           string     `json:"gender,omitempty"`
	PhoneNumber      *string    `json:"phone_number,omitempty"`
	Address          string     `json:"address,omitempty"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	JoinDate         *time.Time `json:"join_date,omitempty"`
	EmploymentStatus string     `json:"employment_status,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	// Kita bisa tambahkan UserInfo (dari user_id) di sini nanti jika perlu
}
