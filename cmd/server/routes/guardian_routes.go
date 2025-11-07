package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/service"

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
	}
}
