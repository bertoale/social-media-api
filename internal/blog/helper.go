package blog

import "go-sosmed/internal/user"

func ToBlogResponse(b *Blog) *BlogResponse {
	return &BlogResponse{
		ID:        b.ID,
		Title:     b.Title,
		Content:   b.Content,
		Image:     b.Image,
		AuthorID:  b.AuthorID,
		Archived:  b.Archived,
		Edited:    b.Edited,
		CreatedAt: b.CreatedAt,
		Author: user.AuthorResponse{
			ID:       b.Author.ID,
			Username: b.Author.Username,
			Avatar:   b.Author.Avatar,
		},
	}
}
