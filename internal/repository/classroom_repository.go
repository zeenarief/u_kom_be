package repository

import (
	"errors"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type ClassroomRepository interface {
	Create(classroom *domain.Classroom) error
	FindAll(academicYearID string) ([]domain.Classroom, error)
	FindByID(id string) (*domain.Classroom, error)
	Update(classroom *domain.Classroom) error
	Delete(id string) error

	// Student Management
	AddStudents(studentClassrooms []domain.StudentClassroom) error
	RemoveStudent(classroomID string, studentID string) error
	IsStudentInClass(classroomID string, studentID string) (bool, error)
}

type classroomRepository struct {
	db *gorm.DB
}

func NewClassroomRepository(db *gorm.DB) ClassroomRepository {
	return &classroomRepository{db: db}
}

func (r *classroomRepository) Create(classroom *domain.Classroom) error {
	return r.db.Create(classroom).Error
}

func (r *classroomRepository) FindAll(academicYearID string) ([]domain.Classroom, error) {
	var classrooms []domain.Classroom
	query := r.db.Preload("AcademicYear").Preload("HomeroomTeacher").
		// Preload count students (Subquery optimization)
		Select("classrooms.*, (SELECT COUNT(*) FROM student_classrooms WHERE student_classrooms.classroom_id = classrooms.id) as total_students")

	if academicYearID != "" {
		query = query.Where("academic_year_id = ?", academicYearID)
	}

	err := query.Find(&classrooms).Error
	return classrooms, err
}

func (r *classroomRepository) FindByID(id string) (*domain.Classroom, error) {
	var classroom domain.Classroom
	err := r.db.Preload("AcademicYear").
		Preload("HomeroomTeacher").
		Preload("StudentClassrooms").         // Load Pivot
		Preload("StudentClassrooms.Student"). // Load Student Data
		First(&classroom, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &classroom, err
}

func (r *classroomRepository) Update(classroom *domain.Classroom) error {
	return r.db.Save(classroom).Error
}

func (r *classroomRepository) Delete(id string) error {
	return r.db.Delete(&domain.Classroom{}, "id = ?", id).Error
}

func (r *classroomRepository) AddStudents(studentClassrooms []domain.StudentClassroom) error {
	return r.db.Create(&studentClassrooms).Error
}

func (r *classroomRepository) RemoveStudent(classroomID string, studentID string) error {
	return r.db.Where("classroom_id = ? AND student_id = ?", classroomID, studentID).Delete(&domain.StudentClassroom{}).Error
}

func (r *classroomRepository) IsStudentInClass(classroomID string, studentID string) (bool, error) {
	var count int64
	err := r.db.Model(&domain.StudentClassroom{}).Where("classroom_id = ? AND student_id = ?", classroomID, studentID).Count(&count).Error
	return count > 0, err
}
