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

func (om *OfficeMemberRepositoryImpl) FindByUserID(ctx context.Context, userID uint) (*entity.OfficeMember, error) {
	var item model.OfficeMember
	if err := om.GetDB(ctx).Preload("User").First(&item, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return item.ToDomain(), nil
}

func (om *OfficeMemberRepositoryImpl) ListByOfficeID(ctx context.Context, officeID uint) ([]entity.OfficeMember, error) {
	var items []model.OfficeMember
	if err := om.GetDB(ctx).Preload("User").Where("office_id = ?", officeID).Find(&items).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	result := make([]entity.OfficeMember, 0, len(items))
	for _, item := range items {
		result = append(result, *item.ToDomain())
	}
	return result, nil
}
