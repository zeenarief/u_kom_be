package handler

import (
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

type TeachingAssignmentHandler struct {
	service service.TeachingAssignmentService
}

func NewTeachingAssignmentHandler(service service.TeachingAssignmentService) *TeachingAssignmentHandler {
	return &TeachingAssignmentHandler{service: service}
}

func (h *TeachingAssignmentHandler) Create(c *gin.Context) {
	var req request.TeachingAssignmentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		if err.Error() == "this subject already has a teacher in this class" {
			BadRequestError(c, "Assignment conflict", err.Error())
			return
		}
		InternalServerError(c, err.Error())
		return
	}
	CreatedResponse(c, "Teaching assignment created successfully", res)
}

func (h *TeachingAssignmentHandler) GetByClassroom(c *gin.Context) {
	classID := c.Query("classroom_id")
	if classID == "" {
		BadRequestError(c, "classroom_id query parameter is required", nil)
		return
	}

	res, err := h.service.GetByClassroom(classID)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Classroom assignments retrieved", res)
}

func (h *TeachingAssignmentHandler) GetByTeacher(c *gin.Context) {
	teacherID := c.Query("teacher_id")
	if teacherID == "" {
		BadRequestError(c, "teacher_id query parameter is required", nil)
		return
	}

	res, err := h.service.GetByTeacher(teacherID)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Teacher assignments retrieved", res)
}

func (h *TeachingAssignmentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Assignment removed", nil)
}
