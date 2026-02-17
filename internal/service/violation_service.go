package service

import (
	"errors"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
	"smart_school_be/internal/utils"
)

type ViolationService interface {
	// Category
	CreateCategory(req request.CreateViolationCategoryRequest) error
	GetCategories() ([]response.ViolationCategoryResponse, error)
	UpdateCategory(id string, req request.UpdateViolationCategoryRequest) error
	DeleteCategory(id string) error

	// Type
	CreateType(req request.CreateViolationTypeRequest) error
	GetTypes(categoryID string) ([]response.ViolationTypeResponse, error)
	UpdateType(id string, req request.UpdateViolationTypeRequest) error
	DeleteType(id string) error

	// Student Violation
	// Student Violation
	RecordViolation(req request.CreateStudentViolationRequest) error
	GetStudentViolations(studentID string, pagination request.PaginationRequest) (*response.PaginatedData, error)
	GetViolationDetail(id string) (*response.StudentViolationDetailResponse, error)
	UpdateViolation(id string, req request.UpdateStudentViolationRequest) error
	DeleteViolation(id string) error
	GetAllViolations(filter string, pagination request.PaginationRequest) (*response.PaginatedData, error)
}

type violationService struct {
	violationRepo repository.ViolationRepository
	studentRepo   repository.StudentRepository
}

func NewViolationService(violationRepo repository.ViolationRepository, studentRepo repository.StudentRepository) ViolationService {
	return &violationService{
		violationRepo: violationRepo,
		studentRepo:   studentRepo,
	}
}

// Category Implementation
func (s *violationService) CreateCategory(req request.CreateViolationCategoryRequest) error {
	category := &domain.ViolationCategory{
		ID:          utils.GenerateUUID(),
		Name:        req.Name,
		Description: req.Description,
	}
	return s.violationRepo.CreateCategory(category)
}

func (s *violationService) GetCategories() ([]response.ViolationCategoryResponse, error) {
	categories, err := s.violationRepo.FindAllCategories()
	if err != nil {
		return nil, err
	}

	var responses []response.ViolationCategoryResponse
	for _, c := range categories {
		responses = append(responses, response.ViolationCategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}
	return responses, nil
}

func (s *violationService) UpdateCategory(id string, req request.UpdateViolationCategoryRequest) error {
	category, err := s.violationRepo.FindCategoryByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	return s.violationRepo.UpdateCategory(category)
}

func (s *violationService) DeleteCategory(id string) error {
	return s.violationRepo.DeleteCategory(id)
}

// Type Implementation
func (s *violationService) CreateType(req request.CreateViolationTypeRequest) error {
	// Optional validaton if category exists, but FK constraint handles it too.
	// Better to check.
	// We can check if category exists via repo.

	violationType := &domain.ViolationType{
		ID:            utils.GenerateUUID(),
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		Description:   req.Description,
		DefaultPoints: req.DefaultPoints,
	}
	return s.violationRepo.CreateType(violationType)
}

func (s *violationService) GetTypes(categoryID string) ([]response.ViolationTypeResponse, error) {
	types, err := s.violationRepo.FindAllTypes(categoryID)
	if err != nil {
		return nil, err
	}

	var responses []response.ViolationTypeResponse
	for _, t := range types {
		resp := response.ViolationTypeResponse{
			ID:            t.ID,
			CategoryID:    t.CategoryID,
			Name:          t.Name,
			Description:   t.Description,
			DefaultPoints: t.DefaultPoints,
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
		if t.Category != nil {
			resp.Category = &response.ViolationCategoryResponse{
				ID:          t.Category.ID,
				Name:        t.Category.Name,
				Description: t.Category.Description,
				CreatedAt:   t.Category.CreatedAt,
				UpdatedAt:   t.Category.UpdatedAt,
			}
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *violationService) UpdateType(id string, req request.UpdateViolationTypeRequest) error {
	violationType, err := s.violationRepo.FindTypeByID(id)
	if err != nil {
		return err
	}
	if violationType == nil {
		return errors.New("violation type not found")
	}

	if req.CategoryID != "" {
		// Validate category if needed
		violationType.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		violationType.Name = req.Name
	}
	if req.Description != "" {
		violationType.Description = req.Description
	}
	if req.DefaultPoints != nil {
		violationType.DefaultPoints = *req.DefaultPoints
	}

	return s.violationRepo.UpdateType(violationType)
}

func (s *violationService) DeleteType(id string) error {
	return s.violationRepo.DeleteType(id)
}

// Student Violation Implementation
func (s *violationService) RecordViolation(req request.CreateStudentViolationRequest) error {
	// Validate student exists
	student, err := s.studentRepo.FindByID(req.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return errors.New("student not found")
	}

	// Validate violation type exists
	violationType, err := s.violationRepo.FindTypeByID(req.ViolationTypeID)
	if err != nil {
		return err
	}
	if violationType == nil {
		return errors.New("violation type not found")
	}

	violation := &domain.StudentViolation{
		ID:              utils.GenerateUUID(),
		StudentID:       req.StudentID,
		ViolationTypeID: req.ViolationTypeID,
		ViolationDate:   req.ViolationDate,
		Points:          violationType.DefaultPoints, // Snapshot points!
		ActionTaken:     req.ActionTaken,
		Notes:           req.Notes,
	}

	return s.violationRepo.RecordViolation(violation)
}

func (s *violationService) GetStudentViolations(studentID string, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	violations, total, err := s.violationRepo.FindStudentViolations(studentID, limit, offset)
	if err != nil {
		return nil, err
	}

	violationResponses := s.mapToResponses(violations)

	paginatedData := response.NewPaginatedData(violationResponses, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

func (s *violationService) GetViolationDetail(id string) (*response.StudentViolationDetailResponse, error) {
	violation, err := s.violationRepo.FindViolationByID(id)
	if err != nil {
		return nil, err
	}
	if violation == nil {
		return nil, errors.New("violation not found")
	}

	resp := &response.StudentViolationDetailResponse{
		ID:              violation.ID,
		StudentID:       violation.StudentID,
		ViolationTypeID: violation.ViolationTypeID,
		ViolationDate:   violation.ViolationDate,
		Points:          violation.Points,
		ActionTaken:     violation.ActionTaken,
		Notes:           violation.Notes,
		CreatedAt:       violation.CreatedAt,
		UpdatedAt:       violation.UpdatedAt,
	}

	if violation.Student != nil {
		resp.StudentName = violation.Student.FullName
		resp.StudentNIM = violation.Student.NIM
		resp.StudentNISN = violation.Student.NISN
		// resp.StudentClass = violation.Student.Classroom.Name // Requires Preload
	}

	if violation.ViolationType != nil {
		resp.ViolationName = violation.ViolationType.Name
		if violation.ViolationType.Category != nil {
			resp.ViolationCategory = violation.ViolationType.Category.Name
		}
	}

	return resp, nil
}

func (s *violationService) UpdateViolation(id string, req request.UpdateStudentViolationRequest) error {
	violation, err := s.violationRepo.FindViolationByID(id)
	if err != nil {
		return err
	}
	if violation == nil {
		return errors.New("violation not found")
	}

	if req.ViolationTypeID != "" {
		// Validasi violation type exists?
		// Better to check if needed or rely on FK error. checking is safer.
		vType, err := s.violationRepo.FindTypeByID(req.ViolationTypeID)
		if err != nil {
			return err
		}
		if vType == nil {
			return errors.New("violation type not found")
		}
		violation.ViolationTypeID = req.ViolationTypeID
		// If points not provided, should we update points to default of new type?
		// Requirement unclear. But usually if points NOT in request, we might want to keep existing OR update to new default.
		// Let's assume: if points IS provided, use it. If NOT provided but Type CHANGED, maybe use default?
		// For safety, let's only update points if explicitly requested.
	}

	if req.ViolationDate != nil {
		violation.ViolationDate = *req.ViolationDate
	}
	if req.Points != nil {
		violation.Points = *req.Points
	}
	if req.ActionTaken != "" {
		violation.ActionTaken = req.ActionTaken
	}
	if req.Notes != "" {
		violation.Notes = req.Notes
	}

	return s.violationRepo.UpdateViolation(violation)
}

func (s *violationService) DeleteViolation(id string) error {
	return s.violationRepo.DeleteViolation(id)
}

func (s *violationService) GetAllViolations(filter string, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	violations, total, err := s.violationRepo.FindAllViolations(filter, limit, offset)
	if err != nil {
		return nil, err
	}

	violationResponses := s.mapToResponses(violations)

	paginatedData := response.NewPaginatedData(violationResponses, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

func (s *violationService) mapToResponses(violations []domain.StudentViolation) []response.StudentViolationListResponse {
	var responses []response.StudentViolationListResponse
	for _, v := range violations {
		resp := response.StudentViolationListResponse{
			ID:            v.ID,
			StudentID:     v.StudentID,
			ViolationDate: v.ViolationDate,
			Points:        v.Points,
			CreatedAt:     v.CreatedAt,
			UpdatedAt:     v.UpdatedAt,
		}

		if v.Student != nil {
			resp.StudentName = v.Student.FullName
		}
		if v.ViolationType != nil {
			resp.ViolationName = v.ViolationType.Name
			if v.ViolationType.Category != nil {
				resp.ViolationCategory = v.ViolationType.Category.Name
			}
		}

		responses = append(responses, resp)
	}
	return responses
}
