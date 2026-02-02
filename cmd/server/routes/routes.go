package routes

import (
	"time"
	"u_kom_be/internal/handler"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

// SetupRoutes registers all application routes
func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authService service.AuthService,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	studentHandler *handler.StudentHandler,
	parentHandler *handler.ParentHandler,
	guardianHandler *handler.GuardianHandler,
	employeeHandler *handler.EmployeeHandler,
	dashboardHandler *handler.DashboardHandler,
	academicYearHandler *handler.AcademicYearHandler,
	classroomHandler *handler.ClassroomHandler,
) {
	// API v1 group
	apiV1 := router.Group("/api/v1")

	// Register all routes
	RegisterAuthRoutes(apiV1, authHandler, authService)
	RegisterUserRoutes(apiV1, userHandler, authService)
	RegisterRoleRoutes(apiV1, roleHandler, authService)
	RegisterPermissionRoutes(apiV1, permissionHandler, authService)
	RegisterStudentRoutes(apiV1, studentHandler, authService)
	RegisterParentRoutes(apiV1, parentHandler, authService)
	RegisterGuardianRoutes(apiV1, guardianHandler, authService)
	RegisterEmployeeRoutes(apiV1, employeeHandler, authService)
	RegisterDashboardRoutes(apiV1, dashboardHandler, authService)
	RegisterAcademicYearRoutes(apiV1, academicYearHandler, authService)
	RegisterClassroomRoutes(apiV1, classroomHandler, authService)

	// Health check route
	apiV1.GET("/health", healthCheck)
}

// healthCheck handles health check requests
func healthCheck(c *gin.Context) {
	handler.SuccessResponse(c, "Server is healthy and running", gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}
