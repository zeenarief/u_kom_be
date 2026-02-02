package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Subject struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Code        string    `gorm:"type:varchar(20);unique;not null" json:"code"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Type        string    `gorm:"type:varchar(50)" json:"type"` // Contoh: "Umum", "Kejuruan", "Muatan Lokal"
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (s *Subject) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateUUID()
	}
	return
}
