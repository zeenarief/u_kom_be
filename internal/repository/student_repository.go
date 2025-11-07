package repository

import (
	"belajar-golang/internal/model/domain"
	"errors"

	"gorm.io/gorm"
)

type StudentRepository interface {
	Create(student *domain.Student) error
	FindByID(id string) (*domain.Student, error)
	FindByNISN(nisn string) (*domain.Student, error)
	FindByNIM(nim string) (*domain.Student, error)
	FindAll() ([]domain.Student, error)
	Update(student *domain.Student) error
	Delete(id string) error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(student *domain.Student) error {
	return r.db.Create(student).Error
}

func (r *studentRepository) FindByID(id string) (*domain.Student, error) {
	var student domain.Student
	// Belum ada relasi, jadi tidak perlu .Preload()
	err := r.db.First(&student, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Data tidak ditemukan, return nil tanpa error
	}
	if err != nil {
		return nil, err // Error GORM lainnya
	}
	return &student, nil
}

func (r *studentRepository) FindByNISN(nisn string) (*domain.Student, error) {
	var student domain.Student
	err := r.db.First(&student, "nisn = ?", nisn).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) FindByNIM(nim string) (*domain.Student, error) {
	var student domain.Student
	err := r.db.First(&student, "nim = ?", nim).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) FindAll() ([]domain.Student, error) {
	var students []domain.Student
	err := r.db.Find(&students).Error
	return students, err
}

func (r *studentRepository) Update(student *domain.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id string) error {
	return r.db.Delete(&domain.Student{}, "id = ?", id).Error
}
