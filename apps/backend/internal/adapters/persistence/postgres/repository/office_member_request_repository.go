package repository

import (
	"cidadon/internal/adapters/persistence/postgres/model"
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/platform/database"
	"context"
	"errors"
	"time"

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
		OfficeID:  officeMemberRequest.OfficeID,
		Email:     officeMemberRequest.Email,
		Token:     officeMemberRequest.Token,
		ExpiresAt: officeMemberRequest.ExpiresAt,
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
	if err := om.GetDB(ctx).Preload("Office.Councillor.User").Where("token = ? AND expires_at > ?", token, time.Now()).First(&officeMemberRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeMemberRequest.ToDomain(), nil
}

func (om *OfficeMemberRequestRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.OfficeMemberRequest, error) {
	var item model.OfficeMemberRequest
	if err := om.GetDB(ctx).First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return item.ToDomain(), nil
}

func (om *OfficeMemberRequestRepositoryImpl) ListByOfficeID(ctx context.Context, officeID uint) ([]entity.OfficeMemberRequest, error) {
	var items []model.OfficeMemberRequest
	if err := om.GetDB(ctx).Where("office_id = ? AND expires_at > ?", officeID, time.Now()).Order("created_at desc").Find(&items).Error; err != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	result := make([]entity.OfficeMemberRequest, 0, len(items))
	for _, item := range items {
		result = append(result, *item.ToDomain())
	}
	return result, nil
}

func (om *OfficeMemberRequestRepositoryImpl) FindByEmailAndOffice(ctx context.Context, email string, officeID uint) (*entity.OfficeMemberRequest, error) {
	var item model.OfficeMemberRequest
	if err := om.GetDB(ctx).Where("email = ? AND office_id = ? AND expires_at > ?", email, officeID, time.Now()).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return item.ToDomain(), nil
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
