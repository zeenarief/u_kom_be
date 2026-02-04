package repository

import (
	"errors"
	"u_kom_be/internal/model/domain"

	"gorm.io/gorm"
)

type StudentRepository interface {
	Create(student *domain.Student) error
	FindByID(id string) (*domain.Student, error)
	FindByNISN(nisn string) (*domain.Student, error)
	FindByNIM(nim string) (*domain.Student, error)
	FindAll(search string) ([]domain.Student, error)
	Update(student *domain.Student) error
	Delete(id string) error
	FindByIDWithParents(id string) (*domain.Student, error)
	SyncParents(studentID string, parents []domain.StudentParent) error
	SetGuardian(studentID string, guardianID *string, guardianType *string) error
	SetUserID(studentID string, userID *string) error
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

func (r *studentRepository) FindAll(search string) ([]domain.Student, error) {
	var students []domain.Student
	query := r.db

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("full_name LIKE ? OR nisn LIKE ? OR city LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	err := query.Find(&students).Error
	return students, err
}

func (r *studentRepository) Update(student *domain.Student) error {
	return r.db.Save(student).Error
}

func (r *studentRepository) Delete(id string) error {
	return r.db.Delete(&domain.Student{}, "id = ?", id).Error
}

// FindByIDWithParents mengambil Student beserta relasi Parents dan data Parent-nya
func (r *studentRepository) FindByIDWithParents(id string) (*domain.Student, error) {
	var student domain.Student
	err := r.db.
		Preload("User").
		Preload("Parents").        // 1. Ambil data dari pivot table (student_parent)
		Preload("Parents.Parent"). // 2. Untuk setiap data pivot, ambil data dari tabel parents
		First(&student, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// SyncParents menghapus semua relasi lama dan membuat yang baru dalam satu transaksi
func (r *studentRepository) SyncParents(studentID string, parents []domain.StudentParent) error {
	// Gunakan Transaksi agar atomik (semua berhasil atau semua gagal)
	return r.db.Transaction(func(tx *gorm.DB) error {

		// 1. Hapus semua relasi parent yang ada untuk student ini
		if err := tx.Where("student_id = ?", studentID).Delete(&domain.StudentParent{}).Error; err != nil {
			return err
		}

		// 2. Tambahkan relasi baru (jika list tidak kosong)
		if len(parents) == 0 {
			return nil // Tidak ada yang perlu ditambahkan, selesai.
		}

		// Pastikan StudentID ter-set untuk setiap entri (meski service harusnya sudah)
		for i := range parents {
			parents[i].StudentID = studentID
		}

		// 3. Buat (batch insert) relasi baru
		if err := tx.Create(&parents).Error; err != nil {
			return err
		}

		return nil
	})
}

// SetGuardian meng-update penanda wali (polymorphic) pada tabel student
func (r *studentRepository) SetGuardian(studentID string, guardianID *string, guardianType *string) error {
	// Jika nil, GORM akan meng-set kolom ke NULL
	// Jika tidak nil, GORM akan meng-set ke nilainya
	return r.db.Model(&domain.Student{}).Where("id = ?", studentID).Updates(map[string]interface{}{
		"guardian_id":   guardianID,
		"guardian_type": guardianType,
	}).Error
}

// SetUserID meng-update kolom user_id untuk student
func (r *studentRepository) SetUserID(studentID string, userID *string) error {
	// GORM akan otomatis meng-set ke NULL jika userID adalah nil
	return r.db.Model(&domain.Student{}).Where("id = ?", studentID).Update("user_id", userID).Error
}
