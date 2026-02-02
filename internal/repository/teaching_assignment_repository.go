package repository

import (
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type TeachingAssignmentRepository interface {
	Create(assignment *domain.TeachingAssignment) error
	FindByClassroomID(classroomID string) ([]domain.TeachingAssignment, error)
	FindByTeacherID(teacherID string) ([]domain.TeachingAssignment, error)
	FindOne(classroomID, subjectID string) (*domain.TeachingAssignment, error)
	Delete(id string) error
}

type teachingAssignmentRepository struct {
	db *gorm.DB
}

func NewTeachingAssignmentRepository(db *gorm.DB) TeachingAssignmentRepository {
	return &teachingAssignmentRepository{db: db}
}

func (r *teachingAssignmentRepository) Create(assignment *domain.TeachingAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *teachingAssignmentRepository) FindByClassroomID(classroomID string) ([]domain.TeachingAssignment, error) {
	var assignments []domain.TeachingAssignment
	err := r.db.Preload("Classroom").Preload("Subject").Preload("Teacher").
		Where("classroom_id = ?", classroomID).
		Find(&assignments).Error
	return assignments, err
}

func (r *teachingAssignmentRepository) FindByTeacherID(teacherID string) ([]domain.TeachingAssignment, error) {
	var assignments []domain.TeachingAssignment
	err := r.db.Preload("Classroom").Preload("Subject").Preload("Teacher").
		Where("teacher_id = ?", teacherID).
		Find(&assignments).Error
	return assignments, err
}

func (r *teachingAssignmentRepository) FindOne(classroomID, subjectID string) (*domain.TeachingAssignment, error) {
	var assignment domain.TeachingAssignment
	err := r.db.Where("classroom_id = ? AND subject_id = ?", classroomID, subjectID).First(&assignment).Error
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *teachingAssignmentRepository) Delete(id string) error {
	return r.db.Delete(&domain.TeachingAssignment{}, "id = ?", id).Error
}
