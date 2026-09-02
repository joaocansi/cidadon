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

type CouncillorRepositoryImpl struct {
	*database.BaseRepository
}

func (c *CouncillorRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*entity.Councillor, error) {
	var item model.Councillor
	if err := c.GetDB(ctx).Preload("User").First(&item, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return item.ToDomain(), nil
}

func (c *CouncillorRepositoryImpl) UpdateByUserID(ctx context.Context, userID uint, data repository.UpdateCouncillorData) (*entity.Councillor, error) {
	var item model.Councillor
	result := c.GetDB(ctx).Model(&item).Where("user_id = ?", userID).Updates(map[string]any{
		"party": data.Party,
		"city":  data.City,
		"state": data.State,
	})
	if result.Error != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, repository.NewDBError(repository.DBErrorNotFound, gorm.ErrRecordNotFound)
	}
	if err := c.GetDB(ctx).Preload("User").First(&item, "user_id = ?", userID).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return item.ToDomain(), nil
}

func NewCouncillorRepository(baseRepository *database.BaseRepository) *CouncillorRepositoryImpl {
	return &CouncillorRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (c *CouncillorRepositoryImpl) Create(ctx context.Context, data repository.CreateCouncillorData) (*entity.Councillor, error) {
	councillorModel := model.Councillor{
		User: &model.User{
			Name:     data.Name,
			Password: data.Password,
			Email:    data.Email,
			Role:     string(entity.CouncillorUser),
		},
		ImageURL: data.ImageURL,
		Party:    data.Party,
		City:     data.City,
		State:    data.State,
	}
	if err := c.GetDB(ctx).Create(&councillorModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return councillorModel.ToDomain(), nil
}
