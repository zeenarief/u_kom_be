package response

import "time"

type UserWithRoleResponse struct {
	ID             string          `json:"id"`
	Username       string          `json:"username"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Roles          []string        `json:"roles"`
	ProfileContext *ProfileContext `json:"profile_context,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type ProfileContext struct {
	Type     string `json:"type"`      // "student", "employee", "parent", "admin"
	EntityID string `json:"entity_id"` // ID from students/employees/parents table
}

type UserWithRolesResponseAndPermissions struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Struct untuk info user yang ringkas
type UserLinkedResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
