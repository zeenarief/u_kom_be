package domain

import (
	"time"
)

type ViolationCategory struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ViolationType struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	CategoryID    string    `gorm:"type:char(36);not null;index" json:"category_id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	DefaultPoints int       `gorm:"not null;default:0" json:"default_points"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relationships
	Category *ViolationCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category,omitempty"`
}

type StudentViolation struct {
	ID              string    `gorm:"type:char(36);primaryKey" json:"id"`
	StudentID       string    `gorm:"type:char(36);not null;index" json:"student_id"`
	ViolationTypeID string    `gorm:"type:char(36);not null;index" json:"violation_type_id"`
	ViolationDate   time.Time `gorm:"type:datetime;not null" json:"violation_date"`
	Points          int       `gorm:"not null" json:"points"`
	ActionTaken     string    `gorm:"type:text" json:"action_taken"`
	Notes           string    `gorm:"type:text" json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relationships
	Student       *Student       `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	ViolationType *ViolationType `gorm:"foreignKey:ViolationTypeID" json:"violation_type,omitempty"`
}
