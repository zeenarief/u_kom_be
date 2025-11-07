package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type Guardian struct {
	ID                    string    `gorm:"primaryKey;type:char(36)" json:"id"`
	FullName              string    `gorm:"type:varchar(100);not null" json:"full_name"`
	NIK                   string    `gorm:"type:text" json:"nik,omitempty"` // Akan dienkripsi
	Gender                string    `gorm:"type:varchar(10)" json:"gender"`
	PhoneNumber           string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"phone_number"`
	Email                 string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Address               string    `gorm:"type:text" json:"address"`
	RT                    string    `gorm:"type:varchar(3)" json:"rt"`
	RW                    string    `gorm:"type:varchar(3)" json:"rw"`
	SubDistrict           string    `gorm:"type:varchar(100)" json:"sub_district"`
	District              string    `gorm:"type:varchar(100)" json:"district"`
	City                  string    `gorm:"type:varchar(100)" json:"city"`
	Province              string    `gorm:"type:varchar(100)" json:"province"`
	PostalCode            string    `gorm:"type:varchar(5)" json:"postal_code"`
	RelationshipToStudent string    `gorm:"type:varchar(50)" json:"relationship_to_student"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Hook BeforeCreate untuk generate UUID, sama seperti model lainnya
func (g *Guardian) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == "" {
		g.ID = utils.GenerateUUID()
	}
	return
}
