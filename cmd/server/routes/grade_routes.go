package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterGradeRoutes(router *gin.RouterGroup, gradeHandler *handler.GradeHandler, authService service.AuthService) {
	gradeGroup := router.Group("/grades")
	gradeGroup.Use(middleware.AuthMiddleware(authService))

	// Assessments
	gradeGroup.POST("/assessments",
		middleware.PermissionMiddleware("assignments.manage", authService),
		gradeHandler.CreateAssessment)

	gradeGroup.PUT("/assessments/:id",
		middleware.PermissionMiddleware("assignments.manage", authService),
		gradeHandler.UpdateAssessment)

	gradeGroup.GET("/assessments/teaching-assignment/:teachingAssignmentID",
		middleware.PermissionMiddleware("assessments.read", authService),
		gradeHandler.GetAssessmentsByTeachingAssignment) // Note: permissions need to be checked

	gradeGroup.GET("/assessments/:id",
		middleware.PermissionMiddleware("assessments.read", authService),
		gradeHandler.GetAssessmentDetail)

	gradeGroup.DELETE("/assessments/:id",
		middleware.PermissionMiddleware("assignments.manage", authService),
		gradeHandler.DeleteAssessment)

	// Scores
	gradeGroup.POST("/scores/bulk",
		middleware.PermissionMiddleware("assignments.manage", authService),
		gradeHandler.SubmitScores)
}
