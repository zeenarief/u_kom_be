package response

import "time"

type UserWithRoleResponse struct {
	ID        string             `json:"id"`
	Username  string             `json:"username"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Roles     []RoleListResponse `json:"roles"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type UserWithRolesResponseAndPermissions struct {
	ID          string             `json:"id"`
	Username    string             `json:"username"`
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	Roles       []RoleListResponse `json:"roles"`
	Permissions []string           `json:"permissions"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}
