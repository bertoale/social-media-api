package middlewares

import (
	"context"
	"fmt"
	"go-sosmed/pkg/cloudinary"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type UploadConfig struct {
	MaxFileSize   int64
	AllowedTypes  []string
	Folder        string
	FileFieldName string
}

func DefaultUploadConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize:   5 * 1024 * 1024,
		AllowedTypes:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		Folder:        "uploads",
		FileFieldName: "image",
	}
}

func UploadSingleFile(cloudinaryService *cloudinary.Service, config *UploadConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if config == nil {
			defaultConfig := DefaultUploadConfig()
			config = &defaultConfig
		}

		if err := c.Request.ParseMultipartForm(config.MaxFileSize); err != nil {
			c.Set("uploadedFile", "")
			c.Next()
			return
		}

		file, header, err := c.Request.FormFile(config.FileFieldName)
		if err != nil {
			c.Set("uploadedFile", "")
			c.Next()
			return
		}
		defer file.Close()

		if header.Size > config.MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Ukuran file maksimal %d MB", config.MaxFileSize/(1024*1024)),
			})
			c.Abort()
			return
		}

		ext := strings.ToLower(filepath.Ext(header.Filename))
		isAllowed := false
		for _, allowedType := range config.AllowedTypes {
			if ext == allowedType {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Tipe file tidak diizinkan. Hanya %v yang diperbolehkan", config.AllowedTypes),
			})
			c.Abort()
			return
		}

		if cloudinaryService == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Cloudinary tidak dikonfigurasi",
			})
			c.Abort()
			return
		}

		url, err := cloudinaryService.Upload(context.Background(), file, header, config.Folder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal upload ke Cloudinary: " + err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("uploadedFile", url)
		c.Next()
	}
}

func UploadAvatar(cloudinaryService *cloudinary.Service) gin.HandlerFunc {
	return UploadSingleFile(cloudinaryService, &UploadConfig{
		MaxFileSize:   2 * 1024 * 1024,
		AllowedTypes:  []string{".jpg", ".jpeg", ".png", ".webp"},
		Folder:        "avatars",
		FileFieldName: "avatar",
	})
}

func UploadPostImage(cloudinaryService *cloudinary.Service) gin.HandlerFunc {
	return UploadSingleFile(cloudinaryService, &UploadConfig{
		MaxFileSize:   5 * 1024 * 1024,
		AllowedTypes:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		Folder:        "posts",
		FileFieldName: "image",
	})
}
