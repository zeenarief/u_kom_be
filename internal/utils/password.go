package utils

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePasswordComplexity mengecek apakah password cukup kuat
func ValidatePasswordComplexity(pass string) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		// hasSpecial = false
	)

	if len(pass) >= 8 { // Naikkan jadi 8 karakter
		hasMinLen = true
	}

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
			// case unicode.IsPunct(char) || unicode.IsSymbol(char):
			//     hasSpecial = true
		}
	}

	if !hasMinLen || !hasUpper || !hasLower || !hasNumber {
		return fmt.Errorf("password harus minimal 8 karakter, mengandung huruf besar, huruf kecil, dan angka")
	}

	return nil
}
