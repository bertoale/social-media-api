package comment

import (
	"go-sosmed/internal/blog"
	"go-sosmed/internal/user"
	"time"
)

type Comment struct {
	ID            uint      `gorm:"primaryKey"`
	BlogID        uint      `gorm:"not null"`
	UserID        uint      `gorm:"not null"`
	ParentID      *uint     `gorm:"index"` // index agar cepat mencari replies
	ReplyToUserID *uint     `gorm:"index"` // ID user yang direply, bisa null
	Content       string    `gorm:"type:text;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	Edited        bool      `gorm:"default:false"`

	// Relations
	Blog        blog.Blog  `gorm:"foreignKey:BlogID"` // harus kapital
	User        user.User  `gorm:"foreignKey:UserID"`
	ReplyToUser *user.User `gorm:"foreignKey:ReplyToUserID"`
	Replies     []Comment  `gorm:"foreignKey:ParentID"`
}

type CommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID          uint                 `json:"id"`
	Content     string               `json:"content"`
	CreatedAt   time.Time            `json:"created_at"`
	Edited      bool                 `json:"edited"`
	User        user.AuthorResponse  `json:"user"`
	ReplyToUser *user.AuthorResponse `json:"reply_to_user,omitempty"`
	Replies     []CommentResponse    `json:"replies,omitempty"`
}

type UpdateCommentRequest struct {
	Content *string `json:"content" form:"content" binding:"required"`
	Edited  *bool   `json:"edited" form:"edited"`
}

type ReplyToUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type ReplyCommentRequest struct {
	Content       string `json:"content" binding:"required"`
	ReplyToUserID uint   `json:"reply_to_user_id"`
}

// //==========================================================//

// type ReplyComment struct {
//     ID        uint      `gorm:"primaryKey"`
//     CommentID uint      `gorm:"not null"` // Komentar yang direply
//     UserID    uint      `gorm:"not null"` // User yang reply
//     Content   string    `gorm:"type:text;not null"`
//     CreatedAt time.Time

//     Comment Comment   `gorm:"foreignKey:CommentID"`
//     User    user.User `gorm:"foreignKey:UserID"`
// }

// type ReplyCommentRequest struct {
//     CommentID uint   `json:"comment_id" binding:"required"`
//     Content   string `json:"content" binding:"required"`
// }

// type ReplyCommentResponse struct {
//     ID        uint   `json:"id"`
//     CommentID uint   `json:"comment_id"`
//     UserID    uint   `json:"user_id"`
//     Content   string `json:"content"`
//     CreatedAt string `json:"created_at"`
// }
