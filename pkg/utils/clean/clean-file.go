package clean

import (
	"os"
	"path/filepath"

	"gorm.io/gorm"

	"go-sosmed/internal/blog"
	"go-sosmed/internal/user"
)

/*
CleanupUnusedUploads
- uploadDir  : path folder uploads (contoh: "./uploads")
- db         : gorm DB
*/
func CleanupUnusedUploads(db *gorm.DB, uploadDir string) error {
	usedFiles, err := getAllUsedFiles(db)
	if err != nil {
		return err
	}

	usedMap := make(map[string]struct{})
	for _, f := range usedFiles {
		usedMap[f] = struct{}{}
	}

	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		// skip file default / penting
		if filename == "default-avatar.png" {
			continue
		}

		dbPath := "/uploads/" + filename
		if _, ok := usedMap[dbPath]; !ok {
			_ = os.Remove(filepath.Join(uploadDir, filename))
		}
	}

	return nil
}

// =====================================
// PRIVATE FUNCTIONS
// =====================================

func getAllUsedFiles(db *gorm.DB) ([]string, error) {
	var files []string

	var avatars []string
	if err := db.Model(&user.User{}).
		Where("avatar != ''").
		Pluck("avatar", &avatars).Error; err != nil {
		return nil, err
	}

	var blogImages []string
	if err := db.Model(&blog.Blog{}).
		Where("image != ''").
		Pluck("image", &blogImages).Error; err != nil {
		return nil, err
	}

	files = append(files, avatars...)
	files = append(files, blogImages...)

	return files, nil
}
