package user

import (
	"fmt"
	"go-sosmed/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req *RegisterRequest) (*UserResponse, error)
	Login(req *LoginRequest) (string, *UserResponse, error)
	UpdateProfile(userID uint, req *UpdateProfileRequest) (*UserResponse, error)
	GetUserByID(userID uint) (*UserResponse, error)
	GenerateToken(user *User) (string, error)
	GetUserByUsername(username string) (*UserResponse, error)
}

type service struct {
	repo Repository
	cfg  *config.Config
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
