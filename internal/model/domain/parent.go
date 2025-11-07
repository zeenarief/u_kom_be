package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type Parent struct {
	ID             string    `gorm:"primaryKey;type:char(36)" json:"id"`
	FullName       string    `gorm:"type:varchar(100);not null" json:"full_name"`
	NIK            string    `gorm:"type:text" json:"nik,omitempty"` // Akan dienkripsi
	Gender         string    `gorm:"type:varchar(10)" json:"gender"`
	PlaceOfBirth   string    `gorm:"type:varchar(100)" json:"place_of_birth"`
	DateOfBirth    time.Time `gorm:"type:date" json:"date_of_birth"`
	LifeStatus     string    `gorm:"type:varchar(10);default:'alive'" json:"life_status"`
	MaritalStatus  string    `gorm:"type:varchar(10)" json:"marital_status"`
	PhoneNumber    string    `gorm:"type:varchar(20);uniqueIndex" json:"phone_number"`
	Email          string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	EducationLevel string    `gorm:"type:varchar(50)" json:"education_level"`
	Occupation     string    `gorm:"type:varchar(100)" json:"occupation"`
	IncomeRange    string    `gorm:"type:varchar(50)" json:"income_range"`
	Address        string    `gorm:"type:text" json:"address"`
	RT             string    `gorm:"type:varchar(3)" json:"rt"`
	RW             string    `gorm:"type:varchar(3)" json:"rw"`
	SubDistrict    string    `gorm:"type:varchar(100)" json:"sub_district"`
	District       string    `gorm:"type:varchar(100)" json:"district"`
	City           string    `gorm:"type:varchar(100)" json:"city"`
	Province       string    `gorm:"type:varchar(100)" json:"province"`
	PostalCode     string    `gorm:"type:varchar(5)" json:"postal_code"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Hook BeforeCreate untuk generate UUID, sama seperti Student/Role
func (p *Parent) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = utils.GenerateUUID()
	}
	return
}
