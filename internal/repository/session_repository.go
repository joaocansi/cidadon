package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewSessionRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("SessionRepository"),
	}
}

func (r *SessionRepositoryImpl) FindByRefreshTokenHash(ctx context.Context, hash string) (*entity.Session, error) {
	var sessionModel model.Session
	err := r.GetDB(ctx).Where("refresh_token_hash = ?", hash).First(&sessionModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrDBNotFound
		}
		return nil, repository.ErrDBInternal
	}
	return sessionModel.ToDomain(), nil
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, session repository.CreateSessionData) (*entity.Session, error) {
	sessionModel := model.Session{
		IpAddress:        session.IpAddress,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		UserAgent:        session.UserAgent,
	}
	err := r.GetDB(ctx).Create(&sessionModel).Error
	if err != nil {
		return nil, repository.ErrDBInternal
	}
	return sessionModel.ToDomain(), err
}
