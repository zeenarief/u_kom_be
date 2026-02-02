package database

import (
	"fmt"
	"log"
	"u_kom_be/internal/model/domain"

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
		// ===== Users =====
		{Name: "users.read", Description: "Read all users"},
		{Name: "users.create", Description: "Create new users"},
		{Name: "users.update", Description: "Update users"},
		{Name: "users.update.others", Description: "Update other users"},
		{Name: "users.change_password.others", Description: "Change password other users"},
		{Name: "users.delete", Description: "Delete users"},
		{Name: "users.manage_roles", Description: "Manage user roles"},
		{Name: "users.manage_permissions", Description: "Manage user permissions"},

		// ===== Roles & Permissions =====
		{Name: "roles.manage", Description: "Manage roles"},
		{Name: "permissions.manage", Description: "Manage permissions"},

		// ===== Profile & Auth =====
		{Name: "profile.read", Description: "Read own profile"},
		{Name: "profile.update", Description: "Update own profile"},
		{Name: "auth.logout", Description: "Logout from system"},

		// ===== Students =====
		{Name: "students.create", Description: "Create new student"},
		{Name: "students.read", Description: "Read students data"},
		{Name: "students.update", Description: "Update student data"},
		{Name: "students.delete", Description: "Delete student"},

		{Name: "students.manage_parents", Description: "Manage student parents relationship"},
		{Name: "students.manage_guardian", Description: "Set or remove student guardian"},
		{Name: "students.manage_account", Description: "Link or unlink student user account"},

		// ===== Parents =====
		{Name: "parents.create", Description: "Create new parent"},
		{Name: "parents.read", Description: "Read parents data"},
		{Name: "parents.update", Description: "Update parent data"},
		{Name: "parents.delete", Description: "Delete parent"},
		{Name: "parents.manage_account", Description: "Link or unlink parent user account"},

		// ===== Guardians =====
		{Name: "guardians.create", Description: "Create new guardian"},
		{Name: "guardians.read", Description: "Read guardians data"},
		{Name: "guardians.update", Description: "Update guardian data"},
		{Name: "guardians.delete", Description: "Delete guardian"},
		{Name: "guardians.manage_account", Description: "Link or unlink guardian user account"},

		// ===== Employees =====
		{Name: "employees.create", Description: "Create new employee"},
		{Name: "employees.read", Description: "Read employees data"},
		{Name: "employees.update", Description: "Update employee data"},
		{Name: "employees.delete", Description: "Delete employee"},
		{Name: "employees.manage_account", Description: "Link or unlink employee user account"},

		// ===== Academic Year =====
		{Name: "academic_years.manage", Description: "Manage academic years data"},
		{Name: "classrooms.manage", Description: "Manage classroom data"},
		{Name: "classrooms.manage_students", Description: "Manage classroom stuident data"},
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
