package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all application routes
func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authService service.AuthService,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	studentHandler *handler.StudentHandler,
) {
	// API v1 group
	apiV1 := router.Group("/api/v1")

	// Register all routes
	RegisterAuthRoutes(apiV1, authHandler, authService)
	RegisterUserRoutes(apiV1, userHandler, authService)
	RegisterRoleRoutes(apiV1, roleHandler, authService)
	RegisterPermissionRoutes(apiV1, permissionHandler, authService)
	RegisterStudentRoutes(apiV1, studentHandler, authService)

	// Health check route
	apiV1.GET("/health", healthCheck)
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	handler.SuccessResponse(c, "Server is healthy and running", gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}
