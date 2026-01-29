package routes

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterParentRoutes(
	router *gin.RouterGroup,
	parentHandler *handler.ParentHandler,
	authService service.AuthService,
) {
	// Grup rute 'parents' akan dilindungi oleh AuthMiddleware
	parents := router.Group("/parents")
	parents.Use(middleware.AuthMiddleware(authService))
	{
		parents.POST("",
			middleware.PermissionMiddleware("parents.create", authService),
			parentHandler.CreateParent)

		parents.GET("",
			middleware.PermissionMiddleware("parents.read", authService),
			parentHandler.GetAllParents)

		parents.GET("/:id",
			middleware.PermissionMiddleware("parents.read", authService),
			parentHandler.GetParentByID)

		parents.PUT("/:id",
			middleware.PermissionMiddleware("parents.update", authService),
			parentHandler.UpdateParent)

		parents.DELETE("/:id",
			middleware.PermissionMiddleware("parents.delete", authService),
			parentHandler.DeleteParent)

		parents.POST("/:id/link-user",
			middleware.PermissionMiddleware("parents.manage_account", authService),
			parentHandler.LinkUser)

		parents.DELETE("/:id/unlink-user",
			middleware.PermissionMiddleware("parents.manage_account", authService),
			parentHandler.UnlinkUser)
	}
}
