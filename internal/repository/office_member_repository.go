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

type OfficeMemberRepositoryImpl struct {
	*database.BaseRepository
}

func NewOfficeMemberRepository(baseRepository *database.BaseRepository) *OfficeMemberRepositoryImpl {
	return &OfficeMemberRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (om *OfficeMemberRepositoryImpl) Create(ctx context.Context, officeMember repository.CreateOfficeMemberData) (*entity.OfficeMember, error) {
	officeMemberModel := model.OfficeMember{
		OfficeID: officeMember.OfficeID,
		User: &model.User{
			Name:     officeMember.Name,
			Email:    officeMember.Email,
			Password: officeMember.Password,
			Role:     string(entity.OfficeMemberUser),
		},
		ImageURL: officeMember.ImageURL,
	}
	if err := om.GetDB(ctx).Create(&officeMemberModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeMemberModel.ToDomain(), nil
}
