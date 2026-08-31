package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type CreateUserData struct {
	Name     string
	Email    string
	Password string
	Role     string
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	Create(ctx context.Context, user CreateUserData) (*entity.User, error)
	DeleteByID(ctx context.Context, id uint) error
}
