package routes

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterSubjectRoutes(router *gin.RouterGroup, h *handler.SubjectHandler, authService service.AuthService) {
	group := router.Group("/subjects")
	group.Use(middleware.AuthMiddleware(authService))
	{
		// Semua user terautentikasi boleh melihat list mapel
		group.GET("", h.FindAll)
		group.GET("/:id", h.FindByID)

		// Hanya user dengan permission yang boleh mengelola (Admin/Kurikulum)
		group.POST("", middleware.PermissionMiddleware("subjects.manage", authService), h.Create)
		group.PUT("/:id", middleware.PermissionMiddleware("subjects.manage", authService), h.Update)
		group.DELETE("/:id", middleware.PermissionMiddleware("subjects.manage", authService), h.Delete)
	}
}
