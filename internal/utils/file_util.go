package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// EnsureDir memastikan folder tujuan ada
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// SaveUploadedFile menyimpan file dengan nama unik dan mengembalikan nama filenya
func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, destFolder string, prefix string) (string, error) {
	// 1. Validasi Ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".pdf" && ext != ".png" {
		return "", fmt.Errorf("invalid file type: only jpg, jpeg, png, and pdf are allowed")
	}

	// 2. Validasi Ukuran (Misal Max 2MB)
	if file.Size > 2*1024*1024 {
		return "", fmt.Errorf("file size exceeds 2MB limit")
	}

	// 3. Buat Folder jika belum ada
	uploadPath := fmt.Sprintf("./storage/uploads/%s", destFolder)
	if err := EnsureDir(uploadPath); err != nil {
		return "", err
	}

	// 4. Generate Nama File Unik (Timestamp + Random/Prefix)
	// Contoh: student_akta_1709999123.pdf
	filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadPath, filename)

	// 5. Simpan File
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	// Mengembalikan path relative untuk disimpan di DB
	// e.g., "students/student_akta_123456.pdf"
	return fmt.Sprintf("%s/%s", destFolder, filename), nil
}

// RemoveFile menghapus file lama saat update/delete
func RemoveFile(filePath string) {
	if filePath == "" {
		return
	}
	// Sesuaikan dengan root storage path Anda
	fullPath := fmt.Sprintf("./storage/uploads/%s", filePath)
	os.Remove(fullPath)
}
