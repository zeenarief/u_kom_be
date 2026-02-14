package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler, authService service.AuthService) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User management dengan permission-based access
		protected.GET("/users",
			middleware.PermissionMiddleware("users.read", authService),
			userHandler.GetAllUsers)

		protected.POST("/users",
			middleware.PermissionMiddleware("users.create", authService),
			userHandler.CreateUser)

		protected.GET("/users/:id",
			middleware.PermissionMiddleware("users.read", authService),
			userHandler.GetUserByID)

		protected.PUT("/users/:id",
			middleware.PermissionMiddleware("users.update", authService),
			userHandler.UpdateUser)

		protected.DELETE("/users/:id",
			middleware.PermissionMiddleware("users.delete", authService),
			userHandler.DeleteUser)

		// RoleIDs and permission management
		protected.POST("/users/:id/sync-roles",
			middleware.PermissionMiddleware("users.manage_roles", authService),
			userHandler.SyncUserRoles)

		protected.POST("/users/:id/sync-permissions",
			middleware.PermissionMiddleware("users.manage_permissions", authService),
			userHandler.SyncUserPermissions)

		protected.GET("/users/:id/permissions",
			middleware.PermissionMiddleware("users.read", authService),
			userHandler.GetUserPermissions)

		protected.POST("/users/:id/change-password",
			middleware.PermissionMiddleware("profile.update", authService),
			userHandler.ChangePassword)

		// User profile
		protected.GET("/profile",
			middleware.PermissionMiddleware("profile.read", authService),
			userHandler.GetProfile)
	}
}
