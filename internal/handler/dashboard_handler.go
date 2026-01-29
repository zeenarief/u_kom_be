package handler

import (
	"github.com/gin-gonic/gin"
	"u_kom_be/internal/service"
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
		InternalServerError(c, err.Error())
		return
	}
	SuccessResponse(c, "Dashboard stats retrieved", stats)
}
