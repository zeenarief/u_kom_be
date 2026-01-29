package routes

import (
	"github.com/gin-gonic/gin"
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"
)

func RegisterDashboardRoutes(router *gin.RouterGroup, handler *handler.DashboardHandler, authService service.AuthService) {
	// Biasanya dashboard read-only bisa diakses semua user login,
	// atau batasi permission "dashboard.read" jika perlu
	router.GET("/dashboard/stats",
		middleware.AuthMiddleware(authService),
		handler.GetStats)
}
