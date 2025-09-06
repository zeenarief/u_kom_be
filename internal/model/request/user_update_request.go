package request

type UserUpdateRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty" binding:"omitempty,email"`
}
