package middleware

import (
	"u_kom_be/internal/handler"
	"u_kom_be/internal/model/domain"
	"u_kom_be/internal/service"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware checks if user has required permission
func PermissionMiddleware(permission string, authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			handler.UnauthorizedError(c, "User not found")
			c.Abort()
			return
		}

		userDomain, ok := user.(*domain.User)
		if !ok {
			handler.InternalServerError(c, "Invalid user object")
			c.Abort()
			return
		}

		// Check if user has the required permission
		if !userDomain.HasPermission(permission) {
			handler.ForbiddenError(c, "You don't have permission to access this resource")
			c.Abort()
			return
		}

		c.Next()
	}
}
