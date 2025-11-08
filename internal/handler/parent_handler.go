package handler

import (
	"belajar-golang/internal/model/request"
	"belajar-golang/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type ParentHandler struct {
	parentService service.ParentService
}

func NewParentHandler(parentService service.ParentService) *ParentHandler {
	return &ParentHandler{parentService: parentService}
}

func (h *ParentHandler) CreateParent(c *gin.Context) {
	var req request.ParentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	parent, err := h.parentService.CreateParent(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Parent creation failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	CreatedResponse(c, "Parent created successfully", parent)
}

func (h *ParentHandler) GetAllParents(c *gin.Context) {
	parents, err := h.parentService.GetAllParents()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Parents retrieved successfully", parents)
}

func (h *ParentHandler) GetParentByID(c *gin.Context) {
	id := c.Param("id")

	parent, err := h.parentService.GetParentByID(id)
	if err != nil {
		if err.Error() == "parent not found" {
			NotFoundError(c, "Parent not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Parent retrieved successfully", parent)
}

func (h *ParentHandler) UpdateParent(c *gin.Context) {
	id := c.Param("id")

	var req request.ParentUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	parent, err := h.parentService.UpdateParent(id, req)
	if err != nil {
		if err.Error() == "parent not found" {
			NotFoundError(c, "Parent not found")
		} else if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Parent update failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Parent updated successfully", parent)
}

func (h *ParentHandler) DeleteParent(c *gin.Context) {
	id := c.Param("id")

	err := h.parentService.DeleteParent(id)
	if err != nil {
		if err.Error() == "parent not found" {
			NotFoundError(c, "Parent not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Parent deleted successfully", nil)
}

// LinkUser menangani POST /parents/:id/link-user
func (h *ParentHandler) LinkUser(c *gin.Context) {
	parentID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload (missing 'user_id')", err.Error())
		return
	}

	err := h.parentService.LinkUser(parentID, req.UserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			NotFoundError(c, err.Error())
		} else if strings.Contains(err.Error(), "already linked") {
			BadRequestError(c, "Link failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Parent linked to user successfully", nil)
}

// UnlinkUser menangani DELETE /parents/:id/unlink-user
func (h *ParentHandler) UnlinkUser(c *gin.Context) {
	parentID := c.Param("id")

	err := h.parentService.UnlinkUser(parentID)
	if err != nil {
		if err.Error() == "parent not found" {
			NotFoundError(c, "Parent not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Parent unlinked from user successfully", nil)
}
