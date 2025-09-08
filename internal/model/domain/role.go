package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          string       `gorm:"primaryKey;type:char(36)" json:"id"`
	Name        string       `gorm:"uniqueIndex;not null" json:"name"`
	Description string       `json:"description"`
	IsDefault   bool         `gorm:"default:false" json:"is_default"`
	Permissions []Permission `gorm:"many2many:role_permission;" json:"permissions,omitempty"`
	Users       []User       `gorm:"many2many:user_role;" json:"users,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = utils.GenerateUUID()
	}
	return
}
