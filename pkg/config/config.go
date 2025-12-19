// Package config handles application configuration management
// Configuration di-load dari file .env dan environment variables
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct menyimpan semua konfigurasi aplikasi
// Semua field adalah string karena dibaca dari environment variables
type Config struct {
	DBHost     string // Database host (default: localhost)
	DBPort     string // Database port (default: 5432 untuk PostgreSQL)
	DBUser     string // Database user
	DBPassword string // Database password
	DBName     string // Database name
	DBSSLMode  string // Database SSL mode (disable/require/verify-ca/verify-full)
	JWTSecret  string // Secret key untuk signing JWT tokens
	JWTExpires string // JWT expiration duration (contoh: 168h = 7 hari)
	Port       string // Port untuk aplikasi web server
	NodeEnv    string // Environment mode (development/production)
	CorsOrigin string // Allowed CORS origin (URL frontend)

	// Mailjet email configuration
	MailjetAPIKey    string // Mailjet API key
	MailjetAPISecret string // Mailjet API secret
	MailjetPort      string // Mailjet SMTP port
	MailjetHost      string // Mailjet SMTP host
	MailSenderEmail  string // Email address untuk sender
	MailSenderName   string // Nama sender yang tampil di email
}

// LoadConfig membaca konfigurasi dari file .env dan environment variables
// Priority: Environment variables > .env file > default values
// Returns: Pointer ke Config struct yang sudah terisi
func LoadConfig() *Config {
	// Load .env file jika ada
	// Jika .env tidak ditemukan, tidak akan error, hanya warning
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, using environment variables")
	}
	// Return Config struct dengan values dari getEnv()
	// getEnv() akan mencari environment variable, jika tidak ada gunakan default value
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "post_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "your_super_secret_jwt_key_post_app_2025"),
		JWTExpires: getEnv("JWT_EXPIRES_IN", "168h"),
		Port:       getEnv("PORT", "5000"),
		NodeEnv:    getEnv("NODE_ENV", "development"),
		CorsOrigin: getEnv("CORS_ORIGIN", "http://localhost:3000"),

		// Mailjet configuration
		MailjetAPIKey:    getEnv("MAILJET_API_KEY", ""),
		MailjetAPISecret: getEnv("MAILJET_API_SECRET", ""),
		MailjetPort:      getEnv("MAILJET_PORT", "587"),
		MailjetHost:      getEnv("MAILJET_HOST", "in-v3.mailjet.com"),
		MailSenderEmail:  getEnv("MAIL_SENDER_EMAIL", "noreply@goevent.com"),
		MailSenderName:   getEnv("MAIL_SENDER_NAME", "GoEvent App"),
	}
}

// getEnv adalah helper function untuk membaca environment variable
// Jika environment variable tidak ada, return default value
// Parameters:
//   - key: Nama environment variable yang ingin dibaca
//   - defaultValue: Nilai default jika environment variable tidak ditemukan
//
// Returns: Value dari environment variable atau default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
