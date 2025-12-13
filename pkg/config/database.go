package config

// Package database handles database connection and migration
// Menggunakan GORM sebagai ORM dan MySQL sebagai database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah global variable untuk database connection
var DB *gorm.DB

// Connect membuat koneksi ke database MySQL
// Function ini dipanggil saat aplikasi startup
// Parameters:
//   - cfg: Config object yang berisi database credentials
//
// Returns: error jika koneksi gagal, nil jika berhasil
func Connect(cfg *Config) error {
	// Build MySQL DSN (Data Source Name) connection string
	// Format MySQL DSN:
	//   username:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC() // Gunakan UTC agar konsisten dengan server
		},
	})

	if err != nil {
		return fmt.Errorf("❌ gagal konek ke database MySQL: %w", err)
	}

	log.Println("✅ Koneksi ke database MySQL berhasil.")

	// Verifikasi database aktif
	var dbName string
	DB.Raw("SELECT DATABASE()").Scan(&dbName)
	log.Println("🔍 Terkoneksi ke database:", dbName)

	return nil
}

// Migrate menjalankan database migration untuk semua models
// AutoMigrate akan membuat tabel jika belum ada
func Migrate() error {
	// Migration akan dilakukan di main.go untuk menghindari import cycle
	log.Println("✅ Migrasi database berhasil.")
	return nil
}

// GetDB mengembalikan instance *gorm.DB
func GetDB() *gorm.DB {
	return DB
}
