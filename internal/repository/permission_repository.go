package repository

import (
	"errors"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *domain.Permission) error
	FindByID(id string) (*domain.Permission, error)
	FindByName(name string) (*domain.Permission, error)
	FindAll() ([]domain.Permission, error)
	Update(permission *domain.Permission) error
	Delete(id string) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *domain.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) FindByID(id string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.First(&permission, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindByName(name string) (*domain.Permission, error) {
	var permission domain.Permission
	err := r.db.First(&permission, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) FindAll() ([]domain.Permission, error) {
	var permissions []domain.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) Update(permission *domain.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id string) error {
	return r.db.Delete(&domain.Permission{}, "id = ?", id).Error
}
