package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(router *gin.RouterGroup, roleHandler *handler.RoleHandler, authService service.AuthService) {
	roles := router.Group("/roles")
	roles.Use(middleware.AuthMiddleware(authService))
	roles.Use(middleware.PermissionMiddleware("roles.manage", authService))

	{
		roles.POST("", roleHandler.CreateRole)
		roles.GET("", roleHandler.GetAllRoles)
		roles.GET("/:id", roleHandler.GetRoleByID)
		roles.PUT("/:id", roleHandler.UpdateRole)
		roles.DELETE("/:id", roleHandler.DeleteRole)
		roles.POST("/:id/sync-permissions", roleHandler.SyncRolePermissions)
	}
}
