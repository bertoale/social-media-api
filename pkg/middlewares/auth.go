package middlewares

import (
	"go-sosmed/internal/session"
	"go-sosmed/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate(cfg *config.Config, sessionService session.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString != "" && strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		} else {
			tokenString = c.GetString("token")
			if tokenString == "" {
				tokenString, _ = c.Cookie("token")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Akses ditolak. Token tidak ditemukan.",
			})
			c.Abort()
			return
		}

		sess, err := sessionService.ValidateSession(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token tidak valid atau kadaluarsa.",
			})
			c.Abort()
			return
		}

		db := config.GetDB()

		var count int64
		if err := db.Table("users").Where("id = ?", sess.UserID).Count(&count).Error; err != nil || count == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User tidak ditemukan.",
			})
			c.Abort()
			return
		}

		var userRole string
		if err := db.Table("users").Where("id = ?", sess.UserID).Select("role").Scan(&userRole).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Gagal mengambil role user.",
			})
			c.Abort()
			return
		}

		c.Set("userID", sess.UserID)
		c.Set("userRole", userRole)

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

		if len(roles) == 0 {
			c.Next()
			return
		}

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
