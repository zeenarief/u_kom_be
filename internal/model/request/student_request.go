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
}

// DTO untuk Update Student
type StudentUpdateRequest struct {
	FullName     string     `json:"full_name"`
	NoKK         string     `json:"no_kk"`
	NIK          string     `json:"nik"`
	NISN         string     `json:"nisn"`
	NIM          string     `json:"nim"`
	Gender       string     `json:"gender"`
	PlaceOfBirth string     `json:"place_of_birth"`
	DateOfBirth  utils.Date `json:"date_of_birth"` // Untuk time.Time, cek `!IsZero()`
	Address      string     `json:"address"`
	RT           string     `json:"rt"`
	RW           string     `json:"rw"`
	SubDistrict  string     `json:"sub_district"`
	District     string     `json:"district"`
	City         string     `json:"city"`
	Province     string     `json:"province"`
	PostalCode   string     `json:"postal_code"`
}
