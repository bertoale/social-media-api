package comment

import "gorm.io/gorm"

type Repository interface {
	//main
	Create(comment *Comment) error
	GetByID(id uint) (*Comment, error)
	GetRootCommentsByPostID(postID uint) ([]Comment, error)
	Update(comment *Comment) error
	Delete(comment *Comment) error
	CountByPostID(postID uint) (int64, error)
	//replies
	GetReplies(parentID uint) ([]Comment, error)
	//utils
	IsOwner(commentID uint, userID uint) (bool, error)
	GetCommentTree(postID uint) ([]Comment, error)
}

type repository struct {
	db *gorm.DB
}

// CountByPostID implements Repository.
func (r *repository) CountByPostID(postID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Comment{}).
		Where("post_id = ?", postID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create implements Repository.
func (r *repository) Create(comment *Comment) error {
	return r.db.Create(comment).Error
}

// Delete implements Repository.
func (r *repository) Delete(comment *Comment) error {
	return r.db.Delete(comment).Error
}

// GetByID implements Repository.
func (r *repository) GetByID(id uint) (*Comment, error) {
	var c Comment

	err := r.db.
		Preload("Post").
		Preload("Post.Author").
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
func (r *repository) GetCommentTree(postID uint) ([]Comment, error) {
	var comments []Comment

	err := r.db.
		Preload("User").
		Preload("ReplyToUser").
		Preload("Replies").
		Preload("Replies.User").
		Preload("Replies.ReplyToUser").
		Where("post_id = ?", postID).
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
		Preload("Post").
		Preload("Post.Author").
		Find(&replies).Error
	if err != nil {
		return nil, err
	}
	return replies, nil
}

// GetRootCommentsByPostID implements Repository.
func (r *repository) GetRootCommentsByPostID(postID uint) ([]Comment, error) {
	var comments []Comment
	err := r.db.
		Where("post_id = ? AND parent_id IS NULL", postID).
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
