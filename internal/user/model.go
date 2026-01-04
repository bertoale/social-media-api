package user

type RoleType = string

const (
	RoleAdmin RoleType = "admin"
	RoleUser  RoleType = "user"
)

type User struct {
	ID       uint     `gorm:"primaryKey"`
	Username string   `gorm:"unique;not null"`
	Email    string   `gorm:"unique;not null"`
	Password string   `gorm:"not null"`
	Bio      string   `gorm:"type:text"`
	Avatar   string   `gorm:"type:text"`
	Role     RoleType `gorm:"default:'user'"`
	//computed fields
	FollowersCount int64 `gorm:"-:migration;<-:false"` // ignored by GORM migrations and write operations
	FollowingCount int64 `gorm:"-:migration;<-:false"` // ignored by GORM migrations and write operations
	IsFollowed     bool  `gorm:"-:migration;<-:false"` // ignored by GORM migrations and write operations
}

type RegisterRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UpdateProfileRequest struct {
	Username *string `json:"username" form:"username"`
	Bio      *string `json:"bio" form:"bio"`
	Avatar   *string `json:"avatar"`
}

type UserResponse struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	Avatar         string `json:"avatar"`
	FollowersCount int64  `json:"followers_count"`
	FollowingCount int64  `json:"following_count"`
	IsFollowed     bool   `json:"is_followed"`
	Role           RoleType `json:"role"`
}

type AuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
