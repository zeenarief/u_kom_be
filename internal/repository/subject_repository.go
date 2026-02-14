package repository

import (
	"errors"
	"smart_school_be/internal/model/domain"

	"gorm.io/gorm"
)

type SubjectRepository interface {
	Create(subject *domain.Subject) error
	FindAll(search string) ([]domain.Subject, error)
	FindByID(id string) (*domain.Subject, error)
	FindByCode(code string) (*domain.Subject, error)
	Update(subject *domain.Subject) error
	Delete(id string) error
}

type subjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepository{db: db}
}

func (r *subjectRepository) Create(subject *domain.Subject) error {
	return r.db.Create(subject).Error
}

func (r *subjectRepository) FindAll(search string) ([]domain.Subject, error) {
	var subjects []domain.Subject
	query := r.db.Order("code asc")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR type LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	err := query.Find(&subjects).Error
	return subjects, err
}

func (r *subjectRepository) FindByID(id string) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.db.First(&subject, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subject, err
}

func (r *subjectRepository) FindByCode(code string) (*domain.Subject, error) {
	var subject domain.Subject
	err := r.db.First(&subject, "code = ?", code).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subject, err
}

func (r *subjectRepository) Update(subject *domain.Subject) error {
	return r.db.Save(subject).Error
}

func (r *subjectRepository) Delete(id string) error {
	return r.db.Delete(&domain.Subject{}, "id = ?", id).Error
}
