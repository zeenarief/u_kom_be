package repository

import (
	"errors"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type GuardianRepository interface {
	Create(guardian *domain.Guardian) error
	FindByID(id string) (*domain.Guardian, error)
	FindByPhone(phone string) (*domain.Guardian, error)
	FindByEmail(email string) (*domain.Guardian, error)
	FindAll(search string) ([]domain.Guardian, error)
	Update(guardian *domain.Guardian) error
	Delete(id string) error
	SetUserID(guardianID string, userID *string) error
	FindByNIKHash(hash string) (*domain.Guardian, error)
}

type guardianRepository struct {
	db *gorm.DB
}

func NewGuardianRepository(db *gorm.DB) GuardianRepository {
	return &guardianRepository{db: db}
}

func (r *guardianRepository) Create(guardian *domain.Guardian) error {
	return r.db.Create(guardian).Error
}

func (r *guardianRepository) FindByID(id string) (*domain.Guardian, error) {
	var guardian domain.Guardian
	err := r.db.Preload("User").First(&guardian, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Data tidak ditemukan
	}
	if err != nil {
		return nil, err // Error GORM lainnya
	}
	return &guardian, nil
}

func (r *guardianRepository) FindByPhone(phone string) (*domain.Guardian, error) {
	var guardian domain.Guardian
	err := r.db.First(&guardian, "phone_number = ?", phone).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &guardian, nil
}

func (r *guardianRepository) FindByEmail(email string) (*domain.Guardian, error) {
	var guardian domain.Guardian
	err := r.db.First(&guardian, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &guardian, nil
}

func (r *guardianRepository) FindAll(search string) ([]domain.Guardian, error) {
	var guardians []domain.Guardian
	query := r.db

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("full_name LIKE ? OR email LIKE ? OR phone_number LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	err := query.Find(&guardians).Error
	return guardians, err
}

func (r *guardianRepository) Update(guardian *domain.Guardian) error {
	return r.db.Save(guardian).Error
}

func (r *guardianRepository) Delete(id string) error {
	return r.db.Delete(&domain.Guardian{}, "id = ?", id).Error
}

// SetUserID meng-update kolom user_id untuk guardian
func (r *guardianRepository) SetUserID(guardianID string, userID *string) error {
	// GORM akan otomatis meng-set ke NULL jika userID adalah nil
	return r.db.Model(&domain.Guardian{}).Where("id = ?", guardianID).Update("user_id", userID).Error
}

func (r *guardianRepository) FindByNIKHash(hash string) (*domain.Guardian, error) {
	var guardian domain.Guardian
	err := r.db.First(&guardian, "nik_hash = ?", hash).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &guardian, nil
}
