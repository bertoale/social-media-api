package post

import (
	"fmt"
	"go-sosmed/pkg/utils"
)

type Service interface {
	Create(req *PostRequest, authorID uint) (*PostResponse, error)
	GetByID(postID uint) (*PostResponse, error)
	Update(userID, postID uint, req *UpdatePostRequest) (*PostResponse, error)
	Delete(postID, UserID uint, userRole string) error
	GetAll() ([]*PostResponse, error)
	GetAllUnarchived() ([]*PostResponse, error)
	GetPostsByAuthor(authorID uint) ([]*PostResponse, error)
	GetPostsByCurrentUser(authorID uint) ([]*PostResponse, error)
	Archive(postID, userID uint) error
	Unarchive(postID, userID uint) error
	GetPostsByFollowing(userID uint) ([]*PostResponse, error)
}

type service struct {
	repo Repository
}

// GetPostsByFollowing implements Service.
func (s *service) GetPostsByFollowing(userID uint) ([]*PostResponse, error) {
	posts, err := s.repo.FindByFollowing(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}
	var responses []*PostResponse
	for _, b := range posts {
		responses = append(responses, ToPostResponse(b))
	}
	return responses, nil
}

// GetPostsByCurrentUser implements Service.
func (s *service) GetPostsByCurrentUser(authorID uint) ([]*PostResponse, error) {
	posts, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}
	var responses []*PostResponse
	for _, b := range posts {
		if b.AuthorID == authorID && !b.Archived {
			responses = append(responses, ToPostResponse(b))
		}
	}
	return responses, nil
}

// GetAllUnarchived implements Service.
func (s *service) GetAllUnarchived() ([]*PostResponse, error) {
	posts, err := s.repo.FindAllUnarchived()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}
	var responses []*PostResponse
	for _, b := range posts {
		responses = append(responses, ToPostResponse(b))
	}
	return responses, nil
}

// Archive implements Service.
func (s *service) Archive(postID uint, userID uint) error {
	post, err := s.repo.FindByID(uint(postID))
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}
	if post.AuthorID != userID {
		return fmt.Errorf("unauthorized to archive this post")
	}
	if post.Archived {
		return fmt.Errorf("post is already archived")
	}
	return s.repo.Archive(postID)
}

// Unarchive implements Service.
func (s *service) Unarchive(postID uint, userID uint) error {
	post, err := s.repo.FindByID(uint(postID))
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}
	if post.AuthorID != userID {
		return fmt.Errorf("unauthorized to unarchive this post")
	}
	if !post.Archived {
		return fmt.Errorf("post is not archived")
	}
	return s.repo.Unarchive(postID)
}

// Create implements Service.
func (s *service) Create(req *PostRequest, authorID uint) (*PostResponse, error) {
	post := &Post{
		Title:    req.Title,
		Content:  req.Content,
		Image:    req.Image,
		AuthorID: authorID,
	}
	if err := s.repo.Create(post); err != nil {
		return nil, err
	}
	return ToPostResponse(post), nil
}

// Delete implements Service.
func (s *service) Delete(postID uint, userID uint, userRole string) error {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return fmt.Errorf("post not found: %w", err)
	}

	if userRole != "admin" && post.AuthorID != userID {
		return fmt.Errorf("unauthorized to delete this post")
	}

	// Hapus gambar jika ada
	if post.Image != "" {
		if err := utils.DeleteFile(post.Image); err != nil {
			fmt.Printf("Warning: failed to delete post image: %v\n", err)
		}
	}

	return s.repo.Delete(postID)
}

// GetAll implements Service.
func (s *service) GetAll() ([]*PostResponse, error) {
	posts, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}
	var responses []*PostResponse
	for _, b := range posts {
		responses = append(responses, ToPostResponse(b))
	}
	return responses, nil
}

// GetPostsByAuthor implements Service.
func (s *service) GetPostsByAuthor(authorID uint) ([]*PostResponse, error) {
	posts, err := s.repo.FindAllUnarchived()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %w", err)
	}
	var responses []*PostResponse
	for _, b := range posts {
		if b.AuthorID == authorID {
			responses = append(responses, ToPostResponse(b))
		}
	}
	return responses, nil
}

// GetByID implements Service.
func (s *service) GetByID(postID uint) (*PostResponse, error) {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	return ToPostResponse(post), nil
}

// Update implements Service.
func (s *service) Update(userID, postID uint, req *UpdatePostRequest) (*PostResponse, error) {
	post, err := s.repo.FindByID(postID)
	if err != nil {
		return nil, fmt.Errorf("post not found: %w", err)
	}
	// Update only fields that are not nil
	if req.Title != nil {
		post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Archived != nil {
		post.Archived = *req.Archived
	}

	if post.AuthorID != userID {
		return nil, fmt.Errorf("unauthorized to update this post")
	}
	post.Edited = true

	if err := s.repo.Update(post); err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}
	return ToPostResponse(post), nil
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
