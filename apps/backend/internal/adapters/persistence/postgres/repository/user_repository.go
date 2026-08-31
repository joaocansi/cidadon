package repository

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
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

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var item model.User
	if err := r.GetDB(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return item.ToDomain(), nil
}

func (r *UserRepositoryImpl) DeleteByID(ctx context.Context, id uint) error {
	result := r.GetDB(ctx).Unscoped().Delete(&model.User{}, id)
	if result.Error != nil {
		return repository.NewDBError(repository.DBErrorInternal, result.Error)
	}
	if result.RowsAffected == 0 {
		return repository.NewDBError(repository.DBErrorNotFound, gorm.ErrRecordNotFound)
	}
	return nil
}
