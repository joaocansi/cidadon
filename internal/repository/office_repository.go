package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"
	"encoding/json"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OfficeRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewOfficeRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *OfficeRepositoryImpl {
	return &OfficeRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("OfficeRepository"),
	}
}

func (o *OfficeRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Office, error) {
	var officeModel model.Office
	err := o.GetDB(ctx).First(&officeModel, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.WrapDB(repository.ErrDBNotFound, "office not found")
		}
		o.Logger.Warnw("failed to get office by id", "error", err)
		return nil, err
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) FindByCouncillorID(ctx context.Context, councillorID uint) (*entity.Office, error) {
	var officeModel model.Office
	err := o.GetDB(ctx).First(&officeModel, "councillor_id = ?", councillorID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.WrapDB(repository.ErrDBNotFound, "office not found")
		}
		o.Logger.Warnw("failed to get office by id", "error", err)
		return nil, err
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) Create(ctx context.Context, data repository.CreateOfficeData) (*entity.Office, error) {
	officeModel := model.Office{
		CouncillorID:   data.CouncillorID,
		Contacts:       o.toRawMessage(data.Contacts),
		SocialNetworks: o.toRawMessage(data.SocialNetworks),
	}
	err := o.GetDB(ctx).Create(&officeModel).Error
	if err != nil {
		o.Logger.Warnw("failed to create office by id", "error", err)
		return nil, repository.ErrDBInternal
	}
	return officeModel.ToDomain(), nil
}

func (o *OfficeRepositoryImpl) Update(ctx context.Context, data repository.UpdateOfficeData) error {
	officeModel := model.Office{
		SocialNetworks: o.toRawMessage(data.SocialNetworks),
		Contacts:       o.toRawMessage(data.Contacts),
	}
	result := o.GetDB(ctx).Model(&model.Office{}).Where("id = ?", data.OfficeID).Updates(&officeModel)
	if result.Error != nil {
		o.Logger.Warnw("failed to update office by id", "error", result.Error)
		return repository.ErrDBInternal
	}
	if result.RowsAffected == 0 {
		return repository.ErrDBNotFound
	}
	return nil
}

func (o *OfficeRepositoryImpl) toRawMessage(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		o.Logger.Warnw("failed to marshal field to json.RawMessage", "error", err)
		return json.RawMessage("[]")
	}
	return b
}
