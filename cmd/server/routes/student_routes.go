package routes

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/middleware"
	"belajar-golang/internal/service"

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
	}
}
