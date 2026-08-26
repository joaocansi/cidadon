package model

import (
	"cidadon/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	RefreshTokenHash string
	UserID           uint
	User             *User
	ExpiresAt        time.Time
	IpAddress        string
	UserAgent        string
	RevokedAt        *time.Time
}

func (s *Session) ToDomain() *entity.Session {
	if s == nil {
		return nil
	}

	return &entity.Session{
		RefreshTokenHash: s.RefreshTokenHash,
		UserID:           s.UserID,
		User:             s.User.ToDomain(),
		ExpiresAt:        s.ExpiresAt,
		IpAddress:        s.IpAddress,
		UserAgent:        s.UserAgent,
		RevokedAt:        s.RevokedAt,
		UpdatedAt:        s.UpdatedAt,
		CreatedAt:        s.CreatedAt,
	}
}
