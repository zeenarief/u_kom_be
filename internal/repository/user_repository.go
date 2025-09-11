package repository

import (
	"belajar-golang/internal/model/domain"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindByIDWithRelations(id string) (*domain.User, error)
	FindByEmailWithRelations(email string) (*domain.User, error)
	FindByUsernameWithRelations(username string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user *domain.User) error
	Delete(id string) error
	UpdateTokenHash(id string, tokenHash string) error
	GetTokenHash(id string) (string, error)
	GetUserWithRolesAndPermissions(id string) (*domain.User, error)
	SyncRoles(id string, roleIDs []string) error
	SyncPermissions(id string, permissionIDs []string) error
	GetDefaultRole() (*domain.Role, error)
	AssignRole(userID string, roleID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetTokenHash(id string) (string, error) {
	var user domain.User
	err := r.db.Select("current_token_hash").First(&user, "id = ?", id).Error
	if err != nil {
		return "", err
	}
	return user.CurrentTokenHash, nil
}

func (r *userRepository) UpdateTokenHash(id string, tokenHash string) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Update("current_token_hash", tokenHash).Error
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) FindByIDWithRelations(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Permissions").
		First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Return nil, nil jika tidak ditemukan
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByEmailWithRelations(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Permissions").
		First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "username = ?", username).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Return nil, nil jika tidak ditemukan
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByUsernameWithRelations(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Permissions").
		First(&user, "username = ?", username).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User

	err := r.db.
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Permissions").
		Find(&users).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) SyncRoles(userID string, roleIDs []string) error {
	// Hapus semua roles user
	if err := r.db.Exec("DELETE FROM user_role WHERE user_id = ?", userID).Error; err != nil {
		return err
	}

	// Tambahkan roles baru
	for _, roleID := range roleIDs {
		if err := r.db.Exec("INSERT INTO user_role (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) SyncPermissions(userID string, permissionIDs []string) error {
	// Hapus semua permissions langsung user
	if err := r.db.Exec("DELETE FROM user_permission WHERE user_id = ?", userID).Error; err != nil {
		return err
	}

	// Tambahkan permissions baru
	for _, permissionID := range permissionIDs {
		if err := r.db.Exec("INSERT INTO user_permission (user_id, permission_id) VALUES (?, ?)", userID, permissionID).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) GetUserWithRolesAndPermissions(userID string) (*domain.User, error) {
	var user domain.User
	err := r.db.
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Permissions").
		First(&user, "id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetDefaultRole() (*domain.Role, error) {
	var role domain.Role
	err := r.db.First(&role, "is_default = ?", true).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *userRepository) AssignRole(userID string, roleID string) error {
	return r.db.Exec("INSERT INTO user_role (user_id, role_id) VALUES (?, ?)", userID, roleID).Error
}
