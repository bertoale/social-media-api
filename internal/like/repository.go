package like

import (
	"go-sosmed/internal/post"

	"gorm.io/gorm"
)

type Repository interface {
	Create(like *Like) error
	Delete(id uint) error
	GetByUserAndPost(userID uint, postID uint) (*Like, error)
	GetPostsLikedByUser(userID uint) ([]post.Post, error)
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

// GetPostsLikedByUser implements Repository.
func (r *repository) GetPostsLikedByUser(
	userID uint,
) ([]post.Post, error) {

	var posts []post.Post

	err := r.db.
		Table("posts").
		Select(`
			posts.*,
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) AS like_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
			TRUE AS is_liked
		`).
		Joins("JOIN likes ON likes.post_id = posts.id").
		Where("likes.user_id = ? AND posts.archived = ?", userID, false).
		Preload("Author").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// GetByUserAndPost implements Repository.
func (r *repository) GetByUserAndPost(userID uint, postID uint) (*Like, error) {
	var like Like
	if err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error; err != nil {
		return nil, err
	}
	return &like, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
