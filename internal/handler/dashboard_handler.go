package handler

import (
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardService.GetStats()
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Dashboard stats retrieved", stats)
}

func (h *DashboardHandler) GetTeacherStats(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		UnauthorizedError(c, "User ID not found in context")
		return
	}

	stats, err := h.dashboardService.GetTeacherStats(userID)
	if err != nil {
		HandleError(c, err)
		return
	}
	SuccessResponse(c, "Teacher stats retrieved", stats)
}
