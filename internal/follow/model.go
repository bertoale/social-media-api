package follow

import "go-sosmed/internal/user"

type Follow struct {
	ID          uint `gorm:"primaryKey"`
	FollowerID  uint `gorm:"not null"`
	FollowingID uint `gorm:"not null"`
	// Relations
	Follower  user.User `gorm:"foreignKey:FollowerID"`
	Following user.User `gorm:"foreignKey:FollowingID"`
}

type FollowRequest struct {
	FollowingID uint `json:"following_id" form:"following_id" binding:"required"`
	FollowerID  uint `json:"follower_id" form:"follower_id" binding:"required"`
}

type FollowerResponse struct {
	ID       uint                `json:"id"`
	Follower user.AuthorResponse `json:"follower"`
}

type FollowingResponse struct {
	ID        uint                `json:"id"`
	Following user.AuthorResponse `json:"following"`
}
