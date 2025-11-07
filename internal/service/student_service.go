package service

import (
	"belajar-golang/internal/converter"
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/repository"
	"belajar-golang/internal/utils"
	"errors"
	"fmt"
)

type StudentService interface {
	CreateStudent(req request.StudentCreateRequest) (*response.StudentDetailResponse, error)
	GetStudentByID(id string) (*response.StudentDetailResponse, error)
	GetAllStudents() ([]response.StudentListResponse, error)
	UpdateStudent(id string, req request.StudentUpdateRequest) (*response.StudentDetailResponse, error)
	DeleteStudent(id string) error
}

type studentService struct {
	studentRepo    repository.StudentRepository
	encryptionUtil utils.EncryptionUtil                // <-- Untuk ENKRIPSI
	converter      converter.StudentConverterInterface // <-- Untuk DEKRIPSI/Response
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	encryptionUtil utils.EncryptionUtil,
	converter converter.StudentConverterInterface,
) StudentService {
	return &studentService{
		studentRepo:    studentRepo,
		encryptionUtil: encryptionUtil,
		converter:      converter,
	}
}

// CreateStudent menangani pembuatan siswa baru
func (s *studentService) CreateStudent(req request.StudentCreateRequest) (*response.StudentDetailResponse, error) {
	// 1. Validasi Duplikat
	if req.NISN != "" {
		if existing, _ := s.studentRepo.FindByNISN(req.NISN); existing != nil {
			return nil, errors.New("nisn already exists")
		}
	}
	if req.NIM != "" {
		if existing, _ := s.studentRepo.FindByNIM(req.NIM); existing != nil {
			return nil, errors.New("nim already exists")
		}
	}

	// 2. Enkripsi Data Sensitif
	encryptedNIK := ""
	if req.NIK != "" {
		var err error
		encryptedNIK, err = s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
	}
	encryptedNoKK := ""
	if req.NoKK != "" {
		var err error
		encryptedNoKK, err = s.encryptionUtil.Encrypt(req.NoKK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt no_kk: %w", err)
		}
	}

	// 3. Buat Domain Object
	student := &domain.Student{
		FullName:     req.FullName,
		NoKK:         encryptedNoKK,
		NIK:          encryptedNIK,
		NISN:         req.NISN,
		NIM:          req.NIM,
		Gender:       req.Gender,
		PlaceOfBirth: req.PlaceOfBirth,
		DateOfBirth:  req.DateOfBirth,
		Address:      req.Address,
		RT:           req.RT,
		RW:           req.RW,
		SubDistrict:  req.SubDistrict,
		District:     req.District,
		City:         req.City,
		Province:     req.Province,
		PostalCode:   req.PostalCode,
	}

	// 4. Panggil Repository
	if err := s.studentRepo.Create(student); err != nil {
		return nil, err
	}

	// 5. Ambil data yang baru dibuat
	createdStudent, err := s.studentRepo.FindByID(student.ID)
	if err != nil {
		return nil, err
	}
	if createdStudent == nil {
		return nil, errors.New("failed to retrieve created student")
	}

	// 6. Konversi ke Response (menggunakan konverter)
	return s.converter.ToStudentDetailResponse(createdStudent), nil
}

// GetStudentByID mengambil satu siswa
func (s *studentService) GetStudentByID(id string) (*response.StudentDetailResponse, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("student not found")
	}
	// Panggil konverter
	return s.converter.ToStudentDetailResponse(student), nil
}

// GetAllStudents mengambil semua siswa
func (s *studentService) GetAllStudents() ([]response.StudentListResponse, error) {
	students, err := s.studentRepo.FindAll()
	if err != nil {
		return nil, err
	}
	// Panggil konverter list
	return s.converter.ToStudentListResponses(students), nil
}

// UpdateStudent memperbarui data siswa
func (s *studentService) UpdateStudent(id string, req request.StudentUpdateRequest) (*response.StudentDetailResponse, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("student not found")
	}

	// Update fields jika disediakan (meniru RoleService)
	if req.FullName != "" {
		student.FullName = req.FullName
	}

	// Validasi duplikat baru
	if req.NISN != "" && req.NISN != student.NISN {
		if existing, _ := s.studentRepo.FindByNISN(req.NISN); existing != nil {
			return nil, errors.New("nisn already exists")
		}
		student.NISN = req.NISN
	}
	if req.NIM != "" && req.NIM != student.NIM {
		if existing, _ := s.studentRepo.FindByNIM(req.NIM); existing != nil {
			return nil, errors.New("nim already exists")
		}
		student.NIM = req.NIM
	}

	// Enkripsi field yang diperbarui
	if req.NIK != "" {
		encryptedNIK, err := s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt nik: %w", err)
		}
		student.NIK = encryptedNIK
	}
	if req.NoKK != "" {
		encryptedNoKK, err := s.encryptionUtil.Encrypt(req.NoKK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt no_kk: %w", err)
		}
		student.NoKK = encryptedNoKK
	}

	// Update field lainnya
	if req.Gender != "" {
		student.Gender = req.Gender
	}
	if req.PlaceOfBirth != "" {
		student.PlaceOfBirth = req.PlaceOfBirth
	}
	// Untuk time.Time, kita cek IsZero()
	if !req.DateOfBirth.IsZero() {
		student.DateOfBirth = req.DateOfBirth
	}
	if req.Address != "" {
		student.Address = req.Address
	}
	// ... (lakukan hal yang sama untuk RT, RW, SubDistrict, Dll.) ...
	if req.RT != "" {
		student.RT = req.RT
	}
	if req.RW != "" {
		student.RW = req.RW
	}
	if req.SubDistrict != "" {
		student.SubDistrict = req.SubDistrict
	}
	if req.District != "" {
		student.District = req.District
	}
	if req.City != "" {
		student.City = req.City
	}
	if req.Province != "" {
		student.Province = req.Province
	}
	if req.PostalCode != "" {
		student.PostalCode = req.PostalCode
	}

	if err := s.studentRepo.Update(student); err != nil {
		return nil, err
	}

	// Ambil data yang sudah diupdate
	updatedStudent, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.converter.ToStudentDetailResponse(updatedStudent), nil
}

// DeleteStudent menghapus siswa
func (s *studentService) DeleteStudent(id string) error {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// Opsional: Tambahkan logika bisnis
	// if student.SomeCondition {
	//    return errors.New("cannot delete this student")
	// }

	return s.studentRepo.Delete(id)
}
