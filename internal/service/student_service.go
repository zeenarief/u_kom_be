package service

import (
	"errors"
	"fmt"
	"strings"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"
)

type StudentService interface {
	CreateStudent(req request.StudentCreateRequest) (*response.StudentDetailResponse, error)
	GetStudentByID(id string) (*response.StudentDetailResponse, error)
	GetAllStudents() ([]response.StudentListResponse, error)
	UpdateStudent(id string, req request.StudentUpdateRequest) (*response.StudentDetailResponse, error)
	DeleteStudent(id string) error
	SyncParents(studentID string, req request.StudentSyncParentsRequest) error
	SetGuardian(studentID string, req request.StudentSetGuardianRequest) error
	RemoveGuardian(studentID string) error // Helper untuk menghapus wali
	LinkUser(studentID string, userID string) error
	UnlinkUser(studentID string) error
}

type studentService struct {
	studentRepo    repository.StudentRepository
	parentRepo     repository.ParentRepository
	guardianRepo   repository.GuardianRepository
	userRepo       repository.UserRepository
	encryptionUtil utils.EncryptionUtil                // <-- Untuk ENKRIPSI
	converter      converter.StudentConverterInterface // <-- Untuk DEKRIPSI/Response
}

func NewStudentService(
	studentRepo repository.StudentRepository,
	parentRepo repository.ParentRepository,
	guardianRepo repository.GuardianRepository,
	userRepo repository.UserRepository,
	encryptionUtil utils.EncryptionUtil,
	converter converter.StudentConverterInterface,
) StudentService {
	return &studentService{
		studentRepo:    studentRepo,
		parentRepo:     parentRepo,
		guardianRepo:   guardianRepo,
		userRepo:       userRepo,
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

// GetStudentByID mengambil satu siswa (termasuk M:N Parents dan 1:1 Guardian)
func (s *studentService) GetStudentByID(id string) (*response.StudentDetailResponse, error) {
	// 1. Ambil data student (beserta kolom guardian_id/type) DAN relasi M:N Parents
	student, err := s.studentRepo.FindByIDWithParents(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("student not found")
	}

	// 2. Panggil konverter
	// Konverter akan menangani:
	// - Data dasar student (NIK, Nama, dll.)
	// - Relasi M:N Parents (yang sudah di-preload)
	responseDTO := s.converter.ToStudentDetailResponse(student)

	// Cek apakah student ini punya user_id (terhubung ke akun)
	if student.User.ID != "" {
		responseDTO.User = &response.UserLinkedResponse{
			ID:       student.User.ID,
			Username: student.User.Username,
			Name:     student.User.Name,
			Email:    student.User.Email,
		}
	} else {
		responseDTO.User = nil
	}

	// 3. (LOGIKA BARU) Ambil data Wali Polimorfik secara manual
	// Kita lakukan di service, bukan di converter, karena butuh I/O (repo)
	if student.GuardianID != nil && student.GuardianType != nil {

		guardianInfo, err := s.fetchGuardianInfo(student.GuardianID, student.GuardianType)
		if err != nil {
			// Log error tapi jangan gagalkan request? Tergantung kebutuhan.
			// Untuk sekarang, kita gagalkan jika data wali korup.
			return nil, fmt.Errorf("failed to fetch guardian info: %w", err)
		}
		// Lampirkan data wali ke DTO
		responseDTO.Guardian = guardianInfo
	}

	// 4. Kembalikan DTO yang sudah lengkap
	return responseDTO, nil
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

// SyncParents menangani logika bisnis untuk sinkronisasi orang tua
func (s *studentService) SyncParents(studentID string, req request.StudentSyncParentsRequest) error {
	// 1. Validasi apakah student-nya ada
	student, err := s.studentRepo.FindByID(studentID) // Cukup FindByID, tidak perlu preload
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	var parentRelations []domain.StudentParent
	parentIDMap := make(map[string]bool) // Untuk cek duplikat parent_id di request

	// 2. Validasi setiap parent_id di request
	for _, p := range req.Parents {
		// Cek duplikat di request
		if parentIDMap[p.ParentID] {
			return fmt.Errorf("duplicate parent_id in request: %s", p.ParentID)
		}
		parentIDMap[p.ParentID] = true

		// Cek apakah parent_id ada di database
		parent, err := s.parentRepo.FindByID(p.ParentID)
		if err != nil {
			return fmt.Errorf("error checking parent: %w", err)
		}
		if parent == nil {
			return fmt.Errorf("parent not found with id: %s", p.ParentID)
		}

		// Jika valid, siapkan data untuk repository
		parentRelations = append(parentRelations, domain.StudentParent{
			StudentID:        studentID, // Repository juga akan set ini, tapi lebih baik eksplisit
			ParentID:         p.ParentID,
			RelationshipType: p.RelationshipType,
		})
	}

	// 3. Panggil Repository untuk melakukan sinkronisasi
	return s.studentRepo.SyncParents(studentID, parentRelations)
}

// SetGuardian memvalidasi dan menetapkan wali polimorfik untuk seorang siswa
func (s *studentService) SetGuardian(studentID string, req request.StudentSetGuardianRequest) error {
	// 1. Validasi apakah student-nya ada
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// 2. Validasi apakah guardian_id yang diberikan ada di tabel yang benar
	switch req.GuardianType {
	case "parent":
		parent, err := s.parentRepo.FindByID(req.GuardianID)
		if err != nil {
			return fmt.Errorf("error checking parent: %w", err)
		}
		if parent == nil {
			return fmt.Errorf("parent not found with id: %s", req.GuardianID)
		}
	case "guardian":
		guardian, err := s.guardianRepo.FindByID(req.GuardianID)
		if err != nil {
			return fmt.Errorf("error checking guardian: %w", err)
		}
		if guardian == nil {
			return fmt.Errorf("guardian not found with id: %s", req.GuardianID)
		}
	default:
		// Sebenarnya sudah ditangani oleh validasi 'oneof' di DTO, tapi
		// ini adalah pengaman tambahan.
		return errors.New("invalid guardian_type")
	}

	// 3. Panggil Repository untuk meng-set datanya
	// Kita teruskan pointer ke string dari request
	return s.studentRepo.SetGuardian(studentID, &req.GuardianID, &req.GuardianType)
}

// RemoveGuardian adalah helper untuk menghapus (me-NULL-kan) wali
func (s *studentService) RemoveGuardian(studentID string) error {
	// 1. Validasi apakah student-nya ada
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// 2. Panggil repository dengan nil untuk menghapus
	return s.studentRepo.SetGuardian(studentID, nil, nil)
}

// fetchGuardianInfo adalah helper internal untuk mengambil data wali berdasarkan tipe polimorfiknya.
func (s *studentService) fetchGuardianInfo(guardianID *string, guardianType *string) (*response.GuardianInfoResponse, error) {
	// Cek jika nil (meskipun GetStudentByID sudah cek, ini pengaman)
	if guardianID == nil || guardianType == nil {
		return nil, nil
	}

	id := *guardianID
	tipe := *guardianType

	switch tipe {
	case "parent":
		parent, err := s.parentRepo.FindByID(id)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, fmt.Errorf("data integrity error: parent guardian with id %s not found", id)
		}

		// Petakan domain.Parent ke response.GuardianInfoResponse
		return &response.GuardianInfoResponse{
			ID:           parent.ID,
			FullName:     parent.FullName,
			PhoneNumber:  parent.PhoneNumber,
			Email:        parent.Email,
			Type:         "parent",
			Relationship: "PARENT", // Kita tidak tahu FATHER/MOTHER, jadi 'PARENT'
		}, nil

	case "guardian":
		guardian, err := s.guardianRepo.FindByID(id)
		if err != nil {
			return nil, err
		}
		if guardian == nil {
			return nil, fmt.Errorf("data integrity error: guardian with id %s not found", id)
		}

		// Petakan domain.Guardian ke response.GuardianInfoResponse
		return &response.GuardianInfoResponse{
			ID:           guardian.ID,
			FullName:     guardian.FullName,
			PhoneNumber:  guardian.PhoneNumber,
			Email:        guardian.Email,
			Type:         "guardian",
			Relationship: guardian.RelationshipToStudent, // cth: 'UNCLE', 'AUNT'
		}, nil
	}

	return nil, fmt.Errorf("unknown guardian_type: %s", tipe)
}

// LinkUser menautkan profil Student ke akun User
func (s *studentService) LinkUser(studentID string, userID string) error {
	// 1. Cek apakah Student ada
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// 2. Cek apakah User ada
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// 3. Cek apakah User tsb sudah ditautkan ke Student LAIN
	// (Kita tidak punya FindByUserID di studentRepo, jadi kita tambahkan
	// atau kita cek di service user)
	// Untuk konsistensi, mari kita asumsikan kita perlu menambahkannya:
	// existingStudent, _ := s.studentRepo.FindByUserID(userID)
	// ... (Jika Anda ingin validasi ini, kita harus menambahkannya ke repo)
	// Untuk saat ini, kita andalkan constraint UNIQUE di database

	// 4. Tautkan akun
	if err := s.studentRepo.SetUserID(studentID, &userID); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return errors.New("this user account is already linked to another student")
		}
		return err
	}
	return nil
}

// UnlinkUser menghapus tautan Student dari akun User
func (s *studentService) UnlinkUser(studentID string) error {
	// 1. Cek apakah Student ada
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.studentRepo.SetUserID(studentID, nil)
}
