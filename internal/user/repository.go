package user

import (
	"cidadon/internal/app/database"
	apperrors "cidadon/internal/app/errors"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewUserRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *RepositoryImpl {
	return &RepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger,
	}
}

func (r *RepositoryImpl) Create(ctx context.Context, user *User) error {
	userModel := user.ToModel()
	err := r.GetDB(ctx).Create(&userModel).Error
	if err != nil {
		return apperrors.WrapDB(err, "failed to create user")
	}
	*user = *userModel.ToDomain()
	return nil
}

func (r *RepositoryImpl) FindByEmail(ctx context.Context, email string) (*User, error) {
	userModel := Model{}
	err := r.GetDB(ctx).First(&userModel, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrDBNotFound
		}
		return nil, apperrors.ErrDBInternal
	}
	return userModel.ToDomain(), nil
}
