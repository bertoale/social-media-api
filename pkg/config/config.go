// Package config handles application configuration management
// Configuration di-load dari file .env dan environment variables
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	JWTSecret  string
	SessionExpires string
	Port       string
	NodeEnv    string
	CorsOrigin string

	MailjetAPIKey    string
	MailjetAPISecret string
	MailjetPort      string
	MailjetHost      string
	MailSenderEmail  string
	MailSenderName   string

	CookieDomain string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryFolder    string
}

// LoadConfig membaca konfigurasi dari file .env dan environment variables
// Priority: Environment variables > .env file > default values
// Returns: Pointer ke Config struct yang sudah terisi
func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "post_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "your_super_secret_jwt_key_post_app_2025"),
		SessionExpires: getEnv("SESSION_EXPIRES", "24h"),
		Port:       getEnv("PORT", "5000"),
		NodeEnv:    getEnv("NODE_ENV", "development"),
		CorsOrigin: getEnv("CORS_ORIGIN", "http://localhost:3000"),

		MailjetAPIKey:    getEnv("MAILJET_API_KEY", ""),
		MailjetAPISecret: getEnv("MAILJET_API_SECRET", ""),
		MailjetPort:      getEnv("MAILJET_PORT", "587"),
		MailjetHost:      getEnv("MAILJET_HOST", "in-v3.mailjet.com"),
		MailSenderEmail:  getEnv("MAIL_SENDER_EMAIL", "noreply@goevent.com"),
		MailSenderName:   getEnv("MAIL_SENDER_NAME", "GoEvent App"),

		CookieDomain: getEnv("COOKIE_DOMAIN", ""),

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryFolder:    getEnv("CLOUDINARY_FOLDER", "social-media"),
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
