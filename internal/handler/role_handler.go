package handler

import (
	"smart_school_be/internal/model/request"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService service.RoleService
}

func NewRoleHandler(roleService service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req request.RoleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	role, err := h.roleService.CreateRole(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	CreatedResponse(c, "RoleIDs created successfully", role)
}

func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleService.GetAllRoles()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Roles retrieved successfully", roles)
}

func (h *RoleHandler) GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "RoleIDs retrieved successfully", role)
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var req request.RoleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	role, err := h.roleService.UpdateRole(id, req)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "RoleIDs updated successfully", role)
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")

	err := h.roleService.DeleteRole(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "RoleIDs deleted successfully", nil)
}

func (h *RoleHandler) SyncRolePermissions(c *gin.Context) {
	roleID := c.Param("id")

	var req request.AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	err := h.roleService.SyncRolePermissions(roleID, req.PermissionNames)
	if err != nil {
		HandleError(c, err)
		return
	}

	SuccessResponse(c, "RoleIDs permissions synced successfully", nil)
}
