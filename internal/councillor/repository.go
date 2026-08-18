package councillor

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

func NewCouncillorRepository(baseRepository *database.BaseRepository, logger *zap.SugaredLogger) *RepositoryImpl {
	return &RepositoryImpl{
		BaseRepository: baseRepository,
		Logger:         logger.Named("CouncillorRepository"),
	}
}

func (c *RepositoryImpl) Create(ctx context.Context, data *Councillor) error {
	model := data.ToModel()
	err := c.GetDB(ctx).Create(model).Error
	if err != nil {
		c.Logger.Errorw("failed to create model", "error", err)
		return apperrors.ErrDBInternal
	}
	*data = *model.ToDomain()
	return nil
}
