package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type ViolationHandler interface {
	// Category
	CreateCategory(c *gin.Context)
	GetCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)

	// Type
	CreateType(c *gin.Context)
	GetTypes(c *gin.Context)
	UpdateType(c *gin.Context)
	DeleteType(c *gin.Context)

	// Student Violation
	RecordViolation(c *gin.Context)
	GetStudentViolations(c *gin.Context)
	DeleteViolation(c *gin.Context)
	GetAllViolations(c *gin.Context)
}

type violationHandler struct {
	violationService service.ViolationService
}

func NewViolationHandler(violationService service.ViolationService) ViolationHandler {
	return &violationHandler{violationService: violationService}
}

// Category
func (h *violationHandler) CreateCategory(c *gin.Context) {
	var req request.CreateViolationCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error(), nil)
		return
	}

	if err := h.violationService.CreateCategory(req); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	CreatedResponse(c, "Violation category created successfully", nil)
}

func (h *violationHandler) GetCategories(c *gin.Context) {
	categories, err := h.violationService.GetCategories()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Data retrieved successfully", categories)
}

func (h *violationHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateViolationCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error(), nil)
		return
	}

	if err := h.violationService.UpdateCategory(id, req); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Violation category updated successfully", nil)
}

func (h *violationHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.violationService.DeleteCategory(id); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Violation category deleted successfully", nil)
}

// Type
func (h *violationHandler) CreateType(c *gin.Context) {
	var req request.CreateViolationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error(), nil)
		return
	}

	if err := h.violationService.CreateType(req); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	CreatedResponse(c, "Violation type created successfully", nil)
}

func (h *violationHandler) GetTypes(c *gin.Context) {
	categoryID := c.Query("category_id")
	types, err := h.violationService.GetTypes(categoryID)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Data retrieved successfully", types)
}

func (h *violationHandler) UpdateType(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdateViolationTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error(), nil)
		return
	}

	if err := h.violationService.UpdateType(id, req); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Violation type updated successfully", nil)
}

func (h *violationHandler) DeleteType(c *gin.Context) {
	id := c.Param("id")
	if err := h.violationService.DeleteType(id); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Violation type deleted successfully", nil)
}

// Student Violation
func (h *violationHandler) RecordViolation(c *gin.Context) {
	var req request.CreateStudentViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error(), nil)
		return
	}

	if err := h.violationService.RecordViolation(req); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	CreatedResponse(c, "Violation recorded successfully", nil)
}

func (h *violationHandler) GetStudentViolations(c *gin.Context) {
	studentID := c.Param("studentID")
	pagination := request.NewPaginationRequest(c.Query("page"), c.Query("limit"))

	violations, err := h.violationService.GetStudentViolations(studentID, pagination)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Data retrieved successfully", violations)
}

func (h *violationHandler) DeleteViolation(c *gin.Context) {
	id := c.Param("id")
	if err := h.violationService.DeleteViolation(id); err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Violation record deleted successfully", nil)
}

func (h *violationHandler) GetAllViolations(c *gin.Context) {
	filter := c.Query("search")
	pagination := request.NewPaginationRequest(c.Query("page"), c.Query("limit"))

	violations, err := h.violationService.GetAllViolations(filter, pagination)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Data retrieved successfully", violations)
}
