package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/service"

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
	}
}
