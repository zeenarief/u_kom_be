package domain

import (
	"belajar-golang/internal/utils"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID               string       `gorm:"primaryKey;type:char(36)" json:"id"`
	Username         string       `gorm:"uniqueIndex;not null" json:"username"`
	Name             string       `gorm:"not null" json:"name"`
	Email            string       `gorm:"uniqueIndex;not null" json:"email"`
	Password         string       `gorm:"not null" json:"-"`
	CurrentTokenHash string       `gorm:"type:varchar(255)" json:"-"`
	Roles            []Role       `gorm:"many2many:user_role;" json:"roles,omitempty"`
	Permissions      []Permission `gorm:"many2many:user_permission;" json:"permissions,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

// HasRole checks if user has a specific role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}

// HasPermission checks if user has a specific permission (from roles + direct permissions)
func (u *User) HasPermission(permissionName string) bool {
	// Check direct permissions
	for _, perm := range u.Permissions {
		if perm.Name == permissionName {
			return true
		}
	}

	// Check permissions from roles
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if perm.Name == permissionName {
				return true
			}
		}
	}
	return false
}

// GetPermissions returns all permissions for user
func (u *User) GetPermissions() []string {
	permissionsMap := make(map[string]bool)

	// Add permissions from roles
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			permissionsMap[perm.Name] = true
		}
	}

	// Add direct permissions
	for _, perm := range u.Permissions {
		permissionsMap[perm.Name] = true
	}

	// Convert to slice
	var permissions []string
	for perm := range permissionsMap {
		permissions = append(permissions, perm)
	}

	return permissions
}

func (u *User) GetRoles() []string {
	var roles []string
	for _, role := range u.Roles {
		roles = append(roles, role.Name)
	}
	return roles
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = utils.GenerateUUID()
	}
	return
}
