package report

import (
	"go-sosmed/internal/post"
	"go-sosmed/internal/user"
	"time"
)

type StatusType = string

const (
	StatusPending  StatusType = "pending"
	StatusReviewed StatusType = "reviewed"
	StatusResolved StatusType = "resolved"
	StatusRejected StatusType = "rejected"
)

type Report struct {
	ID          uint       `gorm:"primaryKey"`
	UserID      uint       `gorm:"not null"`
	PostID      uint       `gorm:"not null"`
	Reason      string     `gorm:"type:text;not null"`
	Description string     `gorm:"type:text"`
	Status      StatusType `gorm:"default:'pending'"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	//relations
	User user.User `gorm:"foreignKey:UserID"`
	Post post.Post `gorm:"foreignKey:PostID"`
}

type ReportRequest struct {
	Reason      string `json:"reason" form:"reason" binding:"required"`
	Description string `json:"description" form:"description"`
}

type ReportResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	PostID      uint      `json:"post_id"`
	Reason      string    `json:"reason"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	AdminID     *uint     `json:"admin_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateReportRequest struct {
	Status string `json:"status" form:"status" binding:"required"`
}
