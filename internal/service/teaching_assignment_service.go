package service

import (
	"smart_school_be/internal/apperrors"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
)

type TeachingAssignmentService interface {
	Create(req request.TeachingAssignmentCreateRequest) (*response.TeachingAssignmentResponse, error)
	GetByClassroom(classroomID string) ([]response.TeachingAssignmentResponse, error)
	GetByTeacher(teacherID string) ([]response.TeachingAssignmentResponse, error)
	Delete(id string) error
}

type teachingAssignmentService struct {
	repo          repository.TeachingAssignmentRepository
	classroomRepo repository.ClassroomRepository
	subjectRepo   repository.SubjectRepository
	employeeRepo  repository.EmployeeRepository
}

func NewTeachingAssignmentService(
	repo repository.TeachingAssignmentRepository,
	cRepo repository.ClassroomRepository,
	sRepo repository.SubjectRepository,
	eRepo repository.EmployeeRepository,
) TeachingAssignmentService {
	return &teachingAssignmentService{
		repo:          repo,
		classroomRepo: cRepo,
		subjectRepo:   sRepo,
		employeeRepo:  eRepo,
	}
}

func (s *teachingAssignmentService) toResponse(d *domain.TeachingAssignment) response.TeachingAssignmentResponse {
	nip := "-"
	if d.Teacher.NIP != nil {
		nip = *d.Teacher.NIP
	}

	return response.TeachingAssignmentResponse{
		ID: d.ID,
		Classroom: response.TeachingAssignmentClassroomResponse{
			ID:    d.Classroom.ID,
			Name:  d.Classroom.Name,
			Level: d.Classroom.Level,
			Major: d.Classroom.Major,
		},
		Subject: response.TeachingAssignmentSubjectResponse{
			ID:   d.Subject.ID,
			Name: d.Subject.Name,
			Code: d.Subject.Code,
		},
		Teacher: response.TeachingAssignmentTeacherResponse{
			ID:  d.Teacher.ID,
			NIP: nip,
			User: response.TeachingAssignmentUserResponse{
				Name: d.Teacher.FullName, // Using FullName from Employee struct which seems to serve as display name
			},
		},
	}
}

func (s *teachingAssignmentService) Create(req request.TeachingAssignmentCreateRequest) (*response.TeachingAssignmentResponse, error) {
	// 1. Validasi FK
	c, _ := s.classroomRepo.FindByID(req.ClassroomID)
	if c == nil {
		return nil, apperrors.NewNotFoundError("classroom not found")
	}
	sub, _ := s.subjectRepo.FindByID(req.SubjectID)
	if sub == nil {
		return nil, apperrors.NewNotFoundError("subject not found")
	}
	emp, _ := s.employeeRepo.FindByID(req.TeacherID)
	if emp == nil {
		return nil, apperrors.NewNotFoundError("teacher not found")
	}

	// 2. Cek Duplikat (Apakah mapel ini sudah ada gurunya di kelas ini?)
	existing, _ := s.repo.FindOne(req.ClassroomID, req.SubjectID)
	if existing != nil {
		return nil, apperrors.NewConflictError("this subject already has a teacher in this class")
	}

	// 3. Simpan
	assignment := &domain.TeachingAssignment{
		ClassroomID: req.ClassroomID,
		SubjectID:   req.SubjectID,
		TeacherID:   req.TeacherID,
	}

	if err := s.repo.Create(assignment); err != nil {
		return nil, err
	}

	// Manual Assign for Response (supaya tidak perlu query ulang)
	assignment.Classroom = *c
	assignment.Subject = *sub
	assignment.Teacher = *emp

	res := s.toResponse(assignment)
	return &res, nil
}

func (s *teachingAssignmentService) GetByClassroom(classroomID string) ([]response.TeachingAssignmentResponse, error) {
	data, err := s.repo.FindByClassroomID(classroomID)
	if err != nil {
		return nil, err
	}

	var result []response.TeachingAssignmentResponse
	for _, d := range data {
		result = append(result, s.toResponse(&d))
	}
	return result, nil
}

func (s *teachingAssignmentService) GetByTeacher(teacherID string) ([]response.TeachingAssignmentResponse, error) {
	data, err := s.repo.FindByTeacherID(teacherID)
	if err != nil {
		return nil, err
	}

	var result []response.TeachingAssignmentResponse
	for _, d := range data {
		result = append(result, s.toResponse(&d))
	}
	return result, nil
}

func (s *teachingAssignmentService) Delete(id string) error {
	return s.repo.Delete(id)
}
