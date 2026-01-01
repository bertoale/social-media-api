package user

import (
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	SearchByUsername(
		keyword string,
		currentUserID uint,
		limit int,
	) ([]*User, error)
	FindByUsernameOrEmail(identifier string) (*User, error)
	FindExploreUsers(currentUserID uint, limit, offset int) ([]User, error)
	FindUserDetailByUsername(username string, currentUserID uint) (*User, error)
	Update(user *User) error
	FindFollowingUsers(
		currentUserID uint,
		limit, offset int,
	) ([]*User, error)
	FindFollowerByUsers(
		currentUserID uint,
		limit, offset int,
	) ([]*User, error)
}

type repository struct {
	db *gorm.DB
}

// SearchByUsername implements Repository.
func (r *repository) SearchByUsername(
	keyword string,
	currentUserID uint,
	limit int,
) ([]*User, error) {

	var users []*User

	err := r.db.
		Table("users").
		Select(`
			users.*,
			EXISTS (
				SELECT 1
				FROM follows
				WHERE follows.follower_id = ?
				AND follows.following_id = users.id
			) AS is_followed,
			(SELECT COUNT(*) FROM follows WHERE follows.following_id = users.id) AS follower_count,
			(SELECT COUNT(*) FROM follows WHERE follows.follower_id = users.id) AS following_count
		`, currentUserID).
		Where("LOWER(users.username) LIKE ?", "%"+strings.ToLower(keyword)+"%").
		Where("users.role != ?", "admin").
		Limit(limit).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindExploreUsers implements Repository.
func (r *repository) FindExploreUsers(currentUserID uint, limit int, offset int) ([]User, error) {
	var users []User
	err := r.db.
		Model(&User{}).
		Select(`
			users.*,
			(SELECT COUNT(*) FROM follows WHERE follows.following_id = users.id) AS followers_count,
			(SELECT COUNT(*) FROM follows WHERE follows.follower_id = users.id) AS following_count,
			EXISTS (
				SELECT 1 FROM follows
				WHERE follows.follower_id = ?
				AND follows.following_id = users.id
			) AS is_followed
		`, currentUserID).
		Where("users.id != ?", currentUserID).
		Where("users.role != ?", "admin").
		Limit(limit).
		Offset(offset).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// FindByUsernameOrEmail implements Repository.
func (r *repository) FindByUsernameOrEmail(identifier string) (*User, error) {
	var user User
	if err := r.db.
		Where("username = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername implements Repository.
func (r *repository) FindByUsername(username string) (*User, error) {
	var user User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create implements Repository.
func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

// FindByEmail implements Repository.
func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID implements Repository.
func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update implements Repository.
func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

// FindUserDetailByID implements Repository.
func (r *repository) FindUserDetailByUsername(username string, currentUserID uint) (*User, error) {
	var user User

	err := r.db.
		Model(&User{}).
		Select(`
			users.*,
			(SELECT COUNT(*) FROM follows WHERE follows.following_id = users.id) AS follower_count,
			(SELECT COUNT(*) FROM follows WHERE follows.follower_id = users.id) AS following_count,
			EXISTS (
				SELECT 1 FROM follows
				WHERE follows.follower_id = ?
				AND follows.following_id = users.id
			) AS is_followed
		`, currentUserID).
		Where("users.username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) FindFollowingUsers(
	currentUserID uint,
	limit, offset int,
) ([]*User, error) {

	var users []*User

	err := r.db.
		Table("users").
		Select(`
			users.*,
			TRUE AS is_followed,
			(SELECT COUNT(*) FROM follows WHERE follows.following_id = users.id) AS follower_count,
			(SELECT COUNT(*) FROM follows WHERE follows.follower_id = users.id) AS following_count
		`).
		Joins(`
			JOIN follows 
			ON follows.following_id = users.id
			AND follows.follower_id = ?
		`, currentUserID).
		Where("users.role != ?", "admin").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) FindFollowerByUsers(
	currentUserID uint,
	limit, offset int,
) ([]*User, error) {

	var users []*User

	err := r.db.
		Table("users").
		Select(`
			users.*,
			TRUE AS is_followed,
			(SELECT COUNT(*) FROM follows WHERE follows.following_id = users.id) AS follower_count,
			(SELECT COUNT(*) FROM follows WHERE follows.follower_id = users.id) AS following_count
		`).
		Joins(`
			JOIN follows 
			ON follows.follower_id = users.id
			AND follows.following_id = ?
		`, currentUserID).
		Where("users.role != ?", "admin").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
