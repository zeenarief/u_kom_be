package response

import (
	"time"
	"u_kom_be/internal/utils"
)

// StudentListResponse adalah DTO untuk tampilan list (ringkas)
type StudentListResponse struct {
	ID       string  `json:"id"`
	FullName string  `json:"full_name"`
	NISN     *string `json:"nisn"`
	NIM      *string `json:"nim"`
	Gender   string  `json:"gender"`
	City     string  `json:"city"`
}

// ParentRelationshipResponse adalah DTO untuk menampilkan relasi orang tua
type ParentRelationshipResponse struct {
	RelationshipType string             `json:"relationship_type"`
	Parent           ParentListResponse `json:"parent_info"` // Kita gunakan ListResponse yang ringkas
}

// GuardianInfoResponse adalah DTO generik untuk menampilkan data wali
// Ini bisa berisi data dari 'parent' ATAU 'guardian'
type GuardianInfoResponse struct {
	ID          string `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	// Field ini memberi tahu dari tabel mana data ini berasal
	Type string `json:"type"` // 'parent' or 'guardian'
	// Field ini memberi tahu apa hubungannya (cth: 'FATHER', 'MOTHER', 'UNCLE')
	Relationship string `json:"relationship"`
}

// StudentDetailResponse adalah DTO untuk tampilan detail (lengkap)
type StudentDetailResponse struct {
	ID           string     `json:"id"`
	FullName     string     `json:"full_name"`
	NoKK         string     `json:"no_kk,omitempty"` // Akan berisi plaintext
	NIK          string     `json:"nik,omitempty"`   // Akan berisi plaintext
	NISN         *string    `json:"nisn"`
	NIM          *string    `json:"nim"`
	Gender       string     `json:"gender"`
	PlaceOfBirth string     `json:"place_of_birth"`
	DateOfBirth  utils.Date `json:"date_of_birth"`
	Address      string     `json:"address"`
	RT           string     `json:"rt"`
	RW           string     `json:"rw"`
	SubDistrict  string     `json:"sub_district"`
	District     string     `json:"district"`
	City         string     `json:"city"`
	Province     string     `json:"province"`
	PostalCode   string     `json:"postal_code"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// Relasi M:N ke Parents (Sudah ada)
	Parents []ParentRelationshipResponse `json:"parents,omitempty"`

	// Guardian adalah relasi polimorfik 1:1. Pointer digunakan agar bisa 'null' jika tidak di-set
	Guardian *GuardianInfoResponse `json:"guardian,omitempty"`

	User *UserLinkedResponse `json:"user"`
}
