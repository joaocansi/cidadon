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

type CouncillorRepositoryImpl struct {
	*database.BaseRepository
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
