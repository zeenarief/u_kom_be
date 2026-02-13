package domain

import (
	"time"
	"u_kom_be/internal/utils"

	"gorm.io/gorm"
)

// Assessment merepresentasikan tugas, ujian, atau quiz
type Assessment struct {
	ID                   string     `gorm:"primaryKey;type:char(36)" json:"id"`
	TeachingAssignmentID string     `gorm:"type:char(36);not null" json:"teaching_assignment_id"`
	Title                string     `gorm:"type:varchar(255);not null" json:"title"`
	Type                 string     `gorm:"type:varchar(50);not null" json:"type"` // ASSIGNMENT, MID_EXAM, FINAL_EXAM, QUIZ
	MaxScore             int        `gorm:"type:int;default:100" json:"max_score"`
	Date                 utils.Date `gorm:"type:date;not null" json:"date"`
	Description          string     `gorm:"type:text" json:"description"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	// Relasi
	TeachingAssignment TeachingAssignment `gorm:"foreignKey:TeachingAssignmentID" json:"teaching_assignment,omitempty"`
	Scores             []StudentScore     `gorm:"foreignKey:AssessmentID" json:"scores,omitempty"`
}

func (a *Assessment) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = utils.GenerateUUID()
	}
	return
}

// StudentScore menyimpan nilai siswa untuk assessment tertentu
type StudentScore struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	AssessmentID string    `gorm:"type:char(36);not null" json:"assessment_id"`
	StudentID    string    `gorm:"type:char(36);not null" json:"student_id"`
	Score        float64   `gorm:"type:float;default:0" json:"score"` // Bisa desimal
	Feedback     string    `gorm:"type:text" json:"feedback"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi
	Assessment Assessment `gorm:"foreignKey:AssessmentID" json:"assessment,omitempty"`
	Student    Student    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}

func (s *StudentScore) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateUUID()
	}
	return
}
