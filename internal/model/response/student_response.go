package response

import "time"

// StudentListResponse adalah DTO untuk tampilan list (ringkas)
type StudentListResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	NISN     string `json:"nisn"`
	NIM      string `json:"nim"`
	Gender   string `json:"gender"`
	City     string `json:"city"`
}

// StudentDetailResponse adalah DTO untuk tampilan detail (lengkap)
type StudentDetailResponse struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	NoKK         string    `json:"no_kk,omitempty"` // Akan berisi plaintext
	NIK          string    `json:"nik,omitempty"`   // Akan berisi plaintext
	NISN         string    `json:"nisn"`
	NIM          string    `json:"nim"`
	Gender       string    `json:"gender"`
	PlaceOfBirth string    `json:"place_of_birth"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	Address      string    `json:"address"`
	RT           string    `json:"rt"`
	RW           string    `json:"rw"`
	SubDistrict  string    `json:"sub_district"`
	District     string    `json:"district"`
	City         string    `json:"city"`
	Province     string    `json:"province"`
	PostalCode   string    `json:"postal_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
