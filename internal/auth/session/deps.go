package session

import "context"

type Repository interface {
	FindByRefreshTokenHash(ctx context.Context, hash string) (Session, error)
	Create(ctx context.Context, session Session) error
}
