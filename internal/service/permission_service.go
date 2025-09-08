package service

import (
	"belajar-golang/internal/model/domain"
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/model/response"
	"belajar-golang/internal/repository"
	"errors"
)

type PermissionService interface {
	CreatePermission(req request.PermissionCreateRequest) (*response.PermissionResponse, error)
	GetPermissionByID(id string) (*response.PermissionResponse, error)
	GetPermissionByName(name string) (*response.PermissionResponse, error)
	GetAllPermissions() ([]response.PermissionResponse, error)
	UpdatePermission(id string, req request.PermissionUpdateRequest) (*response.PermissionResponse, error)
	DeletePermission(id string) error
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{permissionRepo: permissionRepo}
}

func (s *permissionService) CreatePermission(req request.PermissionCreateRequest) (*response.PermissionResponse, error) {
	// Check if permission already exists
	existingPerm, _ := s.permissionRepo.FindByName(req.Name)
	if existingPerm != nil {
		return nil, errors.New("permission already exists")
	}

	permission := &domain.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.permissionRepo.Create(permission); err != nil {
		return nil, err
	}

	return s.convertToResponse(permission), nil
}

func (s *permissionService) GetPermissionByID(id string) (*response.PermissionResponse, error) {
	permission, err := s.permissionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.New("permission not found")
	}
	return s.convertToResponse(permission), nil
}

func (s *permissionService) GetPermissionByName(name string) (*response.PermissionResponse, error) {
	permission, err := s.permissionRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.New("permission not found")
	}
	return s.convertToResponse(permission), nil
}

func (s *permissionService) GetAllPermissions() ([]response.PermissionResponse, error) {
	permissions, err := s.permissionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.PermissionResponse
	for _, perm := range permissions {
		responses = append(responses, *s.convertToResponse(&perm))
	}

	return responses, nil
}

func (s *permissionService) UpdatePermission(id string, req request.PermissionUpdateRequest) (*response.PermissionResponse, error) {
	permission, err := s.permissionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if permission == nil {
		return nil, errors.New("permission not found")
	}

	if req.Name != "" {
		// Check if new name already exists
		if existingPerm, _ := s.permissionRepo.FindByName(req.Name); existingPerm != nil && existingPerm.ID != id {
			return nil, errors.New("permission name already exists")
		}
		permission.Name = req.Name
	}

	if req.Description != "" {
		permission.Description = req.Description
	}

	if err := s.permissionRepo.Update(permission); err != nil {
		return nil, err
	}

	return s.convertToResponse(permission), nil
}

func (s *permissionService) DeletePermission(id string) error {
	// Check if permission exists
	permission, err := s.permissionRepo.FindByID(id)
	if err != nil {
		return err
	}
	if permission == nil {
		return errors.New("permission not found")
	}

	return s.permissionRepo.Delete(id)
}

func (s *permissionService) convertToResponse(permission *domain.Permission) *response.PermissionResponse {
	return &response.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		CreatedAt:   permission.CreatedAt,
		UpdatedAt:   permission.UpdatedAt,
	}
}
