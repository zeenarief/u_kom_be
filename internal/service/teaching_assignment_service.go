package service

import (
	"errors"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
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
		ID:            d.ID,
		ClassroomName: d.Classroom.Name,
		SubjectName:   d.Subject.Name,
		TeacherName:   d.Teacher.FullName,
		TeacherNIP:    nip,
	}
}

func (s *teachingAssignmentService) Create(req request.TeachingAssignmentCreateRequest) (*response.TeachingAssignmentResponse, error) {
	// 1. Validasi FK
	c, _ := s.classroomRepo.FindByID(req.ClassroomID)
	if c == nil {
		return nil, errors.New("classroom not found")
	}
	sub, _ := s.subjectRepo.FindByID(req.SubjectID)
	if sub == nil {
		return nil, errors.New("subject not found")
	}
	emp, _ := s.employeeRepo.FindByID(req.TeacherID)
	if emp == nil {
		return nil, errors.New("teacher not found")
	}

	// 2. Cek Duplikat (Apakah mapel ini sudah ada gurunya di kelas ini?)
	existing, _ := s.repo.FindOne(req.ClassroomID, req.SubjectID)
	if existing != nil {
		return nil, errors.New("this subject already has a teacher in this class")
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
