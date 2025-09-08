package request

type RoleCreateRequest struct {
	Name        string   `json:"name" binding:"required,min=3"`
	Description string   `json:"description,omitempty"`
	IsDefault   bool     `json:"is_default,omitempty"`
	Permissions []string `json:"permissions,omitempty"` // Array of permission names
}

type RoleUpdateRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	IsDefault   *bool    `json:"is_default,omitempty"` // Use pointer untuk differentiate between false and not provided
	Permissions []string `json:"permissions,omitempty"`
}

type AssignPermissionsRequest struct {
	PermissionNames []string `json:"permissions" binding:"required"`
}
