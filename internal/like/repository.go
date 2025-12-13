package like

import (
	"go-sosmed/internal/blog"

	"gorm.io/gorm"
)

type Repository interface {
	Create(like *Like) error
	Delete(id uint) error
	GetByUserAndBlog(userID uint, blogID uint) (*Like, error)
	GetBlogLikedByUser(userID uint) ([]blog.Blog, error)
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(like *Like) error {
	return r.db.Create(like).Error
}

// Delete implements Repository.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Like{}, id).Error
}

// GetBlogLikedByUser implements Repository.
func (r *repository) GetBlogLikedByUser(userID uint) ([]blog.Blog, error) {
	var blogs []blog.Blog
	err := r.db.Joins("JOIN likes ON likes.blog_id = blogs.id").
		Where("likes.user_id = ?", userID).
		Preload("Author").
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// GetByUserAndBlog implements Repository.
func (r *repository) GetByUserAndBlog(userID uint, blogID uint) (*Like, error) {
	var like Like
	if err := r.db.Where("user_id = ? AND blog_id = ?", userID, blogID).First(&like).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
