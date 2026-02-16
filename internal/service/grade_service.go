package service

import (
	"smart_school_be/internal/apperrors"
	"smart_school_be/internal/model/domain"
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/model/response"
	"smart_school_be/internal/repository"
)

type GradeService interface {
	CreateAssessment(req request.AssessmentCreateRequest) (*domain.Assessment, error)
	UpdateAssessment(id string, req request.AssessmentCreateRequest) (*domain.Assessment, error)
	GetAssessmentsByTeachingAssignment(teachingAssignmentID string, pagination request.PaginationRequest) (*response.PaginatedData, error)
	GetAssessmentDetail(id string) (*domain.Assessment, error)
	SubmitScores(req request.BulkScoreRequest) error
	DeleteAssessment(id string) error
}

type gradeService struct {
	gradeRepo repository.GradeRepository
}

func NewGradeService(gradeRepo repository.GradeRepository) GradeService {
	return &gradeService{gradeRepo: gradeRepo}
}

func (s *gradeService) CreateAssessment(req request.AssessmentCreateRequest) (*domain.Assessment, error) {
	if req.MaxScore == 0 {
		req.MaxScore = 100
	}

	assessment := &domain.Assessment{
		TeachingAssignmentID: req.TeachingAssignmentID,
		Title:                req.Title,
		Type:                 req.Type,
		MaxScore:             req.MaxScore,
		Date:                 req.Date,
		Description:          req.Description,
	}

	if err := s.gradeRepo.CreateAssessment(assessment); err != nil {
		return nil, err
	}

	// Reload to get full data (relationships)
	return s.gradeRepo.FindAssessmentByID(assessment.ID)
}

func (s *gradeService) UpdateAssessment(id string, req request.AssessmentCreateRequest) (*domain.Assessment, error) {
	assessment, err := s.gradeRepo.FindAssessmentByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields
	assessment.Title = req.Title
	assessment.Type = req.Type
	assessment.Date = req.Date
	assessment.Description = req.Description

	// Update Max Score if provided
	if req.MaxScore > 0 {
		assessment.MaxScore = req.MaxScore
	}

	if err := s.gradeRepo.UpdateAssessment(assessment); err != nil {
		return nil, err
	}

	return assessment, nil
}

func (s *gradeService) GetAssessmentsByTeachingAssignment(teachingAssignmentID string, pagination request.PaginationRequest) (*response.PaginatedData, error) {
	limit := pagination.GetLimit()
	offset := pagination.GetOffset()

	assessments, total, err := s.gradeRepo.GetAssessmentsByTeachingAssignment(teachingAssignmentID, limit, offset)
	if err != nil {
		return nil, err
	}

	paginatedData := response.NewPaginatedData(assessments, total, pagination.GetPage(), limit)
	return &paginatedData, nil
}

func (s *gradeService) GetAssessmentDetail(id string) (*domain.Assessment, error) {
	return s.gradeRepo.FindAssessmentByID(id)
}

func (s *gradeService) SubmitScores(req request.BulkScoreRequest) error {
	for _, scoreReq := range req.Scores {
		score := &domain.StudentScore{
			AssessmentID: req.AssessmentID,
			StudentID:    scoreReq.StudentID,
			Score:        scoreReq.Score,
			Feedback:     scoreReq.Feedback,
		}

		if err := s.gradeRepo.SaveStudentScore(score); err != nil {
			return err
		}
	}
	return nil
}

func (s *gradeService) DeleteAssessment(id string) error {
	// Check if assessment exists
	assessment, err := s.gradeRepo.FindAssessmentByID(id)
	if err != nil {
		return err
	}
	if assessment == nil {
		return apperrors.NewNotFoundError("assessment not found")
	}

	return s.gradeRepo.DeleteAssessment(id)
}
