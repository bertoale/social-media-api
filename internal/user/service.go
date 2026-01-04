package user

import (
	"fmt"
	"go-sosmed/pkg/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req *RegisterRequest) (*UserResponse, error)
	Login(req *LoginRequest) (string, *UserResponse, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserResponse, error)
	GetUserByID(userID uint) (*UserResponse, error)
	GetExploreUsers(currentUserID uint, limit, offset int) ([]UserResponse, error)
	GetUserDetailByUsername(username string, currentUserID uint) (*UserResponse, error)
	GenerateToken(user *User) (string, error)
	GetUserByUsername(username string) (*UserResponse, error)
	SearchUser(keyword string, currentUserID uint) ([]*UserResponse, error)
	GetCurrentUserDetail(currentUserID uint) (*UserResponse, error)
	GetUserFollowers(
		currentUserID uint,
		limit, offset int,
	) ([]*UserResponse, error)
	GetUserFollowings(
		currentUserID uint,
		limit, offset int,
	) ([]*UserResponse, error)
	Logout() error
}

type service struct {
	repo Repository
	cfg  *config.Config
}

// GetCurrentUserDetail implements Service.
func (s *service) GetCurrentUserDetail(currentUserID uint) (*UserResponse, error) {
	user, err := s.repo.FindCurrentUserDetail(currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user detail: %w", err)
	}
	return ToUserResponse(user), nil
}

// Logout implements Service.
func (s *service) Logout() error {
	return nil
}

// GetUserFollowers implements Service.
func (s *service) GetUserFollowers(currentUserID uint, limit int, offset int) ([]*UserResponse, error) {
	followers, err := s.repo.FindFollowerByUsers(currentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user followers: %w", err)
	}

	var responses []*UserResponse
	for _, u := range followers {
		responses = append(responses, ToUserResponse(u))
	}

	return responses, nil
}

// GetUserFollowings implements Service.
func (s *service) GetUserFollowings(currentUserID uint, limit int, offset int) ([]*UserResponse, error) {
	followings, err := s.repo.FindFollowingUsers(currentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user followings: %w", err)
	}

	var responses []*UserResponse
	for _, u := range followings {
		responses = append(responses, ToUserResponse(u))
	}

	return responses, nil
}

// SearchByUsername implements Service.
func (s *service) SearchUser(keyword string, currentUserID uint) ([]*UserResponse, error) {
	// 1. validasi keyword
	keyword = strings.TrimSpace(keyword)
	if len(keyword) < 1 {
		return []*UserResponse{}, nil
	}

	// 2. ambil data dari repo
	users, err := s.repo.SearchByUsername(keyword, currentUserID, 10)
	if err != nil {
		return nil, err
	}

	// 3. mapping ke response
	var responses []*UserResponse
	for _, u := range users {
		// opsional: skip diri sendiri
		if u.ID == currentUserID {
			continue
		}

		responses = append(responses, ToUserResponse(u))
	}

	return responses, nil
}

// GetUserDetailByID implements Service.
func (s *service) GetUserDetailByUsername(username string, currentUserID uint) (*UserResponse, error) {
	user, err := s.repo.FindUserDetailByUsername(username, currentUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user detail by ID: %w", err)
	}
	return ToUserResponse(user), nil
}

// GetExploreUsers implements Service.
func (s *service) GetExploreUsers(currentUserID uint, limit int, offset int) ([]UserResponse, error) {
	users, err := s.repo.FindExploreUsers(currentUserID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get explore users: %w", err)
	}
	var responses []UserResponse
	for _, u := range users {
		responses = append(responses, *ToUserResponse(&u))
	}
	return responses, nil
}

// GetUserByUsername implements Service.
func (s *service) GetUserByUsername(username string) (*UserResponse, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return ToUserResponse(user), nil
}

type Claims struct {
	ID   uint     `json:"id"`
	Role RoleType `json:"role"`
	jwt.RegisteredClaims
}

func (s *service) GenerateToken(user *User) (string, error) {
	duration, err := time.ParseDuration(s.cfg.JWTExpires)
	if err != nil {
		duration = 168 * time.Hour // default 7 days
	}
	claims := Claims{
		ID:   user.ID,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// GetUserByID implements Service.
func (s *service) GetUserByID(userID uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return ToUserResponse(user), nil
}

// Login implements Service.
func (s *service) Login(req *LoginRequest) (string, *UserResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", nil, fmt.Errorf("invalid password: %w", err)
	}
	token, err := s.GenerateToken(user)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}
	return token, ToUserResponse(user), nil
}

// Register implements Service.
func (s *service) Register(req *RegisterRequest) (*UserResponse, error) {
	existingEmail, _ := s.repo.FindByEmail(req.Email)
	if existingEmail != nil {
		return nil, fmt.Errorf("email already in use")
	}
	existingUsername, _ := s.repo.FindByUsername(req.Username)
	if existingUsername != nil {
		return nil, fmt.Errorf("username already in use")
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	if err := s.repo.Create(u); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return ToUserResponse(u), nil
}

// UpdateProfile implements Service.
func (s *service) UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// update allowed fields
	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Avatar != nil {
		user.Avatar = *req.Avatar
	}
	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	existingUsername, err := s.repo.FindByUsername(user.Username)
	if err == nil && existingUsername != nil && existingUsername.ID != user.ID {
		return nil, fmt.Errorf("username already in use")
	}

	existingEmail, err := s.repo.FindByEmail(user.Email)
	if err == nil && existingEmail != nil && existingEmail.ID != user.ID {
		return nil, fmt.Errorf("email already in use")
	}

	if err := s.repo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return ToUserResponse(user), nil
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{repo: repo, cfg: cfg}
}
