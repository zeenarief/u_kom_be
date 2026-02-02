package handler

import (
	"strings"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service service.ScheduleService
}

func NewScheduleHandler(service service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var req request.ScheduleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "conflict") {
			BadRequestError(c, "Schedule Conflict", err.Error())
			return
		}
		InternalServerError(c, err.Error())
		return
	}
	CreatedResponse(c, "Schedule created successfully", res)
}

func (h *ScheduleHandler) GetByClassroom(c *gin.Context) {
	classID := c.Query("classroom_id")
	if classID == "" {
		BadRequestError(c, "classroom_id required", nil)
		return
	}
	res, err := h.service.GetByClassroom(classID)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Schedules retrieved", res)
}

func (h *ScheduleHandler) GetByTeacher(c *gin.Context) {
	teacherID := c.Query("teacher_id")
	if teacherID == "" {
		BadRequestError(c, "teacher_id required", nil)
		return
	}
	res, err := h.service.GetByTeacher(teacherID)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Schedules retrieved", res)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Schedule removed", nil)
}
