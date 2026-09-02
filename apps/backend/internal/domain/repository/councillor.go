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
	City     string
	State    string
}

type UpdateCouncillorData struct {
	Party string
	City  string
	State string
}

type CouncillorRepository interface {
	Create(ctx context.Context, user CreateCouncillorData) (*entity.Councillor, error)
	FindByUserID(ctx context.Context, userID uint) (*entity.Councillor, error)
	UpdateByUserID(ctx context.Context, userID uint, data UpdateCouncillorData) (*entity.Councillor, error)
}
