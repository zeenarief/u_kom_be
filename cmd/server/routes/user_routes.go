package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user management routes
func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler, authService service.AuthService) {
	// Public user routes (registration might be here or in auth)
	// router.POST("/users", userHandler.CreateUser) // Moved to auth/register

	// Protected user routes
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User management
		protected.GET("/users", userHandler.GetAllUsers)
		protected.POST("/users", userHandler.CreateUser)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
		protected.POST("/users/:id/change-password", userHandler.ChangePassword)

		// User profile
		protected.GET("/profile", userHandler.GetProfile)
	}
}
