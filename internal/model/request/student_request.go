package request

import "u_kom_be/internal/utils"

// DTO untuk Create Student
type StudentCreateRequest struct {
	FullName     string     `json:"full_name" binding:"required"`
	NoKK         string     `json:"no_kk"`
	NIK          string     `json:"nik"`
	NISN         string     `json:"nisn"`
	NIM          string     `json:"nim"`
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
	Status       string     `json:"status"`
	EntryYear    string     `json:"entry_year"`
	ExitYear     string     `json:"exit_year"`
}

// DTO untuk Update Student
type StudentUpdateRequest struct {
	FullName     string     `json:"full_name"`
	NoKK         string     `json:"no_kk"`
	NIK          string     `json:"nik"`
	NISN         string     `json:"nisn"`
	NIM          string     `json:"nim"`
	Gender       *string    `json:"gender"`         // Changed to pointer for nullable
	PlaceOfBirth *string    `json:"place_of_birth"` // Changed to pointer for nullable
	DateOfBirth  utils.Date `json:"date_of_birth"`  // Untuk time.Time, cek `!IsZero()`
	Address      *string    `json:"address"`        // Changed to pointer for nullable
	RT           *string    `json:"rt"`             // Changed to pointer for nullable
	RW           *string    `json:"rw"`             // Changed to pointer for nullable
	SubDistrict  *string    `json:"sub_district"`   // Changed to pointer for nullable
	District     *string    `json:"district"`       // Changed to pointer for nullable
	City         *string    `json:"city"`           // Changed to pointer for nullable
	Province     *string    `json:"province"`       // Changed to pointer for nullable
	PostalCode   *string    `json:"postal_code"`    // Changed to pointer for nullable
	Status       *string    `json:"status"`         // Changed to pointer for nullable
	EntryYear    *string    `json:"entry_year"`     // Changed to pointer for nullable
	ExitYear     *string    `json:"exit_year"`      // Changed to pointer for nullable
}
