package response

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserWithRolesResponse struct {
	ID          string         `json:"id"`
	Username    string         `json:"username"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Roles       []RoleResponse `json:"roles"`
	Permissions []string       `json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
