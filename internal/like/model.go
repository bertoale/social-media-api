package like

import (
	"go-sosmed/internal/blog"
	"go-sosmed/internal/user"
)

type Like struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	BlogID uint `gorm:"not null"`
	// Relations
	User user.User `gorm:"foreignKey:UserID"`
	Blog blog.Blog `gorm:"foreignKey:BlogID"`
}

type LikeRequest struct {
	BlogID uint `json:"blog_id" form:"blog_id" binding:"required"`
}

type LikeResponse struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	BlogID uint `json:"blog"`
}
