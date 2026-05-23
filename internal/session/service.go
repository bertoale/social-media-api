package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"go-sosmed/pkg/config"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionExpired       = errors.New("session expired")
	ErrInvalidSessionToken  = errors.New("invalid session token")
)

type Service interface {
	CreateSession(userID uint, ipAddress, userAgent string) (string, error)
	ValidateSession(token string) (*Session, error)
	RefreshSession(token string, ipAddress, userAgent string) (string, error)
	DeleteSession(token string) error
	DeleteAllUserSessions(userID uint) error
	CleanupExpiredSessions() error
	GetUserSessions(userID uint) ([]*SessionResponse, error)
}

type service struct {
	repo   Repository
	config *config.Config
}

func NewService(repo Repository, cfg *config.Config) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
}

func (s *service) CreateSession(userID uint, ipAddress, userAgent string) (string, error) {
	duration, err := time.ParseDuration(s.config.SessionExpires)
	if err != nil {
		duration = 24 * time.Hour
	}

	token, err := s.generateSecureToken()
	if err != nil {
		return "", err
	}

	session := &Session{
		UserID:       userID,
		Token:        token,
		IPAddress:    ipAddress,
		UserAgent:    userAgent,
		ExpiresAt:    time.Now().Add(duration),
		LastActiveAt: time.Now(),
	}

	if err := s.repo.Create(session); err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) ValidateSession(token string) (*Session, error) {
	if token == "" {
		return nil, ErrInvalidSessionToken
	}

	session, err := s.repo.FindByToken(token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	if session.ExpiresAt.Before(time.Now()) {
		_ = s.repo.Delete(token)
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *service) RefreshSession(token string, ipAddress, userAgent string) (string, error) {
	session, err := s.ValidateSession(token)
	if err != nil {
		return "", err
	}

	session.IPAddress = ipAddress
	session.UserAgent = userAgent
	session.LastActiveAt = time.Now()

	duration, err := time.ParseDuration(s.config.SessionExpires)
	if err != nil {
		duration = 24 * time.Hour
	}
	session.ExpiresAt = time.Now().Add(duration)

	if err := s.repo.UpdateLastActive(session); err != nil {
		return "", err
	}

	return session.Token, nil
}

func (s *service) DeleteSession(token string) error {
	return s.repo.Delete(token)
}

func (s *service) DeleteAllUserSessions(userID uint) error {
	return s.repo.DeleteByUserID(userID)
}

func (s *service) CleanupExpiredSessions() error {
	return s.repo.DeleteExpired()
}

func (s *service) GetUserSessions(userID uint) ([]*SessionResponse, error) {
	sessions, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return ToSessionResponseList(sessions), nil
}

func (s *service) generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
