package converter

import (
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/model/response"
)

func ToUserWithRoleResponse(user *domain.User) *response.UserWithRoleResponse {
	var roles []response.RoleListResponse
	for _, r := range user.Roles {
		roles = append(roles, response.RoleListResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		})
	}
	return &response.UserWithRoleResponse{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserWithRolesResponseAndPermissions(user *domain.User) *response.UserWithRolesResponseAndPermissions {
	var roles []response.RoleListResponse
	for _, r := range user.Roles {
		roles = append(roles, response.RoleListResponse{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
		})
	}
	return &response.UserWithRolesResponseAndPermissions{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		Email:       user.Email,
		Roles:       roles,
		Permissions: user.GetPermissions(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
