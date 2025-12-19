package post

import (
	"go-sosmed/internal/user"
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey"`
	Title     string         `gorm:"not null"`
	Content   string         `gorm:"type:text;not null"`
	Image     string         `gorm:"type:text"`
	AuthorID  uint           `gorm:"not null"`
	Archived  bool           `gorm:"default:false"`
	Edited    bool           `gorm:"default:false"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//relation below
	Author user.User `gorm:"foreignKey:AuthorID"`
}

type PostRequest struct {
	Title   string `json:"title" form:"title" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
	Image   string `json:"image"`
}

type PostResponse struct {
	ID        uint                `json:"id"`
	Title     string              `json:"title"`
	Content   string              `json:"content"`
	Image     string              `json:"image"`
	Archived  bool                `json:"archived"`
	Edited    bool                `json:"edited"`
	AuthorID  uint                `json:"author_id"`
	CreatedAt time.Time           `json:"created_at"`
	Author    user.AuthorResponse `json:"author"`
}

type UpdatePostRequest struct {
	Title    *string `json:"title" form:"title" binding:"omitempty"`
	Content  *string `json:"content" form:"content" binding:"omitempty"`
	Archived *bool   `json:"archived" form:"archived" binding:"omitempty"`
	Edited   *bool   `json:"edited" form:"edited" binding:"omitempty"`
}
