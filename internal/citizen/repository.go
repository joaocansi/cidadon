package citizen

import (
	"cidadon/internal/app/database"
	apperrors "cidadon/internal/app/errors"
	"context"

	"go.uber.org/zap"
)

type RepositoryImpl struct {
	*database.BaseRepository
	Logger *zap.SugaredLogger
}

func NewCitizenRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *RepositoryImpl {
	return &RepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("CitizenRepository"),
	}
}

func (r *RepositoryImpl) Create(ctx context.Context, citizen *Citizen) error {
	model := citizen.ToModel()
	err := r.GetDB(ctx).Create(model).Error
	if err != nil {
		r.Logger.Errorw("error while creating citizen", "error", err)
		return apperrors.ErrDBInternal
	}
	return nil
}
