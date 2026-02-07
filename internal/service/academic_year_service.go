package service

import (
	"time"
	"u_kom_be/internal/apperrors"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"

	"gorm.io/gorm"
)

type AcademicYearService interface {
	Create(req request.AcademicYearCreateRequest) (*response.AcademicYearResponse, error)
	FindAll() ([]response.AcademicYearResponse, error)
	FindByID(id string) (*response.AcademicYearResponse, error)
	Update(id string, req request.AcademicYearUpdateRequest) (*response.AcademicYearResponse, error)
	Delete(id string) error
	Activate(id string) error
}

type academicYearService struct {
	repo repository.AcademicYearRepository
	db   *gorm.DB // Butuh DB instance untuk Transaction start
}

func NewAcademicYearService(repo repository.AcademicYearRepository, db *gorm.DB) AcademicYearService {
	return &academicYearService{repo: repo, db: db}
}

func (s *academicYearService) toResponse(a *domain.AcademicYear) *response.AcademicYearResponse {
	return &response.AcademicYearResponse{
		ID:        a.ID,
		Name:      a.Name,
		Status:    a.Status,
		StartDate: a.StartDate,
		EndDate:   a.EndDate,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (s *academicYearService) Create(req request.AcademicYearCreateRequest) (*response.AcademicYearResponse, error) {
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid start_date format (expected YYYY-MM-DD)")
	}

	endDate, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return nil, apperrors.NewBadRequestError("invalid end_date format (expected YYYY-MM-DD)")
	}

	academicYear := &domain.AcademicYear{
		Name:      req.Name,
		Status:    "INACTIVE", // Default Inactive
		StartDate: startDate,
		EndDate:   endDate,
	}

	if err := s.repo.Create(academicYear); err != nil {
		return nil, err
	}

	return s.toResponse(academicYear), nil
}

func (s *academicYearService) FindAll() ([]response.AcademicYearResponse, error) {
	academicYears, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.AcademicYearResponse
	for _, ay := range academicYears {
		responses = append(responses, *s.toResponse(&ay))
	}
	return responses, nil
}

func (s *academicYearService) FindByID(id string) (*response.AcademicYearResponse, error) {
	ay, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ay == nil {
		return nil, apperrors.NewNotFoundError("academic year not found")
	}
	return s.toResponse(ay), nil
}

func (s *academicYearService) Update(id string, req request.AcademicYearUpdateRequest) (*response.AcademicYearResponse, error) {
	ay, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if ay == nil {
		return nil, apperrors.NewNotFoundError("academic year not found")
	}

	layout := "2006-01-02"

	if req.Name != "" {
		ay.Name = req.Name
	}

	// Cek apakah string tidak kosong sebelum parse
	if req.StartDate != "" {
		parsedStart, err := time.Parse(layout, req.StartDate)
		if err != nil {
			return nil, apperrors.NewBadRequestError("invalid start_date format")
		}
		ay.StartDate = parsedStart
	}

	if req.EndDate != "" {
		parsedEnd, err := time.Parse(layout, req.EndDate)
		if err != nil {
			return nil, apperrors.NewBadRequestError("invalid end_date format")
		}
		ay.EndDate = parsedEnd
	}

	if err := s.repo.Update(ay); err != nil {
		return nil, err
	}

	return s.toResponse(ay), nil
}

func (s *academicYearService) Delete(id string) error {
	ay, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if ay == nil {
		return apperrors.NewNotFoundError("academic year not found")
	}

	// Opsional: Cegah hapus jika status ACTIVE
	if ay.Status == "ACTIVE" {
		return apperrors.NewBadRequestError("cannot delete active academic year")
	}

	return s.repo.Delete(id)
}

func (s *academicYearService) Activate(id string) error {
	// 1. Cek existensi
	ay, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if ay == nil {
		return apperrors.NewNotFoundError("academic year not found")
	}

	// 2. Mulai Transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// A. Set semua ke INACTIVE
		if err := s.repo.ResetAllStatus(tx); err != nil {
			return err
		}

		// B. Set yang dipilih ke ACTIVE
		if err := s.repo.UpdateStatus(tx, id, "ACTIVE"); err != nil {
			return err
		}

		return nil
	})
}
