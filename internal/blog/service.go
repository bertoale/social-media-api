package blog

import (
	"fmt"
	"go-sosmed/pkg/utils"
)

type Service interface {
	Create(req *BlogRequest, authorID uint) (*BlogResponse, error)
	GetByID(BlogID uint) (*BlogResponse, error)
	Update(userID, BlogID uint, req *UpdateBlogRequest) (*BlogResponse, error)
	Delete(blogID, UserID uint) error
	GetAll() ([]*BlogResponse, error)
	GetBlogsByAuthor(authorID uint) ([]*BlogResponse, error)
}

type service struct {
	repo Repository
}

// Create implements Service.
func (s *service) Create(req *BlogRequest, authorID uint) (*BlogResponse, error) {
	blog := &Blog{
		Title:    req.Title,
		Content:  req.Content,
		Image:    req.Image,
		AuthorID: authorID,
	}
	if err := s.repo.Create(blog); err != nil {
		return nil, err
	}
	return ToBlogResponse(blog), nil
}

// Delete implements Service.
func (s *service) Delete(blogID uint, UserID uint) error {
	blog, err := s.repo.FindByID(uint(blogID))
	if err != nil {
		return fmt.Errorf("blog not found: %w", err)
	}
	if blog.AuthorID != UserID {
		return fmt.Errorf("unauthorized to delete this blog")
	}

	// Hapus gambar jika ada
	if blog.Image != "" {
		imagePath := utils.GetFilePath(blog.Image)
		if err := utils.DeleteFile(imagePath); err != nil {
			// Log error tapi tetap lanjutkan delete blog
			fmt.Printf("Warning: failed to delete blog image: %v\n", err)
		}
	}

	return s.repo.Delete(blogID)
}

// GetAll implements Service.
func (s *service) GetAll() ([]*BlogResponse, error) {
	blogs, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve blogs: %w", err)
	}
	var responses []*BlogResponse
	for _, b := range blogs {
		responses = append(responses, ToBlogResponse(b))
	}
	return responses, nil
}

// GetBlogsByAuthor implements Service.
func (s *service) GetBlogsByAuthor(authorID uint) ([]*BlogResponse, error) {
	blogs, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve blogs: %w", err)
	}
	var responses []*BlogResponse
	for _, b := range blogs {
		if b.AuthorID == authorID {
			responses = append(responses, ToBlogResponse(b))
		}
	}
	return responses, nil
}

// GetByID implements Service.
func (s *service) GetByID(BlogID uint) (*BlogResponse, error) {
	blog, err := s.repo.FindByID(BlogID)
	if err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}
	return ToBlogResponse(blog), nil
}

// Update implements Service.
func (s *service) Update(userID, BlogID uint, req *UpdateBlogRequest) (*BlogResponse, error) {
	blog, err := s.repo.FindByID(BlogID)
	if err != nil {
		return nil, fmt.Errorf("blog not found: %w", err)
	}
	// Update only fields that are not nil
	if req.Title != nil {
		blog.Title = *req.Title
	}
	if req.Content != nil {
		blog.Content = *req.Content
	}
	if req.Archived != nil {
		blog.Archived = *req.Archived
	}

	if blog.AuthorID != userID {
		return nil, fmt.Errorf("unauthorized to update this blog")
	}
	blog.Edited = true

	if err := s.repo.Update(blog); err != nil {
		return nil, fmt.Errorf("failed to update blog: %w", err)
	}
	return ToBlogResponse(blog), nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
