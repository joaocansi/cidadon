package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateCouncillorData struct {
	Name     string
	Email    string
	Password string
	Party    string
	ImageURL string
}

type CouncillorRepository interface {
	Create(ctx context.Context, user CreateCouncillorData) (*entity.Councillor, error)
}
