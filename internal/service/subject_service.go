package service

import (
	"errors"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
)

type SubjectService interface {
	Create(req request.SubjectCreateRequest) (*response.SubjectResponse, error)
	FindAll() ([]response.SubjectResponse, error)
	FindByID(id string) (*response.SubjectResponse, error)
	Update(id string, req request.SubjectUpdateRequest) (*response.SubjectResponse, error)
	Delete(id string) error
}

type subjectService struct {
	repo repository.SubjectRepository
}

func NewSubjectService(repo repository.SubjectRepository) SubjectService {
	return &subjectService{repo: repo}
}

func (s *subjectService) toResponse(sub *domain.Subject) *response.SubjectResponse {
	return &response.SubjectResponse{
		ID:          sub.ID,
		Code:        sub.Code,
		Name:        sub.Name,
		Type:        sub.Type,
		Description: sub.Description,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}

func (s *subjectService) Create(req request.SubjectCreateRequest) (*response.SubjectResponse, error) {
	// Validasi Duplikat Kode
	existing, _ := s.repo.FindByCode(req.Code)
	if existing != nil {
		return nil, errors.New("subject code already exists")
	}

	subject := &domain.Subject{
		Code:        req.Code,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
	}

	if err := s.repo.Create(subject); err != nil {
		return nil, err
	}

	return s.toResponse(subject), nil
}

func (s *subjectService) FindAll() ([]response.SubjectResponse, error) {
	subjects, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.SubjectResponse
	for _, sub := range subjects {
		responses = append(responses, *s.toResponse(&sub))
	}
	return responses, nil
}

func (s *subjectService) FindByID(id string) (*response.SubjectResponse, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subject not found")
	}
	return s.toResponse(sub), nil
}

func (s *subjectService) Update(id string, req request.SubjectUpdateRequest) (*response.SubjectResponse, error) {
	sub, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, errors.New("subject not found")
	}

	// Validasi Code jika berubah
	if req.Code != "" && req.Code != sub.Code {
		existing, _ := s.repo.FindByCode(req.Code)
		if existing != nil {
			return nil, errors.New("subject code already exists")
		}
		sub.Code = req.Code
	}

	if req.Name != "" {
		sub.Name = req.Name
	}
	if req.Type != "" {
		sub.Type = req.Type
	}
	sub.Description = req.Description // Bisa dikosongkan

	if err := s.repo.Update(sub); err != nil {
		return nil, err
	}

	return s.toResponse(sub), nil
}

func (s *subjectService) Delete(id string) error {
	// Optional: Cek apakah mapel ini sedang dipakai di jadwal/kelas (Future improvement)
	return s.repo.Delete(id)
}
