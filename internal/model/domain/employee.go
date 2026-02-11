package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Employee struct {
	ID               string      `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           *string     `gorm:"type:char(36);uniqueIndex" json:"user_id"` // Pointer untuk NULL
	FullName         string      `gorm:"type:varchar(100);not null" json:"full_name"`
	NIP              *string     `gorm:"type:varchar(50);uniqueIndex;column:nip" json:"nip"` // Pointer untuk NULL
	JobTitle         *string     `gorm:"type:varchar(100)" json:"job_title"`                 // Changed to pointer for nullable
	NIK              string      `gorm:"type:text" json:"nik,omitempty"`                     // Akan dienkripsi
	NIKHash          string      `gorm:"type:varchar(64);uniqueIndex" json:"-"`              // Blind Index for Unique Check
	Gender           *string     `gorm:"type:varchar(10)" json:"gender"`                     // Changed to pointer for nullable
	PhoneNumber      *string     `gorm:"type:varchar(20);uniqueIndex" json:"phone_number"`   // Pointer untuk NULL
	Address          *string     `gorm:"type:text" json:"address"`                           // Changed to pointer for nullable
	DateOfBirth      *utils.Date `gorm:"type:date" json:"date_of_birth"`
	JoinDate         *utils.Date `gorm:"type:date" json:"join_date"`
	EmploymentStatus *string     `gorm:"type:varchar(20)" json:"employment_status"` // Changed to pointer for nullable
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`

	// Relasi GORM (opsional, tapi bagus untuk Preload jika diperlukan)
	User User `gorm:"foreignKey:UserID;references:ID"`
}

// Hook BeforeCreate untuk generate UUID
func (e *Employee) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = utils.GenerateUUID()
	}
	return
}
