package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// Hapus fungsi ValidateToken dari sini karena sudah pindah ke AuthService
// Biarkan hanya fungsi yang general-purpose di sini

// GenerateTokenFromClaims - fungsi helper untuk generate token dengan claims custom
func GenerateTokenFromClaims(claims jwt.Claims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken - fungsi helper untuk parse token
func ParseToken(tokenString, secret string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
}
