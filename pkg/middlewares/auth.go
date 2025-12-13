package middlewares

import (
	"net/http"
	"strings"

	"go-sosmed/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func Authenticate(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Ambil token dari header Authorization atau cookie
		tokenString := c.GetHeader("Authorization")
		if tokenString != "" && strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		} else {
			tokenString = c.GetString("token")
			if tokenString == "" {
				tokenString, _ = c.Cookie("token")
			}
		}

		// Token tidak ada
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Akses ditolak. Token tidak ditemukan.",
			})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak valid atau kadaluarsa.",
			})
			c.Abort()
			return
		}

		// Validasi user ada di database
		db := config.GetDB()

		var count int64
		if err := db.Table("users").Where("id = ?", claims.ID).Count(&count).Error; err != nil || count == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User tidak ditemukan.",
			})
			c.Abort()
			return
		}

		// Taruh user data di context
		c.Set("userID", claims.ID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

func Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		roleValue, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User belum terautentikasi.",
			})
			c.Abort()
			return
		}

		userRole := roleValue.(string)

		// Jika tidak ada batasan role, semua user boleh
		if len(roles) == 0 {
			c.Next()
			return
		}

		// Cocokkan role
		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"message": "Akses ditolak. Anda tidak memiliki izin yang sesuai.",
		})
		c.Abort()
	}
}
