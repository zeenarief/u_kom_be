package request

import "smart_school_be/internal/utils"

// DTO untuk Create Student
type StudentCreateRequest struct {
	FullName     string     `json:"full_name" form:"full_name" binding:"required"`
	NoKK         string     `json:"no_kk" form:"no_kk"`
	NIK          string     `json:"nik" form:"nik"`
	NISN         string     `json:"nisn" form:"nisn"`
	NIM          string     `json:"nim" form:"nim"`
	Gender       string     `json:"gender" form:"gender"`
	PlaceOfBirth string     `json:"place_of_birth" form:"place_of_birth"`
	DateOfBirth  utils.Date `json:"date_of_birth" form:"date_of_birth"`
	Address      string     `json:"address"  form:"address"`
	RT           string     `json:"rt" form:"rt"`
	RW           string     `json:"rw" form:"rw"`
	SubDistrict  string     `json:"sub_district" form:"sub_district"`
	District     string     `json:"district" form:"district"`
	City         string     `json:"city" form:"city"`
	Province     string     `json:"province" form:"province"`
	PostalCode   string     `json:"postal_code" form:"postal_code"`
	Status       string     `json:"status" form:"status"`
	EntryYear    string     `json:"entry_year" form:"entry_year"`
	ExitYear     string     `json:"exit_year" form:"exit_year"`
}

// DTO untuk Update Student
type StudentUpdateRequest struct {
	FullName     string      `json:"full_name" form:"full_name"`
	NoKK         string      `json:"no_kk" form:"no_kk"`
	NIK          *string     `json:"nik" form:"nik"` // Changed to pointer
	NISN         string      `json:"nisn" form:"nisn"`
	NIM          string      `json:"nim" form:"nim"`
	Gender       *string     `json:"gender" form:"gender"`                 // Changed to pointer for nullable
	PlaceOfBirth *string     `json:"place_of_birth" form:"place_of_birth"` // Changed to pointer for nullable
	DateOfBirth  *utils.Date `json:"date_of_birth" form:"date_of_birth"`   // Changed to pointer
	Address      *string     `json:"address" form:"address"`               // Changed to pointer for nullable
	RT           *string     `json:"rt" form:"rt"`                         // Changed to pointer for nullable
	RW           *string     `json:"rw" form:"rw"`                         // Changed to pointer for nullable
	SubDistrict  *string     `json:"sub_district" form:"sub_district"`     // Changed to pointer for nullable
	District     *string     `json:"district" form:"district"`             // Changed to pointer for nullable
	City         *string     `json:"city" form:"city"`                     // Changed to pointer for nullable
	Province     *string     `json:"province" form:"province"`             // Changed to pointer for nullable
	PostalCode   *string     `json:"postal_code" form:"postal_code"`       // Changed to pointer for nullable
	Status       *string     `json:"status" form:"status"`                 // Changed to pointer for nullable
	EntryYear    *string     `json:"entry_year" form:"entry_year"`         // Changed to pointer for nullable
	ExitYear     *string     `json:"exit_year" form:"exit_year"`           // Changed to pointer for nullable
}
