package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	Roles       []Role    `gorm:"many2many:role_permission;" json:"roles,omitempty"`
	Users       []User    `gorm:"many2many:user_permission;" json:"users,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = utils.GenerateUUID()
	}
	return
}
