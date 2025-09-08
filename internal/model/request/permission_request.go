package request

type PermissionCreateRequest struct {
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description,omitempty"`
}

type PermissionUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
