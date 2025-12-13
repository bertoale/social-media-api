package like

import "go-sosmed/internal/blog"

func ToLikeResponse(l *Like) *LikeResponse {
	return &LikeResponse{
		ID:     l.ID,
		UserID: l.UserID,
		BlogID: l.BlogID,
	}
}

func ToBlogsLikedByUserResponse(blogs []blog.Blog) []*blog.BlogResponse {
	var responses []*blog.BlogResponse
	for _, b := range blogs {
		responses = append(responses, blog.ToBlogResponse(&b))
	}
	return responses
}
