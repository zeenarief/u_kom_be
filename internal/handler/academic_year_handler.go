package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type AcademicYearHandler struct {
	service service.AcademicYearService
}

func NewAcademicYearHandler(service service.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{service: service}
}

func (h *AcademicYearHandler) Create(c *gin.Context) {
	var req request.AcademicYearCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	result, err := h.service.Create(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Academic year created successfully", result)
}

func (h *AcademicYearHandler) FindAll(c *gin.Context) {
	result, err := h.service.FindAll()
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Academic years retrieved successfully", result)
}

func (h *AcademicYearHandler) FindByID(c *gin.Context) {
	id := c.Param("id")
	result, err := h.service.FindByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Academic year retrieved successfully", result)
}

func (h *AcademicYearHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req request.AcademicYearUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	result, err := h.service.Update(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Academic year updated successfully", result)
}

func (h *AcademicYearHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Academic year deleted successfully", nil)
}

func (h *AcademicYearHandler) Activate(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Activate(id); err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Academic year activated successfully", nil)
}
