package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type AttendanceSession struct {
	ID         string    `gorm:"primaryKey;type:char(36)" json:"id"`
	ScheduleID string    `gorm:"type:char(36);not null" json:"schedule_id"`
	Date       time.Time `gorm:"type:date;not null" json:"date"`
	Topic      string    `gorm:"type:text" json:"topic"`
	Notes      string    `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Schedule Schedule           `gorm:"foreignKey:ScheduleID" json:"schedule,omitempty"`
	Details  []AttendanceDetail `gorm:"foreignKey:AttendanceSessionID" json:"details,omitempty"`
}

func (a *AttendanceSession) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	return
}

type AttendanceDetail struct {
	ID                  string    `gorm:"primaryKey;type:char(36)" json:"id"`
	AttendanceSessionID string    `gorm:"type:char(36);not null" json:"attendance_session_id"`
	StudentID           string    `gorm:"type:char(36);not null" json:"student_id"`
	Status              string    `gorm:"type:varchar(10);default:'PRESENT'" json:"status"` // PRESENT, SICK, PERMISSION, ABSENT
	Notes               string    `gorm:"type:varchar(255)" json:"notes"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Relations
	Student Student `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

func (ad *AttendanceDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if ad.ID == "" {
		ad.ID = utils.GenerateUUID()
	}
	return
}
