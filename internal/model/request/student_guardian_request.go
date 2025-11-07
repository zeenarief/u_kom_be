package request

// StudentSetGuardianRequest adalah DTO untuk mengatur wali (polymorphic)
type StudentSetGuardianRequest struct {
	// GuardianID adalah UUID dari tabel 'parents' atau 'guardians'
	GuardianID string `json:"guardian_id" binding:"required"`

	// GuardianType adalah nama tabel: 'parent' atau 'guardian'
	GuardianType string `json:"guardian_type" binding:"required,oneof=parent guardian"`
}
