package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req request.PermissionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	permission, err := h.permissionService.CreatePermission(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "Permission created successfully", permission)
}

func (h *PermissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.permissionService.GetAllPermissions()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Permissions retrieved successfully", permissions)
}

func (h *PermissionHandler) GetPermissionByID(c *gin.Context) {
	id := c.Param("id")

	permission, err := h.permissionService.GetPermissionByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Permission retrieved successfully", permission)
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")

	var req request.PermissionUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	permission, err := h.permissionService.UpdatePermission(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Permission updated successfully", permission)
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")

	err := h.permissionService.DeletePermission(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "Permission deleted successfully", nil)
}
