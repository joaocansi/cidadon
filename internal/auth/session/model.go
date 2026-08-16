package session

import (
	"cidadon/internal/user"
	"time"

	"gorm.io/gorm"
)

type Model struct {
	gorm.Model
	RefreshTokenHash string
	UserID           uint
	User             user.Model
	ExpiresAt        time.Time
	IpAddress        string
	UserAgent        string
	RevokedAt        *time.Time
}

func (*Model) TableName() string {
	return "sessions"
}

func (s *Model) ToDomain() Session {
	return Session{
		RefreshTokenHash: s.RefreshTokenHash,
		UserID:           s.UserID,
		User:             s.User.ToDomain(),
		ExpiresAt:        s.ExpiresAt,
		IpAddress:        s.IpAddress,
		UserAgent:        s.UserAgent,
		RevokedAt:        s.RevokedAt,
	}
}
