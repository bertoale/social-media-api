package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadConfig konfigurasi untuk upload file
type UploadConfig struct {
	MaxFileSize   int64
	AllowedTypes  []string
	UploadDir     string // path di server (./uploads/posts)
	PublicPath    string // path untuk frontend (/uploads/posts)
	FileFieldName string
}

// DefaultUploadConfig memberikan konfigurasi default untuk upload
func DefaultUploadConfig() UploadConfig {
	return UploadConfig{
		MaxFileSize:   5 * 1024 * 1024,
		AllowedTypes:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		UploadDir:     "./uploads",
		PublicPath:    "/uploads",
		FileFieldName: "image",
	}
}

// UploadSingleFile middleware untuk upload single file
// Menyimpan file ke folder uploads dan menyimpan path-nya di context
// Parameter:
//   - config: Konfigurasi upload (opsional, jika nil akan pakai default)
func UploadSingleFile(config *UploadConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Gunakan config default jika tidak disediakan
		if config == nil {
			defaultConfig := DefaultUploadConfig()
			config = &defaultConfig
		}

		// Pastikan folder uploads ada
		if err := ensureUploadDir(config.UploadDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal membuat folder uploads",
			})
			c.Abort()
			return
		}

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(config.MaxFileSize); err != nil {
			// Jika error parsing, cek apakah memang tidak ada file
			// Jika tidak ada file, lanjutkan tanpa upload (opsional)
			c.Set("uploadedFile", "")
			c.Next()
			return
		}

		// Ambil file dari form
		file, header, err := c.Request.FormFile(config.FileFieldName)
		if err != nil {
			// Jika file tidak ada, lanjutkan tanpa upload (opsional)
			// Set context dengan empty string
			c.Set("uploadedFile", "")
			c.Next()
			return
		}
		defer file.Close()

		// Validasi ukuran file
		if header.Size > config.MaxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Ukuran file maksimal %d MB", config.MaxFileSize/(1024*1024)),
			})
			c.Abort()
			return
		}

		// Validasi tipe file
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

		// Generate nama file unik menggunakan timestamp
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		// Sanitize filename untuk keamanan
		sanitizedFilename := sanitizeFilename(header.Filename)
		filename := fmt.Sprintf("%s-%s", timestamp, sanitizedFilename)
		filePath := filepath.Join(config.UploadDir, filename)

		// simpan file
		if err := c.SaveUploadedFile(header, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Gagal menyimpan file: " + err.Error(),
			})
			c.Abort()
			return
		}

		// path RELATIF untuk DB & frontend
		publicURL := filepath.ToSlash(
			filepath.Join(config.PublicPath, filename),
		)

		c.Set("uploadedFile", publicURL)
		c.Set("uploadedFileName", filename)
		c.Set("uploadedFilePath", filePath)

		c.Next()
	}
}

func UploadAvatar() gin.HandlerFunc {
	return UploadSingleFile(&UploadConfig{
		MaxFileSize:   2 * 1024 * 1024,
		AllowedTypes:  []string{".jpg", ".jpeg", ".png", ".webp"},
		UploadDir:     "./uploads/avatars",
		PublicPath:    "/uploads/avatars",
		FileFieldName: "avatar",
	})
}

func UploadPostImage() gin.HandlerFunc {
	return UploadSingleFile(&UploadConfig{
		MaxFileSize:   5 * 1024 * 1024,
		AllowedTypes:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		UploadDir:     "./uploads/posts",
		PublicPath:    "/uploads/posts",
		FileFieldName: "image",
	})
}

// ensureUploadDir memastikan folder uploads ada, jika tidak ada akan dibuat
func ensureUploadDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// sanitizeFilename membersihkan nama file dari karakter yang tidak aman
func sanitizeFilename(filename string) string {
	// Hapus karakter yang berbahaya
	replacer := strings.NewReplacer(
		"..", "",
		"/", "",
		"\\", "",
		":", "",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	return replacer.Replace(filename)
}
