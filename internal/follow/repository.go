package follow

import "gorm.io/gorm"

type Repository interface {
	Create(follow *Follow) error
	Delete(id uint) error
	FindByUserAndFollowed(userID uint, followedID uint) (*Follow, error)
	FindFollowersByUserID(userID uint) ([]*Follow, error)
	FindFollowingByUserID(userID uint) ([]*Follow, error)
}

type repository struct {
	db *gorm.DB
}

// FindFollowingByUserID implements Repository.
func (r *repository) FindFollowingByUserID(userID uint) ([]*Follow, error) {
	var follows []*Follow
	if err := r.db.Preload("Following").
		Where("follower_id = ?", userID).
		Find(&follows).Error; err != nil {
		return nil, err
	}
	return follows, nil
}

// FindFollowersByUserID implements Repository.
func (r *repository) FindFollowersByUserID(userID uint) ([]*Follow, error) {
	var follows []*Follow
	if err := r.db.Preload("Follower").
		Where("following_id = ?", userID).
		Find(&follows).Error; err != nil {
		return nil, err
	}
	return follows, nil
}

// Create implements Repository.
func (r *repository) Create(follow *Follow) error {
	return r.db.Create(follow).Error
}

// Delete implements Repository.
func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Follow{}, id).Error
}

// FindByUserAndFollowed implements Repository.
func (r *repository) FindByUserAndFollowed(followerID uint, followingID uint) (*Follow, error) {
	var follow Follow

	err := r.db.
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		First(&follow).Error

	return &follow, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
