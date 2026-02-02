package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Schedule struct {
	ID                   string    `gorm:"primaryKey;type:char(36)" json:"id"`
	TeachingAssignmentID string    `gorm:"type:char(36);not null" json:"teaching_assignment_id"`
	DayOfWeek            int       `gorm:"type:tinyint;not null" json:"day_of_week"` // 1=Senin, 7=Minggu
	StartTime            string    `gorm:"type:time;not null" json:"start_time"`     // "07:00"
	EndTime              string    `gorm:"type:time;not null" json:"end_time"`       // "08:40"
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	// Relasi (Preload)
	TeachingAssignment TeachingAssignment `gorm:"foreignKey:TeachingAssignmentID" json:"teaching_assignment,omitempty"`
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateUUID()
	}
	return
}
