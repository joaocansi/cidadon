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

type OfficeMemberRequestRepositoryImpl struct {
	*database.BaseRepository
}

func NewOfficeMemberRequestRepository(baseRepository *database.BaseRepository) *OfficeMemberRequestRepositoryImpl {
	return &OfficeMemberRequestRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (om *OfficeMemberRequestRepositoryImpl) Create(ctx context.Context, officeMemberRequest repository.CreateOfficeMemberRequestData) (*entity.OfficeMemberRequest, error) {
	officeMemberRequestModel := model.OfficeMemberRequest{
		OfficeID: officeMemberRequest.OfficeID,
		Email:    officeMemberRequest.Email,
		Token:    officeMemberRequest.Token,
	}
	if err := om.GetDB(ctx).Create(&officeMemberRequestModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeMemberRequestModel.ToDomain(), nil
}

func (om *OfficeMemberRequestRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.OfficeMemberRequest, error) {
	var officeMemberRequest model.OfficeMemberRequest
	if err := om.GetDB(ctx).Where("token = ?", token).First(&officeMemberRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeMemberRequest.ToDomain(), nil
}

func (om *OfficeMemberRequestRepositoryImpl) Delete(ctx context.Context, id uint) error {
	tx := om.GetDB(ctx).Delete(&model.OfficeMemberRequest{}, id)
	if tx.Error != nil {
		return repository.NewDBError(repository.DBErrorInternal, tx.Error)
	}
	if tx.RowsAffected == 0 {
		return repository.NewDBError(repository.DBErrorNotFound, gorm.ErrRecordNotFound)
	}
	return nil
}
