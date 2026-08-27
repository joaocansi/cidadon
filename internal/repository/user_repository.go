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

type UserRepositoryImpl struct {
	*database.BaseRepository
}

func NewUserRepository(baseRepository *database.BaseRepository) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user repository.CreateUserData) (*entity.User, error) {
	userModel := model.User{
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Role:     user.Role,
	}
	if err := r.GetDB(ctx).Create(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return userModel.ToDomain(), nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User
	if err := r.GetDB(ctx).First(&userModel, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return userModel.ToDomain(), nil
}
