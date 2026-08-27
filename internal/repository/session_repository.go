package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	*database.BaseRepository
}

func NewSessionRepository(baseRepository *database.BaseRepository) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (r *SessionRepositoryImpl) FindByRefreshTokenHash(ctx context.Context, hash string) (*entity.Session, error) {
	var sessionModel model.Session
	err := r.GetDB(ctx).Where("refresh_token_hash = ?", hash).First(&sessionModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
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
	if err := r.GetDB(ctx).Create(&sessionModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return sessionModel.ToDomain(), nil
}
