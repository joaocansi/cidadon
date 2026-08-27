package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfficeRepositoryImpl struct {
	*database.BaseRepository
}

func NewOfficeRepository(baseRepository *database.BaseRepository) *OfficeRepositoryImpl {
	return &OfficeRepositoryImpl{
		BaseRepository: baseRepository,
	}
}

func (o *OfficeRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Office, error) {
	var officeModel model.Office
	if err := o.GetDB(ctx).First(&officeModel, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) FindByCouncillorID(ctx context.Context, councillorID uint) (*entity.Office, error) {
	var officeModel model.Office
	if err := o.GetDB(ctx).First(&officeModel, "councillor_id = ?", councillorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.NewDBError(repository.DBErrorNotFound, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) Create(ctx context.Context, data repository.CreateOfficeData) (*entity.Office, error) {
	officeModel := model.Office{
		CouncillorID:   data.CouncillorID,
		Contacts:       o.toRawMessage(data.Contacts),
		SocialNetworks: o.toRawMessage(data.SocialNetworks),
	}
	if err := o.GetDB(ctx).Create(&officeModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, repository.NewDBError(repository.DBErrorConflict, err)
		}
		return nil, repository.NewDBError(repository.DBErrorInternal, err)
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) UpdateByCouncillorID(ctx context.Context, councillorID uint, data repository.UpdateOfficeData) (*entity.Office, error) {
	officeModel := model.Office{
		SocialNetworks: o.toRawMessage(data.SocialNetworks),
		Contacts:       o.toRawMessage(data.Contacts),
	}

	result := o.GetDB(ctx).Model(&officeModel).
		Clauses(clause.Returning{}).
		Where("councillor_id = ?", councillorID).
		Updates(&officeModel)

	if result.Error != nil {
		return nil, repository.NewDBError(repository.DBErrorInternal, result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, repository.NewDBError(repository.DBErrorNotFound, fmt.Errorf("office not found"))
	}

	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) toRawMessage(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		return json.RawMessage("[]")
	}
	return b
}
