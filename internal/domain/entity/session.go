package domain

import (
	"time"
)

type Session struct {
	ID               uint       `json:"id"`
	RefreshTokenHash string     `json:"refresh_token_hash"`
	UserID           uint       `json:"user_id"`
	User             *User      `json:"user,omitempty"`
	ExpiresAt        time.Time  `json:"expires_at"`
	IpAddress        string     `json:"ip_address"`
	UserAgent        string     `json:"user_agent"`
	RevokedAt        *time.Time `json:"revoked_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (s *Session) IsExpired() bool {
	if s == nil {
		return true
	}
	return time.Now().After(s.ExpiresAt)
}
