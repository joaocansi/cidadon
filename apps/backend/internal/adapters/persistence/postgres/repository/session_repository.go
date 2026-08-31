package repository

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	*database.BaseRepository
}

func (r *SessionRepositoryImpl) RevokeByRefreshTokenHash(ctx context.Context, hash string) error {
	result := r.GetDB(ctx).Model(&model.Session{}).Where("refresh_token_hash = ? AND revoked_at IS NULL", hash).Update("revoked_at", time.Now())
	if result.Error != nil {
		return repository.NewDBError(repository.DBErrorInternal, result.Error)
	}
	return nil
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
		ExpiresAt:        session.ExpiresAt,
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
