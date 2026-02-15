package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"
	"time"

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
	subjectHandler *handler.SubjectHandler,
	teachingAssignmentHandler *handler.TeachingAssignmentHandler,
	scheduleHandler *handler.ScheduleHandler,
	attendanceHandler *handler.AttendanceHandler,
	gradeHandler *handler.GradeHandler,
	violationHandler handler.ViolationHandler,
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
	RegisterSubjectRoutes(apiV1, subjectHandler, authService)
	RegisterTeachingAssignmentRoutes(apiV1, teachingAssignmentHandler, authService)
	RegisterScheduleRoutes(apiV1, scheduleHandler, authService)
	RegisterAttendanceRoutes(apiV1, attendanceHandler, authService)
	RegisterGradeRoutes(apiV1, gradeHandler, authService)
	RegisterViolationRoutes(apiV1, violationHandler, authService)

	protected := apiV1.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))

	protected.GET("/files/:folder/:filename", authHandler.ServeFile)

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
