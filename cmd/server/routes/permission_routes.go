package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(router *gin.RouterGroup, permissionHandler *handler.PermissionHandler, authService service.AuthService) {
	permissions := router.Group("/permissions")
	permissions.Use(middleware.AuthMiddleware(authService))
	permissions.Use(middleware.PermissionMiddleware("permissions.manage", authService))

	{
		permissions.POST("", permissionHandler.CreatePermission)
		permissions.GET("", permissionHandler.GetAllPermissions)
		permissions.GET("/:id", permissionHandler.GetPermissionByID)
		permissions.PUT("/:id", permissionHandler.UpdatePermission)
		permissions.DELETE("/:id", permissionHandler.DeletePermission)
	}
}
