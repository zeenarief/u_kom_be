package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

type Classroom struct {
	ID                string    `gorm:"primaryKey;type:char(36)" json:"id"`
	AcademicYearID    string    `gorm:"type:char(36);not null" json:"academic_year_id"`
	HomeroomTeacherID *string   `gorm:"type:char(36)" json:"homeroom_teacher_id"`
	Name              string    `gorm:"type:varchar(50);not null" json:"name"`
	Level             string    `gorm:"type:varchar(10);not null" json:"level"`
	Major             string    `gorm:"type:varchar(50)" json:"major"`
	Description       string    `gorm:"type:text" json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	TotalStudents int64 `gorm:"->"` // Menampung hasil subquery COUNT

	// Relations
	AcademicYear    AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year,omitempty"`
	HomeroomTeacher *Employee    `gorm:"foreignKey:HomeroomTeacherID" json:"homeroom_teacher,omitempty"`

	// HasMany relation via pivot struct
	StudentClassrooms []StudentClassroom `gorm:"foreignKey:ClassroomID" json:"student_classrooms,omitempty"`
}

func (c *Classroom) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = utils.GenerateUUID()
	}
	return
}

// Struct Pivot Eksplisit
type StudentClassroom struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	ClassroomID string    `gorm:"type:char(36);not null" json:"classroom_id"`
	StudentID   string    `gorm:"type:char(36);not null" json:"student_id"`
	Status      string    `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Student   Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Classroom Classroom `gorm:"foreignKey:ClassroomID" json:"classroom,omitempty"`
}

func (sc *StudentClassroom) BeforeCreate(tx *gorm.DB) (err error) {
	if sc.ID == "" {
		sc.ID = utils.GenerateUUID()
	}
	return
}
