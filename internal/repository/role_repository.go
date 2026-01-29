package repository

import (
	"errors"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *domain.Role) error
	FindByID(id string) (*domain.Role, error)
	FindByName(name string) (*domain.Role, error)
	FindAll() ([]domain.Role, error)
	Update(role *domain.Role) error
	Delete(id string) error
	SyncPermissions(roleID string, permissionIDs []string) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *domain.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) FindByID(id string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindByName(name string) (*domain.Role, error) {
	var role domain.Role
	err := r.db.Preload("Permissions").First(&role, "name = ?", name).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) FindAll() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Update(role *domain.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id string) error {
	return r.db.Delete(&domain.Role{}, "id = ?", id).Error
}

func (r *roleRepository) SyncPermissions(roleID string, permissionIDs []string) error {
	// Hapus semua permissions role
	if err := r.db.Exec("DELETE FROM role_permission WHERE role_id = ?", roleID).Error; err != nil {
		return err
	}

	// Tambahkan permissions baru
	for _, permissionID := range permissionIDs {
		if err := r.db.Exec("INSERT INTO role_permission (role_id, permission_id) VALUES (?, ?)", roleID, permissionID).Error; err != nil {
			return err
		}
	}

	return nil
}
