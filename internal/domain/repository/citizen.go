package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateCitizenData struct {
	Name     string
	Email    string
	Password string
}

type CitizenRepository interface {
	Create(ctx context.Context, citizen CreateCitizenData) (*entity.Citizen, error)
}
