package like

import (
	"go-sosmed/internal/post"
)

func ToLikeResponse(l *Like) *LikeResponse {
	return &LikeResponse{
		ID:     l.ID,
		UserID: l.UserID,
		PostID: l.PostID,
	}
}

func ToPostsLikedByUserResponse(posts []post.Post) []*post.PostResponse {
	var responses []*post.PostResponse
	for _, b := range posts {
		responses = append(responses, post.ToPostResponse(&b))
	}
	return responses
}
