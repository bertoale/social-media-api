package like

import (
	"errors"
	"fmt"
	"go-sosmed/internal/post"
	"go-sosmed/internal/user"

	"gorm.io/gorm"
)

type Service interface {
	LikePost(userID, postID uint) error
	UnlikePost(userID, postID uint) error
	IsPostLiked(userID, postID uint) (bool, error)
	GetPostsLikedByUser(userID uint) ([]post.PostResponse, error)
}

type service struct {
	repo Repository
}

func (s *service) GetPostsLikedByUser(userID uint) ([]post.PostResponse, error) {
	posts, err := s.repo.GetPostsLikedByUser(userID)
	if err != nil {
		return nil, err
	}

	var resp []post.PostResponse
	for _, b := range posts {
		resp = append(resp, post.PostResponse{
			ID:      b.ID,
			Title:   b.Title,
			Content: b.Content,
			Image:   b.Image,
			Author: user.AuthorResponse{
				ID:       b.Author.ID,
				Username: b.Author.Username,
				Avatar:   b.Author.Avatar,
			},
		})
	}

	return resp, nil
}

func (s *service) IsPostLiked(userID uint, postID uint) (bool, error) {
	_, err := s.repo.GetByUserAndPost(userID, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *service) LikePost(userID uint, postID uint) error {
	existing, err := s.repo.GetByUserAndPost(userID, postID)
	if err == nil && existing != nil {
		return fmt.Errorf("already liked")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed checking existing like: %w", err)
	}
	like := &Like{
		UserID: userID,
		PostID: postID,
	}
	return s.repo.Create(like)

}

func (s *service) UnlikePost(userID uint, postID uint) error {
	like, err := s.repo.GetByUserAndPost(userID, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("not liked yet")
		}
		return fmt.Errorf("failed retrieving like: %w", err)
	}

	if err := s.repo.Delete(like.ID); err != nil {
		return fmt.Errorf("failed to unlike: %w", err)
	}

	return nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
