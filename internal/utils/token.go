package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken creates a hash of the token for storage
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
