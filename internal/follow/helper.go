package follow

import "go-sosmed/internal/user"

func ToFollowerResponse(f *Follow) *FollowerResponse {
	return &FollowerResponse{
		ID: f.ID,
		Follower: user.AuthorResponse{
			ID:       f.Follower.ID,
			Username: f.Follower.Username,
			Avatar:   f.Follower.Avatar,
		},
	}
}

func ToFollowingResponse(f *Follow) *FollowingResponse {
	return &FollowingResponse{
		ID: f.ID,
		Following: user.AuthorResponse{
			ID:       f.Following.ID,
			Username: f.Following.Username,
			Avatar:   f.Following.Avatar,
		},
	}
}
