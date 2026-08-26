package repository

import (
	"context"
	"time"
)

type CreateCitizenData struct {
	Name     string
	Email    string
	Password string
}

type CreateCitizenResult struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CitizenRepository interface {
	Create(ctx context.Context, citizen CreateCitizenData) (CreateCitizenResult, error)
}
