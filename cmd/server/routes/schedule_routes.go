package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterScheduleRoutes(router *gin.RouterGroup, h *handler.ScheduleHandler, authService service.AuthService) {
	group := router.Group("/schedules")
	group.Use(middleware.AuthMiddleware(authService))
	{
		// Kurikulum / Admin mengatur jadwal
		group.POST("", middleware.PermissionMiddleware("schedules.manage", authService), h.Create)
		group.DELETE("/:id", middleware.PermissionMiddleware("schedules.manage", authService), h.Delete)

		// Read Access
		group.GET("/by-class", h.GetByClassroom)
		group.GET("/by-teacher", h.GetByTeacher)
		group.GET("/by-teaching-assignment/:id", h.GetByTeachingAssignment)
		group.GET("/teacher/today", h.GetTodaySchedule)
	}
}
