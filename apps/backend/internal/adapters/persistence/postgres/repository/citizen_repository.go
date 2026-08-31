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

type CitizenRepository interface {
	Create(ctx context.Context, citizen *entity.Citizen) error
}

type CitizenRepositoryImpl struct {
	*database.BaseRepository
}

func NewCitizenRepository(baseRepository *database.BaseRepository) *CitizenRepositoryImpl {
	return &CitizenRepositoryImpl{
		BaseRepository: baseRepository,
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
		City:  citizen.City,
		State: citizen.State,
	}
	if err := r.GetDB(ctx).Create(&citizenModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return citizenModel.ToDomain(), nil
}
