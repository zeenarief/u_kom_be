package routes

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterAttendanceRoutes(router *gin.RouterGroup, h *handler.AttendanceHandler, authService service.AuthService) {
	group := router.Group("/attendances")
	group.Use(middleware.AuthMiddleware(authService))
	{
		// Submit Absen
		// Format: middleware(permission, service), handler
		group.POST("",
			middleware.PermissionMiddleware("attendance.submit", authService),
			h.Submit,
		)

		// Get Detail
		// Sebaiknya juga dilindungi permission read
		group.GET("/:id",
			// middleware.PermissionMiddleware("attendance.read", authService), // Opsional
			h.GetDetail,
		)

		// Get History
		group.GET("/history",
			// middleware.PermissionMiddleware("attendance.read", authService), // Opsional
			h.GetHistory,
		)

		group.GET("/check", h.CheckSession) // GET /api/v1/attendances/check?schedule_id=...&date=...
	}
}
