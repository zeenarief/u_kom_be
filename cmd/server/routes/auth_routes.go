package routes

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler, authService service.AuthService) {
	// Public auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected auth routes
	protectedAuth := router.Group("/auth")
	protectedAuth.Use(middleware.AuthMiddleware(authService))
	{
		protectedAuth.POST("/logout", authHandler.Logout)
	}
}
