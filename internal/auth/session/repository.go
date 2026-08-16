package session

import (
	"cidadon/internal/app/shared/database"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewSessionRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *RepositoryImpl {
	return &RepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("SessionRepository"),
	}
}

func (r *RepositoryImpl) FindByRefreshTokenHash(ctx context.Context, hash string) (Session, error) {
	var session Model
	err := r.GetDB(ctx).Where("refresh_token_hash = ?", hash).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Session{}, database.ErrDBNotFound
		}
		return Session{}, database.ErrDBInternal
	}
	return session.ToDomain(), nil
}

func (r *RepositoryImpl) Create(ctx context.Context, session Session) error {
	err := r.GetDB(ctx).Create(&session).Error
	if err != nil {
		return database.ErrDBInternal
	}
	return nil
}
