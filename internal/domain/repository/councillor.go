package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CouncillorRepository interface {
	Create(ctx context.Context, user *entity.Councillor) error
}
