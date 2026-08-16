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
	DB := repo.DB
	if ctx.Value(dbContextKey{}) != nil {
		DB = ctx.Value(dbContextKey{}).(*gorm.DB)
	}
	return DB
}
