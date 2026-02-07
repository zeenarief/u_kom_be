package handler

import (
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: s}
}

// CreateEmployee menangani POST /employees
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req request.EmployeeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	employee, err := h.employeeService.CreateEmployee(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Employee created successfully", employee)
}

// GetAllEmployees menangani GET /employees
func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	searchQuery := c.Query("q")
	employees, err := h.employeeService.GetAllEmployees(searchQuery)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Employees retrieved successfully", employees)
}

// GetEmployeeByID menangani GET /employees/:id
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	id := c.Param("id")

	employee, err := h.employeeService.GetEmployeeByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Employee retrieved successfully", employee)
}

// UpdateEmployee menangani PUT /employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")

	var req request.EmployeeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	employee, err := h.employeeService.UpdateEmployee(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Employee updated successfully", employee)
}

// DeleteEmployee menangani DELETE /employees/:id
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	err := h.employeeService.DeleteEmployee(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Employee deleted successfully", nil)
}

// --- Handler untuk Penautan Akun ---

// LinkUser menangani POST /employees/:id/link-user
func (h *EmployeeHandler) LinkUser(c *gin.Context) {
	employeeID := c.Param("id")

	// DTO Request didefinisikan inline
	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload (missing 'user_id')", err.Error())
		return
	}

	err := h.employeeService.LinkUser(employeeID, req.UserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Employee linked to user successfully", nil)
}

// UnlinkUser menangani DELETE /employees/:id/unlink-user
func (h *EmployeeHandler) UnlinkUser(c *gin.Context) {
	employeeID := c.Param("id")

	err := h.employeeService.UnlinkUser(employeeID)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Employee unlinked from user successfully", nil)
}
