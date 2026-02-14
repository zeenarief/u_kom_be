package response

import (
	"time"
	"smart_school_be/internal/utils"
)

// ParentListResponse adalah DTO untuk tampilan list (ringkas)
// TIDAK mengandung NIK atau alamat lengkap
type ParentListResponse struct {
	ID          string  `json:"id"`
	FullName    string  `json:"full_name"`
	Gender      *string `json:"gender"`
	LifeStatus  *string `json:"life_status"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	Occupation  *string `json:"occupation"`
}

// ParentDetailResponse adalah DTO untuk tampilan detail (lengkap)
type ParentDetailResponse struct {
	ID             string      `json:"id"`
	FullName       string      `json:"full_name"`
	NIK            *string     `json:"nik,omitempty"` // Akan berisi plaintext
	Gender         *string     `json:"gender"`
	PlaceOfBirth   *string     `json:"place_of_birth"`
	DateOfBirth    *utils.Date `json:"date_of_birth"`
	LifeStatus     *string     `json:"life_status"`
	MaritalStatus  *string     `json:"marital_status"`
	PhoneNumber    *string     `json:"phone_number"`
	Email          *string     `json:"email"`
	EducationLevel *string     `json:"education_level"`
	Occupation     *string     `json:"occupation"`
	IncomeRange    *string     `json:"income_range"`
	Address        *string     `json:"address"`
	RT             *string     `json:"rt"`
	RW             *string     `json:"rw"`
	SubDistrict    *string     `json:"sub_district"`
	District       *string     `json:"district"`
	City           *string     `json:"city"`
	Province       *string     `json:"province"`
	PostalCode     *string     `json:"postal_code"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`

	User *UserLinkedResponse `json:"user"`
}
