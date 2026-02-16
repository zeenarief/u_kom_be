package repository

import (
	"errors"
	"smart_school_be/internal/model/domain"

	"gorm.io/gorm"
)

type ParentRepository interface {
	Create(parent *domain.Parent) error
	FindByID(id string) (*domain.Parent, error)
	FindByPhone(phone string) (*domain.Parent, error)
	FindByEmail(email string) (*domain.Parent, error)
	FindAll(search string, limit, offset int) ([]domain.Parent, int64, error)
	Update(parent *domain.Parent) error
	Delete(id string) error
	SetUserID(parentID string, userID *string) error
	FindByNIKHash(hash string) (*domain.Parent, error)
	FindByUserID(userID string) (*domain.Parent, error)
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
	err := r.db.Preload("User").First(&parent, "id = ?", id).Error
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

func (r *parentRepository) FindAll(search string, limit, offset int) ([]domain.Parent, int64, error) {
	var parents []domain.Parent
	var total int64
	query := r.db.Model(&domain.Parent{})

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("full_name LIKE ? OR email LIKE ? OR phone_number LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("full_name ASC").
		Limit(limit).Offset(offset).
		Find(&parents).Error
	return parents, total, err
}

func (r *parentRepository) Update(parent *domain.Parent) error {
	return r.db.Save(parent).Error
}

func (r *parentRepository) Delete(id string) error {
	return r.db.Delete(&domain.Parent{}, "id = ?", id).Error
}

// SetUserID meng-update kolom user_id untuk parent
func (r *parentRepository) SetUserID(parentID string, userID *string) error {
	// GORM akan otomatis meng-set ke NULL jika userID adalah nil
	return r.db.Model(&domain.Parent{}).Where("id = ?", parentID).Update("user_id", userID).Error
}

func (r *parentRepository) FindByNIKHash(hash string) (*domain.Parent, error) {
	var parent domain.Parent
	err := r.db.First(&parent, "nik_hash = ?", hash).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}

func (r *parentRepository) FindByUserID(userID string) (*domain.Parent, error) {
	var parent domain.Parent
	err := r.db.First(&parent, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &parent, nil
}
