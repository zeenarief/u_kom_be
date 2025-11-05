package database

import (
	"belajar-golang/internal/model/domain"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// SeedData menjalankan seeding data awal
func SeedData(db *gorm.DB) error {
	log.Println("Seeding initial data...")

	// Mulai transaction
	return db.Transaction(func(tx *gorm.DB) error {
		// Seed permissions - SESUAIKAN dengan SQL seed
		if err := seedPermissions(tx); err != nil {
			return err
		}

		// Seed roles - SESUAIKAN dengan SQL seed
		if err := seedRoles(tx); err != nil {
			return err
		}

		// Seed admin user - SESUAIKAN dengan SQL seed
		if err := seedAdminUser(tx); err != nil {
			return err
		}

		log.Println("Database seeding completed successfully")
		return nil
	})
}

func seedPermissions(db *gorm.DB) error {
	// SESUAIKAN dengan permissions di SQL seed
	permissions := []domain.Permission{
		{Name: "users.read", Description: "Read all users"},
		{Name: "users.create", Description: "Create new users"},
		{Name: "users.update", Description: "Update users"},
		{Name: "users.update.others", Description: "Update other users"},
		{Name: "users.change_password.others", Description: "Change password other users"},
		{Name: "users.delete", Description: "Delete users"},
		{Name: "users.manage_roles", Description: "Manage user roles"},
		{Name: "users.manage_permissions", Description: "Manage user permissions"},
		{Name: "roles.manage", Description: "Manage roles"},
		{Name: "permissions.manage", Description: "Manage permissions"},
		{Name: "profile.read", Description: "Read own profile"},
		{Name: "profile.update", Description: "Update own profile"},
		{Name: "auth.logout", Description: "Logout from system"},
	}

	for _, permission := range permissions {
		var existing domain.Permission
		if err := db.Where("name = ?", permission.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					return fmt.Errorf("failed to create permission %s: %w", permission.Name, err)
				}
				log.Printf("Created permission: %s", permission.Name)
			} else {
				return err
			}
		}
	}

	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []domain.Role{
		{
			Name:        "admin",
			Description: "Administrator role",
			IsDefault:   false,
		},
		{
			Name:        "user",
			Description: "Regular user role",
			IsDefault:   true,
		},
	}

	for i, role := range roles {
		var existing domain.Role
		if err := db.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Untuk role admin, assign semua permissions
				if role.Name == "admin" {
					var allPermissions []domain.Permission
					if err := db.Find(&allPermissions).Error; err != nil {
						return err
					}
					roles[i].Permissions = allPermissions
				}

				if err := db.Create(&roles[i]).Error; err != nil {
					return fmt.Errorf("failed to create role %s: %w", role.Name, err)
				}
				log.Printf("Created role: %s", role.Name)
			} else {
				return err
			}
		}
	}

	return nil
}

func seedAdminUser(db *gorm.DB) error {
	// Cari role admin
	var adminRole domain.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("failed to find admin role: %w", err)
	}

	adminUser := domain.User{
		Username: "admin",
		Name:     "Super Admin",
		Email:    "admin@example.com",
		Password: "$2a$10$Y4ZQaUO.VTUMoYJJSU3VYe2UIRfDg./SqdbQ71E8gm2CHavcUMx42", // password dari SQL seed
		Roles:    []domain.Role{adminRole},
	}

	var existingUser domain.User
	if err := db.Where("email = ?", adminUser.Email).First(&existingUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&adminUser).Error; err != nil {
				return fmt.Errorf("failed to create admin user: %w", err)
			}
			log.Printf("Created admin user: %s", adminUser.Email)
		} else {
			return err
		}
	}

	return nil
}
