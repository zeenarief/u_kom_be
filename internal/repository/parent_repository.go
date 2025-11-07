package repository

import (
	"belajar-golang/internal/model/domain"
	"errors"

	"gorm.io/gorm"
)

type ParentRepository interface {
	Create(parent *domain.Parent) error
	FindByID(id string) (*domain.Parent, error)
	FindByPhone(phone string) (*domain.Parent, error)
	FindByEmail(email string) (*domain.Parent, error)
	FindAll() ([]domain.Parent, error)
	Update(parent *domain.Parent) error
	Delete(id string) error
}

type parentRepository struct {
	db *gorm.DB
}

func NewParentRepository(db *gorm.DB) ParentRepository {
	return &parentRepository{db: db}
}

func (r *parentRepository) Create(parent *domain.Parent) error {
	return r.db.Create(parent).Error
}

func (r *parentRepository) FindByID(id string) (*domain.Parent, error) {
	var parent domain.Parent
	// Belum ada relasi, jadi tidak perlu .Preload()
	err := r.db.First(&parent, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Data tidak ditemukan
	}
	if err != nil {
		return nil, err // Error GORM lainnya
	}
	return &parent, nil
}

func (r *parentRepository) FindByPhone(phone string) (*domain.Parent, error) {
	var parent domain.Parent
	err := r.db.First(&parent, "phone_number = ?", phone).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func (r *parentRepository) FindByEmail(email string) (*domain.Parent, error) {
	var parent domain.Parent
	err := r.db.First(&parent, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func (r *parentRepository) FindAll() ([]domain.Parent, error) {
	var parents []domain.Parent
	err := r.db.Find(&parents).Error
	return parents, err
}

func (r *parentRepository) Update(parent *domain.Parent) error {
	return r.db.Save(parent).Error
}

func (r *parentRepository) Delete(id string) error {
	return r.db.Delete(&domain.Parent{}, "id = ?", id).Error
}
