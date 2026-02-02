package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type AcademicYear struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Status    string    `gorm:"type:varchar(20);default:'INACTIVE';index" json:"status"` // ACTIVE, INACTIVE
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Hook BeforeCreate untuk generate UUID
func (a *AcademicYear) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	return
}
