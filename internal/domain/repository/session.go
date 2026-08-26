package repository

import (
	"cidadon/internal/domain/entity"
	"context"
	"time"
)

type CreateSessionData struct {
	RefreshTokenHash string
	UserID           uint
	ExpiresAt        time.Time
	IpAddress        string
	UserAgent        string
}

type SessionRepository interface {
	FindByRefreshTokenHash(ctx context.Context, hash string) (*entity.Session, error)
	Create(ctx context.Context, session CreateSessionData) (*entity.Session, error)
}
