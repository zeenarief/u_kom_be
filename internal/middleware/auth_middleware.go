package middleware

import (
	"strings"
	"smart_school_be/internal/handler"
	"smart_school_be/internal/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handler.UnauthorizedError(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handler.UnauthorizedError(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			handler.UnauthorizedError(c, "Invalid token")
			c.Abort()
			return
		}

		// Dapatkan user lengkap dari database
		user, err := authService.GetUserWithPermissions(userID) // Anda perlu menambahkan method ini di AuthService
		if err != nil {
			handler.InternalServerError(c, "Failed to get user data")
			c.Abort()
			return
		}

		// Set userID dan user object dalam context
		c.Set("user_id", userID)
		c.Set("user", user) // ✅ Ini yang penting ditambahkan

		c.Next()
	}
}
