package handler

import (
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

type GradeHandler struct {
	service service.GradeService
}

func NewGradeHandler(service service.GradeService) *GradeHandler {
	return &GradeHandler{service: service}
}

func (h *GradeHandler) CreateAssessment(c *gin.Context) {
	var req request.AssessmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.CreateAssessment(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Assessment created successfully", res)
}

func (h *GradeHandler) UpdateAssessment(c *gin.Context) {
	id := c.Param("id")
	var req request.AssessmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.UpdateAssessment(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Assessment updated successfully", res)
}

func (h *GradeHandler) GetAssessmentsByTeachingAssignment(c *gin.Context) {
	teachingAssignmentID := c.Param("teachingAssignmentID")
	if teachingAssignmentID == "" {
		BadRequestError(c, "teachingAssignmentID is required", nil)
		return
	}

	res, err := h.service.GetAssessmentsByTeachingAssignment(teachingAssignmentID)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Assessments retrieved successfully", res)
}

func (h *GradeHandler) GetAssessmentDetail(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.GetAssessmentDetail(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Assessment detail retrieved successfully", res)
}

func (h *GradeHandler) SubmitScores(c *gin.Context) {
	var req request.BulkScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.service.SubmitScores(req); err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Scores submitted successfully", nil)
}
