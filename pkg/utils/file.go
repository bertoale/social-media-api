package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DeleteFile menghapus file dari sistem jika ada
// Parameter:
//   - filePath: path relatif atau absolut dari file yang akan dihapus
//
// Returns: error jika gagal menghapus file
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	// Jika path adalah URL lengkap (misal: http://localhost:5000/uploads/file.jpg)
	// Extract hanya nama file-nya
	if strings.Contains(filePath, "/uploads/") {
		parts := strings.Split(filePath, "/uploads/")
		if len(parts) > 1 {
			filePath = filepath.Join("uploads", parts[1])
		}
	}

	// Cek apakah file ada
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File tidak ada, tidak perlu dihapus
	}

	// Hapus file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// FileExists mengecek apakah file ada di sistem
// Parameter:
//   - filePath: path dari file yang akan dicek
//
// Returns: true jika file ada, false jika tidak
func FileExists(filePath string) bool {
	if filePath == "" {
		return false
	}

	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetFilePath mengubah URL menjadi path lokal
// Parameter:
//   - fileURL: URL dari file (misal: http://localhost:5000/uploads/file.jpg)
//
// Returns: path lokal (misal: uploads/file.jpg)
func GetFilePath(fileURL string) string {
	if fileURL == "" {
		return ""
	}

	// Jika sudah path lokal, return langsung
	if !strings.HasPrefix(fileURL, "http") {
		return fileURL
	}

	// Extract path dari URL
	if strings.Contains(fileURL, "/uploads/") {
		parts := strings.Split(fileURL, "/uploads/")
		if len(parts) > 1 {
			return filepath.Join("uploads", parts[1])
		}
	}

	return fileURL
}
