package database

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{
		DB: db,
	}
}

func (repo *BaseRepository) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(dbContextKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return repo.DB.WithContext(ctx)
}
