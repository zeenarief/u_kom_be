package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Guardian struct {
	ID                    string    `gorm:"primaryKey;type:char(36)" json:"id"`
	FullName              string    `gorm:"type:varchar(100);not null" json:"full_name"`
	NIK                   *string   `gorm:"type:text" json:"nik,omitempty"`        // Akan dienkripsi
	NIKHash               *string   `gorm:"type:varchar(64);uniqueIndex" json:"-"` // Blind Index for Unique Check
	Gender                *string   `gorm:"type:varchar(10)" json:"gender"`
	PhoneNumber           *string   `gorm:"type:varchar(20);uniqueIndex" json:"phone_number"`
	Email                 *string   `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Address               *string   `gorm:"type:text" json:"address"`
	RT                    *string   `gorm:"type:varchar(3)" json:"rt"`
	RW                    *string   `gorm:"type:varchar(3)" json:"rw"`
	SubDistrict           *string   `gorm:"type:varchar(100)" json:"sub_district"`
	District              *string   `gorm:"type:varchar(100)" json:"district"`
	City                  *string   `gorm:"type:varchar(100)" json:"city"`
	Province              *string   `gorm:"type:varchar(100)" json:"province"`
	PostalCode            *string   `gorm:"type:varchar(5)" json:"postal_code"`
	RelationshipToStudent *string   `gorm:"type:varchar(50)" json:"relationship_to_student"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	UserID *string `gorm:"type:char(36);uniqueIndex" json:"user_id"`
	User   User    `gorm:"foreignKey:UserID;references:ID"`
}

// Hook BeforeCreate untuk generate UUID, sama seperti model lainnya
func (g *Guardian) BeforeCreate(tx *gorm.DB) (err error) {
	if g.ID == "" {
		g.ID = utils.GenerateUUID()
	}
	return
}
