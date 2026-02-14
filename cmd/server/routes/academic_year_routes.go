package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterAcademicYearRoutes(router *gin.RouterGroup, handler *handler.AcademicYearHandler, authService service.AuthService) {
	group := router.Group("/academic-years")
	group.Use(middleware.AuthMiddleware(authService))
	{
		// Read (Mungkin bisa diakses semua user yang login)
		group.GET("", handler.FindAll)
		group.GET("/:id", handler.FindByID)

		// Write (Hanya Admin/Operator) - Sesuaikan permission stringnya
		group.POST("", middleware.PermissionMiddleware("academic_years.manage", authService), handler.Create)
		group.PUT("/:id", middleware.PermissionMiddleware("academic_years.manage", authService), handler.Update)
		group.DELETE("/:id", middleware.PermissionMiddleware("academic_years.manage", authService), handler.Delete)

		// Activate
		group.PATCH("/:id/activate", middleware.PermissionMiddleware("academic_years.manage", authService), handler.Activate)
	}
}
