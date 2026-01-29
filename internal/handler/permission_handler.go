package handler

import (
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/service"

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
		if err.Error() == "permission already exists" {
			BadRequestError(c, "Permission creation failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
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
		if err.Error() == "permission not found" {
			NotFoundError(c, "Permission not found")
		} else {
			InternalServerError(c, err.Error())
		}
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
		if err.Error() == "permission not found" {
			NotFoundError(c, "Permission not found")
		} else if err.Error() == "permission name already exists" {
			BadRequestError(c, "Permission update failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Permission updated successfully", permission)
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	id := c.Param("id")

	err := h.permissionService.DeletePermission(id)
	if err != nil {
		if err.Error() == "permission not found" {
			NotFoundError(c, "Permission not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Permission deleted successfully", nil)
}
