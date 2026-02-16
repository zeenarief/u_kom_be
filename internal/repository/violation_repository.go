package repository

import (
	"errors"
	"smart_school_be/internal/model/domain"

	"gorm.io/gorm"
)

type ViolationRepository interface {
	// Category
	CreateCategory(category *domain.ViolationCategory) error
	FindAllCategories() ([]domain.ViolationCategory, error)
	FindCategoryByID(id string) (*domain.ViolationCategory, error)
	UpdateCategory(category *domain.ViolationCategory) error
	DeleteCategory(id string) error

	// Type
	CreateType(violationType *domain.ViolationType) error
	FindAllTypes(categoryID string) ([]domain.ViolationType, error)
	FindTypeByID(id string) (*domain.ViolationType, error)
	UpdateType(violationType *domain.ViolationType) error
	DeleteType(id string) error

	// Student Violation
	RecordViolation(violation *domain.StudentViolation) error
	FindStudentViolations(studentID string, limit, offset int) ([]domain.StudentViolation, int64, error)
	FindViolationByID(id string) (*domain.StudentViolation, error)
	DeleteViolation(id string) error
	FindAllViolations(filter string, limit, offset int) ([]domain.StudentViolation, int64, error)
}

type violationRepository struct {
	db *gorm.DB
}

func NewViolationRepository(db *gorm.DB) ViolationRepository {
	return &violationRepository{db: db}
}

// Category Implementation
func (r *violationRepository) CreateCategory(category *domain.ViolationCategory) error {
	return r.db.Create(category).Error
}

func (r *violationRepository) FindAllCategories() ([]domain.ViolationCategory, error) {
	var categories []domain.ViolationCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *violationRepository) FindCategoryByID(id string) (*domain.ViolationCategory, error) {
	var category domain.ViolationCategory
	err := r.db.First(&category, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *violationRepository) UpdateCategory(category *domain.ViolationCategory) error {
	return r.db.Save(category).Error
}

func (r *violationRepository) DeleteCategory(id string) error {
	return r.db.Delete(&domain.ViolationCategory{}, "id = ?", id).Error
}

// Type Implementation
func (r *violationRepository) CreateType(violationType *domain.ViolationType) error {
	return r.db.Create(violationType).Error
}

func (r *violationRepository) FindAllTypes(categoryID string) ([]domain.ViolationType, error) {
	var types []domain.ViolationType
	query := r.db.Preload("Category")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	err := query.Find(&types).Error
	return types, err
}

func (r *violationRepository) FindTypeByID(id string) (*domain.ViolationType, error) {
	var violationType domain.ViolationType
	err := r.db.Preload("Category").First(&violationType, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &violationType, nil
}

func (r *violationRepository) UpdateType(violationType *domain.ViolationType) error {
	return r.db.Save(violationType).Error
}

func (r *violationRepository) DeleteType(id string) error {
	return r.db.Delete(&domain.ViolationType{}, "id = ?", id).Error
}

// Student Violation Implementation
func (r *violationRepository) RecordViolation(violation *domain.StudentViolation) error {
	return r.db.Create(violation).Error
}

func (r *violationRepository) FindStudentViolations(studentID string, limit, offset int) ([]domain.StudentViolation, int64, error) {
	var violations []domain.StudentViolation
	var total int64
	query := r.db.Model(&domain.StudentViolation{}).Where("student_id = ?", studentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("ViolationType").
		Preload("ViolationType.Category").
		Preload("Student").
		Order("violation_date DESC").
		Limit(limit).Offset(offset).
		Find(&violations).Error
	return violations, total, err
}

func (r *violationRepository) FindViolationByID(id string) (*domain.StudentViolation, error) {
	var violation domain.StudentViolation
	err := r.db.Preload("ViolationType").
		Preload("ViolationType.Category").
		Preload("Student").
		First(&violation, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &violation, nil
}

func (r *violationRepository) DeleteViolation(id string) error {
	return r.db.Delete(&domain.StudentViolation{}, "id = ?", id).Error
}
func (r *violationRepository) FindAllViolations(filter string, limit, offset int) ([]domain.StudentViolation, int64, error) {
	var violations []domain.StudentViolation
	var total int64
	query := r.db.Model(&domain.StudentViolation{})

	if filter != "" {
		searchPattern := "%" + filter + "%"
		// Join to filter by student name or violation name
		query = query.Joins("JOIN students s ON s.id = student_violations.student_id").
			Joins("JOIN violation_types vt ON vt.id = student_violations.violation_type_id").
			Where("s.full_name LIKE ? OR vt.name LIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("ViolationType").
		Preload("ViolationType.Category").
		Preload("Student").
		Order("violation_date DESC").
		Limit(limit).Offset(offset).
		Find(&violations).Error
	return violations, total, err
}
