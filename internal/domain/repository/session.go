package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type SessionRepository interface {
	FindByRefreshTokenHash(ctx context.Context, hash string) (entity.Session, error)
	Create(ctx context.Context, session entity.Session) error
}
