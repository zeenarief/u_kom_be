package request

type UserCreateRequest struct {
	Username string   `json:"username" binding:"required,min=3"`
	Name     string   `json:"name" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	RoleIDs  []string `json:"role_ids,omitempty"`
}

type UserUpdateRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty" binding:"omitempty,email"`
	RoleIDs []string `json:"role_ids,omitempty"`

	Password string   `json:"password,omitempty" binding:"omitempty,min=6"`

}
