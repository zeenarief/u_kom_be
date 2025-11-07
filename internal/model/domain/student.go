package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	FullName     string    `gorm:"type:varchar(100);not null" json:"full_name"`
	NoKK         string    `gorm:"type:text" json:"no_kk,omitempty"` // akan dienkripsi
	NIK          string    `gorm:"type:text" json:"nik,omitempty"`   // akan dienkripsi
	NISN         string    `gorm:"type:varchar(20);uniqueIndex" json:"nisn"`
	NIM          string    `gorm:"type:varchar(20);uniqueIndex" json:"nim"`
	Gender       string    `gorm:"type:varchar(10)" json:"gender"`
	PlaceOfBirth string    `gorm:"type:varchar(100)" json:"place_of_birth"`
	DateOfBirth  time.Time `gorm:"type:date" json:"date_of_birth"`
	Address      string    `gorm:"type:text" json:"address"`
	RT           string    `gorm:"type:varchar(3)" json:"rt"`
	RW           string    `gorm:"type:varchar(3)" json:"rw"`
	SubDistrict  string    `gorm:"type:varchar(100)" json:"sub_district"`
	District     string    `gorm:"type:varchar(100)" json:"district"`
	City         string    `gorm:"type:varchar(100)" json:"city"`
	Province     string    `gorm:"type:varchar(100)" json:"province"`
	PostalCode   string    `gorm:"type:varchar(5)" json:"postal_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi ke tabel pivot StudentParent
	Parents []StudentParent `gorm:"foreignKey:StudentID" json:"parents,omitempty"`

	GuardianID   *string `gorm:"type:char(36);index:idx_student_guardian" json:"guardian_id"`
	GuardianType *string `gorm:"type:varchar(20);index:idx_student_guardian" json:"guardian_type"`
}

// Hook BeforeCreate untuk generate UUID
func (s *Student) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateUUID()
	}
	return
}
