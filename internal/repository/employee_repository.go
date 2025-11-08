package repository

import (
	"belajar-golang/internal/model/domain"
	"errors"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	Create(employee *domain.Employee) error
	FindByID(id string) (*domain.Employee, error)
	FindByNIP(nip string) (*domain.Employee, error)
	FindByPhone(phone string) (*domain.Employee, error)
	FindByUserID(userID string) (*domain.Employee, error)
	FindAll() ([]domain.Employee, error)
	Update(employee *domain.Employee) error
	Delete(id string) error
	SetUserID(employeeID string, userID *string) error // Untuk link/unlink user
}

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) FindByID(id string) (*domain.Employee, error) {
	var employee domain.Employee
	// FindByID tidak melakukan Preload User secara default
	err := r.db.First(&employee, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Data tidak ditemukan
	}
	if err != nil {
		return nil, err // Error GORM lainnya
	}
	return &employee, nil
}

func (r *employeeRepository) FindByNIP(nip string) (*domain.Employee, error) {
	var employee domain.Employee
	err := r.db.First(&employee, "nip = ?", nip).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindByPhone(phone string) (*domain.Employee, error) {
	var employee domain.Employee
	err := r.db.First(&employee, "phone_number = ?", phone).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindByUserID(userID string) (*domain.Employee, error) {
	var employee domain.Employee
	err := r.db.First(&employee, "user_id = ?", userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepository) FindAll() ([]domain.Employee, error) {
	var employees []domain.Employee
	// FindAll tidak melakukan Preload User untuk performa
	err := r.db.Find(&employees).Error
	return employees, err
}

func (r *employeeRepository) Update(employee *domain.Employee) error {
	// .Save akan mengupdate semua field, termasuk yang pointer (NULL atau bernilai)
	return r.db.Save(employee).Error
}

func (r *employeeRepository) Delete(id string) error {
	return r.db.Delete(&domain.Employee{}, "id = ?", id).Error
}

func (r *employeeRepository) SetUserID(employeeID string, userID *string) error {
	// Menggunakan .Update untuk mengubah satu kolom
	// GORM akan otomatis meng-set ke NULL jika userID adalah nil
	return r.db.Model(&domain.Employee{}).Where("id = ?", employeeID).Update("user_id", userID).Error
}
