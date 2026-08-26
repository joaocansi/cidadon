package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OfficeMemberRequestRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewOfficeMemberRequestRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *OfficeMemberRequestRepositoryImpl {
	return &OfficeMemberRequestRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("OfficeMemberRequestRepository"),
	}
}

func (om *OfficeMemberRequestRepositoryImpl) Create(ctx context.Context, officeMemberRequest repository.CreateOfficeMemberRequestData) (*entity.OfficeMemberRequest, error) {
	officeMemberRequestModel := model.OfficeMemberRequest{
		OfficeID: officeMemberRequest.OfficeID,
		Email:    officeMemberRequest.Email,
		Token:    officeMemberRequest.Token,
	}
	err := om.GetDB(ctx).Create(&officeMemberRequestModel).Error
	if err != nil {
		om.Logger.Warnw("failed to create office member request", "error", err)
		return nil, repository.ErrDBInternal
	}
	return officeMemberRequestModel.ToDomain(), nil
}

func (om *OfficeMemberRequestRepositoryImpl) FindByToken(ctx context.Context, token string) (*entity.OfficeMemberRequest, error) {
	var officeMemberRequest model.OfficeMemberRequest
	err := om.GetDB(ctx).Where("token = ?", token).First(&officeMemberRequest).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repository.ErrDBNotFound
		}
		om.Logger.Warnw("failed to find office member request", "error", err)
		return nil, repository.ErrDBInternal
	}
	return officeMemberRequest.ToDomain(), nil
}
