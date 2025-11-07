package response

import "time"

// GuardianListResponse adalah DTO untuk tampilan list (ringkas)
// TIDAK mengandung NIK atau alamat lengkap
type GuardianListResponse struct {
	ID                    string `json:"id"`
	FullName              string `json:"full_name"`
	PhoneNumber           string `json:"phone_number"`
	Email                 string `json:"email"`
	RelationshipToStudent string `json:"relationship_to_student"`
}

// GuardianDetailResponse adalah DTO untuk tampilan detail (lengkap)
type GuardianDetailResponse struct {
	ID                    string    `json:"id"`
	FullName              string    `json:"full_name"`
	NIK                   string    `json:"nik,omitempty"` // Akan berisi plaintext
	Gender                string    `json:"gender"`
	PhoneNumber           string    `json:"phone_number"`
	Email                 string    `json:"email"`
	Address               string    `json:"address"`
	RT                    string    `json:"rt"`
	RW                    string    `json:"rw"`
	SubDistrict           string    `json:"sub_district"`
	District              string    `json:"district"`
	City                  string    `json:"city"`
	Province              string    `json:"province"`
	PostalCode            string    `json:"postal_code"`
	RelationshipToStudent string    `json:"relationship_to_student"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
