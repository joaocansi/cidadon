package repository

import (
	"cidadon/internal/domain/entity"
	"cidadon/internal/domain/repository"
	"cidadon/internal/infrastructure/database"
	"cidadon/internal/model"
	"context"

	"go.uber.org/zap"
)

type OfficeMemberRepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewOfficeMemberRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *OfficeMemberRepositoryImpl {
	return &OfficeMemberRepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("OfficeMemberRepository"),
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
	err := om.GetDB(ctx).Create(&officeMemberModel).Error
	if err != nil {
		om.Logger.Warnw("failed to create office member", "error", err)
		return nil, repository.ErrDBInternal
	}
	return officeMemberModel.ToDomain(), nil
}
