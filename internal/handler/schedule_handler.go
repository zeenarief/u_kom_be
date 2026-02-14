package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

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
		HandleError(c, err)
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
		HandleError(c, err)
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
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Schedules retrieved", res)
}

func (h *ScheduleHandler) GetByTeachingAssignment(c *gin.Context) {
	taID := c.Param("id")
	if taID == "" {
		BadRequestError(c, "teaching_assignment_id required", nil)
		return
	}
	res, err := h.service.GetByTeachingAssignment(taID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Schedules retrieved", res)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	// 	id := c.Param("id")
	if err := h.service.Delete(c.Param("id")); err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Schedule removed", nil)
}

func (h *ScheduleHandler) GetTodaySchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		UnauthorizedError(c, "User ID not found in context")
		return
	}

	res, err := h.service.GetTodaySchedule(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Today's schedule retrieved", res)
}
