package repository

import (
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type GradeRepository interface {
	// Assessments
	CreateAssessment(assessment *domain.Assessment) error
	UpdateAssessment(assessment *domain.Assessment) error
	FindAssessmentByID(id string) (*domain.Assessment, error)
	GetAssessmentsByTeachingAssignment(teachingAssignmentID string) ([]domain.Assessment, error)

	// Scores
	SaveStudentScore(score *domain.StudentScore) error
	GetScoresByAssessmentID(assessmentID string) ([]domain.StudentScore, error)
	GetScoreByAssessmentAndStudent(assessmentID, studentID string) (*domain.StudentScore, error) // Helper to check existence
}

type gradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) GradeRepository {
	return &gradeRepository{db: db}
}

func (r *gradeRepository) CreateAssessment(assessment *domain.Assessment) error {
	return r.db.Create(assessment).Error
}

func (r *gradeRepository) UpdateAssessment(assessment *domain.Assessment) error {
	return r.db.Save(assessment).Error
}

func (r *gradeRepository) FindAssessmentByID(id string) (*domain.Assessment, error) {
	var assessment domain.Assessment
	err := r.db.Preload("TeachingAssignment").
		Preload("TeachingAssignment.Subject").
		Preload("TeachingAssignment.Classroom").
		Preload("Scores").
		Preload("Scores.Student").
		First(&assessment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &assessment, nil
}

func (r *gradeRepository) GetAssessmentsByTeachingAssignment(teachingAssignmentID string) ([]domain.Assessment, error) {
	var assessments []domain.Assessment
	err := r.db.Where("teaching_assignment_id = ?", teachingAssignmentID).
		Order("date DESC").
		Find(&assessments).Error
	return assessments, err
}

func (r *gradeRepository) SaveStudentScore(score *domain.StudentScore) error {
	// Upsert: If ID exists Update, else Create.
	// However, usually for bulk input we check existence by assessment_id + student_id
	var existing domain.StudentScore
	err := r.db.Where("assessment_id = ? AND student_id = ?", score.AssessmentID, score.StudentID).First(&existing).Error

	if err == nil {
		// Update existing
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"score":      score.Score,
			"feedback":   score.Feedback,
			"updated_at": score.UpdatedAt, // Let GORM handle time, or pass it explicitly
		}).Error
	}

	// Create new
	return r.db.Create(score).Error
}

func (r *gradeRepository) GetScoresByAssessmentID(assessmentID string) ([]domain.StudentScore, error) {
	var scores []domain.StudentScore
	err := r.db.Preload("Student").
		Where("assessment_id = ?", assessmentID).
		Find(&scores).Error
	return scores, err
}

func (r *gradeRepository) GetScoreByAssessmentAndStudent(assessmentID, studentID string) (*domain.StudentScore, error) {
	var score domain.StudentScore
	err := r.db.Where("assessment_id = ? AND student_id = ?", assessmentID, studentID).First(&score).Error
	if err != nil {
		return nil, err
	}
	return &score, nil
}
