package session

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(session *Session) error
	FindByToken(token string) (*Session, error)
	FindByUserID(userID uint) ([]*Session, error)
	UpdateLastActive(session *Session) error
	Delete(token string) error
	DeleteExpired() error
	DeleteByUserID(userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(session *Session) error {
	return r.db.Create(session).Error
}

func (r *repository) FindByToken(token string) (*Session, error) {
	var session Session
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) FindByUserID(userID uint) ([]*Session, error) {
	var sessions []*Session
	err := r.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions).Error
	return sessions, err
}

func (r *repository) UpdateLastActive(session *Session) error {
	session.LastActiveAt = time.Now()
	return r.db.Save(session).Error
}

func (r *repository) Delete(token string) error {
	return r.db.Where("token = ?", token).Delete(&Session{}).Error
}

func (r *repository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
}

func (r *repository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&Session{}).Error
}
