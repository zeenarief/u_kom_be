package request

import "u_kom_be/internal/utils"

// DTO untuk Create Parent
type ParentCreateRequest struct {
	FullName       string      `json:"full_name" binding:"required"`
	NIK            *string     `json:"nik"` // Akan dienkripsi
	Gender         *string     `json:"gender"`
	PlaceOfBirth   *string     `json:"place_of_birth"`
	DateOfBirth    *utils.Date `json:"date_of_birth"`
	LifeStatus     *string     `json:"life_status" binding:"omitempty,oneof=alive deceased"`
	MaritalStatus  *string     `json:"marital_status" binding:"omitempty,oneof=married divorced widowed"`
	PhoneNumber    *string     `json:"phone_number" binding:"omitempty"`
	Email          *string     `json:"email" binding:"omitempty,email"`
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
}

// DTO untuk Update Parent
type ParentUpdateRequest struct {
	FullName       string      `json:"full_name"`
	NIK            *string     `json:"nik"`
	Gender         *string     `json:"gender"`
	PlaceOfBirth   *string     `json:"place_of_birth"`
	DateOfBirth    *utils.Date `json:"date_of_birth"`
	LifeStatus     *string     `json:"life_status" binding:"omitempty,oneof=alive deceased"`
	MaritalStatus  *string     `json:"marital_status" binding:"omitempty,oneof=married divorced widowed"`
	PhoneNumber    *string     `json:"phone_number"`
	Email          *string     `json:"email" binding:"omitempty,email"`
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
}
