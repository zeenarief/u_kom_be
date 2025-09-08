package response

import "time"

type RoleDetailResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description,omitempty"`
	IsDefault   bool                 `json:"is_default"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type RoleListResponse struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
