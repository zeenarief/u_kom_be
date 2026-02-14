package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterEmployeeRoutes(
	router *gin.RouterGroup,
	employeeHandler *handler.EmployeeHandler,
	authService service.AuthService,
) {
	// Grup rute 'employees' akan dilindungi oleh AuthMiddleware
	employees := router.Group("/employees")
	employees.Use(middleware.AuthMiddleware(authService))
	{
		employees.POST("",
			middleware.PermissionMiddleware("employees.create", authService),
			employeeHandler.CreateEmployee)

		employees.GET("",
			middleware.PermissionMiddleware("employees.read", authService),
			employeeHandler.GetAllEmployees)

		employees.GET("/:id",
			middleware.PermissionMiddleware("employees.read", authService),
			employeeHandler.GetEmployeeByID)

		employees.PUT("/:id",
			middleware.PermissionMiddleware("employees.update", authService),
			employeeHandler.UpdateEmployee)

		employees.DELETE("/:id",
			middleware.PermissionMiddleware("employees.delete", authService),
			employeeHandler.DeleteEmployee)

		// Endpoint untuk menautkan/melepas tautan akun user
		employees.POST("/:id/link-user",
			middleware.PermissionMiddleware("employees.manage_account", authService),
			employeeHandler.LinkUser)

		employees.DELETE("/:id/unlink-user",
			middleware.PermissionMiddleware("employees.manage_account", authService),
			employeeHandler.UnlinkUser)
	}
}
