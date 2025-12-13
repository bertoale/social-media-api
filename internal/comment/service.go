package comment

import (
	"fmt"
)

type Service interface {
	//main
	CreateComment(comment *Comment) (*Comment, error)
	ReplyToComment(userID uint, parentID uint, blogID uint, content string) (*Comment, error)
	UpdateComment(userID uint, commentID uint, req UpdateCommentRequest) (*Comment, error)
	DeleteComment(userID uint, commentID uint) error
	GetCommentTree(blogID uint) ([]Comment, error)
	GetReplies(commentID uint) ([]Comment, error)
	GetByID(commentID uint) (*Comment, error)
}

type service struct {
	repo Repository
}

// GetByID implements Service.
func (s *service) GetByID(commentID uint) (*Comment, error) {
	return s.repo.GetByID(commentID)
}

func (s *service) CreateComment(comment *Comment) (*Comment, error) {
	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return s.repo.GetByID(comment.ID) 
}

func (s *service) DeleteComment(userID uint, commentID uint) error {
	isOwner, err := s.repo.IsOwner(commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized")
	}

	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	return s.repo.Delete(comment)
}

func (s *service) GetCommentTree(blogID uint) ([]Comment, error) {
	return s.repo.GetCommentTree(blogID)
}

func (s *service) GetReplies(commentID uint) ([]Comment, error) {
	return s.repo.GetReplies(commentID)
}

func (s *service) ReplyToComment(userID uint, targetID uint, blogID uint, content string) (*Comment, error) {

	target, err := s.repo.GetByID(targetID)
	if err != nil {
		return nil, fmt.Errorf("target comment not found")
	}

	var parentID uint
	if target.ParentID == nil {
		// reply ke root → parent = root
		parentID = target.ID
	} else {
		// reply ke reply → tetap parent = root
		parentID = *target.ParentID
	}

	replyToUserID := target.UserID

	reply := &Comment{
		Content:       content,
		UserID:        userID,
		BlogID:        blogID,
		ParentID:      &parentID,
		ReplyToUserID: &replyToUserID,
	}

	if err := s.repo.Create(reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *service) UpdateComment(userID uint, commentID uint, req UpdateCommentRequest) (*Comment, error) {
	isOwner, err := s.repo.IsOwner(commentID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized")
	}

	comment, err := s.repo.GetByID(commentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}

	if req.Content != nil {
		comment.Content = *req.Content
	}

	if req.Edited != nil {
		comment.Edited = *req.Edited
	} else {
		comment.Edited = true
	}

	if err := s.repo.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}
	return comment, nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
