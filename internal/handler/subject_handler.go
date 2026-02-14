package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type SubjectHandler struct {
	service service.SubjectService
}

func NewSubjectHandler(service service.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

func (h *SubjectHandler) Create(c *gin.Context) {
	var req request.SubjectCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Create(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Subject created successfully", res)
}

func (h *SubjectHandler) FindAll(c *gin.Context) {
	searchQuery := c.Query("q")
	res, err := h.service.FindAll(searchQuery)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Subjects retrieved successfully", res)
}

func (h *SubjectHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.FindByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Subject detail retrieved successfully", res)
}

func (h *SubjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.SubjectUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.service.Update(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Subject updated successfully", res)
}

func (h *SubjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Subject deleted successfully", nil)
}
