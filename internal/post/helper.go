package post

import "go-sosmed/internal/user"

func ToPostResponse(b *Post) *PostResponse {
	return &PostResponse{
		ID:           b.ID,
		Title:        b.Title,
		Content:      b.Content,
		Image:        b.Image,
		AuthorID:     b.AuthorID,
		Archived:     b.Archived,
		LikeCount:    int(b.LikeCount),
		CommentCount: int(b.CommentCount),
		IsLiked:      b.IsLiked,
		Edited:       b.Edited,
		CreatedAt:    b.CreatedAt,
		Author: user.AuthorResponse{
			ID:       b.Author.ID,
			Username: b.Author.Username,
			Avatar:   b.Author.Avatar,
		},
	}
}
