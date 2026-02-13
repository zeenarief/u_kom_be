package service

import (
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/repository"
)

type GradeService interface {
	CreateAssessment(req request.AssessmentCreateRequest) (*domain.Assessment, error)
	GetAssessmentsByTeachingAssignment(teachingAssignmentID string) ([]domain.Assessment, error)
	GetAssessmentDetail(id string) (*domain.Assessment, error)
	SubmitScores(req request.BulkScoreRequest) error
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

	return assessment, nil
}

func (s *gradeService) GetAssessmentsByTeachingAssignment(teachingAssignmentID string) ([]domain.Assessment, error) {
	return s.gradeRepo.GetAssessmentsByTeachingAssignment(teachingAssignmentID)
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
