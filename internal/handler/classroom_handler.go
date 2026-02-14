package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type ClassroomHandler struct {
	service service.ClassroomService
}

func NewClassroomHandler(service service.ClassroomService) *ClassroomHandler {
	return &ClassroomHandler{service: service}
}

func (h *ClassroomHandler) Create(c *gin.Context) {
	var req request.ClassroomCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Classroom created successfully", res)
}

func (h *ClassroomHandler) FindAll(c *gin.Context) {
	// Filter by academic_year_id via query param
	ayID := c.Query("academic_year_id")

	res, err := h.service.FindAll(ayID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Classrooms retrieved successfully", res)
}

func (h *ClassroomHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.FindByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Classroom detail retrieved successfully", res)
}

func (h *ClassroomHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.ClassroomUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Update(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Classroom updated successfully", res)
}

func (h *ClassroomHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Classroom deleted successfully", nil)
}

func (h *ClassroomHandler) AddStudents(c *gin.Context) {
	id := c.Param("id")
	var req request.AddStudentsToClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	if err := h.service.AddStudents(id, req); err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Students added to classroom successfully", nil)
}

func (h *ClassroomHandler) RemoveStudent(c *gin.Context) {
	id := c.Param("id")
	studentID := c.Param("studentID")

	if err := h.service.RemoveStudent(id, studentID); err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Student removed from classroom successfully", nil)
}
