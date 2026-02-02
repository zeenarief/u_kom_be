package repository

import (
	"errors"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type AcademicYearRepository interface {
	Create(academicYear *domain.AcademicYear) error
	FindAll() ([]domain.AcademicYear, error)
	FindByID(id string) (*domain.AcademicYear, error)
	Update(academicYear *domain.AcademicYear) error
	Delete(id string) error

	// Transactional methods for activation logic
	ResetAllStatus(tx *gorm.DB) error
	UpdateStatus(tx *gorm.DB, id string, status string) error
}

type academicYearRepository struct {
	db *gorm.DB
}

func NewAcademicYearRepository(db *gorm.DB) AcademicYearRepository {
	return &academicYearRepository{db: db}
}

func (r *academicYearRepository) Create(academicYear *domain.AcademicYear) error {
	return r.db.Create(academicYear).Error
}

func (r *academicYearRepository) FindAll() ([]domain.AcademicYear, error) {
	var academicYears []domain.AcademicYear
	// Order by start_date descending (terbaru diatas)
	err := r.db.Order("start_date desc").Find(&academicYears).Error
	return academicYears, err
}

func (r *academicYearRepository) FindByID(id string) (*domain.AcademicYear, error) {
	var academicYear domain.AcademicYear
	err := r.db.First(&academicYear, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &academicYear, err
}

func (r *academicYearRepository) Update(academicYear *domain.AcademicYear) error {
	return r.db.Save(academicYear).Error
}

func (r *academicYearRepository) Delete(id string) error {
	return r.db.Delete(&domain.AcademicYear{}, "id = ?", id).Error
}

// ResetAllStatus men-set semua status menjadi INACTIVE
func (r *academicYearRepository) ResetAllStatus(tx *gorm.DB) error {
	// Gunakan tx (transaction DB) jika ada, jika tidak gunakan r.db
	conn := r.db
	if tx != nil {
		conn = tx
	}
	return conn.Model(&domain.AcademicYear{}).Where("1=1").Update("status", "INACTIVE").Error
}

// UpdateStatus mengubah status record spesifik
func (r *academicYearRepository) UpdateStatus(tx *gorm.DB, id string, status string) error {
	conn := r.db
	if tx != nil {
		conn = tx
	}
	return conn.Model(&domain.AcademicYear{}).Where("id = ?", id).Update("status", status).Error
}
