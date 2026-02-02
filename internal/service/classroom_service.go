package service

import (
	"errors"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"

	"gorm.io/gorm"
)

type ClassroomService interface {
	Create(req request.ClassroomCreateRequest) (*response.ClassroomResponse, error)
	FindAll(academicYearID string) ([]response.ClassroomResponse, error)
	FindByID(id string) (*response.ClassroomDetailResponse, error)
	Update(id string, req request.ClassroomUpdateRequest) (*response.ClassroomResponse, error)
	Delete(id string) error

	AddStudents(id string, req request.AddStudentsToClassRequest) error
	RemoveStudent(id string, studentID string) error
}

type classroomService struct {
	repo             repository.ClassroomRepository
	academicYearRepo repository.AcademicYearRepository // Untuk validasi
	employeeRepo     repository.EmployeeRepository     // Untuk validasi
	studentRepo      repository.StudentRepository      // Untuk validasi
	db               *gorm.DB
}

func NewClassroomService(
	repo repository.ClassroomRepository,
	ayRepo repository.AcademicYearRepository,
	empRepo repository.EmployeeRepository,
	stdRepo repository.StudentRepository,
	db *gorm.DB,
) ClassroomService {
	return &classroomService{
		repo:             repo,
		academicYearRepo: ayRepo,
		employeeRepo:     empRepo,
		studentRepo:      stdRepo,
		db:               db,
	}
}

// Helper untuk mapping response
func (s *classroomService) toResponse(c *domain.Classroom) *response.ClassroomResponse {
	teacherName := "-"
	// Tambahkan variabel teacherID
	var teacherID *string

	if c.HomeroomTeacher != nil {
		teacherName = c.HomeroomTeacher.FullName
		teacherID = &c.HomeroomTeacher.ID // Ambil ID
	} else if c.HomeroomTeacherID != nil {
		// Fallback jika preload teacher gagal tapi ID ada
		teacherID = c.HomeroomTeacherID
	}

	return &response.ClassroomResponse{
		ID:          c.ID,
		Name:        c.Name,
		Level:       c.Level,
		Major:       c.Major,
		Description: c.Description,
		AcademicYear: response.AcademicYearResponse{
			ID:   c.AcademicYear.ID,
			Name: c.AcademicYear.Name,
			// ... pastikan field lain mappingnya benar
		},
		HomeroomTeacherID:   teacherID, // MAP ID DI SINI
		HomeroomTeacherName: teacherName,
		TotalStudents:       c.TotalStudents, // AMBIL DARI DOMAIN (Hasil Query Repository)
		CreatedAt:           c.CreatedAt,
	}
}

func (s *classroomService) Create(req request.ClassroomCreateRequest) (*response.ClassroomResponse, error) {
	// 1. Validasi Academic Year
	ay, err := s.academicYearRepo.FindByID(req.AcademicYearID)
	if err != nil || ay == nil {
		return nil, errors.New("academic year not found")
	}

	// 2. Validasi Employee (Jika ada)
	if req.HomeroomTeacherID != nil {
		emp, err := s.employeeRepo.FindByID(*req.HomeroomTeacherID)
		if err != nil || emp == nil {
			return nil, errors.New("homeroom teacher not found")
		}
	}

	classroom := &domain.Classroom{
		AcademicYearID:    req.AcademicYearID,
		HomeroomTeacherID: req.HomeroomTeacherID,
		Name:              req.Name,
		Level:             req.Level,
		Major:             req.Major,
		Description:       req.Description,
	}

	if err := s.repo.Create(classroom); err != nil {
		return nil, err
	}

	// Load ulang untuk dapat relasi
	createdClass, _ := s.repo.FindByID(classroom.ID)
	return s.toResponse(createdClass), nil
}

func (s *classroomService) FindAll(academicYearID string) ([]response.ClassroomResponse, error) {
	classrooms, err := s.repo.FindAll(academicYearID)
	if err != nil {
		return nil, err
	}

	var responses []response.ClassroomResponse
	for _, c := range classrooms {
		// Hacky way to get count from subquery scan (requires GORM hook or map decoding usually)
		// For simplicity, we assume repo handled basic struct
		res := s.toResponse(&c)
		responses = append(responses, *res)
	}
	return responses, nil
}

func (s *classroomService) FindByID(id string) (*response.ClassroomDetailResponse, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("classroom not found")
	}

	baseRes := s.toResponse(c)

	// Map Students
	var studentResponses []response.StudentInClassResponse
	for _, sc := range c.StudentClassrooms {
		nisn := "-"
		if sc.Student.NISN != nil {
			nisn = *sc.Student.NISN
		}

		studentResponses = append(studentResponses, response.StudentInClassResponse{
			ID:       sc.Student.ID,
			FullName: sc.Student.FullName,
			NISN:     nisn,
			Gender:   sc.Student.Gender,
			Status:   sc.Status,
		})
	}

	baseRes.TotalStudents = int64(len(studentResponses))

	return &response.ClassroomDetailResponse{
		ClassroomResponse: *baseRes,
		Students:          studentResponses,
	}, nil
}

func (s *classroomService) Update(id string, req request.ClassroomUpdateRequest) (*response.ClassroomResponse, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("classroom not found")
	}

	if req.Name != "" {
		c.Name = req.Name
	}
	if req.Level != "" {
		c.Level = req.Level
	}
	if req.Major != "" {
		c.Major = req.Major
	}
	c.Description = req.Description // Boleh kosong

	// Handle pointer update untuk Wali Kelas
	if req.HomeroomTeacherID != nil {
		// 1. Validasi: Pastikan Guru ada di DB (kecuali string kosong untuk menghapus)
		if *req.HomeroomTeacherID != "" {
			emp, err := s.employeeRepo.FindByID(*req.HomeroomTeacherID)
			if err != nil || emp == nil {
				return nil, errors.New("homeroom teacher not found")
			}
		}

		// 2. Assign ID Baru
		c.HomeroomTeacherID = req.HomeroomTeacherID

		// 3. FIX: Kosongkan struct relasi agar GORM mengupdate berdasarkan ID baru
		c.HomeroomTeacher = nil
	}

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}

	// Load ulang untuk memastikan response mendapatkan data relasi terbaru
	updated, _ := s.repo.FindByID(id)

	// FIX: Kembalikan 'updated' object, bukan 'c'
	// (Walaupun kode lama sudah pakai 'updated', pastikan ini konsisten)
	return s.toResponse(updated), nil
}

func (s *classroomService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *classroomService) AddStudents(id string, req request.AddStudentsToClassRequest) error {
	// 1. Cek Kelas
	c, err := s.repo.FindByID(id)
	if err != nil || c == nil {
		return errors.New("classroom not found")
	}

	var newRelations []domain.StudentClassroom

	for _, studentID := range req.StudentIDs {
		// Validasi Student Exists
		student, _ := s.studentRepo.FindByID(studentID)
		if student == nil {
			continue // Skip student invalid
		}

		// Validasi Duplikat di kelas yang sama
		exists, _ := s.repo.IsStudentInClass(id, studentID)
		if exists {
			continue // Skip jika sudah ada
		}

		newRelations = append(newRelations, domain.StudentClassroom{
			ClassroomID: id,
			StudentID:   studentID,
			Status:      "ACTIVE",
		})
	}

	if len(newRelations) > 0 {
		return s.repo.AddStudents(newRelations)
	}

	return nil // Tidak ada yang ditambahkan (mungkin duplikat semua)
}

func (s *classroomService) RemoveStudent(id string, studentID string) error {
	return s.repo.RemoveStudent(id, studentID)
}
