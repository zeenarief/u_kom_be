package middleware

import (
	"belajar-golang/internal/handler"
	"belajar-golang/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			handler.UnauthorizedError(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			handler.UnauthorizedError(c, "Authorization header format must be Bearer {token}")
			c.Abort()
			return
		}

		tokenString := parts[1]
		userID, err := authService.ValidateToken(tokenString)
		if err != nil {
			//c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			handler.UnauthorizedError(c, "Invalid token")
			c.Abort()
			return
		}

		// Set userID dalam context untuk digunakan di handler
		c.Set("userID", userID)
		c.Next()
	}
}
