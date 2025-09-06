package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       string `gorm:"primaryKey;type:char(36)" json:"id"`
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"_"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = utils.GenerateUUID()
	}
	return
}
