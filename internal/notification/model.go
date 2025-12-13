package notification

import (
	"go-sosmed/internal/user"
	"time"
)

type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"` // User yang menerima notifikasi
	Type      string `gorm:"not null"` // follow, like, comment, reply
	ActorID   uint   `gorm:"not null"` // User yang melakukan aksi
	TargetID  uint   // ID blog/comment yang terkait (optional)
	Message   string `gorm:"type:text"`
	IsRead    bool   `gorm:"default:false"`
	CreatedAt time.Time

	User  user.User `gorm:"foreignKey:UserID"`
	Actor user.User `gorm:"foreignKey:ActorID"`
}

type NotificationRequest struct {
	Type     string `json:"type" binding:"required"`
	ActorID  uint   `json:"actor_id" binding:"required"`
	TargetID uint   `json:"target_id"`
	Message  string `json:"message"`
}

type NotificationResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Type      string `json:"type"`
	ActorID   uint   `json:"actor_id"`
	TargetID  uint   `json:"target_id"`
	Message   string `json:"message"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}
