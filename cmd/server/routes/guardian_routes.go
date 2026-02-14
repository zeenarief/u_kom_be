package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterGuardianRoutes(
	router *gin.RouterGroup,
	guardianHandler *handler.GuardianHandler,
	authService service.AuthService,
) {
	// Grup rute 'guardians' akan dilindungi oleh AuthMiddleware
	guardians := router.Group("/guardians")
	guardians.Use(middleware.AuthMiddleware(authService))
	{
		guardians.POST("",
			middleware.PermissionMiddleware("guardians.create", authService),
			guardianHandler.CreateGuardian)

		guardians.GET("",
			middleware.PermissionMiddleware("guardians.read", authService),
			guardianHandler.GetAllGuardians)

		guardians.GET("/:id",
			middleware.PermissionMiddleware("guardians.read", authService),
			guardianHandler.GetGuardianByID)

		guardians.PUT("/:id",
			middleware.PermissionMiddleware("guardians.update", authService),
			guardianHandler.UpdateGuardian)

		guardians.DELETE("/:id",
			middleware.PermissionMiddleware("guardians.delete", authService),
			guardianHandler.DeleteGuardian)

		guardians.POST("/:id/link-user",
			middleware.PermissionMiddleware("guardians.manage_account", authService),
			guardianHandler.LinkUser)

		guardians.DELETE("/:id/unlink-user",
			middleware.PermissionMiddleware("guardians.manage_account", authService),
			guardianHandler.UnlinkUser)
	}
}
