package comment

import (
	"errors"
	"fmt"
	"go-sosmed/internal/post"
)

type Service interface {
	//main
	CreateComment(comment *Comment) (*Comment, error)
	ReplyToComment(userID uint, parentID uint, postID uint, content string) (*Comment, error)
	UpdateComment(userID uint, commentID uint, req UpdateCommentRequest) (*Comment, error)
	DeleteComment(userID uint, commentID uint) error
	GetCommentTree(postID uint) ([]Comment, error)
	GetReplies(commentID uint) ([]Comment, error)
	GetByID(commentID uint) (*Comment, error)
	GetCommentCount(postID uint) (int64, error)
}

type service struct {
	commentRepo Repository
	postRepo    post.Repository
}

// GetCommentCount implements Service.
func (s *service) GetCommentCount(postID uint) (int64, error) {
	return s.commentRepo.CountByPostID(postID)
}

// GetByID implements Service.
func (s *service) GetByID(commentID uint) (*Comment, error) {
	return s.commentRepo.GetByID(commentID)
}

func (s *service) CreateComment(comment *Comment) (*Comment, error) {
	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return s.commentRepo.GetByID(comment.ID)
}

func (s *service) DeleteComment(userID uint, commentID uint) error {
	isOwner, err := s.commentRepo.IsOwner(commentID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized")
	}

	comment, err := s.commentRepo.GetByID(commentID)
	if err != nil {
		return fmt.Errorf("comment not found: %w", err)
	}

	return s.commentRepo.Delete(comment)
}

func (s *service) GetCommentTree(postID uint) ([]Comment, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if post.DeletedAt.Valid {
		return nil, errors.New("post not found")
	}

	comments, err := s.commentRepo.GetCommentTree(postID)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *service) GetReplies(commentID uint) ([]Comment, error) {
	return s.commentRepo.GetReplies(commentID)
}

func (s *service) ReplyToComment(userID uint, targetID uint, postID uint, content string) (*Comment, error) {

	target, err := s.commentRepo.GetByID(targetID)
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
		PostID:        postID,
		ParentID:      &parentID,
		ReplyToUserID: &replyToUserID,
	}

	if err := s.commentRepo.Create(reply); err != nil {
		return nil, err
	}

	return reply, nil
}

func (s *service) UpdateComment(userID uint, commentID uint, req UpdateCommentRequest) (*Comment, error) {
	isOwner, err := s.commentRepo.IsOwner(commentID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized")
	}

	comment, err := s.commentRepo.GetByID(commentID)
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

	if err := s.commentRepo.Update(comment); err != nil {
		return nil, fmt.Errorf("failed to update comment: %w", err)
	}
	return comment, nil
}

func NewService(commentRepo Repository, postRepo post.Repository) Service {
	return &service{
		commentRepo: commentRepo,
		postRepo:    postRepo}
}
