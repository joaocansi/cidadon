package user

import (
	"cidadon/internal/app/shared/database"
	"context"

	"go.uber.org/zap"
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
		return err
	}
	*user = *userModel.ToDomain()
	return nil
}

func (r *RepositoryImpl) FindByEmail(ctx context.Context, email string) (*User, error) {
	userModel := Model{}
	err := r.GetDB(ctx).First(&userModel, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return userModel.ToDomain(), nil
}
