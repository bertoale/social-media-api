package like

import (
	"errors"
	"fmt"
	"go-sosmed/internal/blog"
	"go-sosmed/internal/user"

	"gorm.io/gorm"
)

type Service interface {
	LikeBlog(userID, blogID uint) error
	UnlikeBlog(userID, blogID uint) error
	IsBlogLiked(userID, blogID uint) (bool, error)
	GetBlogsLikedByUser(userID uint) ([]blog.BlogResponse, error)
}

type service struct {
	repo Repository
}

func (s *service) GetBlogsLikedByUser(userID uint) ([]blog.BlogResponse, error) {
	blogs, err := s.repo.GetBlogLikedByUser(userID)
	if err != nil {
		return nil, err
	}

	var resp []blog.BlogResponse
	for _, b := range blogs {
		resp = append(resp, blog.BlogResponse{
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

func (s *service) IsBlogLiked(userID uint, blogID uint) (bool, error) {
	_, err := s.repo.GetByUserAndBlog(userID, blogID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *service) LikeBlog(userID uint, blogID uint) error {
	existing, err := s.repo.GetByUserAndBlog(userID, blogID)
	if err == nil && existing != nil {
		return fmt.Errorf("already liked")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed checking existing like: %w", err)
	}
	like := &Like{
		UserID: userID,
		BlogID: blogID,
	}
	return s.repo.Create(like)

}

func (s *service) UnlikeBlog(userID uint, blogID uint) error {
	like, err := s.repo.GetByUserAndBlog(userID, blogID)
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
