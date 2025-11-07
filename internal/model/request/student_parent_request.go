package request

// ParentRelationshipRequest adalah DTO internal untuk sync
type ParentRelationshipRequest struct {
	ParentID         string `json:"parent_id" binding:"required"`
	RelationshipType string `json:"relationship_type" binding:"required"` // Cth: 'FATHER', 'MOTHER'
}

// StudentSyncParentsRequest adalah DTO utama untuk request
type StudentSyncParentsRequest struct {
	Parents []ParentRelationshipRequest `json:"parents" binding:"required,dive"`
}
