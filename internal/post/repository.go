package post

import "gorm.io/gorm"

type Repository interface {
	Create(post *Post) error
	FindByID(id uint) (*Post, error)
	FindDetailByID(id, userID uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint) error
	FindAll() ([]*Post, error)
	FindAllUnarchived() ([]*Post, error)
	Archive(id uint) error
	Unarchive(id uint) error
	FindByFollowing(userID uint) ([]*Post, error)
	FindByCurrentUser(
		authorID uint,
		currentUserID uint,
	) ([]*Post, error)
	FindPostsLikedByUser(userID uint) ([]Post, error)
}

type repository struct {
	db *gorm.DB
}

// GetPostsLikedByUser implements Repository.
func (r *repository) FindPostsLikedByUser(userID uint) ([]Post, error) {
	var posts []Post

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
		// Order("likes.created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// FindByCurrentUser implements Repository.
func (r *repository) FindByCurrentUser(
	authorID uint,
	currentUserID uint,
) ([]*Post, error) {

	var posts []*Post

	err := r.db.
		Table("posts").
		Select(`
			posts.*,
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) AS like_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
			EXISTS (
				SELECT 1 FROM likes
				WHERE likes.post_id = posts.id
				AND likes.user_id = ?
			) AS is_liked
		`, currentUserID).
		Where("posts.author_id = ?", authorID).
		Preload("Author").
		Order("posts.created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// FindByFollowing implements Repository.
func (r *repository) FindByFollowing(userID uint) ([]*Post, error) {
	var posts []*Post

	err := r.db.
		Model(&Post{}).
		Select(`
			posts.*,
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) AS like_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
			EXISTS (
				SELECT 1 FROM likes
				WHERE likes.post_id = posts.id
				AND likes.user_id = ?
			) AS is_liked
		`, userID).
		Joins("JOIN follows ON follows.following_id = posts.author_id").
		Where("follows.follower_id = ? AND posts.archived = ?", userID, false).
		Preload("Author").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// FindAllUnarchived implements Repository.
func (r *repository) FindAllUnarchived() ([]*Post, error) {
	var posts []*Post
	if err := r.db.Preload("Author").Where("archived = ?", false).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// Unarchive implements Repository.
func (r *repository) Unarchive(id uint) error {
	return r.db.Model(&Post{}).Where("id = ?", id).Update("archived", false).Error
}

// Archive implements Repository.
func (r *repository) Archive(id uint) error {
	return r.db.Model(&Post{}).Where("id = ?", id).Update("archived", true).Error
}

// Create implements Repository.
func (r *repository) Create(post *Post) error {
	return r.db.Create(post).Error
}

// Delete implements Repository.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Post{}, id).Error
}

// FindAll implements Repository.
func (r *repository) FindAll() ([]*Post, error) {
	var posts []*Post
	if err := r.db.Preload("Author").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*Post, error) {
	var post Post
	if err := r.db.Preload("Author").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// FindDetailByID implements Repository.
func (r *repository) FindDetailByID(id, userID uint) (*Post, error) {
	var post Post

	err := r.db.
		Model(&Post{}).
		Select(`
			posts.*,
			(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) AS like_count,
			(SELECT COUNT(*) FROM comments WHERE comments.post_id = posts.id) AS comment_count,
			EXISTS (
				SELECT 1 FROM likes
				WHERE likes.post_id = posts.id
				AND likes.user_id = ?
			) AS is_liked
		`, userID).
		Preload("Author").
		First(&post, id).Error

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Update implements Repository.
func (r *repository) Update(post *Post) error {
	return r.db.Save(post).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
