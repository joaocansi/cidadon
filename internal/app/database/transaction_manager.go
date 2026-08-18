package database

import (
	"context"

	"gorm.io/gorm"
)

type dbContextKey struct{}

type TransactionManager interface {
	WithTransaction(ctx context.Context, fc func(ctx context.Context) error) error
}

type TransactionManagerImpl struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *TransactionManagerImpl {
	return &TransactionManagerImpl{
		db: db,
	}
}

// WithTransaction starts a transaction and pass the DB instance do context so repositories can access.
// Nested calls reuse the ongoing transaction instead of opening a new one.
func (t *TransactionManagerImpl) WithTransaction(ctx context.Context, fc func(ctx context.Context) error) error {
	if _, ok := ctx.Value(dbContextKey{}).(*gorm.DB); ok {
		return fc(ctx)
	}

	return t.db.Transaction(func(tx *gorm.DB) error {
		return fc(context.WithValue(ctx, dbContextKey{}, tx))
	})
}
