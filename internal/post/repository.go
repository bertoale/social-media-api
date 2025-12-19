package post

import "gorm.io/gorm"

type Repository interface {
	Create(post *Post) error
	FindByID(id uint) (*Post, error)
	Update(post *Post) error
	Delete(id uint) error
	FindAll() ([]*Post, error)
	FindAllUnarchived() ([]*Post, error)
	Archive(id uint) error
	Unarchive(id uint) error
	FindByFollowing(userID uint) ([]*Post, error)
}

type repository struct {
	db *gorm.DB
}

// FindByFollowing implements Repository.
func (r *repository) FindByFollowing(userID uint) ([]*Post, error) {
	var posts []*Post
	if err := r.db.
		Preload("Author").
		Joins("JOIN follows ON follows.following_id = posts.author_id").
Where("follows.follower_id = ? AND posts.archived = ?", userID, false).
		Find(&posts).Error; err != nil {
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

// Update implements Repository.
func (r *repository) Update(post *Post) error {
	return r.db.Save(post).Error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
