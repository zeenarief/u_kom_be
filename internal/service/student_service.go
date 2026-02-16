package service

import (
	"bytes"
	"errors"
	"fmt"
	"smart_school_be/internal/apperrors"
	"smart_school_be/internal/converter"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
	"smart_school_be/internal/utils"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
	"github.com/xuri/excelize/v2"
)

type StudentService interface {
	CreateStudent(req request.StudentCreateRequest, files StudentFiles) (*response.StudentDetailResponse, error)
	GetStudentByID(id string) (*response.StudentDetailResponse, error)
	GetAllStudents(search string, classroomID string, pagination request.PaginationRequest) (*response.PaginatedData, error)
	UpdateStudent(id string, req request.StudentUpdateRequest, files StudentFiles) (*response.StudentDetailResponse, error)
	DeleteStudent(id string) error
	SyncParents(studentID string, req request.StudentSyncParentsRequest) error
	SetGuardian(studentID string, req request.StudentSetGuardianRequest) error
	RemoveGuardian(studentID string) error // Helper untuk menghapus wali
	LinkUser(studentID string, userID string) error
	UnlinkUser(studentID string) error
	ExportStudentsToExcel() (*bytes.Buffer, error)
	ExportStudentsToPdf() (*bytes.Buffer, error)
	ExportStudentBiodata(id string) (*bytes.Buffer, error)
}

type StudentFiles struct {
	BirthCertificateFile        string
	FamilyCardFile              string
	ParentStatementFile         string
	StudentStatementFile        string
	HealthInsuranceFile         string
	DiplomaCertificateFile      string
	GraduationCertificateFile   string
	FinancialHardshipLetterFile string
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
func (s *studentService) CreateStudent(req request.StudentCreateRequest, files StudentFiles) (*response.StudentDetailResponse, error) {
	// Helpers untuk konversi string kosong ke nil pointer
	toPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}
	toDatePtr := func(d utils.Date) *utils.Date {
		if d.IsZero() {
			return nil
		}
		return &d
	}

	// 1. Validasi Duplikat
	var nisn *string
	if req.NISN != "" {
		nisn = &req.NISN
		if existing, _ := s.studentRepo.FindByNISN(req.NISN); existing != nil {
			return nil, apperrors.NewConflictError("NISN already exists")
		}
	}

	var nim *string
	if req.NIM != "" {
		nim = &req.NIM
		if existing, _ := s.studentRepo.FindByNIM(req.NIM); existing != nil {
			return nil, apperrors.NewConflictError("NIM already exists")
		}
	}

	// 2. Enkripsi Data Sensitif & Hash NIK
	var encryptedNIK *string
	var nikHash *string

	if req.NIK != "" {
		// a. Hash & Check Unique
		hash, err := s.encryptionUtil.Hash(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to hash NIK: %w", err)
		}

		existing, err := s.studentRepo.FindByNIKHash(hash)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, apperrors.NewConflictError("NIK already exists")
		}
		nikHash = &hash

		// b. Encrypt
		encrypted, err := s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt NIK: %w", err)
		}
		encryptedNIK = &encrypted
	}

	encryptedNoKK := ""
	if req.NoKK != "" {
		var err error
		encryptedNoKK, err = s.encryptionUtil.Encrypt(req.NoKK)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt NoKK: %w", err)
		}
	}

	// 3. Buat Domain Object
	student := &domain.Student{
		FullName:                    req.FullName,
		NoKK:                        encryptedNoKK,
		NIK:                         encryptedNIK,
		NIKHash:                     nikHash,
		NISN:                        nisn,
		NIM:                         nim,
		Gender:                      req.Gender,
		PlaceOfBirth:                toPtr(req.PlaceOfBirth),
		DateOfBirth:                 toDatePtr(req.DateOfBirth),
		Address:                     toPtr(req.Address),
		RT:                          toPtr(req.RT),
		RW:                          toPtr(req.RW),
		SubDistrict:                 toPtr(req.SubDistrict),
		District:                    toPtr(req.District),
		City:                        toPtr(req.City),
		Province:                    toPtr(req.Province),
		PostalCode:                  toPtr(req.PostalCode),
		Status:                      req.Status,
		EntryYear:                   toPtr(req.EntryYear),
		ExitYear:                    toPtr(req.ExitYear),
		BirthCertificateFile:        toPtr(files.BirthCertificateFile),
		FamilyCardFile:              toPtr(files.FamilyCardFile),
		ParentStatementFile:         toPtr(files.ParentStatementFile),
		StudentStatementFile:        toPtr(files.StudentStatementFile),
		HealthInsuranceFile:         toPtr(files.HealthInsuranceFile),
		DiplomaCertificateFile:      toPtr(files.DiplomaCertificateFile),
		GraduationCertificateFile:   toPtr(files.GraduationCertificateFile),
		FinancialHardshipLetterFile: toPtr(files.FinancialHardshipLetterFile),
	}

	// Set default status if empty
	if student.Status == "" {
		student.Status = "ACTIVE"
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
		return nil, apperrors.NewInternalError("Failed to retrieve created student")
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
		return nil, apperrors.NewNotFoundError("Student not found")
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
			return nil, fmt.Errorf("Failed to fetch guardian info: %w", err)
		}
		// Lampirkan data wali ke DTO
		responseDTO.Guardian = guardianInfo
	}

	// 4. Kembalikan DTO yang sudah lengkap
	return responseDTO, nil
}

// GetAllStudents mengambil semua siswa dengan pagination
func (s *studentService) GetAllStudents(search string, classroomID string, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	students, total, err := s.studentRepo.FindAll(search, classroomID, limit, offset)
	if err != nil {
		return nil, err
	}
	// Panggil konverter list
	data := s.converter.ToStudentListResponses(students)
	paginatedData := response.NewPaginatedData(data, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

// UpdateStudent memperbarui data siswa
func (s *studentService) UpdateStudent(id string, req request.StudentUpdateRequest, files StudentFiles) (*response.StudentDetailResponse, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, apperrors.NewNotFoundError("Student not found")
	}

	// Update fields jika disediakan (meniru RoleService)
	if req.FullName != "" {
		student.FullName = req.FullName
	}

	// Validasi duplikat baru
	if req.NISN == "" {
		// request eksplisit ingin mengosongkan
		student.NISN = nil
	} else {
		// ada value baru
		if student.NISN == nil || *student.NISN != req.NISN {
			if existing, _ := s.studentRepo.FindByNISN(req.NISN); existing != nil {
				return nil, apperrors.NewConflictError("NISN already exists")
			}
			student.NISN = &req.NISN
		}
	}
	if req.NIM == "" {
		student.NIM = nil
	} else {
		if student.NIM == nil || *student.NIM != req.NIM {
			if existing, _ := s.studentRepo.FindByNIM(req.NIM); existing != nil {
				return nil, apperrors.NewConflictError("NIM already exists")
			}
			student.NIM = &req.NIM
		}
	}

	// Enkripsi field yang diperbarui
	if req.NIK != nil {
		// Logika update NIK:
		// Jika empty string ("") -> Clear NIK
		// Jika ada isi -> Update (Hash+Encrypt)
		if *req.NIK == "" {
			student.NIK = nil
			student.NIKHash = nil
		} else {
			// Hash & Check
			newHash, err := s.encryptionUtil.Hash(*req.NIK)
			if err != nil {
				return nil, fmt.Errorf("Failed to hash NIK: %w", err)
			}

			// Cek keunikan jika hash berubah atau sebelumnya null
			isDifferent := student.NIKHash == nil || *student.NIKHash != newHash
			if isDifferent {
				existing, err := s.studentRepo.FindByNIKHash(newHash)
				if err != nil {
					return nil, err
				}
				if existing != nil && existing.ID != id {
					return nil, apperrors.NewConflictError("NIK already exists")
				}
			}

			student.NIKHash = &newHash

			// Enkripsi
			encryptedNIK, err := s.encryptionUtil.Encrypt(*req.NIK)
			if err != nil {
				return nil, fmt.Errorf("Failed to encrypt NIK: %w", err)
			}
			student.NIK = &encryptedNIK
		}
	}
	// NoKK - bisa di-null dengan mengirim empty string
	if req.NoKK == "" {
		student.NoKK = "" // Set ke empty untuk null
	} else {
		encryptedNoKK, err := s.encryptionUtil.Encrypt(req.NoKK)
		if err != nil {
			return nil, fmt.Errorf("Failed to encrypt NoKK: %w", err)
		}
		student.NoKK = encryptedNoKK
	}

	// Update field lainnya - direct assign untuk pointer fields
	student.PlaceOfBirth = req.PlaceOfBirth
	student.DateOfBirth = req.DateOfBirth
	student.Address = req.Address
	student.RT = req.RT
	student.RW = req.RW
	student.SubDistrict = req.SubDistrict
	student.District = req.District
	student.City = req.City
	student.Province = req.Province
	student.PostalCode = req.PostalCode
	student.EntryYear = req.EntryYear
	student.ExitYear = req.ExitYear

	// Non-pointer fields need dereference
	if req.Gender != nil {
		student.Gender = *req.Gender
	}
	if req.Status != nil {
		student.Status = *req.Status
	}

	// Update File Paths
	// Helper to update file path if provided
	updateFile := func(currentPath **string, newPath string) {
		if newPath != "" {
			if *currentPath != nil {
				utils.RemoveFile(**currentPath)
			}
			*currentPath = &newPath
		}
	}

	updateFile(&student.BirthCertificateFile, files.BirthCertificateFile)
	updateFile(&student.FamilyCardFile, files.FamilyCardFile)
	updateFile(&student.ParentStatementFile, files.ParentStatementFile)
	updateFile(&student.StudentStatementFile, files.StudentStatementFile)
	updateFile(&student.HealthInsuranceFile, files.HealthInsuranceFile)
	updateFile(&student.DiplomaCertificateFile, files.DiplomaCertificateFile)
	updateFile(&student.GraduationCertificateFile, files.GraduationCertificateFile)
	updateFile(&student.FinancialHardshipLetterFile, files.FinancialHardshipLetterFile)

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
		return apperrors.NewNotFoundError("student not found")
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
			return apperrors.NewBadRequestError(fmt.Sprintf("duplicate parent_id in request: %s", p.ParentID))
		}
		parentIDMap[p.ParentID] = true

		// Cek apakah parent_id ada di database
		parent, err := s.parentRepo.FindByID(p.ParentID)
		if err != nil {
			return fmt.Errorf("error checking parent: %w", err)
		}
		if parent == nil {
			return apperrors.NewNotFoundError(fmt.Sprintf("Parent not found with id: %s", p.ParentID))
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
			return apperrors.NewNotFoundError(fmt.Sprintf("Parent not found with id: %s", req.GuardianID))
		}
	case "guardian":
		guardian, err := s.guardianRepo.FindByID(req.GuardianID)
		if err != nil {
			return fmt.Errorf("error checking guardian: %w", err)
		}
		if guardian == nil {
			return apperrors.NewNotFoundError(fmt.Sprintf("Guardian not found with id: %s", req.GuardianID))
		}
	default:
		// Sebenarnya sudah ditangani oleh validasi 'oneof' di DTO, tapi
		// ini adalah pengaman tambahan.
		return apperrors.NewBadRequestError("invalid guardian_type")
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
			return nil, fmt.Errorf("data integrity error: Parent guardian with id %s not found", id)
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
			return nil, fmt.Errorf("data integrity error: Guardian with id %s not found", id)
		}

		// Petakan domain.Guardian ke response.GuardianInfoResponse
		return &response.GuardianInfoResponse{
			ID:          guardian.ID,
			FullName:    guardian.FullName,
			PhoneNumber: guardian.PhoneNumber,
			Email:       guardian.Email,
			Type:        "guardian",
			// Jika RelationshipToStudent nil, default ke empty string atau "-"
			Relationship: func() string {
				if guardian.RelationshipToStudent != nil {
					return *guardian.RelationshipToStudent
				}
				return ""
			}(),
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
		return apperrors.NewNotFoundError("User not found")
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
			return apperrors.NewConflictError("this user account is already linked to another student")
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
		return errors.New("Student not found")
	}

	// 2. Hapus tautan (set user_id ke NULL)
	return s.studentRepo.SetUserID(studentID, nil)
}

func (s *studentService) ExportStudentsToExcel() (*bytes.Buffer, error) {
	// 1. Ambil semua data siswa
	// 1. Ambil semua data siswa (tanpa limit/offset, atau set limit besar)
	// Gunakan limit besar untuk export
	students, _, err := s.studentRepo.FindAll("", "", 10000, 0)
	if err != nil {
		return nil, err
	}

	// 2. Buat File Excel Baru
	f := excelize.NewFile()
	sheetName := "Data Siswa"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	// Hapus sheet default "Sheet1" jika tidak dipakai
	f.DeleteSheet("Sheet1")

	// 3. Buat Header
	headers := []string{"No", "NISN", "NIM", "Nama Lengkap", "Jenis Kelamin", "Tempat Lahir", "Tanggal Lahir", "Alamat", "Status", "Tahun Masuk", "Tahun Keluar"}
	for i, header := range headers {
		// Konversi koordinat (0,0 -> A1, 1,0 -> B1)
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style Header (Bold, Kuning) - Opsional biar cantik
	style, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#FFFF00"}, Pattern: 1},
	})
	f.SetCellStyle(sheetName, "A1", "L1", style)

	// 4. Isi Data
	for i, student := range students {
		row := i + 2 // Mulai dari baris ke-2

		// Format Tanggal
		dob := ""
		if !student.DateOfBirth.IsZero() {
			dob = student.DateOfBirth.Format("2006-01-02")
		}

		address := utils.JoinAddress(
			student.Address,
			student.RT,
			student.RW,
			student.SubDistrict,
			student.District,
			student.City,
			student.Province,
			student.PostalCode,
		)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), utils.SafeString(student.NISN))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), utils.SafeString(student.NIM))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), student.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), student.Gender)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), utils.SafeString(student.PlaceOfBirth))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), dob)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), address)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), student.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), student.EntryYear)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), student.ExitYear)
	}

	// 5. Simpan ke Buffer (Memory)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func (s *studentService) ExportStudentsToPdf() (*bytes.Buffer, error) {
	// Gunakan limit besar untuk export
	students, _, err := s.studentRepo.FindAll("", "", 10000, 0)
	if err != nil {
		return nil, err
	}

	// 1. Init PDF (Landscape, mm, A4)
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	// 2. Judul
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "LAPORAN DATA SISWA", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// 3. Header Tabel
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240) // Abu-abu muda

	// Lebar Kolom: No, NISN, Nama, Gender, Alamat
	widths := []float64{10, 30, 60, 20, 40, 80}
	headers := []string{"No", "NISN", "Nama Lengkap", "JK", "Tgl Lahir", "Alamat"}

	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1) // Pindah baris

	// 4. Isi Data
	pdf.SetFont("Arial", "", 9)
	pdf.SetFillColor(255, 255, 255) // Putih

	for i, student := range students {
		// Format Tanggal
		dob := "-"
		if !student.DateOfBirth.IsZero() {
			dob = student.DateOfBirth.Format("02-01-2006")
		}

		// Convert Gender
		gender := "L"
		if student.Gender == "female" {
			gender = "P"
		}

		pdf.CellFormat(widths[0], 8, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 8, student.NISNValue(), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[2], 8, student.FullName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, gender, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[4], 8, dob, "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[5], 8, utils.SafeString(student.Address), "1", 0, "L", false, 0, "") // Alamat mungkin terpotong jika panjang, nanti bisa pakai MultiCell
		pdf.Ln(-1)
	}

	// 5. Output ke Buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return &buf, nil
}

func (s *studentService) ExportStudentBiodata(id string) (*bytes.Buffer, error) {
	// 1. Ambil data lengkap (termasuk parents & guardian)
	student, err := s.studentRepo.FindByIDWithParents(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, apperrors.NewNotFoundError("Student not found")
	}

	// 2. Buat PDF instance baru
	pdf := fpdf.New("P", "mm", "A4", "")

	// 3. Tambahkan halaman baru terlebih dahulu
	pdf.AddPage()

	// 4. Import kop surat PDF menggunakan gofpdi (CARA YANG BENAR)
	// Pastikan path file benar. Jika dijalankan dari root project, biasanya "assets/..."
	kopSuratPath := "assets/kop_surat_a4.pdf"

	// ImportPage(pdfInstance, pathFile, halamanKe, boxType)
	// Fungsi ini otomatis me-link template ke object 'pdf' Anda
	tpl := gofpdi.ImportPage(pdf, kopSuratPath, 1, "/MediaBox")

	// Gambar template ke halaman
	// UseImportedTemplate(pdfInstance, tplId, x, y, width, height)
	// x=0, y=0, w=210 (Full A4 Width), h=0 (Auto height/Full Height)
	gofpdi.UseImportedTemplate(pdf, tpl, 0, 0, 210, 0)

	// Format: SetMargins(kiri, atas, kanan)
	// Satuan dalam mm.
	// 15mm = 1.5cm (supaya lebih masuk 0.5cm dari default 1cm)
	// Margin atas diset 10mm (standar), karena kita pakai SetY(45) manual untuk halaman 1.
	pdf.SetMargins(15, 10, 15)

	// 5. Set posisi Y mulai dari 50mm (5cm) agar tidak menimpa kop surat
	pdf.SetY(45)

	pdf.SetFont("Arial", "BU", 14)
	pdf.CellFormat(0, 8, "LEMBAR DATA DIRI SANTRI", "", 1, "C", false, 0, "")

	pdf.Ln(5)

	// Helper function untuk baris data: [Label : Value]
	printRow := func(label, value string) {
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(40, 7, label)                     // Lebar Label
		pdf.Cell(5, 7, ":")                        // Titik dua
		pdf.SetFont("Arial", "B", 10)              // Value agak tebal
		pdf.MultiCell(0, 7, value, "", "L", false) // MultiCell biar kalau panjang dia wrap ke bawah
	}

	// --- A. DATA PRIBADI ---
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(0, 8, " A. DATA PRIBADI", "1", 1, "L", true, 0, "")
	pdf.Ln(2)

	printRow("Nama Lengkap", student.FullName)
	printRow("NISN", student.NISNValue())
	printRow("NIM", student.NIMValue())

	// Format Tanggal & Gender
	dob := "-"
	if !student.DateOfBirth.IsZero() {
		dob = student.DateOfBirth.Format("02 January 2006")
	}
	placeDate := fmt.Sprintf("%s, %s", utils.SafeString(student.PlaceOfBirth), dob)
	printRow("Tempat, Tgl Lahir", placeDate)

	gender := "Laki-laki"
	if student.Gender == "female" {
		gender = "Perempuan"
	}
	printRow("Jenis Kelamin", gender)

	// Alamat lengkap
	fullAddress := utils.JoinAddress(
		student.Address,
		student.RT,
		student.RW,
		student.SubDistrict,
		student.District,
	)
	printRow("Detail Alamat", fullAddress)
	printRow("Kota/Kab", utils.SafeString(student.City))
	printRow("Provinsi", utils.SafeString(student.Province))

	pdf.Ln(5)

	// --- B. DATA ORANG TUA ---
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, " B. DATA ORANG TUA / WALI", "1", 1, "L", true, 0, "")
	pdf.Ln(2)

	if len(student.Parents) > 0 {
		for i, p := range student.Parents {
			if p.Parent.ID != "" {
				// Header Orang Tua

				relation := ""

				if p.RelationshipType != "" {
					switch p.RelationshipType {
					case "FATHER":
						relation = "Ayah"
					case "MOTHER":
						relation = "Ibu"
					default:
						relation = p.RelationshipType
					}
				}

				parentLabel := fmt.Sprintf("Orang Tua %d - %s", i+1, relation)
				pdf.SetFont("Arial", "B", 10)
				pdf.Cell(0, 7, parentLabel)
				pdf.Ln(7)

				// Nama
				printRow("Nama Lengkap", p.Parent.FullName)

				// No Telepon
				phone := "-"
				if p.Parent.PhoneNumber != nil {
					phone = *p.Parent.PhoneNumber
				}
				printRow("No. Telepon", phone)

				// Email
				email := "-"
				if p.Parent.Email != nil {
					email = *p.Parent.Email
				}
				printRow("Email", email)

				// Pendidikan
				education := "-"
				if p.Parent.EducationLevel != nil {
					education = *p.Parent.EducationLevel
				}
				printRow("Pendidikan", education)

				// Pekerjaan
				occupation := "-"
				if p.Parent.Occupation != nil {
					occupation = *p.Parent.Occupation
				}
				printRow("Pekerjaan", occupation)

				// Penghasilan
				income := "-"
				if p.Parent.IncomeRange != nil {
					income = *p.Parent.IncomeRange
				}
				printRow("Penghasilan", income)

				// Alamat Lengkap
				parentAddress := utils.JoinAddress(
					p.Parent.Address,
					p.Parent.RT,
					p.Parent.RW,
					p.Parent.SubDistrict,
					p.Parent.District,
					p.Parent.City,
					p.Parent.Province,
				)
				if parentAddress != "" {
					printRow("Alamat", parentAddress)
				}

				pdf.Ln(3) // Spacing antar orang tua
			}
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "Belum ada data orang tua yang terhubung.")
		pdf.Ln(5)
	}

	// Jika ada Wali (Guardian)
	if student.GuardianID != nil && student.GuardianType != nil {
		pdf.Ln(3)

		// Fetch guardian info menggunakan helper yang sudah ada
		guardianInfo, err := s.fetchGuardianInfo(student.GuardianID, student.GuardianType)
		if err == nil && guardianInfo != nil {
			// Header Wali
			pdf.SetFont("Arial", "B", 10)
			pdf.Cell(0, 7, "Wali Murid")
			pdf.Ln(7)

			// Nama
			printRow("Nama Lengkap", guardianInfo.FullName)

			// No Telepon
			phone := "-"
			if guardianInfo.PhoneNumber != nil && *guardianInfo.PhoneNumber != "" {
				phone = *guardianInfo.PhoneNumber
			}
			printRow("No. Telepon", phone)

			// Email
			email := "-"
			if guardianInfo.Email != nil && *guardianInfo.Email != "" {
				email = *guardianInfo.Email
			}
			printRow("Email", email)

			// Relasi/Hubungan
			relationship := "-"
			if guardianInfo.Relationship != "" {
				relationship = guardianInfo.Relationship
			}
			printRow("Hubungan", relationship)

			// Untuk Guardian, kita perlu fetch data lengkap dari repository
			// karena GuardianInfoResponse tidak memiliki semua field
			switch *student.GuardianType {
			case "guardian":
				guardian, err := s.guardianRepo.FindByID(*student.GuardianID)
				if err == nil && guardian != nil {
					// Alamat Lengkap
					guardianAddress := utils.JoinAddress(
						guardian.Address,
						guardian.RT,
						guardian.RW,
						guardian.SubDistrict,
						guardian.District,
					)
					if guardianAddress != "" {
						printRow("Alamat", guardianAddress)

						city := "-"
						if guardian.City != nil {
							city = *guardian.City
						}
						printRow("Kota/Kab", city)

						province := "-"
						if guardian.Province != nil {
							province = *guardian.Province
						}
						printRow("Provinsi", province)
					}
				}
			case "parent":
				// Jika wali adalah parent, fetch dari parent repo
				parent, err := s.parentRepo.FindByID(*student.GuardianID)
				if err == nil && parent != nil {
					// Pendidikan
					education := "-"
					if parent.EducationLevel != nil {
						education = *parent.EducationLevel
					}
					printRow("Pendidikan", education)

					// Pekerjaan
					occupation := "-"
					if parent.Occupation != nil {
						occupation = *parent.Occupation
					}
					printRow("Pekerjaan", occupation)

					// Penghasilan
					income := "-"
					if parent.IncomeRange != nil {
						income = *parent.IncomeRange
					}
					printRow("Penghasilan", income)

					// Alamat Lengkap
					guardianAddress := utils.JoinAddress(
						parent.Address,
						parent.RT,
						parent.RW,
						parent.SubDistrict,
						parent.District,
					)
					if guardianAddress != "" {
						printRow("Alamat", guardianAddress)

						city := "-"
						if parent.City != nil {
							city = *parent.City
						}
						printRow("Kota/Kab", city)

						province := "-"
						if parent.Province != nil {
							province = *parent.Province
						}
						printRow("Provinsi", province)
					}
				}
			}
		}
	}

	pdf.Ln(10)

	// --- TANDA TANGAN ---
	// Posisi kanan bawah
	currentY := pdf.GetY()
	if currentY > 250 { // Kalau halaman mau habis, tambah halaman baru
		pdf.AddPage()
		currentY = pdf.GetY()
	}

	pdf.SetX(120)
	pdf.SetFont("Arial", "", 10)

	// PERBAIKAN: Ganti pdf.Cell menjadi pdf.CellFormat
	pdf.CellFormat(0, 5, fmt.Sprintf("Surakarta, %s", time.Now().Format("02 January 2006")), "", 1, "C", false, 0, "")

	pdf.SetX(120)
	// Ganti pdf.Cell menjadi pdf.CellFormat
	pdf.CellFormat(0, 5, "Mengetahui,", "", 1, "C", false, 0, "")

	pdf.Ln(20) // Spasi Tanda Tangan

	pdf.SetX(120)
	pdf.SetFont("Arial", "B", 10)
	// PERBAIKAN: Ganti pdf.Cell menjadi pdf.CellFormat
	pdf.CellFormat(0, 5, "( ..................................... )", "", 1, "C", false, 0, "")

	// Output
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return &buf, nil
}
