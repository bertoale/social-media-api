package report

import (
	"go-sosmed/internal/blog"
	"go-sosmed/internal/user"
	"time"
)

type Report struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	BlogID    uint   `gorm:"not null"`
	Reason    string `gorm:"type:text;not null"`
	Status    string `gorm:"default:'pending'"` // pending, reviewed, resolved
	AdminID   *uint  `gorm:"default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	//relations
	User  user.User  `gorm:"foreignKey:UserID"`
	Blog  blog.Blog  `gorm:"foreignKey:BlogID"`
	Admin *user.User `gorm:"foreignKey:AdminID"`
}

type ReportRequest struct {
	BlogID uint   `json:"blog_id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

type ReportResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	BlogID    uint      `json:"blog_id"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	AdminID   *uint     `json:"admin_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminReviewRequest struct {
	Status  string `json:"status" binding:"required"` // reviewed, resolved
	AdminID uint   `json:"admin_id" binding:"required"`
}
