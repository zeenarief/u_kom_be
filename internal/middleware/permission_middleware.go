package middleware

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/service"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware checks if user has required permission
func PermissionMiddleware(permission string, authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			handler.UnauthorizedError(c, "User ID not found")
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
			handler.InternalServerError(c, "Invalid user ID")
			c.Abort()
			return
		}

		// Get user with roles and permissions
		user, err := authService.GetUserWithPermissions(userIDStr)
		if err != nil {
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user permissions"})
			handler.InternalServerError(c, "Failed to get user permissions")
			c.Abort()
			return
		}

		// Check if user has the required permission
		if !user.HasPermission(permission) {
			//c.JSON(http.StatusForbidden, gin.H{
			//	"error":               "Access denied",
			//	"message":             "You don't have permission to access this resource",
			//	"required_permission": permission,
			//	"your_permissions":    user.GetPermissions(),
			//})
			handler.ForbiddenError(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
