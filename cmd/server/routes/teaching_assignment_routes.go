package routes

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterTeachingAssignmentRoutes(router *gin.RouterGroup, h *handler.TeachingAssignmentHandler, authService service.AuthService) {
	group := router.Group("/assignments")
	group.Use(middleware.AuthMiddleware(authService))
	{
		// Kurikulum / Admin mengelola assignment
		group.POST("", middleware.PermissionMiddleware("assignments.manage", authService), h.Create)
		group.DELETE("/:id", middleware.PermissionMiddleware("assignments.manage", authService), h.Delete)

		// Read Access (Guru mungkin perlu lihat jadwalnya sendiri)
		// Kita bisa buat permission khusus atau buka untuk authenticated users
		group.GET("/by-class", h.GetByClassroom)
		group.GET("/by-teacher", h.GetByTeacher)
	}
}
