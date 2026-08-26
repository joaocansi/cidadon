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

type UserRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewUserRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user repository.CreateUserData) (*entity.User, error) {
	userModel := model.User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Role:     user.Role,
	}
	err := r.GetDB(ctx).Create(&userModel).Error
	if err != nil {
		r.Logger.Errorw("failed to create user", err)
		return nil, repository.ErrDBInternal
	}
	return userModel.ToDomain(), nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User
	err := r.GetDB(ctx).First(&userModel, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrDBNotFound
		}
		r.Logger.Errorw("failed to find user by email", "email", email)
		return nil, err
	}
	return userModel.ToDomain(), nil
}
