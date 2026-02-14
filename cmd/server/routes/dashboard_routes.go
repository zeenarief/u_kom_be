package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(router *gin.RouterGroup, handler *handler.DashboardHandler, authService service.AuthService) {
	// Biasanya dashboard read-only bisa diakses semua user login,
	// atau batasi permission "dashboard.read" jika perlu
	router.GET("/dashboard/stats",
		middleware.AuthMiddleware(authService),
		handler.GetStats)

	router.GET("/dashboard/teacher/stats",
		middleware.AuthMiddleware(authService),
		handler.GetTeacherStats)
}
