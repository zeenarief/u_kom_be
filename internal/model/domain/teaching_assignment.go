package domain

import (
	"time"
	"smart_school_be/internal/utils"

	"gorm.io/gorm"
)

type TeachingAssignment struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	ClassroomID string    `gorm:"type:char(36);not null" json:"classroom_id"`
	SubjectID   string    `gorm:"type:char(36);not null" json:"subject_id"`
	TeacherID   string    `gorm:"type:char(36);not null" json:"teacher_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations (Preload)
	Classroom Classroom `gorm:"foreignKey:ClassroomID" json:"classroom,omitempty"`
	Subject   Subject   `gorm:"foreignKey:SubjectID" json:"subject,omitempty"`
	Teacher   Employee  `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

func (t *TeachingAssignment) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = utils.GenerateUUID()
	}
	return
}
