package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"

	"go.uber.org/zap"
)

type CitizenRepository interface {
	Create(ctx context.Context, citizen *entity.Citizen) error
}

type CitizenRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewCitizenRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *CitizenRepositoryImpl {
	return &CitizenRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("CitizenRepository"),
	}
}

func (r *CitizenRepositoryImpl) Create(ctx context.Context, citizen repository.CreateCitizenData) (*entity.Citizen, error) {
	citizenModel := model.Citizen{
		User: &model.User{
			Email:    citizen.Email,
			Password: citizen.Password,
			Name:     citizen.Name,
			Role:     string(entity.CitizenUser),
		},
	}
	if err := r.GetDB(ctx).Create(&citizenModel).Error; err != nil {
		r.Logger.Errorw("error while creating citizen", "error", err)
		return nil, repository.ErrDBInternal
	}
	return citizenModel.ToDomain(), nil
}
