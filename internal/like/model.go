package like

import (
	"go-sosmed/internal/post"
	"go-sosmed/internal/user"
)

type Like struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	PostID uint `gorm:"not null"`
	// Relations
	User user.User `gorm:"foreignKey:UserID"`
	Post post.Post `gorm:"foreignKey:PostID"`
}

type LikeResponse struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	PostID uint `json:"post_id"`
}
