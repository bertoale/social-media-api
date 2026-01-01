package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func normalizeUploadPath(input string) string {
	if input == "" {
		return ""
	}

	// Hilangkan domain jika URL
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		if idx := strings.Index(input, "/uploads/"); idx != -1 {
			input = input[idx:]
		}
	}

	// Jika path publik
	if strings.HasPrefix(input, "/uploads/") {
		return filepath.Join(".", input)
	}

	return input
}

// DeleteFile menghapus file dari sistem jika ada
// Parameter:
//   - filePath: path relatif atau absolut dari file yang akan dihapus
//
// Returns: error jika gagal menghapus file
func DeleteFile(filePath string) error {
	if filePath == "" {
		return nil
	}

	localPath := normalizeUploadPath(filePath)

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(localPath); err != nil {
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

	localPath := normalizeUploadPath(filePath)
	_, err := os.Stat(localPath)

	return err == nil
}

// GetFilePath mengubah URL menjadi path lokal
// Parameter:
//   - fileURL: URL dari file (misal: http://localhost:5000/uploads/file.jpg)
//
// Returns: path lokal (misal: uploads/file.jpg)
// func GetFilePath(fileURL string) string {
// 	if fileURL == "" {
// 		return ""
// 	}

// 	// Jika sudah path lokal, return langsung
// 	if !strings.HasPrefix(fileURL, "http") {
// 		return fileURL
// 	}

// 	// Extract path dari URL
// 	if strings.Contains(fileURL, "/uploads/") {
// 		parts := strings.Split(fileURL, "/uploads/")
// 		if len(parts) > 1 {
// 			return filepath.Join("uploads", parts[1])
// 		}
// 	}

// 	return fileURL
// }
