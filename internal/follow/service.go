package follow

import (
	"errors"

	"gorm.io/gorm"
)

type Service interface {
	FollowUser(followerID, followingID uint) error
	UnfollowUser(followerID, followingID uint) error
	GetFollowers(userID uint) ([]*FollowerResponse, error)
	GetFollowing(userID uint) ([]*FollowingResponse, error)
}

type service struct {
	repo Repository
}

// FollowUser implements Service.
func (s *service) FollowUser(followerID uint, followingID uint) error {

	// Validate: cannot follow yourself
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	existing, err := s.repo.FindByUserAndFollowed(followerID, followingID)
	// Jika data sudah ada → sudah follow
	if err == nil && existing != nil {
		return errors.New("already following this user")
	}

	// Jika error tetapi bukan "record not found" → DB error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	follow := &Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	// Lebih aman: DB akan menolak jika duplicate index
	return s.repo.Create(follow)
}

// GetFollowers implements Service.
func (s *service) GetFollowers(userID uint) ([]*FollowerResponse, error) {
	post, err := s.repo.FindFollowersByUserID(userID)
	if err != nil {
		return nil, err
	}
	var followers []*FollowerResponse
	for _, f := range post {
		followers = append(followers, ToFollowerResponse(f))
	}
	return followers, nil
}

// GetFollowing implements Service.
func (s *service) GetFollowing(userID uint) ([]*FollowingResponse, error) {
	follows, err := s.repo.FindFollowingByUserID(userID)
	if err != nil {
		return nil, err
	}
	var following []*FollowingResponse
	for _, f := range follows {
		following = append(following, ToFollowingResponse(f))
	}
	return following, nil
}

// UnfollowUser implements Service.
func (s *service) UnfollowUser(followerID uint, followingID uint) error {
	existing, err := s.repo.FindByUserAndFollowed(followerID, followingID)
	// if data not found → not following
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("you are not following this user")
		}
		return err
	}
	// Data found, proceed to delete
	return s.repo.Delete(existing.ID)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
