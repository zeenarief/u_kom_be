package routes

import (
	"smart_school_be/internal/handler"
	"smart_school_be/internal/middleware"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterFinanceRoutes(router *gin.RouterGroup, handler *handler.FinanceHandler, authService service.AuthService) {
	group := router.Group("/finance")
	group.Use(middleware.AuthMiddleware(authService))

	// Donations
	group.POST("/donations", middleware.PermissionMiddleware("finance_donations.create", authService), handler.CreateDonation)
	group.GET("/donations", middleware.PermissionMiddleware("finance_donations.read", authService), handler.GetDonations)
	group.GET("/donations/:id", middleware.PermissionMiddleware("finance_donations.read", authService), handler.GetDonationByID)
	group.PUT("/donations/:id", middleware.PermissionMiddleware("finance_donations.update", authService), handler.UpdateDonation)

	// Donors
	group.GET("/donors", middleware.PermissionMiddleware("finance_donors.read", authService), handler.GetDonors)
	group.GET("/donors/:id", middleware.PermissionMiddleware("finance_donors.read", authService), handler.GetDonorByID)
	group.PUT("/donors/:id", middleware.PermissionMiddleware("finance_donors.update", authService), handler.UpdateDonor)
}
