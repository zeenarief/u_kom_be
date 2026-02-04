package handler

import (
	"strings"
	"u_kom_be/internal/model/request"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

type GuardianHandler struct {
	guardianService service.GuardianService
}

func NewGuardianHandler(guardianService service.GuardianService) *GuardianHandler {
	return &GuardianHandler{guardianService: guardianService}
}

func (h *GuardianHandler) CreateGuardian(c *gin.Context) {
	var req request.GuardianCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	guardian, err := h.guardianService.CreateGuardian(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Guardian creation failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	CreatedResponse(c, "Guardian created successfully", guardian)
}

func (h *GuardianHandler) GetAllGuardians(c *gin.Context) {
	searchQuery := c.Query("q")
	guardians, err := h.guardianService.GetAllGuardians(searchQuery)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	SuccessResponse(c, "Guardians retrieved successfully", guardians)
}

func (h *GuardianHandler) GetGuardianByID(c *gin.Context) {
	id := c.Param("id")

	guardian, err := h.guardianService.GetGuardianByID(id)
	if err != nil {
		if err.Error() == "guardian not found" {
			NotFoundError(c, "Guardian not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Guardian retrieved successfully", guardian)
}

func (h *GuardianHandler) UpdateGuardian(c *gin.Context) {
	id := c.Param("id")

	var req request.GuardianUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload", err.Error())
		return
	}

	guardian, err := h.guardianService.UpdateGuardian(id, req)
	if err != nil {
		if err.Error() == "guardian not found" {
			NotFoundError(c, "Guardian not found")
		} else if strings.Contains(err.Error(), "already exists") {
			BadRequestError(c, "Guardian update failed", err.Error())
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Guardian updated successfully", guardian)
}

func (h *GuardianHandler) DeleteGuardian(c *gin.Context) {
	id := c.Param("id")

	err := h.guardianService.DeleteGuardian(id)
	if err != nil {
		if err.Error() == "guardian not found" {
			NotFoundError(c, "Guardian not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Guardian deleted successfully", nil)
}

// LinkUser menangani POST /guardians/:id/link-user
func (h *GuardianHandler) LinkUser(c *gin.Context) {
	guardianID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, "Invalid request payload (missing 'user_id')", err.Error())
		return
	}

	err := h.guardianService.LinkUser(guardianID, req.UserID)
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

	SuccessResponse(c, "Guardian linked to user successfully", nil)
}

// UnlinkUser menangani DELETE /guardians/:id/unlink-user
func (h *GuardianHandler) UnlinkUser(c *gin.Context) {
	guardianID := c.Param("id")

	err := h.guardianService.UnlinkUser(guardianID)
	if err != nil {
		if err.Error() == "guardian not found" {
			NotFoundError(c, "Guardian not found")
		} else {
			InternalServerError(c, err.Error())
		}
		return
	}

	SuccessResponse(c, "Guardian unlinked from user successfully", nil)
}
