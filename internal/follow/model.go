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

type FollowerResponse struct {
	ID       uint                `json:"id"`
	Follower user.AuthorResponse `json:"follower"`
}

type FollowingResponse struct {
	ID        uint                `json:"id"`
	Following user.AuthorResponse `json:"following"`
}
