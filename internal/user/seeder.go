package user

import (
	"go-sosmed/pkg/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password: %v", err)
	}
	return string(hashed)
}

func SeedAdminUser() {
	db := config.GetDB()

	admin := User{
		Username: "admin",
		Email:    "admin@email.com",
		Password: HashPassword("adminadmin"),
		Role:     RoleAdmin,
	}

	var count int64
	db.Model(&User{}).Where("username = ? OR email = ?", admin.Username, admin.Email).Count(&count)

	if count == 0 {
		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Gagal membuat user admin: %v", err)
		} else {
			log.Println("User admin berhasil dibuat.")
		}
	} else {
		log.Println("User admin sudah ada.")
	}
}
