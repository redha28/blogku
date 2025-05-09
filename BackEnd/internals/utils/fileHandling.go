package utils

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	fp "path/filepath"

	"github.com/gin-gonic/gin"
)

// Utils provides utility methods
type Utils struct{}

// NewUtils creates a new Utils instance
func NewUtils() *Utils {
	return &Utils{}
}

// UploadedFile contains information about an uploaded file
type UploadedFile struct {
	FileName string
	FilePath string
	FileSize int64
	FileType string
}

// FileHandling handles file upload with validation and cleanup
func (u *Utils) FileHandling(ctx *gin.Context, file *multipart.FileHeader, slug string, oldFilename string) (filename, filepath string, err error) {
	ext := fp.Ext(file.Filename)
	allowedExt := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	if !allowedExt[ext] {
		return "", "", fmt.Errorf("file extension not allowed")
	}

	// Create new filename with timestamp
	filename = fmt.Sprintf("%s_image%s", slug, ext)
	filepath = fp.Join("public", "uploads", filename)

	// Create directory if it doesn't exist
	uploadDir := fp.Join("public", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Save the new file
	if err = ctx.SaveUploadedFile(file, filepath); err != nil {
		return "", "", err
	}

	// Delete old file if exists
	if oldFilename != "" {
		oldPath := fp.Join("public", "uploads", oldFilename)
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			log.Println("[WARNING] Failed to delete old file:", err)
		}
	}

	return filename, filepath, nil
}
