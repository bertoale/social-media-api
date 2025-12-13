package comment

import "go-sosmed/internal/user"

func ToCommentResponse(c *Comment) CommentResponse {
	resp := CommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		Edited:    c.Edited,
		User: user.AuthorResponse{
			ID:       c.User.ID,
			Username: c.User.Username,
			Avatar:   c.User.Avatar,
		},
	}

	if c.ReplyToUser != nil {
		resp.ReplyToUser = &user.AuthorResponse{
			ID:       c.ReplyToUser.ID,
			Username: c.ReplyToUser.Username,
		}
	}

	for _, r := range c.Replies {
		resp.Replies = append(resp.Replies, ToCommentResponse(&r))
	}

	return resp
}
