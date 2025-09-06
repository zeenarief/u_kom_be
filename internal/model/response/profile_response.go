package response

import "time"

type ProfileResponse struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ProfileComplete bool      `json:"profile_complete"`
	AvatarURL       string    `json:"avatar_url,omitempty"`
}
