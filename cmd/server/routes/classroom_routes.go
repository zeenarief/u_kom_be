package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterClassroomRoutes(router *gin.RouterGroup, h *handler.ClassroomHandler, authService service.AuthService) {
	group := router.Group("/classrooms")
	group.Use(middleware.AuthMiddleware(authService))
	{
		group.GET("", h.FindAll)
		group.GET("/:id", h.FindByID)

		// Manage Classroom
		group.POST("", middleware.PermissionMiddleware("classrooms.manage", authService), h.Create)
		group.PUT("/:id", middleware.PermissionMiddleware("classrooms.manage", authService), h.Update)
		group.DELETE("/:id", middleware.PermissionMiddleware("classrooms.manage", authService), h.Delete)

		// Manage Students in Classroom
		group.POST("/:id/students", middleware.PermissionMiddleware("classrooms.manage_students", authService), h.AddStudents)
		group.DELETE("/:id/students/:studentID", middleware.PermissionMiddleware("classrooms.manage_students", authService), h.RemoveStudent)
	}
}
