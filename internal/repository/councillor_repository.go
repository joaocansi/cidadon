package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"

	"go.uber.org/zap"
)

type CouncillorRepository interface {
	Create(ctx context.Context, user *entity.Councillor) error
}

type CouncillorRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewCouncillorRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *CouncillorRepositoryImpl {
	return &CouncillorRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("CouncillorRepository"),
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
	}
	err := c.GetDB(ctx).Create(&councillorModel).Error
	if err != nil {
		c.Logger.Errorw("failed to create model", "error", err)
		return nil, repository.ErrDBInternal
	}
	return councillorModel.ToDomain(), nil
}
