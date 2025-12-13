package comment

import "gorm.io/gorm"

type Repository interface {
	//main
	Create(comment *Comment) error
	GetByID(id uint) (*Comment, error)
	GetRootCommentsByBlogID(blogID uint) ([]Comment, error)
	Update(comment *Comment) error
	Delete(comment *Comment) error
	//replies
	GetReplies(parentID uint) ([]Comment, error)
	//utils
	IsOwner(commentID uint, userID uint) (bool, error)
	GetCommentTree(blogID uint) ([]Comment, error)
}

type repository struct {
	db *gorm.DB
}

// Create implements Repository.
func (r *repository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

// Delete implements Repository.
func (r *repository) Delete(comment *Comment) error {
	// hapus semua child (replies)
	var replies []Comment
	r.db.Where("parent_id = ?", comment.ID).Find(&replies)

	for _, reply := range replies {
		if err := r.Delete(&reply); err != nil {
			return err
		}
	}

	// hapus comment utama
	return r.db.Delete(comment).Error
}

// GetByID implements Repository.
func (r *repository) GetByID(id uint) (*Comment, error) {
	var c Comment

	err := r.db.
		Preload("User").
		Preload("ReplyToUser").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Replies.User").
		Preload("Replies.ReplyToUser").
		First(&c, id).Error

	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetCommentTree implements Repository.
func (r *repository) GetCommentTree(blogID uint) ([]Comment, error) {
	var comments []Comment

	err := r.db.
		Preload("User").
		Preload("ReplyToUser").
		Preload("Replies").
		Preload("Replies.User").
		Preload("Replies.ReplyToUser").
		Where("blog_id = ?", blogID).
		Where("parent_id IS NULL").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

// GetReplies implements Repository.
func (r *repository) GetReplies(parentID uint) ([]Comment, error) {
	var replies []Comment
	err := r.db.
		Where("parent_id = ?", parentID).
		Preload("User").
		Preload("ReplyToUser").
		Preload("Blog").
		Preload("Blog.Author").
		Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

// GetRootCommentsByBlogID implements Repository.
func (r *repository) GetRootCommentsByBlogID(blogID uint) ([]Comment, error) {
	var comments []Comment
	err := r.db.
		Where("blog_id = ? AND parent_id IS NULL", blogID).
		Preload("User").
		Preload("Replies").
		Find(&comments).Error

	if err != nil {
		return nil, err
	}
	return comments, nil
}

// IsOwner implements Repository.
func (r *repository) IsOwner(commentID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Comment{}).
		Where("id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Update implements Repository.
func (r *repository) Update(comment *Comment) error {
	return r.db.Save(comment).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}
