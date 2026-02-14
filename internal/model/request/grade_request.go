package request

import "smart_school_be/internal/utils"

type AssessmentCreateRequest struct {
	TeachingAssignmentID string     `json:"teaching_assignment_id" binding:"required"`
	Title                string     `json:"title" binding:"required"`
	Type                 string     `json:"type" binding:"required"` // ASSIGNMENT, MID_EXAM, FINAL_EXAM, QUIZ
	MaxScore             int        `json:"max_score"`
	Date                 utils.Date `json:"date" binding:"required"`
	Description          string     `json:"description"`
}

type StudentScoreRequest struct {
	StudentID string  `json:"student_id" binding:"required"`
	Score     float64 `json:"score"`
	Feedback  string  `json:"feedback"`
}

type BulkScoreRequest struct {
	AssessmentID string                `json:"assessment_id" binding:"required"`
	Scores       []StudentScoreRequest `json:"scores" binding:"required"`
}
