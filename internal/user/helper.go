package user

func ToUserResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		Bio:            u.Bio,
		Avatar:         u.Avatar,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
		IsFollowed:     u.IsFollowed,
		Role:           u.Role,
	}
}
