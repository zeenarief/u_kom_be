package database

import (
	"fmt"
	"log"
	"u_kom_be/internal/config"
	"u_kom_be/internal/model/domain"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// RunSQLMigrations menjalankan file migrasi SQL dari folder /migrations
func RunSQLMigrations(cfg *config.Config) error {
	// Format DSN untuk golang-migrate sedikit berbeda
	// mysql://user:password@tcp(host:port)/dbname?multiStatements=true
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Tentukan lokasi folder migrasi
	// Kita akan copy folder ini ke dalam Docker image
	migrationPath := "file://migrations"

	log.Println("Connecting to migration source and database...")
	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	log.Println("Running SQL migrations (UP)...")
	// Jalankan migrasi 'up'
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations up: %w", err)
	}

	log.Println("SQL migrations completed successfully.")
	return nil
}

// AutoMigrate melakukan migrasi otomatis untuk semua model
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// HANYA untuk tabel yang belum ada di SQL migration
	// atau untuk development saja
	err := db.AutoMigrate(
	// &domain.User{},  // Comment jika sudah pakai SQL migration
	// &domain.Role{},  // Comment jika sudah pakai SQL migration
	// &domain.Permission{}, // Comment jika sudah pakai SQL migration
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Atau buat function yang hanya untuk development
func AutoMigrateForDev(db *gorm.DB) error {
	log.Println("Running DEV database migrations...")

	err := db.AutoMigrate(
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		// &domain.AcademicYear{},
		// &domain.AttendanceSession{},
		// &domain.AttendanceDetail{},
		// &domain.Classroom{},

		// &domain.Guardian{},
		// &domain.Parent{},
		// &domain.Schedule{},
		&domain.Student{},
		&domain.StudentParent{},
		// &domain.Subject{},
		// &domain.TeachingAssignment{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("DEV Database migrations completed successfully")
	return nil
}
