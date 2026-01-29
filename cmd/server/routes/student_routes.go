package routes

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/middleware"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterStudentRoutes(
	router *gin.RouterGroup,
	studentHandler *handler.StudentHandler,
	authService service.AuthService,
) {
	// Grup rute 'students' akan dilindungi oleh AuthMiddleware
	students := router.Group("/students")
	students.Use(middleware.AuthMiddleware(authService))
	{
		students.POST("",
			middleware.PermissionMiddleware("students.create", authService),
			studentHandler.CreateStudent)

		students.GET("",
			middleware.PermissionMiddleware("students.read", authService),
			studentHandler.GetAllStudents)

		students.GET("/:id",
			middleware.PermissionMiddleware("students.read", authService),
			studentHandler.GetStudentByID)

		students.PUT("/:id",
			middleware.PermissionMiddleware("students.update", authService),
			studentHandler.UpdateStudent)

		students.DELETE("/:id",
			middleware.PermissionMiddleware("students.delete", authService),
			studentHandler.DeleteStudent)

		students.POST("/:id/sync-parents",
			middleware.PermissionMiddleware("students.manage_parents", authService),
			studentHandler.SyncParents)

		// Rute 1:1 Polymorphic Guardian (Set/Update)
		students.PUT("/:id/set-guardian",
			middleware.PermissionMiddleware("students.manage_guardian", authService),
			studentHandler.SetGuardian)

		// Rute 1:1 Polymorphic Guardian (Remove)
		students.DELETE("/:id/remove-guardian",
			middleware.PermissionMiddleware("students.manage_guardian", authService),
			studentHandler.RemoveGuardian)

		students.POST("/:id/link-user",
			middleware.PermissionMiddleware("students.manage_account", authService),
			studentHandler.LinkUser)

		students.DELETE("/:id/unlink-user",
			middleware.PermissionMiddleware("students.manage_account", authService),
			studentHandler.UnlinkUser)
	}
}
