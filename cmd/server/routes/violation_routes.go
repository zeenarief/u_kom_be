package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterViolationRoutes(router *gin.RouterGroup, violationHandler handler.ViolationHandler, authService service.AuthService) {
	violationGroup := router.Group("/violations")
	violationGroup.Use(middleware.AuthMiddleware(authService))
	{
		// Categories
		violationGroup.POST("/categories", middleware.PermissionMiddleware("violation_category.create", authService), violationHandler.CreateCategory)
		violationGroup.GET("/categories", middleware.PermissionMiddleware("violation_category.read", authService), violationHandler.GetCategories)
		violationGroup.PUT("/categories/:id", middleware.PermissionMiddleware("violation_category.update", authService), violationHandler.UpdateCategory)
		violationGroup.DELETE("/categories/:id", middleware.PermissionMiddleware("violation_category.delete", authService), violationHandler.DeleteCategory)

		// Types
		violationGroup.POST("/types", middleware.PermissionMiddleware("violation_type.create", authService), violationHandler.CreateType)
		violationGroup.GET("/types", middleware.PermissionMiddleware("violation_type.read", authService), violationHandler.GetTypes)
		violationGroup.PUT("/types/:id", middleware.PermissionMiddleware("violation_type.update", authService), violationHandler.UpdateType)
		violationGroup.DELETE("/types/:id", middleware.PermissionMiddleware("violation_type.delete", authService), violationHandler.DeleteType)

		// Student Violations
		violationGroup.POST("/record", middleware.PermissionMiddleware("violation_record.create", authService), violationHandler.RecordViolation)
		violationGroup.GET("/record/:id", middleware.PermissionMiddleware("violation_record.read", authService), violationHandler.GetStudentViolationDetail)
		violationGroup.PUT("/record/:id", middleware.PermissionMiddleware("violation_record.update", authService), violationHandler.UpdateViolation)
		violationGroup.GET("/student/:studentID", middleware.PermissionMiddleware("violation_record.read", authService), violationHandler.GetStudentViolations)
		violationGroup.DELETE("/record/:id", middleware.PermissionMiddleware("violation_record.delete", authService), violationHandler.DeleteViolation)
		violationGroup.GET("/all", middleware.PermissionMiddleware("violation_record.read_all", authService), violationHandler.GetAllViolations)
	}
}
