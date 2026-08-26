package repository

import (
	"cidadon/internal/domain/entity"
	"context"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}
