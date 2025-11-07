package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// EncryptionUtil adalah interface untuk helper enkripsi/dekripsi
type EncryptionUtil interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

type encryptionUtil struct {
	key []byte
}

// NewEncryptionUtil membuat instance baru dari EncryptionUtil
// Pastikan 'key' adalah 32 byte (untuk AES-256)
func NewEncryptionUtil(key string) (EncryptionUtil, error) {
	byteKey := []byte(key)
	if len(byteKey) != 32 {
		return nil, errors.New("encryption key must be 32 bytes long")
	}
	return &encryptionUtil{key: byteKey}, nil
}

func (e *encryptionUtil) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	// Kita gabungkan nonce dan ciphertext lalu encode ke Base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *encryptionUtil) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
