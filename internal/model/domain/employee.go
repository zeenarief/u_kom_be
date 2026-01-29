package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Employee struct {
	ID               string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           *string    `gorm:"type:char(36);uniqueIndex" json:"user_id"` // Pointer untuk NULL
	FullName         string     `gorm:"type:varchar(100);not null" json:"full_name"`
	NIP              *string    `gorm:"type:varchar(50);uniqueIndex;column:nip" json:"nip"` // Pointer untuk NULL
	JobTitle         string     `gorm:"type:varchar(100);not null" json:"job_title"`
	NIK              string     `gorm:"type:text" json:"nik,omitempty"` // Akan dienkripsi
	Gender           string     `gorm:"type:varchar(10)" json:"gender"`
	PhoneNumber      *string    `gorm:"type:varchar(20);uniqueIndex" json:"phone_number"` // Pointer untuk NULL
	Address          string     `gorm:"type:text" json:"address"`
	DateOfBirth      *time.Time `gorm:"type:date" json:"date_of_birth"` // Pointer untuk NULL
	JoinDate         *time.Time `gorm:"type:date" json:"join_date"`     // Pointer untuk NULL
	EmploymentStatus string     `gorm:"type:varchar(20)" json:"employment_status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

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
