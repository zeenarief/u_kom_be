package service

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
	"u_kom_be/internal/apperrors"
	"u_kom_be/internal/converter"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
	"u_kom_be/internal/utils"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"
)

type StudentService interface {
	CreateStudent(req request.StudentCreateRequest) (*response.StudentDetailResponse, error)
	GetStudentByID(id string) (*response.StudentDetailResponse, error)
	GetAllStudents(search string) ([]response.StudentListResponse, error)
	UpdateStudent(id string, req request.StudentUpdateRequest) (*response.StudentDetailResponse, error)
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
			return nil, fmt.Errorf("Failed to hash NIK: %w", err)
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
			return nil, fmt.Errorf("Failed to encrypt NIK: %w", err)
		}
		encryptedNIK = &encrypted
	}

	encryptedNoKK := ""
	if req.NoKK != "" {
		var err error
		encryptedNoKK, err = s.encryptionUtil.Encrypt(req.NoKK)
		if err != nil {
			return nil, fmt.Errorf("Failed to encrypt NoKK: %w", err)
		}
	}

	// 3. Buat Domain Object
	student := &domain.Student{
		FullName:     req.FullName,
		NoKK:         encryptedNoKK,
		NIK:          encryptedNIK,
		NIKHash:      nikHash,
		NISN:         nisn,
		NIM:          nim,
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
		Status:       req.Status,
		EntryYear:    req.EntryYear,
		ExitYear:     req.ExitYear,
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

// GetAllStudents mengambil semua siswa
func (s *studentService) GetAllStudents(search string) ([]response.StudentListResponse, error) {
	students, err := s.studentRepo.FindAll(search)
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
	// Enkripsi field yang diperbarui
	if req.NIK != "" {
		// Cek apakah NIK berubah atau datanya belum ada (meski user kirim NIK sama, kita proses saja)
		// Optimalisasi: cek hash dulu? Tapi request raw string, di DB terenkripsi/hash.
		// Kita hitung hash dari request.
		newHash, err := s.encryptionUtil.Hash(req.NIK)
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
		encryptedNIK, err := s.encryptionUtil.Encrypt(req.NIK)
		if err != nil {
			return nil, fmt.Errorf("Failed to encrypt NIK: %w", err)
		}
		student.NIK = &encryptedNIK
	}
	if req.NoKK != "" {
		encryptedNoKK, err := s.encryptionUtil.Encrypt(req.NoKK)
		if err != nil {
			return nil, fmt.Errorf("Failed to encrypt NoKK: %w", err)
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
	if req.Status != "" {
		student.Status = req.Status
	}
	if req.EntryYear != "" {
		student.EntryYear = req.EntryYear
	}
	if req.ExitYear != "" {
		student.ExitYear = req.ExitYear
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
		return apperrors.NewNotFoundError("Student not found")
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
		return errors.New("Student not found")
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
			return fmt.Errorf("Error checking parent: %w", err)
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
		return errors.New("Student not found")
	}

	// 2. Validasi apakah guardian_id yang diberikan ada di tabel yang benar
	switch req.GuardianType {
	case "parent":
		parent, err := s.parentRepo.FindByID(req.GuardianID)
		if err != nil {
			return fmt.Errorf("Error checking parent: %w", err)
		}
		if parent == nil {
			return apperrors.NewNotFoundError(fmt.Sprintf("Parent not found with id: %s", req.GuardianID))
		}
	case "guardian":
		guardian, err := s.guardianRepo.FindByID(req.GuardianID)
		if err != nil {
			return fmt.Errorf("Error checking guardian: %w", err)
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
		return errors.New("Student not found")
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
			return nil, fmt.Errorf("Data integrity error: Parent guardian with id %s not found", id)
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
			return nil, fmt.Errorf("Data integrity error: Guardian with id %s not found", id)
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

	return nil, fmt.Errorf("Unknown guardian_type: %s", tipe)
}

// LinkUser menautkan profil Student ke akun User
func (s *studentService) LinkUser(studentID string, userID string) error {
	// 1. Cek apakah Student ada
	student, err := s.studentRepo.FindByID(studentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("Student not found")
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
	students, err := s.studentRepo.FindAll("")
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
			&student.Address,
			&student.RT,
			&student.RW,
			&student.SubDistrict,
			&student.District,
			&student.City,
			&student.Province,
			&student.PostalCode,
		)

		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), utils.SafeString(student.NISN))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), utils.SafeString(student.NIM))
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), student.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), student.Gender)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), student.PlaceOfBirth)
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
	students, err := s.studentRepo.FindAll("")
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
		pdf.CellFormat(widths[5], 8, student.Address, "1", 0, "L", false, 0, "") // Alamat mungkin terpotong jika panjang, nanti bisa pakai MultiCell
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
	// Kita gunakan GetStudentByID yg sudah ada logic lengkapnya (parents, guardian info, dll)
	// Tapi karena GetStudentByID mengembalikan Response DTO, dan kita butuh Domain Object untuk report,
	// Lebih aman kita panggil Repo langsung:
	student, err := s.studentRepo.FindByIDWithParents(id)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, apperrors.NewNotFoundError("Student not found")
	}

	// 2. Init PDF (Portrait, A4)
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// --- HEADER KOP ---
	// Cek apakah ada logo
	logoPath := "assets/logo_sekolah.png"
	// Kita coba load image, jika tidak ada, skip
	// ImageOptions(src, x, y, width, height, flow, options, link, linkStr)
	// x=10, y=10, w=25, h=0 (auto keep aspect ratio)
	pdf.ImageOptions(logoPath, 10, 10, 25, 0, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// Geser Text agak ke kanan jika ada logo (atau selalu geser biar rapi)
	// Tapi karena "C" (Center) itu relatif terhadap page width, kita bisa mainkan Margin atau Cell width.
	// Cara umum KOP: Logo kiri absolute, Teks Center di page.

	pdf.SetFont("Arial", "", 13)
	pdf.SetY(10)
	pdf.CellFormat(0, 6, "YAYASAN MAJLIS TALIM NURUL HUDA KARTASURA", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 6, "PONDOK PESANTREN NURUL HUDA KARTASURA", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, "(WISMA ASUHAN YATIM NURUL HUDA KARTASURA)", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6, "Gg. Anggrek, Bakalan 02/02, Pucangan, Kartasura, Sukoharjo", "", 1, "C", false, 0, "")

	pdf.Ln(5)

	// Garis pembatas (Double line atau tebal)
	pdf.SetLineWidth(0.5)
	pdf.Line(10, 32, 200, 32)
	pdf.SetLineWidth(0.2) // Balikin ke default
	pdf.Ln(5)             // Space setelah garis

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
	placeDate := fmt.Sprintf("%s, %s", student.PlaceOfBirth, dob)
	printRow("Tempat, Tgl Lahir", placeDate)

	gender := "Laki-laki"
	if student.Gender == "female" {
		gender = "Perempuan"
	}
	printRow("Jenis Kelamin", gender)

	// Alamat lengkap
	fullAddress := fmt.Sprintf("%s RT %s / RW %s, Kel. %s, Kec. %s", student.Address, student.RT, student.RW, student.SubDistrict, student.District)
	printRow("Detail Alamat", fullAddress)
	printRow("Kota/Kab", student.City)
	printRow("Provinsi", student.Province)

	pdf.Ln(5)

	// --- B. DATA ORANG TUA ---
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, " B. DATA ORANG TUA / WALI", "1", 1, "L", true, 0, "")
	pdf.Ln(2)

	if len(student.Parents) > 0 {
		for i, p := range student.Parents {
			// Kita perlu fetch data parent detail karena di student.Parents cuma ada ID & Relasi (tergantung preload)
			// Asumsi FindByIDWithParents sudah preload 'Parent' objectnya
			// Jika struktur domain Anda: Student struct { Parents []StudentParent } dan StudentParent punya { Parent Parent }

			// Sederhananya, kita ambil parent pertama sebagai contoh,
			// atau jika Anda sudah implementasi preload yang benar:
			parentLabel := fmt.Sprintf("Orang Tua %d (%s)", i+1, p.RelationshipType)
			// Disini kita butuh akses ke nama parent.
			// Jika preload di repo belum deep, mungkin nama parent kosong.
			// Asumsikan sudah ada:
			if p.Parent.ID != "" {
				printRow(parentLabel, p.Parent.FullName)
				phone := "-"
				if p.Parent.PhoneNumber != nil {
					phone = *p.Parent.PhoneNumber
				}
				printRow("   No. HP", phone)
			}
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.Cell(0, 10, "Belum ada data orang tua yang terhubung.")
		pdf.Ln(5)
	}

	// Jika ada Wali (Guardian)
	if student.GuardianID != nil {
		// Logic fetch guardian info sama seperti di GetStudentByID
		// Untuk ringkasnya di PDF, Anda bisa fetch ulang atau pakai helper fetchGuardianInfo
		// ...
		printRow("Wali Murid", "Terdata (Lihat detail di sistem)")
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
