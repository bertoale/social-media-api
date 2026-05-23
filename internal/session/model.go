package session

import (
	"time"
)

// Session represents a user session stored in database
type Session struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	Token        string    `gorm:"unique;not null;index"`
	IPAddress    string    `gorm:"type:varchar(45)"` // IPv6 max length
	UserAgent    string    `gorm:"type:text"`
	ExpiresAt    time.Time `gorm:"not null;index"`
	LastActiveAt time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

// TableName specifies the table name for the Session model
func (Session) TableName() string {
	return "sessions"
}

// SessionResponse represents the session data returned to client
type SessionResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	ExpiresAt    time.Time `json:"expires_at"`
	LastActiveAt time.Time `json:"last_active_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToSessionResponse converts Session to SessionResponse
func ToSessionResponse(s *Session) *SessionResponse {
	return &SessionResponse{
		ID:           s.ID,
		UserID:       s.UserID,
		IPAddress:    s.IPAddress,
		UserAgent:    s.UserAgent,
		ExpiresAt:    s.ExpiresAt,
		LastActiveAt: s.LastActiveAt,
		CreatedAt:    s.CreatedAt,
	}
}

// ToSessionResponseList converts a slice of Session to SessionResponse slice
func ToSessionResponseList(sessions []*Session) []*SessionResponse {
	responses := make([]*SessionResponse, len(sessions))
	for i, s := range sessions {
		responses[i] = ToSessionResponse(s)
	}
	return responses
}
