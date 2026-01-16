package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(router *gin.RouterGroup, handler *handler.DashboardHandler, authService service.AuthService) {
	// Biasanya dashboard read-only bisa diakses semua user login,
	// atau batasi permission "dashboard.read" jika perlu
	router.GET("/dashboard/stats",
		middleware.AuthMiddleware(authService),
		handler.GetStats)
}
