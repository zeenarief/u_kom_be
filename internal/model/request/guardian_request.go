package request

// DTO untuk Create Guardian
type GuardianCreateRequest struct {
	FullName              string  `json:"full_name" binding:"required"`
	NIK                   *string `json:"nik"` // Akan dienkripsi
	Gender                *string `json:"gender"`
	PhoneNumber           *string `json:"phone_number" binding:"omitempty"` // Wajib ada sesuai skema -> Changed to optional
	Email                 *string `json:"email" binding:"omitempty,email"`
	Address               *string `json:"address"`
	RT                    *string `json:"rt"`
	RW                    *string `json:"rw"`
	SubDistrict           *string `json:"sub_district"`
	District              *string `json:"district"`
	City                  *string `json:"city"`
	Province              *string `json:"province"`
	PostalCode            *string `json:"postal_code"`
	RelationshipToStudent *string `json:"relationship_to_student"`
}

// DTO untuk Update Guardian
type GuardianUpdateRequest struct {
	FullName              string  `json:"full_name"`
	NIK                   *string `json:"nik"`
	Gender                *string `json:"gender"`
	PhoneNumber           *string `json:"phone_number"` // Tidak required di update
	Email                 *string `json:"email" binding:"omitempty,email"`
	Address               *string `json:"address"`
	RT                    *string `json:"rt"`
	RW                    *string `json:"rw"`
	SubDistrict           *string `json:"sub_district"`
	District              *string `json:"district"`
	City                  *string `json:"city"`
	Province              *string `json:"province"`
	PostalCode            *string `json:"postal_code"`
	RelationshipToStudent *string `json:"relationship_to_student"`
}
