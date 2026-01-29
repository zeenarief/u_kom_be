package service

import (
	"errors"
	"fmt"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/model/response"
	"u_kom_be/internal/repository"
)

type RoleService interface {
	CreateRole(req request.RoleCreateRequest) (*response.RoleDetailResponse, error)
	GetRoleByID(id string) (*response.RoleDetailResponse, error)
	GetRoleByName(name string) (*response.RoleDetailResponse, error)
	GetAllRoles() ([]response.RoleDetailResponse, error)
	UpdateRole(id string, req request.RoleUpdateRequest) (*response.RoleDetailResponse, error)
	DeleteRole(id string) error
	SyncRolePermissions(roleID string, permissionNames []string) error
}

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *roleService) CreateRole(req request.RoleCreateRequest) (*response.RoleDetailResponse, error) {
	// Check if role already exists
	existingRole, _ := s.roleRepo.FindByName(req.Name)
	if existingRole != nil {
		return nil, errors.New("role already exists")
	}

	role := &domain.Role{
		Name:        req.Name,
		Description: req.Description,
		IsDefault:   req.IsDefault,
	}

	// Create role
	if err := s.roleRepo.Create(role); err != nil {
		return nil, err
	}

	// Sync permissions if provided
	if len(req.Permissions) > 0 {
		if err := s.SyncRolePermissions(role.ID, req.Permissions); err != nil {
			// Rollback role creation if permission sync fails
			s.roleRepo.Delete(role.ID)
			return nil, fmt.Errorf("failed to assign permissions: %v", err)
		}
	}

	// Get created role with permissions
	createdRole, err := s.roleRepo.FindByID(role.ID)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(createdRole), nil
}

func (s *roleService) SyncRolePermissions(roleID string, permissionNames []string) error {
	var permissionIDs []string

	for _, permName := range permissionNames {
		permission, err := s.permissionRepo.FindByName(permName)
		if err != nil {
			return fmt.Errorf("permission not found: %s", permName)
		}
		if permission == nil {
			return fmt.Errorf("permission not found: %s", permName)
		}
		permissionIDs = append(permissionIDs, permission.ID)
	}

	return s.roleRepo.SyncPermissions(roleID, permissionIDs)
}

func (s *roleService) GetRoleByID(id string) (*response.RoleDetailResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return s.convertToResponse(role), nil
}

func (s *roleService) GetRoleByName(name string) (*response.RoleDetailResponse, error) {
	role, err := s.roleRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}
	return s.convertToResponse(role), nil
}

func (s *roleService) GetAllRoles() ([]response.RoleDetailResponse, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.RoleDetailResponse
	for _, role := range roles {
		responses = append(responses, *s.convertToResponse(&role))
	}

	return responses, nil
}

func (s *roleService) UpdateRole(id string, req request.RoleUpdateRequest) (*response.RoleDetailResponse, error) {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	// Update fields if provided
	if req.Name != "" {
		// Check if new name already exists
		if existingRole, _ := s.roleRepo.FindByName(req.Name); existingRole != nil && existingRole.ID != id {
			return nil, errors.New("role name already exists")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if req.IsDefault != nil {
		role.IsDefault = *req.IsDefault
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, err
	}

	// Sync permissions if provided
	if req.Permissions != nil {
		if err := s.SyncRolePermissions(id, req.Permissions); err != nil {
			return nil, err
		}
	}

	// Get updated role
	updatedRole, err := s.roleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.convertToResponse(updatedRole), nil
}

func (s *roleService) DeleteRole(id string) error {
	// Prevent deletion of default roles if needed
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Optional: Add business logic to prevent deletion of certain roles
	if role.IsDefault {
		return errors.New("cannot delete default role")
	}

	return s.roleRepo.Delete(id)
}

func (s *roleService) convertToResponse(role *domain.Role) *response.RoleDetailResponse {
	var permissionResponses []response.PermissionResponse
	for _, perm := range role.Permissions {
		permissionResponses = append(permissionResponses, response.PermissionResponse{
			ID:          perm.ID,
			Name:        perm.Name,
			Description: perm.Description,
			CreatedAt:   perm.CreatedAt,
			UpdatedAt:   perm.UpdatedAt,
		})
	}

	return &response.RoleDetailResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsDefault:   role.IsDefault,
		Permissions: permissionResponses,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
